# API Examples

Examples below follow current gRPC contract in `proto/*/v1/*.proto`.

## Create Role (`RoleService/CreateRole`)

Request:

```json
{
  "role": {
    "name": "admin",
    "description": "Platform administrator",
    "module_id": "2cce26f5-7266-49cf-b487-84e6fa0e41e1",
    "permissions": [
      {
        "resource_id": "9ac1f2f0-9f2d-4f95-89dd-49d4b8f8b990",
        "action_id": "f4ec4db5-7702-47c2-aefe-c8f35f8770ee"
      }
    ]
  },
  "user_id": "b1aa92ce-fd81-4787-90b9-95f9c3639fce"
}
```

Success response:

```json
{
  "role": {
    "role_id": "f9dcb09e-9c79-4d9e-9db0-d7a2f2f0f6f2",
    "name": "admin",
    "description": "Platform administrator",
    "module_id": "2cce26f5-7266-49cf-b487-84e6fa0e41e1",
    "created_user_id": "b1aa92ce-fd81-4787-90b9-95f9c3639fce",
    "created_at": "2026-04-25T10:00:00Z"
  },
  "errors": []
}
```

## List Roles (`RoleService/GetRoles`)

Request:

```json
{
  "pagination": { "page": 1, "limit": 10 },
  "sorts": [{ "field": "created_at", "order": "SORT_ORDER_DESC" }],
  "filters": [
    {
      "field": "name",
      "operator": "FILTER_OPERATOR_LIKE",
      "string_value": "admin"
    }
  ]
}
```

## Create User (`UserService/CreateUser`)

Request:

```json
{
  "user": {
    "username": "john_doe",
    "email": "john@example.com",
    "password": "123456",
    "phone": "0901234567",
    "status": "USER_STATUS_ACTIVE",
    "role_ids": ["f9dcb09e-9c79-4d9e-9db0-d7a2f2f0f6f2"]
  },
  "user_id": "b1aa92ce-fd81-4787-90b9-95f9c3639fce"
}
```

Validation error response:

```json
{
  "errors": [
    {
      "code": 2,
      "message": "username must be between 3 and 50 characters and contain only letters, numbers, '_', '.', '-'",
      "extra": {
        "field": "username"
      }
    },
    {
      "code": 3,
      "message": "email is required",
      "extra": {
        "field": "email"
      }
    }
  ]
}
```

## List Users (`UserService/GetUsers`)

Request:

```json
{
  "pagination": { "page": 1, "limit": 10 },
  "sorts": [{ "field": "created_at", "order": "SORT_ORDER_DESC" }],
  "filters": [
    {
      "field": "username",
      "operator": "FILTER_OPERATOR_LIKE",
      "string_value": "john"
    },
    {
      "field": "email",
      "operator": "FILTER_OPERATOR_LIKE",
      "string_value": "@example.com"
    }
  ]
}
```

## Standard error format

Most RPC errors are returned in this structure:

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

For full request/response samples per RPC, see:

- [docs/api_test_results.md](./api_test_results.md)
- [proto/user/v1/user.proto](../proto/user/v1/user.proto)
- [proto/role/v1/role.proto](../proto/role/v1/role.proto)
