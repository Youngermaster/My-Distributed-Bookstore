# User Service

## Overview

Handles user authentication, authorization, profile management, and address management using JWT tokens and role-based access control (RBAC).

## Technology Stack

- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **ORM**: GORM
- **Database**: PostgreSQL 15
- **Auth**: JWT
- **Ports**: HTTP: 8082, gRPC: 50052

## Responsibilities

- User registration and authentication
- JWT token generation and validation
- User profile management
- Address management (shipping/billing)
- Role-based access control (RBAC)
- Session management
- Password reset functionality

## Database Schema

**users**: id, email, password_hash, full_name, phone, created_at, updated_at, last_login_at
**roles**: id, name, permissions (JSONB)
**user_roles**: user_id, role_id
**addresses**: id, user_id, address_line1, address_line2, city, state, postal_code, country, is_default
**sessions**: id, user_id, token_hash, ip_address, user_agent, expires_at, created_at

## REST API Endpoints

```
POST   /api/v1/auth/register      # User registration
POST   /api/v1/auth/login         # Login (returns JWT)
POST   /api/v1/auth/logout        # Logout
POST   /api/v1/auth/refresh       # Refresh JWT token
POST   /api/v1/auth/forgot-password
POST   /api/v1/auth/reset-password
GET    /api/v1/users/me           # Get current user profile
PUT    /api/v1/users/me           # Update profile
GET    /api/v1/users/me/addresses # List addresses
POST   /api/v1/users/me/addresses # Add address
PUT    /api/v1/users/me/addresses/:id
DELETE /api/v1/users/me/addresses/:id
```

## gRPC Methods

```protobuf
rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
rpc GetUser(GetUserRequest) returns (GetUserResponse);
rpc CheckPermission(CheckPermissionRequest) returns (CheckPermissionResponse);
```

## Events Published

- `user.registered` - New user signed up
- `user.logged_in` - User logged in
- `user.profile_updated` - Profile changed
- `user.password_reset` - Password was reset

## Environment Variables

```bash
HTTP_PORT=8082
GRPC_PORT=50052
DATABASE_URL=postgresql://bookstore:password@postgres:5432/users_db
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h
BCRYPT_COST=10
```

## Next Steps

- [ ] Implement user models and migrations
- [ ] Create authentication handlers
- [ ] Implement JWT middleware
- [ ] Add password hashing with bcrypt
- [ ] Implement RBAC
- [ ] Add email verification
- [ ] Implement 2FA (optional)
- [ ] Add rate limiting for auth endpoints
- [ ] Write comprehensive tests
