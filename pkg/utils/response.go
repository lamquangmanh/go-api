package utils

import (
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	"go-api/pkg/constants"

	errorv1 "github.com/lamquangmanh/protobuf/gen/go/proto/error/v1"
)

// ErrMsg builds a single ErrorItem from an ErrorDef constant.
func ErrMsg(def constants.ErrorDef, extra ...map[string]any) *errorv1.ErrorItem {
	var st *structpb.Struct
	if len(extra) > 0 && extra[0] != nil {
		s, err := structpb.NewStruct(extra[0])
		if err == nil {
			st = s
		}
	}

	return &errorv1.ErrorItem{
		Code:    int32(def.Code),
		Message: def.Message,
		ExtraData:   st,
	}
}

func ResponseError(code codes.Code, errors []*errorv1.ErrorItem) (any, error) {
	errList := &errorv1.ErrorList{
		Errors: errors,
	}
	data, err := json.Marshal(errList)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshal error response")
	}
	return nil, status.Error(code, string(data))
}
