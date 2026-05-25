package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"go-api/internal/config"
	"go-api/internal/repository"
	"go-api/pkg/constants"
	"go-api/pkg/utils"

	authv1 "github.com/lamquangmanh/protobuf/gen/go/proto/auth/v1"
	errorv1 "github.com/lamquangmanh/protobuf/gen/go/proto/error/v1"
)

type AuthService struct {
	q *repository.Queries
	cfg *config.Config
}

// NewAuthService creates a business service for auth operations.
func NewAuthService(q *repository.Queries, cfg *config.Config) *AuthService {
	return &AuthService{q: q, cfg: cfg}
}

type LoginInput struct {
	Email    string
	Password string
}

type authTokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"type"`
	Issued int64  `json:"iat"`
	Expiry int64  `json:"exp"`
}

func (s *AuthService) tokenSecret() []byte {
	return []byte(s.cfg.JWT.Secret)
}

func (s *AuthService) verifyToken(tokenStr string) (*authTokenClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	// Reconstruct unsigned token
	unsignedToken := parts[0] + "." + parts[1]

	// Decode and verify signature
	signatureBytes, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.New("invalid token signature encoding")
	}

	mac := hmac.New(sha256.New, s.tokenSecret())
	_, _ = mac.Write([]byte(unsignedToken))
	expectedSignature := mac.Sum(nil)

	if !hmac.Equal(signatureBytes, expectedSignature) {
		return nil, errors.New("invalid token signature")
	}

	// Decode payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid token payload encoding")
	}

	var claims authTokenClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("invalid token payload")
	}

	return &claims, nil
}

func (s *AuthService) issueToken(userID, email, tokenType string, ttl time.Duration) (string, error) {
	claims := authTokenClaims{
		UserID: userID,
		Email:  email,
		Type:   tokenType,
		Issued: time.Now().Unix(),
		Expiry: time.Now().Add(ttl).Unix(),
	}

	headerJSON, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", err
	}
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadJSON)
	unsignedToken := encodedHeader + "." + encodedPayload

	mac := hmac.New(sha256.New, s.tokenSecret())
	_, _ = mac.Write([]byte(unsignedToken))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return unsignedToken + "." + signature, nil
}

func (s *AuthService) findActiveUserByEmail(ctx context.Context, email string) (*repository.User, []*errorv1.ErrorItem) {
	if s.q == nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrAuthInternal.Messagef("queries are not configured"))}
	}

	user, err := s.q.GetUserByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrUserNotFound.Messagef(email))}
		}
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrAuthInternal.Messagef(err.Error()))}
	}

	if user.Status != repository.UserStatusACTIVE {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrUserNotFound.Messagef(email))}
	}

	return user, nil
}

// Login validates credentials and returns access/refresh tokens.
func (s *AuthService) Login(ctx context.Context, in LoginInput) (*authv1.LoginResponse, []*errorv1.ErrorItem) {
	email := strings.TrimSpace(in.Email)
	password := strings.TrimSpace(in.Password)

	if !utils.ValidateEmail(email) {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrEmailInvalid)}
	}
	if password == "" {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrPasswordRequired)}
	}

	user, err := s.findActiveUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidEmailOrPassword)}
	}

	accessToken, tokenErr := s.issueToken(user.UserID.String(), user.Email, "access", 15*time.Minute)
	if tokenErr != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrAuthInternal.Messagef(tokenErr.Error()))}
	}
	refreshToken, tokenErr := s.issueToken(user.UserID.String(), user.Email, "refresh", 7*24*time.Hour)
	if tokenErr != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrAuthInternal.Messagef(tokenErr.Error()))}
	}

	return &authv1.LoginResponse{
		Auth: &authv1.Auth{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		Errors: nil,
	}, nil
}

// GetMe retrieves the current user information by user ID.
func (s *AuthService) GetMe(ctx context.Context, userID string) (*authv1.User, []*errorv1.ErrorItem) {
	if s.q == nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrAuthInternal.Messagef("queries are not configured"))}
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrAuthInternal.Messagef("invalid user id format"))}
	}

	user, err := s.q.GetUser(ctx, parsedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrUserNotFound.Messagef(userID))}
		}
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrAuthInternal.Messagef(err.Error()))}
	}

	if user.Status != repository.UserStatusACTIVE {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrUserNotFound.Messagef(userID))}
	}

	var avatar *string
	if user.Avatar.Valid {
		avatar = &user.Avatar.String
	}
	var phone *string
	if user.Phone.Valid {
		phone = &user.Phone.String
	}

	return &authv1.User{
		UserId:   user.UserID.String(),
		Email:    user.Email,
		Username: user.UserName,
		Avatar:   avatar,
		Phone:    phone,
		Status:   string(user.Status),
	}, nil
}

// Verify validates a JWT token and returns user information if valid.
func (s *AuthService) Verify(ctx context.Context, token string) (*authv1.User, []*errorv1.ErrorItem) {
	if s.q == nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrAuthInternal.Messagef("queries are not configured"))}
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrAuthInternal.Messagef("token is required"))}
	}

	// Parse and verify token
	claims, verifyErr := s.verifyToken(token)
	if verifyErr != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrAuthInternal.Messagef(verifyErr.Error()))}
	}

	// Check token expiration
	now := time.Now().Unix()
	if claims.Expiry < now {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrAuthInternal.Messagef("token expired"))}
	}

	// Get user from database
	parsedID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrAuthInternal.Messagef("invalid user id in token"))}
	}

	user, err := s.q.GetUser(ctx, parsedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrUserNotFound.Messagef(claims.UserID))}
		}
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrAuthInternal.Messagef(err.Error()))}
	}

	if user.Status != repository.UserStatusACTIVE {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrUserNotFound.Messagef(claims.UserID))}
	}

	var avatar *string
	if user.Avatar.Valid {
		avatar = &user.Avatar.String
	}
	var phone *string
	if user.Phone.Valid {
		phone = &user.Phone.String
	}

	return &authv1.User{
		UserId:   user.UserID.String(),
		Email:    user.Email,
		Username: user.UserName,
		Avatar:   avatar,
		Phone:    phone,
		Status:   string(user.Status),
	}, nil
}
