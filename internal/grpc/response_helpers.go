package grpc

import (
	basepb "go-api/pkg/api/basepb"
	"go-api/pkg/constants"

	"google.golang.org/grpc/status"
)

// errMsg builds a single ErrorMessage from an ErrorDef constant.
func errMsg(def constants.ErrorDef) *basepb.ErrorMessage {
	return &basepb.ErrorMessage{Code: int32(def.Code), Message: def.Message}
}

// errMsgFromErr builds an ErrorMessage by extracting the gRPC status message from an error.
func errMsgFromErr(err error, fallback constants.ErrorDef) *basepb.ErrorMessage {
	if err == nil {
		return errMsg(fallback)
	}
	st, ok := status.FromError(err)
	if ok {
		return &basepb.ErrorMessage{Code: int32(fallback.Code), Message: st.Message()}
	}
	return &basepb.ErrorMessage{Code: int32(fallback.Code), Message: err.Error()}
}

// errMsgsFromValidation converts service validation error maps into ErrorMessage slice.
func errMsgsFromValidation(valErrs []map[string]string) []*basepb.ErrorMessage {
	items := make([]*basepb.ErrorMessage, 0, len(valErrs))
	for _, ve := range valErrs {
		code := int32(0)
		if c, ok := ve["code"]; ok {
			for _, ch := range c {
				code = code*10 + int32(ch-'0')
			}
		}
		items = append(items, &basepb.ErrorMessage{Code: code, Message: ve["error"]})
	}
	return items
}

// errorsList returns a repeated ErrorMessage field value.
func errorsList(msgs ...*basepb.ErrorMessage) []*basepb.ErrorMessage {
	return msgs
}

// updateOK returns a successful UpdateSuccess response.
func updateOK() *basepb.UpdateSuccess {
	return &basepb.UpdateSuccess{Success: true}
}

// updateErr returns an UpdateSuccess carrying error messages in the body.
func updateErr(msgs ...*basepb.ErrorMessage) *basepb.UpdateSuccess {
	return &basepb.UpdateSuccess{Errors: msgs}
}

// deleteOK returns a successful DeleteSuccess response.
func deleteOK() *basepb.DeleteSuccess {
	return &basepb.DeleteSuccess{Success: true}
}

// deleteErr returns a DeleteSuccess carrying error messages in the body.
func deleteErr(msgs ...*basepb.ErrorMessage) *basepb.DeleteSuccess {
	return &basepb.DeleteSuccess{Errors: msgs}
}
