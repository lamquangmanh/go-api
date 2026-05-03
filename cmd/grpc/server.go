package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-api/internal/config"
	appdb "go-api/internal/db"
	appgrpc "go-api/internal/grpc"
	applogger "go-api/internal/logger"
	"go-api/internal/repository"
	appsvc "go-api/internal/service"
	actionpb "go-api/pkg/api/actionpb"
	authpb "go-api/pkg/api/authpb"
	menupb "go-api/pkg/api/menupb"
	modulepb "go-api/pkg/api/modulepb"
	permissionpb "go-api/pkg/api/permissionpb"
	productpb "go-api/pkg/api/productpb"
	resourcepb "go-api/pkg/api/resourcepb"
	rolepb "go-api/pkg/api/rolepb"
	userpb "go-api/pkg/api/userpb"

	"google.golang.org/grpc"
	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}

func run() error {
	// Load configuration from configs/config.yml (path relative to working directory)
	// Fall back to GOAPI_CONFIG env var for overriding the path in deployment
	cfgPath := os.Getenv("GOAPI_CONFIG")
	if cfgPath == "" {
		cfgPath = "configs/config.yml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Initialize structured application logger from config.
	logger, err := applogger.New(cfg.Logger)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	logger.Info("config loaded", "config_path", cfgPath)

	// Connect to PostgreSQL using the DSN built from config
	ctx := context.Background()
	pool, err := appdb.NewPool(ctx, cfg.Database.DSN())
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	logger.Info("database connected", "host", cfg.Database.Host, "port", cfg.Database.Port, "name", cfg.Database.Name)

	// Initialize the sqlc query layer bound to the connection pool
	queries := repository.New(pool)

	// Create gRPC listener on configured port
	grpcAddr := fmt.Sprintf(":%d", cfg.GRPC.Port)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	// Create gRPC server with a unary logging interceptor.
	// Every unary request will emit a structured log entry.
	s := grpc.NewServer(grpc.UnaryInterceptor(applogger.UnaryServerInterceptor(logger)))

	// Wire up UserService: business logic layer + thin gRPC transport layer.
	uss := appsvc.NewUserService(queries, pool)
	us := appgrpc.NewUserHandler(uss)
	userpb.RegisterUserServiceServer(s, us)

	// Wire up RoleService with separated business/service layer.
	rss := appsvc.NewRoleService(queries)
	rs := appgrpc.NewRoleHandler(rss)
	rolepb.RegisterRoleServiceServer(s, rs)

	// Wire up AuthService: authentication and token management
	as := appgrpc.NewAuthHandler()
	authpb.RegisterAuthServiceServer(s, as)

	// Wire up MenuService: hierarchical menu structure for UI
	ms := appgrpc.NewMenuHandler()
	menupb.RegisterMenuServiceServer(s, ms)

	// Wire up ResourceService: API resource and action management
	ress := appsvc.NewResourceService(queries, pool)
	res := appgrpc.NewResourceHandler(ress)
	resourcepb.RegisterResourceServiceServer(s, res)

	// Wire up ProductService
	pros := appsvc.NewProductService(queries)
	pro := appgrpc.NewProductHandler(pros)
	productpb.RegisterProductServiceServer(s, pro)

	// Wire up ModuleService
	mods := appsvc.NewModuleService(queries)
	mod := appgrpc.NewModuleHandler(mods, pros)
	modulepb.RegisterModuleServiceServer(s, mod)

	// Wire up ActionService
	acts := appsvc.NewActionService(queries)
	act := appgrpc.NewActionHandler(acts)
	actionpb.RegisterActionServiceServer(s, act)

	// Wire up PermissionService: role-based permission management
	ps := appgrpc.NewPermissionHandler()
	permissionpb.RegisterPermissionServiceServer(s, ps)

	// Register gRPC health service (used by load balancers and k8s probes)
	hs := health.NewServer()
	healthpb.RegisterHealthServer(s, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	hs.SetServingStatus(userpb.UserService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	hs.SetServingStatus(rolepb.RoleService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	hs.SetServingStatus(authpb.AuthService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	hs.SetServingStatus(menupb.MenuService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	hs.SetServingStatus(resourcepb.ResourceService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	hs.SetServingStatus(permissionpb.PermissionService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	hs.SetServingStatus(productpb.ProductService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	hs.SetServingStatus(modulepb.ModuleService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	hs.SetServingStatus(actionpb.ActionService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)

	// Register reflection service (enables dynamic service discovery via grpcurl, Postman, etc.)
	reflection.Register(s)

	// Start HTTP server for Kubernetes probes on configured port.
	// - /livez: process liveness only
	// - /readyz: dependency readiness (database ping)
	// - /healthz: full health (currently same checks as /readyz)
	httpAddr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	httpSrv := &http.Server{
		Addr: httpAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writeProbe := func(statusCode int, body string) {
				w.WriteHeader(statusCode)
				_, _ = w.Write([]byte(body))
			}

			switch r.URL.Path {
			case "/livez":
				writeProbe(http.StatusOK, "live")
				return
			case "/readyz", "/healthz":
				checkCtx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
				defer cancel()
				if err := pool.Ping(checkCtx); err != nil {
					writeProbe(http.StatusServiceUnavailable, "not ready")
					return
				}
				writeProbe(http.StatusOK, "ok")
				return
			default:
				http.NotFound(w, r)
			}
		}),
	}

	go func() {
		logger.Info("http health server started", "address", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("http server error", "error", err)
		}
	}()

	// Serve gRPC in a goroutine so we can handle shutdown signals
	grpcErrCh := make(chan error, 1)
	go func() {
		logger.Info("grpc server started", "address", lis.Addr().String())
		grpcErrCh <- s.Serve(lis)
	}()

	// Handle OS shutdown signals (SIGINT, SIGTERM) for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-grpcErrCh:
		return fmt.Errorf("grpc server stopped: %w", err)
	case sig := <-sigCh:
		logger.Info("shutdown signal received", "signal", sig.String())
		// Gracefully stop gRPC (waits for in-flight RPCs to complete)
		s.GracefulStop()
		// Gracefully shut down HTTP server with a 5s timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpSrv.Shutdown(shutdownCtx)
		return nil
	}
}
