# Go Naming Conventions (DO / DON'T)

## 1) Package Names

### DO

- Use short, lowercase package names.
- Use a single word when possible.
- Keep names meaningful in import usage.

```go
package config
package repository
package service
```

### DON'T

- Don't use uppercase letters.
- Don't use underscores or hyphens.
- Don't use generic names like `util` if a specific name is possible.

```go
package Config      // bad
package user_service // bad
package my-utils     // bad
```

---

## 2) File Names

### DO

- Use lowercase with underscores for multi-word files.
- Name by responsibility.
- Keep suffixes consistent by layer.

```text
user_handler.go
role_service.go
postgres.go
config_test.go
```

### DON'T

- Don't use camelCase or PascalCase for file names.
- Don't use vague names like `common.go` unless truly shared.
- Don't put multiple unrelated responsibilities in one file.

```text
UserHandler.go   // bad
userHandler.go   // bad
misc.go          // bad
```

---

## 3) Variable Names

### DO

- Use short but clear names.
- Use `camelCase` for local variables.
- Use common abbreviations (`ctx`, `err`, `id`) where idiomatic.

```go
userID := req.GetId()
createdAt := time.Now()
validationErr := validate(input)
```

### DON'T

- Don't use unclear abbreviations.
- Don't use all-caps (except constants).
- Don't use one-letter names unless loop index or tiny scope.

```go
u := request.GetUser() // bad (unclear in normal scope)
USER_ID := req.GetId() // bad
x := getUser()         // bad (unclear)
```

---

## 4) Function & Method Names

### DO

- Use `PascalCase` for exported functions.
- Use `camelCase` for unexported functions.
- Use verb-based names for actions.

```go
func CreateUser(ctx context.Context, req *Request) error
func validateEmail(email string) bool
func parseOptionalUUID(raw string) (uuid.UUID, error)
```

### DON'T

- Don't use snake_case in function names.
- Don't use names that hide behavior.

```go
func create_user() {}      // bad
func Handle() {}           // bad (too generic)
func DoStuff() {}          // bad
```

---

## 5) Struct, Interface, and Type Names

### DO

- Use `PascalCase`.
- Use noun names for structs/types.
- Use behavior-based names for interfaces.
- Keep interface names small and specific.

```go
type UserService struct {}
type RoleRepository interface {
	GetByID(ctx context.Context, id string) (*Role, error)
}

type Validator interface {
	Validate(any) error
}
```

### DON'T

- Don't prefix types with package name redundantly.
- Don't use `IUserService` style interface naming.

```go
type ServiceUser struct {}      // bad naming style

type IUserRepository interface { // bad (Java/C# style)
	FindByID(ctx context.Context, id string) error
}
```

---

## 6) Constants

### DO

- Use `PascalCase` for exported constants.
- Use `camelCase` for unexported constants.
- Group related constants with `const (...)`.

```go
const (
	DefaultPageLimit = 20
	MaxPageLimit     = 100
)

const defaultTimeoutSeconds = 30
```

### DON'T

- Don't use ALL_CAPS style from other languages.

```go
const MAX_LIMIT = 100 // bad in Go style
```

---

## 7) Acronyms & Initialisms

### DO

- Use standard Go initialism style in identifiers.
- Keep initialisms uppercase in names.

```go
userID
apiURL
httpClient
parseUUID
```

### DON'T

- Don't mix inconsistent capitalization.

```go
userId     // bad
apiUrl     // bad
parseUuid  // bad
```

---

## 8) Error Variables

### DO

- Use `err` as local error variable.
- Prefix sentinel errors with `Err`.

```go
if err := svc.CreateUser(ctx, req); err != nil {
	return err
}

var ErrUserNotFound = errors.New("user not found")
```

### DON'T

- Don't use `e`, `errorObj`, or unclear names for common error flow.

```go
if e != nil { return e } // bad style in Go
```

---

## 9) Receiver Names

### DO

- Use short receiver names (1-2 letters), usually from type name.
- Keep receiver names consistent for same type.

```go
func (s *UserService) CreateUser(...) error { ... }
func (r *RoleRepository) GetByID(...) (*Role, error) { ... }
```

### DON'T

- Don't use `this` or very long receiver names.

```go
func (this *UserService) CreateUser(...) error { ... } // bad
func (userService *UserService) CreateUser(...) error { ... } // usually too long
```

---

## 10) Directory & Layer Naming

### DO

- Use folder names reflecting architecture/layer.
- Keep naming consistent across services.

```text
internal/grpc/
internal/service/
internal/repository/
pkg/constants/
pkg/utils/
```

### DON'T

- Don't mix naming styles between layers.
- Don't use ambiguous folder names when role is clear.

```text
internal/Service/   // bad
pkg/helpers/        // bad if purpose is actually validation
```

---

## Quick Checklist

### DO

- Prefer clarity over brevity.
- Follow Go idioms (`camelCase`, `PascalCase`, `err`, `ctx`).
- Keep naming consistent by module/layer.

### DON'T

- Don't mix naming styles across files.
- Don't use snake_case for Go identifiers.
- Don't use Java/C# naming habits (`IService`, `this`, `MAX_LIMIT`).
