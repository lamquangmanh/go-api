package constants

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorDef defines a reusable gRPC error with a code and a default message.
type ErrorDef struct {
	Code     int32
	GrpcCode codes.Code
	Message  string
	Extra    any // Optional field for additional context
}

// Status returns the error as a gRPC status error.
func (e ErrorDef) Status() error {
	return status.Error(e.GrpcCode, e.Message)
}

// Statusf returns a gRPC status error using Message as a printf-style template.
func (e ErrorDef) Statusf(args ...any) error {
	return status.Errorf(e.GrpcCode, e.Message, args...)
}

func (e ErrorDef) Messagef(args ...any) ErrorDef {
	e.Message = fmt.Sprintf(e.Message, args...)
	return e
}

var (
	// Errors of Product Service
	ErrProductIDRequired          = ErrorDef{Code: 33, GrpcCode: codes.InvalidArgument, Message: "product_id is required"}
	ErrInvalidProductIDFormat     = ErrorDef{Code: 34, GrpcCode: codes.InvalidArgument, Message: "invalid product id format"}
	ErrProductNotFound            = ErrorDef{Code: 36, GrpcCode: codes.NotFound, Message: "product %s not found"}
	ErrGetProductInternal         = ErrorDef{Code: 38, GrpcCode: codes.Internal, Message: "get product: %v"}
	ErrCreateProductInternal      = ErrorDef{Code: 37, GrpcCode: codes.Internal, Message: "create product: %v"}
	ErrProductPayloadRequired     = ErrorDef{Code: 97, GrpcCode: codes.InvalidArgument, Message: "product is required"}
	ErrUpdateProductInternal      = ErrorDef{Code: 39, GrpcCode: codes.Internal, Message: "update product: %v"}
	ErrProductNameRequired        = ErrorDef{Code: 35, GrpcCode: codes.InvalidArgument, Message: "name is required"}
	ErrDeleteProductInternal      = ErrorDef{Code: 40, GrpcCode: codes.Internal, Message: "delete product: %v"}
	ErrListProductsInternal       = ErrorDef{Code: 41, GrpcCode: codes.Internal, Message: "list products: %v"}
	ErrMissingProductsQueryConfig = ErrorDef{Code: 42, GrpcCode: codes.Internal, Message: "missing query config for products"}
	ErrInvalidProductSort         = ErrorDef{Code: 43, GrpcCode: codes.InvalidArgument, Message: "%v"}
	ErrInvalidProductFilter       = ErrorDef{Code: 44, GrpcCode: codes.InvalidArgument, Message: "%v"}

	// ------------

	// Errors of User Service
	ErrUsernameRequired      = ErrorDef{Code: 1, GrpcCode: codes.InvalidArgument, Message: "username is required"}
	ErrUsernameInvalid       = ErrorDef{Code: 2, GrpcCode: codes.InvalidArgument, Message: "username must be between 3 and 50 characters and contain only letters, numbers, '_', '.', '-'"}
	ErrEmailRequired         = ErrorDef{Code: 3, GrpcCode: codes.InvalidArgument, Message: "email is required"}
	ErrEmailInvalid          = ErrorDef{Code: 4, GrpcCode: codes.InvalidArgument, Message: "email is invalid or exceeds 254 characters"}
	ErrUserIDRequired        = ErrorDef{Code: 5, GrpcCode: codes.InvalidArgument, Message: "id is required"}
	ErrInvalidUserIDFormat   = ErrorDef{Code: 6, GrpcCode: codes.InvalidArgument, Message: "invalid user id format"}
	ErrEmailAlreadyExists    = ErrorDef{Code: 7, GrpcCode: codes.AlreadyExists, Message: "email already exists"}
	ErrUsernameAlreadyExists = ErrorDef{Code: 8, GrpcCode: codes.AlreadyExists, Message: "username already exists"}
	ErrUserNotFound          = ErrorDef{Code: 9, GrpcCode: codes.NotFound, Message: "user %s not found"}
	ErrInvalidRoleIDFormatIn = ErrorDef{Code: 10, GrpcCode: codes.InvalidArgument, Message: "invalid role_id format: %s"}
	ErrBeginTransaction      = ErrorDef{Code: 11, GrpcCode: codes.Internal, Message: "begin transaction: %v"}
	ErrValidateRoles         = ErrorDef{Code: 12, GrpcCode: codes.Internal, Message: "validate roles: %v"}
	ErrCreateUserInternal    = ErrorDef{Code: 13, GrpcCode: codes.Internal, Message: "create user: %v"}
	ErrAddUserRoleMapping    = ErrorDef{Code: 14, GrpcCode: codes.Internal, Message: "add user role mapping: %v"}
	ErrCommitTransaction     = ErrorDef{Code: 15, GrpcCode: codes.Internal, Message: "commit transaction: %v"}
	ErrGetUserInternal       = ErrorDef{Code: 16, GrpcCode: codes.Internal, Message: "get user: %v"}
	ErrUpdateUserInternal    = ErrorDef{Code: 17, GrpcCode: codes.Internal, Message: "update user: %v"}
	ErrClearUserRoleMappings = ErrorDef{Code: 18, GrpcCode: codes.Internal, Message: "clear user role mappings: %v"}
	ErrDeleteUserInternal    = ErrorDef{Code: 19, GrpcCode: codes.Internal, Message: "delete user: %v"}
	ErrListUsersInternal     = ErrorDef{Code: 20, GrpcCode: codes.Internal, Message: "list users: %v"}

	// Errors of Role Service
	ErrRoleIDRequired        = ErrorDef{Code: 21, GrpcCode: codes.InvalidArgument, Message: "id is required"}
	ErrInvalidRoleIDFormat   = ErrorDef{Code: 22, GrpcCode: codes.InvalidArgument, Message: "invalid role id format"}
	ErrRoleNameRequired      = ErrorDef{Code: 23, GrpcCode: codes.InvalidArgument, Message: "name is required"}
	ErrRoleNameInvalid       = ErrorDef{Code: 24, GrpcCode: codes.InvalidArgument, Message: "name must be between 2 and 100 characters"}
	ErrRoleIDsNotExist       = ErrorDef{Code: 25, GrpcCode: codes.InvalidArgument, Message: "one or more role_ids do not exist"}
	ErrRoleNameAlreadyExists = ErrorDef{Code: 26, GrpcCode: codes.AlreadyExists, Message: "role name already exists"}
	ErrRoleNotFound          = ErrorDef{Code: 27, GrpcCode: codes.NotFound, Message: "role %s not found"}
	ErrCreateRoleInternal    = ErrorDef{Code: 28, GrpcCode: codes.Internal, Message: "create role: %v"}
	ErrGetRoleInternal       = ErrorDef{Code: 29, GrpcCode: codes.Internal, Message: "get role: %v"}
	ErrUpdateRoleInternal    = ErrorDef{Code: 30, GrpcCode: codes.Internal, Message: "update role: %v"}
	ErrDeleteRoleInternal    = ErrorDef{Code: 31, GrpcCode: codes.Internal, Message: "delete role: %v"}
	ErrListRolesInternal     = ErrorDef{Code: 32, GrpcCode: codes.Internal, Message: "list roles: %v"}

	// Errors of Module Service
	ErrModuleIDRequired          = ErrorDef{Code: 45, GrpcCode: codes.InvalidArgument, Message: "id is required"}
	ErrInvalidModuleIDFormat     = ErrorDef{Code: 46, GrpcCode: codes.InvalidArgument, Message: "invalid module id format"}
	ErrInvalidModuleProductID    = ErrorDef{Code: 47, GrpcCode: codes.InvalidArgument, Message: "invalid product_id format"}
	ErrModuleNameRequired        = ErrorDef{Code: 48, GrpcCode: codes.InvalidArgument, Message: "name is required"}
	ErrModuleNotFound            = ErrorDef{Code: 49, GrpcCode: codes.NotFound, Message: "module %s not found"}
	ErrCreateModuleInternal      = ErrorDef{Code: 50, GrpcCode: codes.Internal, Message: "create module: %v"}
	ErrGetModuleInternal         = ErrorDef{Code: 51, GrpcCode: codes.Internal, Message: "get module: %v"}
	ErrUpdateModuleInternal      = ErrorDef{Code: 52, GrpcCode: codes.Internal, Message: "update module: %v"}
	ErrDeleteModuleInternal      = ErrorDef{Code: 53, GrpcCode: codes.Internal, Message: "delete module: %v"}
	ErrListModulesInternal       = ErrorDef{Code: 54, GrpcCode: codes.Internal, Message: "list modules: %v"}
	ErrMissingModulesQueryConfig = ErrorDef{Code: 55, GrpcCode: codes.Internal, Message: "missing query config for modules"}
	ErrInvalidModuleSort         = ErrorDef{Code: 56, GrpcCode: codes.InvalidArgument, Message: "%v"}
	ErrInvalidModuleFilter       = ErrorDef{Code: 57, GrpcCode: codes.InvalidArgument, Message: "%v"}

	// Errors of Action Service
	ErrActionIDRequired          = ErrorDef{Code: 58, GrpcCode: codes.InvalidArgument, Message: "id is required"}
	ErrInvalidActionIDFormat     = ErrorDef{Code: 59, GrpcCode: codes.InvalidArgument, Message: "invalid action id format"}
	ErrInvalidActionResourceID   = ErrorDef{Code: 60, GrpcCode: codes.InvalidArgument, Message: "invalid resource_id format"}
	ErrActionNameRequired        = ErrorDef{Code: 61, GrpcCode: codes.InvalidArgument, Message: "name is required"}
	ErrActionNotFound            = ErrorDef{Code: 62, GrpcCode: codes.NotFound, Message: "action %s not found"}
	ErrCreateActionInternal      = ErrorDef{Code: 63, GrpcCode: codes.Internal, Message: "create action: %v"}
	ErrGetActionInternal         = ErrorDef{Code: 64, GrpcCode: codes.Internal, Message: "get action: %v"}
	ErrUpdateActionInternal      = ErrorDef{Code: 65, GrpcCode: codes.Internal, Message: "update action: %v"}
	ErrDeleteActionInternal      = ErrorDef{Code: 66, GrpcCode: codes.Internal, Message: "delete action: %v"}
	ErrListActionsInternal       = ErrorDef{Code: 67, GrpcCode: codes.Internal, Message: "list actions: %v"}
	ErrMissingActionsQueryConfig = ErrorDef{Code: 68, GrpcCode: codes.Internal, Message: "missing query config for actions"}
	ErrInvalidActionSort         = ErrorDef{Code: 69, GrpcCode: codes.InvalidArgument, Message: "%v"}
	ErrInvalidActionFilter       = ErrorDef{Code: 70, GrpcCode: codes.InvalidArgument, Message: "%v"}

	// Errors of Resource Service
	ErrResourceIDRequired           = ErrorDef{Code: 71, GrpcCode: codes.InvalidArgument, Message: "id is required"}
	ErrInvalidResourceIDFormat      = ErrorDef{Code: 72, GrpcCode: codes.InvalidArgument, Message: "invalid resource id format"}
	ErrInvalidResourceModuleID      = ErrorDef{Code: 73, GrpcCode: codes.InvalidArgument, Message: "invalid module_id format"}
	ErrResourceNameRequired         = ErrorDef{Code: 74, GrpcCode: codes.InvalidArgument, Message: "name is required"}
	ErrResourceNotFound             = ErrorDef{Code: 75, GrpcCode: codes.NotFound, Message: "resource %s not found"}
	ErrCreateResourceInternal       = ErrorDef{Code: 76, GrpcCode: codes.Internal, Message: "create resource: %v"}
	ErrGetResourceInternal          = ErrorDef{Code: 77, GrpcCode: codes.Internal, Message: "get resource: %v"}
	ErrUpdateResourceInternal       = ErrorDef{Code: 78, GrpcCode: codes.Internal, Message: "update resource: %v"}
	ErrDeleteResourceInternal       = ErrorDef{Code: 79, GrpcCode: codes.Internal, Message: "delete resource: %v"}
	ErrListResourcesInternal        = ErrorDef{Code: 80, GrpcCode: codes.Internal, Message: "list resources: %v"}
	ErrCreateResourceActionInternal = ErrorDef{Code: 81, GrpcCode: codes.Internal, Message: "create action for resource: %v"}
	ErrListOldActionsInternal       = ErrorDef{Code: 82, GrpcCode: codes.Internal, Message: "list old actions: %v"}
	ErrDeleteOldActionInternal      = ErrorDef{Code: 83, GrpcCode: codes.Internal, Message: "delete old action: %v"}
	ErrInsertActionInternal         = ErrorDef{Code: 84, GrpcCode: codes.Internal, Message: "insert action: %v"}
	ErrMissingResourcesQueryConfig  = ErrorDef{Code: 85, GrpcCode: codes.Internal, Message: "missing query config for resources"}
	ErrInvalidResourceSort          = ErrorDef{Code: 86, GrpcCode: codes.InvalidArgument, Message: "%v"}
	ErrInvalidResourceFilter        = ErrorDef{Code: 87, GrpcCode: codes.InvalidArgument, Message: "%v"}

	// Additional Errors of Role Service
	ErrMissingRolesQueryConfig = ErrorDef{Code: 88, GrpcCode: codes.Internal, Message: "missing query config for roles"}
	ErrInvalidRoleSort         = ErrorDef{Code: 89, GrpcCode: codes.InvalidArgument, Message: "%v"}
	ErrInvalidRoleFilter       = ErrorDef{Code: 90, GrpcCode: codes.InvalidArgument, Message: "%v"}

	// Additional Errors of User Service
	ErrMissingUsersQueryConfig = ErrorDef{Code: 91, GrpcCode: codes.Internal, Message: "missing query config for users"}
	ErrInvalidUserSort         = ErrorDef{Code: 92, GrpcCode: codes.InvalidArgument, Message: "%v"}
	ErrInvalidUserFilter       = ErrorDef{Code: 93, GrpcCode: codes.InvalidArgument, Message: "%v"}
	ErrPasswordRequired        = ErrorDef{Code: 94, GrpcCode: codes.InvalidArgument, Message: "password is required"}

	// Errors of gRPC Handler payload checks
	ErrUserPayloadRequired = ErrorDef{Code: 95, GrpcCode: codes.InvalidArgument, Message: "user is required"}
	ErrRolePayloadRequired = ErrorDef{Code: 96, GrpcCode: codes.InvalidArgument, Message: "role is required"}

	ErrModulePayloadRequired   = ErrorDef{Code: 98, GrpcCode: codes.InvalidArgument, Message: "module is required"}
	ErrResourcePayloadRequired = ErrorDef{Code: 99, GrpcCode: codes.InvalidArgument, Message: "resource is required"}
	ErrActionPayloadRequired   = ErrorDef{Code: 100, GrpcCode: codes.InvalidArgument, Message: "action is required"}
)

// var (
// 	// Errors of User Service
// 	ErrUsernameRequired      = ErrorDef{Code: 1, GrpcCode: codes.InvalidArgument, Message: "username is required"}
// 	ErrUsernameInvalid       = ErrorDef{Code: 2, GrpcCode: codes.InvalidArgument, Message: "username must be between 3 and 50 characters and contain only letters, numbers, '_', '.', '-'"}
// 	ErrEmailRequired         = ErrorDef{Code: 3, GrpcCode: codes.InvalidArgument, Message: "email is required"}
// 	ErrEmailInvalid          = ErrorDef{Code: 4, GrpcCode: codes.InvalidArgument, Message: "email is invalid or exceeds 254 characters"}
// 	ErrUserIDRequired        = ErrorDef{Code: 5, GrpcCode: codes.InvalidArgument, Message: "id is required"}
// 	ErrInvalidUserIDFormat   = ErrorDef{Code: 6, GrpcCode: codes.InvalidArgument, Message: "invalid user id format"}
// 	ErrEmailAlreadyExists    = ErrorDef{Code: 7, GrpcCode: codes.AlreadyExists, Message: "email already exists"}
// 	ErrUsernameAlreadyExists = ErrorDef{Code: 8, GrpcCode: codes.AlreadyExists, Message: "username already exists"}
// 	ErrUserNotFound          = ErrorDef{Code: 9, GrpcCode: codes.NotFound, Message: "user %s not found"}
// 	ErrInvalidRoleIDFormatIn = ErrorDef{Code: 10, GrpcCode: codes.InvalidArgument, Message: "invalid role_id format: %s"}
// 	ErrBeginTransaction      = ErrorDef{Code: 11, GrpcCode: codes.Internal, Message: "begin transaction: %v"}
// 	ErrValidateRoles         = ErrorDef{Code: 12, GrpcCode: codes.Internal, Message: "validate roles: %v"}
// 	ErrCreateUserInternal    = ErrorDef{Code: 13, GrpcCode: codes.Internal, Message: "create user: %v"}
// 	ErrAddUserRoleMapping    = ErrorDef{Code: 14, GrpcCode: codes.Internal, Message: "add user role mapping: %v"}
// 	ErrCommitTransaction     = ErrorDef{Code: 15, GrpcCode: codes.Internal, Message: "commit transaction: %v"}
// 	ErrGetUserInternal       = ErrorDef{Code: 16, GrpcCode: codes.Internal, Message: "get user: %v"}
// 	ErrUpdateUserInternal    = ErrorDef{Code: 17, GrpcCode: codes.Internal, Message: "update user: %v"}
// 	ErrClearUserRoleMappings = ErrorDef{Code: 18, GrpcCode: codes.Internal, Message: "clear user role mappings: %v"}
// 	ErrDeleteUserInternal    = ErrorDef{Code: 19, GrpcCode: codes.Internal, Message: "delete user: %v"}
// 	ErrListUsersInternal     = ErrorDef{Code: 20, GrpcCode: codes.Internal, Message: "list users: %v"}

// 	// Errors of Role Service
// 	ErrRoleIDRequired         = ErrorDef{Code: 21, GrpcCode: codes.InvalidArgument, Message: "id is required"}
// 	ErrInvalidRoleIDFormat    = ErrorDef{Code: 22, GrpcCode: codes.InvalidArgument, Message: "invalid role id format"}
// 	ErrRoleNameRequired       = ErrorDef{Code: 23, GrpcCode: codes.InvalidArgument, Message: "name is required"}
// 	ErrRoleNameInvalid        = ErrorDef{Code: 24, GrpcCode: codes.InvalidArgument, Message: "name must be between 2 and 100 characters"}
// 	ErrRoleIDsNotExist        = ErrorDef{Code: 25, GrpcCode: codes.InvalidArgument, Message: "one or more role_ids do not exist"}
// 	ErrRoleNameAlreadyExists  = ErrorDef{Code: 26, GrpcCode: codes.AlreadyExists, Message: "role name already exists"}
// 	ErrRoleNotFound           = ErrorDef{Code: 27, GrpcCode: codes.NotFound, Message: "role %s not found"}
// 	ErrCreateRoleInternal     = ErrorDef{Code: 28, GrpcCode: codes.Internal, Message: "create role: %v"}
// 	ErrGetRoleInternal        = ErrorDef{Code: 29, GrpcCode: codes.Internal, Message: "get role: %v"}
// 	ErrUpdateRoleInternal     = ErrorDef{Code: 30, GrpcCode: codes.Internal, Message: "update role: %v"}
// 	ErrDeleteRoleInternal     = ErrorDef{Code: 31, GrpcCode: codes.Internal, Message: "delete role: %v"}
// 	ErrListRolesInternal      = ErrorDef{Code: 32, GrpcCode: codes.Internal, Message: "list roles: %v"}

// 	// Errors of Product Service
// 	ErrProductIDRequired          = ErrorDef{Code: 33, GrpcCode: codes.InvalidArgument, Message: "id is required"}
// 	ErrInvalidProductIDFormat     = ErrorDef{Code: 34, GrpcCode: codes.InvalidArgument, Message: "invalid product id format"}
// 	ErrProductNameRequired        = ErrorDef{Code: 35, GrpcCode: codes.InvalidArgument, Message: "name is required"}
// 	ErrProductNotFound            = ErrorDef{Code: 36, GrpcCode: codes.NotFound, Message: "product %s not found"}
// 	ErrCreateProductInternal      = ErrorDef{Code: 37, GrpcCode: codes.Internal, Message: "create product: %v"}
// 	ErrGetProductInternal         = ErrorDef{Code: 38, GrpcCode: codes.Internal, Message: "get product: %v"}
// 	ErrUpdateProductInternal      = ErrorDef{Code: 39, GrpcCode: codes.Internal, Message: "update product: %v"}
// 	ErrDeleteProductInternal      = ErrorDef{Code: 40, GrpcCode: codes.Internal, Message: "delete product: %v"}
// 	ErrListProductsInternal       = ErrorDef{Code: 41, GrpcCode: codes.Internal, Message: "list products: %v"}
// 	ErrMissingProductsQueryConfig = ErrorDef{Code: 42, GrpcCode: codes.Internal, Message: "missing query config for products"}
// 	ErrInvalidProductSort         = ErrorDef{Code: 43, GrpcCode: codes.InvalidArgument, Message: "%v"}
// 	ErrInvalidProductFilter       = ErrorDef{Code: 44, GrpcCode: codes.InvalidArgument, Message: "%v"}

// 	// Errors of Module Service
// 	ErrModuleIDRequired          = ErrorDef{Code: 45, GrpcCode: codes.InvalidArgument, Message: "id is required"}
// 	ErrInvalidModuleIDFormat     = ErrorDef{Code: 46, GrpcCode: codes.InvalidArgument, Message: "invalid module id format"}
// 	ErrInvalidModuleProductID    = ErrorDef{Code: 47, GrpcCode: codes.InvalidArgument, Message: "invalid product_id format"}
// 	ErrModuleNameRequired        = ErrorDef{Code: 48, GrpcCode: codes.InvalidArgument, Message: "name is required"}
// 	ErrModuleNotFound            = ErrorDef{Code: 49, GrpcCode: codes.NotFound, Message: "module %s not found"}
// 	ErrCreateModuleInternal      = ErrorDef{Code: 50, GrpcCode: codes.Internal, Message: "create module: %v"}
// 	ErrGetModuleInternal         = ErrorDef{Code: 51, GrpcCode: codes.Internal, Message: "get module: %v"}
// 	ErrUpdateModuleInternal      = ErrorDef{Code: 52, GrpcCode: codes.Internal, Message: "update module: %v"}
// 	ErrDeleteModuleInternal      = ErrorDef{Code: 53, GrpcCode: codes.Internal, Message: "delete module: %v"}
// 	ErrListModulesInternal       = ErrorDef{Code: 54, GrpcCode: codes.Internal, Message: "list modules: %v"}
// 	ErrMissingModulesQueryConfig = ErrorDef{Code: 55, GrpcCode: codes.Internal, Message: "missing query config for modules"}
// 	ErrInvalidModuleSort         = ErrorDef{Code: 56, GrpcCode: codes.InvalidArgument, Message: "%v"}
// 	ErrInvalidModuleFilter       = ErrorDef{Code: 57, GrpcCode: codes.InvalidArgument, Message: "%v"}

// 	// Errors of Action Service
// 	ErrActionIDRequired          = ErrorDef{Code: 58, GrpcCode: codes.InvalidArgument, Message: "id is required"}
// 	ErrInvalidActionIDFormat     = ErrorDef{Code: 59, GrpcCode: codes.InvalidArgument, Message: "invalid action id format"}
// 	ErrInvalidActionResourceID   = ErrorDef{Code: 60, GrpcCode: codes.InvalidArgument, Message: "invalid resource_id format"}
// 	ErrActionNameRequired        = ErrorDef{Code: 61, GrpcCode: codes.InvalidArgument, Message: "name is required"}
// 	ErrActionNotFound            = ErrorDef{Code: 62, GrpcCode: codes.NotFound, Message: "action %s not found"}
// 	ErrCreateActionInternal      = ErrorDef{Code: 63, GrpcCode: codes.Internal, Message: "create action: %v"}
// 	ErrGetActionInternal         = ErrorDef{Code: 64, GrpcCode: codes.Internal, Message: "get action: %v"}
// 	ErrUpdateActionInternal      = ErrorDef{Code: 65, GrpcCode: codes.Internal, Message: "update action: %v"}
// 	ErrDeleteActionInternal      = ErrorDef{Code: 66, GrpcCode: codes.Internal, Message: "delete action: %v"}
// 	ErrListActionsInternal       = ErrorDef{Code: 67, GrpcCode: codes.Internal, Message: "list actions: %v"}
// 	ErrMissingActionsQueryConfig = ErrorDef{Code: 68, GrpcCode: codes.Internal, Message: "missing query config for actions"}
// 	ErrInvalidActionSort         = ErrorDef{Code: 69, GrpcCode: codes.InvalidArgument, Message: "%v"}
// 	ErrInvalidActionFilter       = ErrorDef{Code: 70, GrpcCode: codes.InvalidArgument, Message: "%v"}

// 	// Errors of Resource Service
// 	ErrResourceIDRequired           = ErrorDef{Code: 71, GrpcCode: codes.InvalidArgument, Message: "id is required"}
// 	ErrInvalidResourceIDFormat      = ErrorDef{Code: 72, GrpcCode: codes.InvalidArgument, Message: "invalid resource id format"}
// 	ErrInvalidResourceModuleID      = ErrorDef{Code: 73, GrpcCode: codes.InvalidArgument, Message: "invalid module_id format"}
// 	ErrResourceNameRequired         = ErrorDef{Code: 74, GrpcCode: codes.InvalidArgument, Message: "name is required"}
// 	ErrResourceNotFound             = ErrorDef{Code: 75, GrpcCode: codes.NotFound, Message: "resource %s not found"}
// 	ErrCreateResourceInternal       = ErrorDef{Code: 76, GrpcCode: codes.Internal, Message: "create resource: %v"}
// 	ErrGetResourceInternal          = ErrorDef{Code: 77, GrpcCode: codes.Internal, Message: "get resource: %v"}
// 	ErrUpdateResourceInternal       = ErrorDef{Code: 78, GrpcCode: codes.Internal, Message: "update resource: %v"}
// 	ErrDeleteResourceInternal       = ErrorDef{Code: 79, GrpcCode: codes.Internal, Message: "delete resource: %v"}
// 	ErrListResourcesInternal        = ErrorDef{Code: 80, GrpcCode: codes.Internal, Message: "list resources: %v"}
// 	ErrCreateResourceActionInternal = ErrorDef{Code: 81, GrpcCode: codes.Internal, Message: "create action for resource: %v"}
// 	ErrListOldActionsInternal       = ErrorDef{Code: 82, GrpcCode: codes.Internal, Message: "list old actions: %v"}
// 	ErrDeleteOldActionInternal      = ErrorDef{Code: 83, GrpcCode: codes.Internal, Message: "delete old action: %v"}
// 	ErrInsertActionInternal         = ErrorDef{Code: 84, GrpcCode: codes.Internal, Message: "insert action: %v"}
// 	ErrMissingResourcesQueryConfig  = ErrorDef{Code: 85, GrpcCode: codes.Internal, Message: "missing query config for resources"}
// 	ErrInvalidResourceSort          = ErrorDef{Code: 86, GrpcCode: codes.InvalidArgument, Message: "%v"}
// 	ErrInvalidResourceFilter        = ErrorDef{Code: 87, GrpcCode: codes.InvalidArgument, Message: "%v"}

// 	// Additional Errors of Role Service
// 	ErrMissingRolesQueryConfig = ErrorDef{Code: 88, GrpcCode: codes.Internal, Message: "missing query config for roles"}
// 	ErrInvalidRoleSort         = ErrorDef{Code: 89, GrpcCode: codes.InvalidArgument, Message: "%v"}
// 	ErrInvalidRoleFilter       = ErrorDef{Code: 90, GrpcCode: codes.InvalidArgument, Message: "%v"}

// 	// Additional Errors of User Service
// 	ErrMissingUsersQueryConfig = ErrorDef{Code: 91, GrpcCode: codes.Internal, Message: "missing query config for users"}
// 	ErrInvalidUserSort         = ErrorDef{Code: 92, GrpcCode: codes.InvalidArgument, Message: "%v"}
// 	ErrInvalidUserFilter       = ErrorDef{Code: 93, GrpcCode: codes.InvalidArgument, Message: "%v"}
// 	ErrPasswordRequired        = ErrorDef{Code: 94, GrpcCode: codes.InvalidArgument, Message: "password is required"}

// 	// Errors of gRPC Handler payload checks
// 	ErrUserPayloadRequired     = ErrorDef{Code: 95, GrpcCode: codes.InvalidArgument, Message: "user is required"}
// 	ErrRolePayloadRequired     = ErrorDef{Code: 96, GrpcCode: codes.InvalidArgument, Message: "role is required"}
// 	ErrProductPayloadRequired  = ErrorDef{Code: 97, GrpcCode: codes.InvalidArgument, Message: "product is required"}
// 	ErrModulePayloadRequired   = ErrorDef{Code: 98, GrpcCode: codes.InvalidArgument, Message: "module is required"}
// 	ErrResourcePayloadRequired = ErrorDef{Code: 99, GrpcCode: codes.InvalidArgument, Message: "resource is required"}
// 	ErrActionPayloadRequired   = ErrorDef{Code: 100, GrpcCode: codes.InvalidArgument, Message: "action is required"}
// )
