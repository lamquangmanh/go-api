# gRPC Testing

gRPC reflection is enabled, so you can test with `grpcurl` or Postman without importing proto files.

## grpcurl

Install:

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

Health check:

```bash
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
```

List services:

```bash
grpcurl -plaintext localhost:50051 list
```

List methods of UserService:

```bash
grpcurl -plaintext localhost:50051 list proto.user.v1.UserService
```

Create user (`proto.user.v1.UserService/CreateUser`):

```bash
grpcurl -plaintext \
  -d '{
    "user": {
      "username": "alice123",
      "email": "alice@example.com",
      "password": "123456",
      "phone": "0901234567",
      "status": "USER_STATUS_ACTIVE",
      "role_ids": []
    },
    "user_id": "b1aa92ce-fd81-4787-90b9-95f9c3639fce"
  }' \
  localhost:50051 proto.user.v1.UserService/CreateUser
```

Get users (`proto.user.v1.UserService/GetUsers`):

```bash
grpcurl -plaintext \
  -d '{
    "pagination": {"page": 1, "limit": 10},
    "sorts": [{"field": "created_at", "order": "SORT_ORDER_DESC"}],
    "filters": [{"field": "username", "operator": "FILTER_OPERATOR_LIKE", "string_value": "alice"}]
  }' \
  localhost:50051 proto.user.v1.UserService/GetUsers
```

Create role (`proto.role.v1.RoleService/CreateRole`):

```bash
grpcurl -plaintext \
  -d '{
    "role": {
      "name": "admin",
      "description": "Platform administrator",
      "module_id": "2cce26f5-7266-49cf-b487-84e6fa0e41e1",
      "permissions": []
    },
    "user_id": "b1aa92ce-fd81-4787-90b9-95f9c3639fce"
  }' \
  localhost:50051 proto.role.v1.RoleService/CreateRole
```

Get roles (`proto.role.v1.RoleService/GetRoles`):

```bash
grpcurl -plaintext \
  -d '{
    "pagination": {"page": 1, "limit": 10},
    "sorts": [{"field": "created_at", "order": "SORT_ORDER_DESC"}],
    "filters": [{"field": "name", "operator": "FILTER_OPERATOR_LIKE", "string_value": "admin"}]
  }' \
  localhost:50051 proto.role.v1.RoleService/GetRoles
```

Common error format:

```json
{
  "errors": [
    {
      "code": 9,
      "message": "user 00000000-0000-0000-0000-000000000000 not found"
    }
  ]
}
```

Run full smoke test script:

```bash
bash scripts/test_all_apis.sh
```

Generated file:

- `docs/api_test_results.md`

If port `50051` is busy:

```bash
lsof -i :50051
kill <PID>
```

## Postman

1. Create a new gRPC request
2. Server: `localhost:50051`
3. Select `Plain Text` (no TLS)
4. Pick method from discovered services (reflection)
5. Send JSON request body
