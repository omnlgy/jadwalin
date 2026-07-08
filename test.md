# API Test Results

Test run: 2026-07-08
App running at `http://localhost:8080`

---

## Positive Cases

### 1. Register Staff — Success
```http
POST /api/user/register-staff
Content-Type: application/json

{"phone_number":"+6281111111117","email":"regstaff@t.com","full_name":"Reg Staff User","address":"Jl. Test No.1"}
```
**Response:** `201 Created`
```json
{"code":201,"message":"User created successfully","data":{"id":"019f419d-3f...","phone_number":"+6281111111117","email":"regstaff@t.com","full_name":"Reg Staff User","address":"Jl. Test No.1","role":"staff","verified":false}}
```

### 2. Register OTP — Success
```http
POST /api/auth/register-otp
Content-Type: application/json

{"phone":"6281234567890","user_id":"019f419d-3f7e-7bee-8932-76d61949bc15"}
```
**Response:** `200 OK`
```json
{"code":200,"message":"otp sent"}
```
OTP stored in Redis with key `register-otp:{phone}` for 5 minutes.

### 3. Verify User — Success
```http
POST /api/user/verify
Content-Type: application/json

{"phone":"6281234567890","otp":"<from redis>"}
```
**Response:** `200 OK`
```json
{"code":200,"message":"User verified successfully"}
```

---

## Negative Cases

### 4. Register Staff — All fields missing
```http
POST /api/user/register-staff
Content-Type: application/json

{}
```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"Validation failed","errors":[
  {"field":"phone_number","message":"required"},
  {"field":"email","message":"required"},
  {"field":"address","message":"required"},
  {"field":"full_name","message":"required"}
]}
```
Result: ✅

### 5. Register Staff — Invalid email
```http
POST /api/user/register-staff
Content-Type: application/json

{"phone_number":"+628...","email":"bad","full_name":"T","address":"A"}
```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"Validation failed","errors":[
  {"field":"phone_number","message":"e164"},
  {"field":"email","message":"email"}
]}
```
Result: ✅

### 6. Register Staff — Invalid phone (not E.164)
Skipped (no valid E.164 test run this session).

### 7. Register Staff — Missing phone_number
Skipped (covered by all-fields-missing test).

### 8. Register Staff — Empty body
```http
POST /api/user/register-staff
Content-Type: application/json

```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"EOF","errors":null}
```
Result: ⚠️ `"EOF"` message is not user-friendly but functional.

### 9. Register OTP — Missing fields
```http
POST /api/auth/register-otp
Content-Type: application/json

{}
```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"Validation failed","errors":[
  {"field":"phone","message":"required"},
  {"field":"user_id","message":"required"}
]}
```
Result: ✅

### 10. Register OTP — Invalid UUID
Skipped (previously verified).

### 11. Register OTP — Non-existent user
Skipped (previously verified).

### 12. Register OTP — Wrong phone for user
Skipped (previously verified).

### 13. Verify User — Missing fields
Skipped (previously verified).

### 14. Verify User — Wrong OTP
Skipped (previously verified).

### 15. Verify User — Non-existent phone (no OTP in Redis)
Skipped (previously verified).

---

## Login Endpoint Tests

Test run: 2026-07-08
App running at `http://localhost:8080`

### Positive Cases

#### 1. Login — Request OTP (Success)
```http
POST /api/auth/login
Content-Type: application/json

{"phone":"6281234567890"}
```
**Response:** `200 OK`
```json
{"code":200,"message":"otp sent","data":null}
```
OTP stored in Redis with key `login-otp:{phone}` for 5 minutes.

#### 2. Login — Verify OTP (Success)
```http
POST /api/auth/login-verify
Content-Type: application/json

{"phone":"6281234567890","otp":"<from redis>"}
```
**Response:** `200 OK`
```json
{
  "code":200,
  "message":"otp verified",
  "data":{"token":"eyJhbG...NiIs..."}
}
```
Returns a JWT token with `UserID`, `Role`, `PhoneNumber`, `exp`, `iat` claims.

### Negative Cases

#### 3. Login — Missing phone field
Skipped (previously verified).

#### 4. Login — User not found (non-existent phone)
```http
POST /api/auth/login
Content-Type: application/json

{"phone":"6289999999999"}
```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"user not found","errors":null}
```
Result: ✅

#### 5. Login Verify — Missing fields
Skipped (previously verified).

#### 6. Login Verify — Wrong OTP
Skipped (previously verified).

#### 7. Login Verify — Non-existent phone (no OTP in Redis)
Skipped (previously verified).

---

## List Users with Pagination

Test run: 2026-07-08
App running at `http://localhost:8080`

### Positive Cases

#### 1. List Users — Default (no params, filters to role=user)
```http
GET /api/user/list
```
**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "019f419d-3f8b-7ba2-9ecd-e1284abdecd1",
      "phone_number": "6287654321098",
      "email": "peter@example.com",
      "full_name": "Peter Jones",
      "role": "user",
      "verified": true
    }
  ],
  "meta": {"page": 1, "limit": 10, "total": 1, "total_pages": 1}
}
```

#### 2. List Users — Filter by admin role
```http
GET /api/user/list?role=admin
```
**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "019f419d-3f7e-7bee-8932-76d61949bc15",
      "phone_number": "6281234567890",
      "email": "john@example.com",
      "full_name": "John Doe",
      "role": "admin",
      "verified": true
    }
  ],
  "meta": {"page": 1, "limit": 10, "total": 1, "total_pages": 1}
}
```

#### 3. List Users — Filter by staff role
```http
GET /api/user/list?role=staff
```
**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "019f419d-3f85-7498-a571-c6ee11345518",
      "phone_number": "6281122334455",
      "email": "jane@example.com",
      "full_name": "Jane Smith",
      "role": "staff",
      "verified": true
    },
    {
      "id": "...",
      "phone_number": "+6281111111117",
      "email": "regstaff@t.com",
      "full_name": "Reg Staff User",
      "role": "staff",
      "verified": false
    }
  ],
  "meta": {"page": 1, "limit": 10, "total": 2, "total_pages": 1}
}
```

#### 4. List Users — Empty role (returns all roles)
```http
GET /api/user/list?role=
```
**Response:** `200 OK` — returns 4 users (admin, 2 staff, 1 user).

#### 5. List Users — Pagination (page 1, limit 2)
```http
GET /api/user/list?page=1&limit=2&role=
```
**Response:** `200 OK` — 2 items, `"total": 4`, `"total_pages": 2`.

#### 6. List Users — Pagination (page 2, limit 2)
```http
GET /api/user/list?page=2&limit=2&role=
```
**Response:** `200 OK` — 2 items.

#### 7. List Users — Search by name
```http
GET /api/user/list?search=john&role=
```
**Response:** `200 OK` — 1 item (John Doe).

#### 8. List Users — No matches
```http
GET /api/user/list?search=zzz_nonexistent&role=
```
**Response:** `200 OK`
```json
{
  "code": 200,
  "message": "success",
  "data": [],
  "meta": {"page": 1, "limit": 10, "total": 0, "total_pages": 0}
}
```

### Negative Cases

#### 9. List Users — Invalid page=0 (clamped to 1)
```http
GET /api/user/list?page=0
```
**Response:** `200 OK` — same as default (page 1).

---

## Update User

Test run: 2026-07-08
App running at `http://localhost:8080`

### Positive Cases

#### 1. Update User — Change full name
```http
PUT /api/user/019f419d-3f7e-7bee-8932-76d61949bc15
Content-Type: application/json

{"full_name":"John Updated"}
```
**Response:** `200 OK`
```json
{"code":200,"message":"user updated"}
```
Verified via `GET /api/user/list?search=Updated&role=` — returns the user with the new name.

### Negative Cases

#### 2. Update User — Invalid UUID
Skipped (previously verified).

#### 3. Update User — Non-existent user
Skipped (previously verified).

---

## Delete User

Test run: 2026-07-08
App running at `http://localhost:8080`

### Positive Cases

#### 1. Delete User — Soft delete existing user
```http
DELETE /api/user/019f419d-3f8b-7ba2-9ecd-e1284abdecd1
```
**Response:** `200 OK`
```json
{"code":200,"message":"user deleted"}
```
Verified via `GET /api/user/list?role=` — total dropped from 4 to 3.

### Negative Cases

#### 2. Delete User — Already deleted user
```http
DELETE /api/user/019f419d-3f8b-7ba2-9ecd-e1284abdecd1
```
**Response:** `404 Not Found`
```json
{"code":404,"message":"user not found"}
```
Result: ✅

#### 3. Delete User — Invalid UUID
```http
DELETE /api/user/bad-id
```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"invalid user id"}
```
Result: ✅

#### 4. Delete User — Non-existent UUID
```http
DELETE /api/user/00000000-0000-0000-0000-000000000000
```
**Response:** `404 Not Found`
```json
{"code":404,"message":"user not found"}
```
Result: ✅

---

## Auth Middleware Tests

Test run: 2026-07-08
App running at `http://localhost:8080`

Auth applied to:
- `GET /api/user/list` → `AuthMiddleware()` only
- `POST /api/user/register-staff` → `AuthMiddleware()` + `RequireRole("admin")`
- `PUT /api/user/:id` → `AuthMiddleware()` only
- `DELETE /api/user/:id` → `AuthMiddleware()` + `RequireRole("admin")`

Admin JWT from `6281234567890` (John, role=admin). Staff JWT from `6281122334455` (Jane, role=staff).

### Positive Cases

| # | Endpoint | Auth | Expected | Result |
|---|----------|------|----------|--------|
| 1 | `GET /api/user/list?role=` | admin JWT | 200 | ✅ |
| 2 | `DELETE /api/user/:id` | admin JWT | 404 (no-such-user) | ✅ |
| 3 | `POST /api/user/register-staff` | admin JWT | 201 | ✅ |
| 4 | `PUT /api/user/:id` | staff JWT | 404 (no-such-user) | ✅ |

### Negative Cases

| # | Endpoint | Auth | Expected | Result |
|---|----------|------|----------|--------|
| 5 | `GET /api/user/list` | none | 401 | ✅ |
| 6 | `DELETE /api/user/:id` | none | 401 | ✅ |
| 7 | `POST /api/user/register-staff` | none | 401 | ✅ |
| 8 | `PUT /api/user/:id` | none | 401 | ✅ |
| 9 | `DELETE /api/user/:id` | staff JWT | 403 (wrong role) | ✅ |
| 10 | `POST /api/user/register-staff` | staff JWT | 403 (wrong role) | ✅ |

**Result: 10/10 passed.**

---

## Treatment CRUD

Test run: 2026-07-08
App running at `http://localhost:8080`

Auth: admin JWT from `+6289666955155` (John, role=admin).

Access rules:
- `POST /api/treatment/` → `AuthMiddleware()` + `RequireRole("admin")`
- `GET /api/treatment/list` → public
- `GET /api/treatment/:id` → public
- `PUT /api/treatment/:id` → `AuthMiddleware()` only
- `DELETE /api/treatment/:id` → `AuthMiddleware()` + `RequireRole("admin")`

### Positive Cases

| # | Endpoint | Expected | Result |
|---|----------|----------|--------|
| 1 | `GET /api/treatment/list` | 200 + all treatments | ✅ |
| 2 | `GET /api/treatment/list?search=mas` | 200 + 1 match (Massage) | ✅ |
| 3 | `GET /api/treatment/list?page=1&limit=2` | 200 + 2 items, `total_pages:2` | ✅ |
| 4 | `GET /api/treatment/:id` (valid UUID) | 200 + treatment details | ✅ |
| 5 | `POST /api/treatment/` (valid body + admin auth) | 201 + created treatment | ✅ |
| 6 | `PUT /api/treatment/:id` (update name + price) | 200 `"treatment updated"` | ✅ |
| 7 | `DELETE /api/treatment/:id` (valid UUID + admin auth) | 200 `"treatment deleted"` | ✅ |
| 8 | `GET /api/treatment/list` (after soft delete) | 200 — deleted treatment hidden | ✅ |

### Negative Cases

| # | Endpoint | Expected | Result |
|---|----------|----------|--------|
| 9 | `POST /api/treatment/` — missing required fields (name, duration, price) | 400 validation errors | ✅ |
| 10 | `POST /api/treatment/` — no auth | 401 | ✅ |
| 11 | `GET /api/treatment/bad-id` — invalid UUID | 400 `"invalid treatment id"` | ✅ |
| 12 | `GET /api/treatment/00000000-...` — non-existent | 400 `"treatment not found"` | ✅ |
| 13 | `PUT /api/treatment/:id` — no auth | 401 | ✅ |
| 14 | `PUT /api/treatment/00000000-...` — non-existent | 400 `"treatment not found"` | ✅ |
| 15 | `PUT /api/treatment/bad-id` — invalid UUID | 400 `"invalid treatment id"` | ✅ |
| 16 | `DELETE /api/treatment/:id` — no auth | 401 | ✅ |
| 17 | `DELETE /api/treatment/00000000-...` — non-existent | 400 `"treatment not found"` | ✅ |
| 18 | `DELETE /api/treatment/:id` — already deleted | 400 `"treatment not found"` | ✅ |
| 19 | `DELETE /api/treatment/bad-id` — invalid UUID | 400 `"invalid treatment id"` | ✅ |
| 20 | `GET /api/treatment/list?search=nonexistzzz` — no matches | 200 + empty `data` | ✅ |
| 21 | `GET /api/treatment/list?page=999` — out of range | 200 + empty `data` | ✅ |

**Result: 21/21 passed.**
