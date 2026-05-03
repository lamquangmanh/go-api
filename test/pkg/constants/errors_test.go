package constants_test

import (
	"testing"

	"go-api/pkg/constants"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorDefStatusHelpers(t *testing.T) {
	e := constants.ErrorDef{Code: 1, GrpcCode: codes.InvalidArgument, Message: "hello %s"}

	err := e.Status()
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("unexpected status code: %v", status.Code(err))
	}
	if status.Convert(err).Message() != "hello %s" {
		t.Fatalf("unexpected message: %s", status.Convert(err).Message())
	}

	err = e.Statusf("world")
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("unexpected status code: %v", status.Code(err))
	}
	if status.Convert(err).Message() != "hello world" {
		t.Fatalf("unexpected formatted message: %s", status.Convert(err).Message())
	}

	if got := e.Messagef("gopher"); got.Message != "hello gopher" {
		t.Fatalf("unexpected Messagef result: %s", got.Message)
	}
}
