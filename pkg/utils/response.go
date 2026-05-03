package utils

import (
	basepb "go-api/pkg/api/basepb"
	"go-api/pkg/constants"

	"google.golang.org/protobuf/types/known/structpb"
)

// ErrMsg builds a single ErrorMessage from an ErrorDef constant.
func ErrMsg(def constants.ErrorDef, extra ...map[string]any) *basepb.ErrorMessage {
	var st *structpb.Struct
	if len(extra) > 0 && extra[0] != nil {
		s, err := structpb.NewStruct(extra[0])
		if err == nil {
			st = s
		}
	}

	return &basepb.ErrorMessage{
		Code:    int32(def.Code),
		Message: def.Message,
		Extra:   st,
	}
}

// ErrMsgFromErr builds an ErrorMessage by extracting the gRPC status message from an error.
// func ErrMsgFromErr(err error, fallback constants.ErrorDef) *basepb.ErrorMessage {
// 	if err == nil {
// 		return ErrMsg(fallback, nil)
// 	}
// 	st, ok := status.FromError(err)
// 	if ok {
// 		return &basepb.ErrorMessage{Code: int32(fallback.Code), Message: st.Message()}
// 	}
// 	return &basepb.ErrorMessage{Code: int32(fallback.Code), Message: err.Error()}
// }

// ErrMsgsFromValidation converts service validation error maps into ErrorMessage slice.
// func ErrMsgsFromValidation(valErrs []map[string]string) []*basepb.ErrorMessage {
// 	items := make([]*basepb.ErrorMessage, 0, len(valErrs))
// 	for _, ve := range valErrs {
// 		code := int32(0)
// 		if c, ok := ve["code"]; ok {
// 			for _, ch := range c {
// 				code = code*10 + int32(ch-'0')
// 			}
// 		}
// 		items = append(items, &basepb.ErrorMessage{Code: code, Message: ve["error"]})
// 	}
// 	return items
// }

// ErrorsList returns a repeated ErrorMessage field value.
// func ErrorsList(msgs ...*basepb.ErrorMessage) []*basepb.ErrorMessage {
// 	return msgs
// }

// UpdateOK returns a successful UpdateSuccess response.
func UpdateOK() *basepb.UpdateSuccess {
	return &basepb.UpdateSuccess{Success: true}
}

// UpdateErr returns an UpdateSuccess carrying error messages in the body.
func UpdateErr(msgs ...*basepb.ErrorMessage) *basepb.UpdateSuccess {
	return &basepb.UpdateSuccess{Errors: msgs}
}

// DeleteOK returns a successful DeleteSuccess response.
func DeleteOK() *basepb.DeleteSuccess {
	return &basepb.DeleteSuccess{Success: true}
}

// DeleteErr returns a DeleteSuccess carrying error messages in the body.
func DeleteErr(msgs ...*basepb.ErrorMessage) *basepb.DeleteSuccess {
	return &basepb.DeleteSuccess{Errors: msgs}
}
