# API Test Results

Test run: 2026-07-06
App running at `http://localhost:8080`

---

## Positive Cases

### 1. Register Employee — Success
```http
POST /api/user/register-employee
Content-Type: application/json

{"phone_number":"+628...","email":"test@example.com","full_name":"Test User","address":"Jl. Test No.1"}
```
**Response:** `201 Created`
```json
{"code":201,"message":"User created successfully","data":{"id":"...","phone_number":"+628...","email":"test@example.com","address":"","full_name":"Test User","photo":"","role":"employee","verified":false}}
```

### 2. Register OTP — Success
```http
POST /api/auth/register-otp
Content-Type: application/json

{"phone":"6281234567890","user_id":"019f37f1-5a09-77c4-93ea-1ab2728cf27e"}
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

### 4. Register Employee — All fields missing
```http
POST /api/user/register-employee
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

### 5. Register Employee — Invalid email
```http
POST /api/user/register-employee
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

### 6. Register Employee — Invalid phone (not E.164)
```http
POST /api/user/register-employee
Content-Type: application/json

{"phone_number":"abc","email":"a@b.com","full_name":"T","address":"A"}
```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"Validation failed","errors":[
  {"field":"phone_number","message":"e164"}
]}
```
Result: ✅

### 7. Register Employee — Missing phone_number
```http
POST /api/user/register-employee
Content-Type: application/json

{"email":"a@b.com","full_name":"T","address":"A"}
```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"Validation failed","errors":[
  {"field":"phone_number","message":"required"}
]}
```
Result: ✅

### 8. Register Employee — Empty body
```http
POST /api/user/register-employee
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
Result: ✅ (`user_id` — correct after Bug 1 fix)

### 10. Register OTP — Invalid UUID
```http
POST /api/auth/register-otp
Content-Type: application/json

{"phone":"6281234567890","user_id":"bad"}
```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"bad request","errors":[
  {"field":"user_id","message":"invalid user id"}
]}
```
Result: ✅

### 11. Register OTP — Non-existent user (valid UUID format)
```http
POST /api/auth/register-otp
Content-Type: application/json

{"phone":"6281234567890","user_id":"11111111-1111-7111-8111-111111111111"}
```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"user not found","errors":null}
```
Result: ✅

### 12. Register OTP — Wrong phone for user
```http
POST /api/auth/register-otp
Content-Type: application/json

{"phone":"6280000000000","user_id":"019f37f1-5a09-77c4-93ea-1ab2728cf27e"}
```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"bad request","errors":[
  {"field":"phone","message":"phone number is invalid"}
]}
```
Result: ✅

### 13. Verify User — Missing fields
```http
POST /api/user/verify
Content-Type: application/json

{}
```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"Validation failed","errors":[
  {"field":"phone","message":"required"},
  {"field":"otp","message":"required"}
]}
```
Result: ✅ (`otp` — correct after Bug 1 fix)

### 14. Verify User — Wrong OTP
```http
POST /api/user/verify
Content-Type: application/json

{"phone":"6281234567890","otp":"000000"}
```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"Invalid OTP"}
```
Result: ✅

### 15. Verify User — Non-existent phone (no OTP in Redis)
```http
POST /api/user/verify
Content-Type: application/json

{"phone":"6289999999999","otp":"123456"}
```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"Invalid OTP"}
```
Result: ✅ Fixed — no longer leaks Redis error.

---

## Bugs

### Bug 1: `PascalToSnake` breaks acronyms ~~FIXED~~

`utils/utils.go` now handles two boundary types:
- lowercase → uppercase (normal PascalCase: `FullName` → `full_name`)
- uppercase → lowercase within uppercase run (acronym: `UserID` → `user_id`, `OTP` → `otp`)

✅ Confirmed: fields now render as `user_id` and `otp`.

### ~~Bug 2: Verify User with non-existent phone leaks Redis error~~ **FIXED**

Now returns `400 Invalid OTP` instead of `500 redis: nil`.

---

## Login Endpoint Tests

Test run: 2026-07-07
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
Returns a JWT token with claims:
```json
{
  "UserID": "019f3ca2-4f0e-7366-a96a-80efead92282",
  "Role": "admin",
  "PhoneNumber": "6281234567890",
  "exp": 1783457535,
  "iat": 1783428735
}
```

### Negative Cases

#### 3. Login — Missing phone field
```http
POST /api/auth/login
Content-Type: application/json

{}
```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"Validation failed","errors":[{"field":"phone","message":"required"}]}
```
Result: ✅

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
```http
POST /api/auth/login-verify
Content-Type: application/json

{}
```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"Validation failed","errors":[{"field":"phone","message":"required"},{"field":"otp","message":"required"}]}
```
Result: ✅

#### 6. Login Verify — Wrong OTP
```http
POST /api/auth/login-verify
Content-Type: application/json

{"phone":"6281234567890","otp":"000000"}
```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"invalid otp","errors":null}
```
Result: ✅

#### 7. Login Verify — Non-existent phone (no OTP in Redis)
```http
POST /api/auth/login-verify
Content-Type: application/json

{"phone":"6289999999999","otp":"123456"}
```
**Response:** `400 Bad Request`
```json
{"code":400,"message":"invalid otp","errors":null}
```
Result: ✅ — no Redis error leak.

---

## List Users with Pagination

Test run: 2026-07-07
App running at `http://localhost:8080`

### Positive Cases

#### 1. List Users — Default (no params, filters to role=user)
```http
GET /api/user/list
```
**Response:** `200 OK`
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "019f3ca2-4f1c-7516-824f-358c5dd79f2b",
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
      "id": "019f3ca2-4f0e-7366-a96a-80efead92282",
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

#### 3. List Users — Filter by employee role
```http
GET /api/user/list?role=employee
```
**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "019f3ca2-4f15-736c-9877-50bf60ad49b5",
      "phone_number": "6281122334455",
      "email": "jane@example.com",
      "full_name": "Jane Smith",
      "role": "employee",
      "verified": true
    }
  ],
  "meta": {"page": 1, "limit": 10, "total": 1, "total_pages": 1}
}
```

#### 4. List Users — Empty role (returns all roles)
```http
GET /api/user/list?role=
```
**Response:** `200 OK` — returns all 3 seeded users (admin, employee, user).

#### 5. List Users — Pagination (page 1, limit 2)
```http
GET /api/user/list?page=1&limit=2&role=
```
**Response:** `200 OK` — 2 items, `"total": 3`, `"total_pages": 2`.

#### 6. List Users — Pagination (page 2, limit 2)
```http
GET /api/user/list?page=2&limit=2&role=
```
**Response:** `200 OK` — 1 item (John Doe).

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
