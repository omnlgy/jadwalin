# API Test Results

Test run: 2026-07-09
App running at `http://localhost:8080`

---

## User / Auth Module

| # | Method | Endpoint | Request / Body | Description | Expected | Actual | Result |
|---|--------|----------|----------------|-------------|----------|--------|--------|
| 1 | POST | `/api/user/register-staff` | `{phone_number, email, full_name, address}` | Register staff with all required fields | `code:201, message:"User created successfully"` | `code:201, message:"User created successfully"` | ✅ |
| 2 | POST | `/api/auth/register-otp` | `{phone, user_id}` | Send OTP for user verification | `code:200, message:"otp sent"` | `code:200, message:"otp sent"` | ✅ |
| 3 | POST | `/api/user/verify` | `{phone, otp}` | Verify user with correct OTP | `code:200, message:"User verified successfully"` | `code:200, message:"User verified successfully"` | ✅ |
| 4 | POST | `/api/user/register-staff` | `{}` (empty) | Register staff with empty body → validation errors | `code:400, errors: [phone_number required, email required, ...]` | `code:400, errors: [phone_number required, email required, ...]` | ✅ |
| 5 | POST | `/api/user/register-staff` | `{email:"bad", phone_number:"+628...", ...}` | Register staff with invalid email and phone format | `code:400, errors: [e164, email]` | `code:400, errors: [e164, email]` | ✅ |
| 6 | POST | `/api/user/register-staff` | no auth header | Register staff without authentication → unauthorized | `code:401` | `code:401` | ✅ |
| 7 | POST | `/api/auth/register-otp` | `{}` (empty) | Send OTP with empty body → validation errors | `code:400, errors: [phone required, user_id required]` | `code:400, errors: [phone required, user_id required]` | ✅ |
| 8 | POST | `/api/auth/register-otp` | invalid UUID in user_id | Send OTP with invalid user_id format | `code:400` | `code:400` | ✅ |
| 9 | POST | `/api/auth/register-otp` | non-existent user_id | Send OTP for non-existent user | `code:400` | `code:400` | ✅ |
| 10 | POST | `/api/auth/register-otp` | wrong phone for user | Send OTP with phone not matching user | `code:400` | `code:400` | ✅ |

### Login

| # | Method | Endpoint | Request / Body | Description | Expected | Actual | Result |
|---|--------|----------|----------------|-------------|----------|--------|--------|
| 11 | POST | `/api/auth/login` | `{phone: "6281234567890"}` | Request login OTP for existing user | `code:200, message:"otp sent"` | `code:200, message:"otp sent"` | ✅ |
| 12 | POST | `/api/auth/login-verify` | `{phone, otp}` | Verify login OTP → get JWT token | `code:200, token: "eyJhbG..."` | `code:200, token: "eyJhbG..."` | ✅ |
| 13 | POST | `/api/auth/login` | `{phone: "6289999999999"}` (non-existent) | Request login OTP for non-existent phone → user not found | `code:400, message:"user not found"` | `code:400, message:"user not found"` | ✅ |

### List Users

| # | Method | Endpoint | Request / Body | Description | Expected | Actual | Result |
|---|--------|----------|----------------|-------------|----------|--------|--------|
| 14 | GET | `/api/user/list` | no params | List users with default role=user filter | `code:200, 1 user (Peter Jones)` | `code:200, 1 user (Peter Jones)` | ✅ |
| 15 | GET | `/api/user/list?role=admin` | filter by admin role | List users filtered by admin role | `code:200, 1 admin (John Doe)` | `code:200, 1 admin (John Doe)` | ✅ |
| 16 | GET | `/api/user/list?role=staff` | filter by staff role | List users filtered by staff role | `code:200, 2 staff` | `code:200, 2 staff` | ✅ |
| 17 | GET | `/api/user/list?role=` | empty role (all) | List all roles with empty role param | `code:200, 4 users` | `code:200, 4 users` | ✅ |
| 18 | GET | `/api/user/list?page=1&limit=2&role=` | pagination params | List users with pagination (page 1, limit 2) | `code:200, items:2, total:4, total_pages:2` | `code:200, items:2, total:4, total_pages:2` | ✅ |
| 19 | GET | `/api/user/list?page=2&limit=2&role=` | page 2 | List users page 2 of pagination | `code:200, 2 items` | `code:200, 2 items` | ✅ |
| 20 | GET | `/api/user/list?search=john&role=` | search by name | Search users by name (John) | `code:200, 1 item (John Doe)` | `code:200, 1 item (John Doe)` | ✅ |
| 21 | GET | `/api/user/list?search=zzz_nonexistent&role=` | no matches | Search users with non-matching name → empty result | `code:200, data:[]` | `code:200, data:[]` | ✅ |
| 22 | GET | `/api/user/list?page=0` | invalid page | List users with page=0 (clamped to 1) | `code:200, same as page 1` | `code:200, same as page 1` | ✅ |

### Update User

| # | Method | Endpoint | Request / Body | Description | Expected | Actual | Result |
|---|--------|----------|----------------|-------------|----------|--------|--------|
| 23 | PUT | `/api/user/:id` | `{full_name: "John Updated"}` + admin auth | Update user's full name | `code:200, message:"user updated"` | `code:200, message:"user updated"` | ✅ |
| 24 | PUT | `/api/user/bad-id` | `{full_name: "..."}` + admin auth | Update user with invalid UUID format | `code:400` | `code:400` | ✅ |
| 25 | PUT | `/api/user/00000000-...` | non-existent user + admin auth | Update non-existent user | `code:404` | `code:404` | ✅ |

### Delete User

| # | Method | Endpoint | Request / Body | Description | Expected | Actual | Result |
|---|--------|----------|----------------|-------------|----------|--------|--------|
| 26 | DELETE | `/api/user/:id` | valid UUID + admin auth | Soft delete existing user | `code:200, message:"user deleted"` | `code:200, message:"user deleted"` | ✅ |
| 27 | DELETE | `/api/user/:id` | already deleted + admin auth | Delete already soft-deleted user → not found | `code:404, message:"user not found"` | `code:404, message:"user not found"` | ✅ |
| 28 | DELETE | `/api/user/bad-id` | invalid UUID + admin auth | Delete user with invalid UUID format | `code:400, message:"invalid user id"` | `code:400, message:"invalid user id"` | ✅ |
| 29 | DELETE | `/api/user/00000000-...` | non-existent UUID + admin auth | Delete non-existent user (valid UUID but no record) | `code:404, message:"user not found"` | `code:404, message:"user not found"` | ✅ |

### Auth Middleware

| # | Method | Endpoint | Auth | Description | Expected | Actual | Result |
|---|--------|----------|------|-------------|----------|--------|--------|
| 30 | GET | `/api/user/list?role=` | admin JWT | Admin can list users | `code:200` | `code:200` | ✅ |
| 31 | DELETE | `/api/user/:id` | admin JWT | Admin can delete users (no-such-user is expected since already deleted) | `code:404` (no-such-user) | `code:404` | ✅ |
| 32 | POST | `/api/user/register-staff` | admin JWT | Admin can register staff | `code:201` | `code:201` | ✅ |
| 33 | PUT | `/api/user/:id` | staff JWT | Staff can update users (no-such-user expected) | `code:404` (no-such-user) | `code:404` | ✅ |
| 34 | GET | `/api/user/list` | none | List users without auth → unauthorized | `code:401` | `code:401` | ✅ |
| 35 | DELETE | `/api/user/:id` | none | Delete user without auth → unauthorized | `code:401` | `code:401` | ✅ |
| 36 | POST | `/api/user/register-staff` | none | Register staff without auth → unauthorized | `code:401` | `code:401` | ✅ |
| 37 | PUT | `/api/user/:id` | none | Update user without auth → unauthorized | `code:401` | `code:401` | ✅ |
| 38 | DELETE | `/api/user/:id` | staff JWT (wrong role) | Delete user with staff role (not admin) → forbidden | `code:403` | `code:403` | ✅ |
| 39 | POST | `/api/user/register-staff` | staff JWT (wrong role) | Register staff with staff role (not admin) → forbidden | `code:403` | `code:403` | ✅ |

---

## Treatment CRUD

Auth: admin JWT from `+6289666955155` (John, role=admin).

Access rules:
- `POST /api/treatment/` → `AuthMiddleware()` + `RequireRole("admin")`
- `GET /api/treatment/list` → public
- `GET /api/treatment/:id` → public
- `PUT /api/treatment/:id` → `AuthMiddleware()` only
- `DELETE /api/treatment/:id` → `AuthMiddleware()` + `RequireRole("admin")`

### Positive Cases

| # | Method | Endpoint | Request / Body | Description | Expected | Actual | Result |
|---|--------|----------|----------------|-------------|----------|--------|--------|
| 1 | GET | `/api/treatment/list` | — | List all treatments (public) | `code:200, 3 treatments` | `code:200, 3 treatments` | ✅ |
| 2 | GET | `/api/treatment/list?search=mas` | `search=mas` | Search treatments by name (partial match) | `code:200, 1 match (Massage)` | `code:200, 1 match (Massage)` | ✅ |
| 3 | GET | `/api/treatment/list?page=1&limit=2` | `page=1, limit=2` | List treatments with pagination | `code:200, items:2, total_pages:2` | `code:200, items:2, total_pages:2` | ✅ |
| 4 | GET | `/api/treatment/:id` | path: valid treatment UUID | Get treatment by ID (public) | `code:200, treatment details` | `code:200, treatment details` | ✅ |
| 5 | POST | `/api/treatment/` | `{name, description, duration, price}` + admin auth | Create new treatment (admin only) | `code:201, created treatment` | `code:201, created treatment` | ✅ |
| 6 | PUT | `/api/treatment/:id` | `{name, price}` + auth | Update treatment name and price | `code:200, message:"treatment updated"` | `code:200, message:"treatment updated"` | ✅ |
| 7 | DELETE | `/api/treatment/:id` | path: valid UUID + admin auth | Soft delete treatment (admin only) | `code:200, message:"treatment deleted"` | `code:200, message:"treatment deleted"` | ✅ |
| 8 | GET | `/api/treatment/list` | after soft delete | Verify deleted treatment hidden from list | `code:200, deleted treatment hidden` | `code:200, deleted treatment hidden` | ✅ |

### Negative Cases

| # | Method | Endpoint | Request / Body | Description | Expected | Actual | Result |
|---|--------|----------|----------------|-------------|----------|--------|--------|
| 9 | POST | `/api/treatment/` | missing name, duration, price | Create treatment with missing required fields → validation error | `code:400, validation errors` | `code:400, validation errors` | ✅ |
| 10 | POST | `/api/treatment/` | no auth | Create treatment without auth → unauthorized | `code:401` | `code:401` | ✅ |
| 11 | GET | `/api/treatment/bad-id` | invalid UUID | Get treatment with invalid UUID format → bad request | `code:400, message:"invalid treatment id"` | `code:400, message:"invalid treatment id"` | ✅ |
| 12 | GET | `/api/treatment/00000000-...` | non-existent UUID | Get non-existent treatment → not found | `code:400, message:"treatment not found"` | `code:400, message:"treatment not found"` | ✅ |
| 13 | PUT | `/api/treatment/:id` | no auth | Update treatment without auth → unauthorized | `code:401` | `code:401` | ✅ |
| 14 | PUT | `/api/treatment/00000000-...` | non-existent UUID + auth | Update non-existent treatment → not found | `code:400, message:"treatment not found"` | `code:400, message:"treatment not found"` | ✅ |
| 15 | PUT | `/api/treatment/bad-id` | invalid UUID + auth | Update treatment with invalid UUID format → bad request | `code:400, message:"invalid treatment id"` | `code:400, message:"invalid treatment id"` | ✅ |
| 16 | DELETE | `/api/treatment/:id` | no auth | Delete treatment without auth → unauthorized | `code:401` | `code:401` | ✅ |
| 17 | DELETE | `/api/treatment/00000000-...` | non-existent UUID + admin auth | Delete non-existent treatment → not found | `code:400, message:"treatment not found"` | `code:400, message:"treatment not found"` | ✅ |
| 18 | DELETE | `/api/treatment/:id` | already deleted + admin auth | Delete already deleted treatment → not found | `code:400, message:"treatment not found"` | `code:400, message:"treatment not found"` | ✅ |
| 19 | DELETE | `/api/treatment/bad-id` | invalid UUID + admin auth | Delete treatment with invalid UUID format → bad request | `code:400, message:"invalid treatment id"` | `code:400, message:"invalid treatment id"` | ✅ |
| 20 | GET | `/api/treatment/list?search=nonexistzzz` | no matches | Search treatments with non-matching keyword → empty result | `code:200, data:[]` | `code:200, data:[]` | ✅ |
| 21 | GET | `/api/treatment/list?page=999` | page out of range | List treatments with out-of-range page → empty result | `code:200, data:[]` | `code:200, data:[]` | ✅ |

**Result: 21/21 passed.**

---

## StaffSkill CRUD

Auth: admin JWT from `+6289666955155` (John, role=admin).
Seeded data: 3 skills for Jane Smith (Haircut Premium, Massage, Facial Treatment).

Access rules:
- `POST /api/staff-skills/` → `AuthMiddleware()` + `RequireRole("admin")`
- `GET /api/staff-skills/list` → public
- `GET /api/staff-skills/:id` → public
- `GET /api/staff-skills/staff/:userId` → public
- `GET /api/staff-skills/treatment/:treatmentId` → public
- `DELETE /api/staff-skills/:id` → `AuthMiddleware()` + `RequireRole("admin")`

### Positive Cases

| # | Method | Endpoint | Request / Body | Description | Expected | Actual | Result |
|---|--------|----------|----------------|-------------|----------|--------|--------|
| 1 | GET | `/api/staff-skills/list` | — | List all staff skills (public) | `code:200, data:3 items` | `code:200, data:3 items` | ✅ |
| 2 | GET | `/api/staff-skills/list?page=1&limit=2` | `page=1, limit=2` | List skills with pagination | `code:200, items:2, total:3` | `code:200, items:2, total:3` | ✅ |
| 3 | GET | `/api/staff-skills/:id` | path: skill ID | Get skill by ID with joined user/treatment names (public) | `code:200, user:"Jane Smith"` | `code:200, user:"Jane Smith"` | ✅ |
| 4 | GET | `/api/staff-skills/staff/:userId` | path: Jane's user ID | List all skills for a specific staff (public) | `code:200, 3 skills` | `code:200, 3 skills` | ✅ |
| 5 | GET | `/api/staff-skills/treatment/:treatmentId` | path: Haircut Premium ID | List all staff who can perform a treatment (public) | `code:200, 1 staff` | `code:200, 1 staff` | ✅ |
| 6 | POST | `/api/staff-skills/` | `{user_id, treatment_id}` + admin token | Assign skill (treatment) to staff (admin only) | `code:201, user:"Jane Smith"` | `code:201, user:"Jane Smith"` | ✅ |
| 7 | GET | `/api/staff-skills/list?search=jane` | `search=jane` | Search skills by staff name (ILIKE on joined user) | `code:200, 3 items` | `code:200, 3 items` | ✅ |
| 8 | GET | `/api/staff-skills/list?search=zzzzz` | `search=zzzzz` | Search skills with non-matching name → empty result | `code:200, 0 items` | `code:200, 0 items` | ✅ |
| 9 | GET | `/api/staff-skills/list?page=99` | `page=99` | List skills with out-of-range page → empty result | `code:200, 0 items` | `code:200, 0 items` | ✅ |

### Negative Cases

| # | Method | Endpoint | Request / Body | Description | Expected | Actual | Result |
|---|--------|----------|----------------|-------------|----------|--------|--------|
| 10 | POST | `/api/staff-skills/` | no auth | Assign skill without auth → unauthorized | `code:401` | `code:401` | ✅ |
| 11 | POST | `/api/staff-skills/` | `{}` (empty) | Assign skill with empty body → bad request | `code:400` | `code:400` | ✅ |
| 12 | POST | `/api/staff-skills/` | `{user_id: "bad"}` (invalid UUID) | Assign skill with invalid UUID format → bad request | `code:400` | `code:400` | ✅ |
| 13 | DELETE | `/api/staff-skills/:id` | no auth | Unassign skill without auth → unauthorized | `code:401` | `code:401` | ✅ |
| 14 | DELETE | `/api/staff-skills/:id` | non-existent UUID + admin auth | Unassign non-existent skill → not found | `code:400, msg:"staff skill not found"` | `code:400, msg:"staff skill not found"` | ✅ |
| 15 | DELETE | `/api/staff-skills/bad-id` | invalid UUID + admin auth | Unassign skill with invalid UUID format → bad request | `code:400` | `code:400` | ✅ |
| 16 | GET | `/api/staff-skills/00000000-...` | non-existent UUID | Get non-existent skill → not found | `code:400` | `code:400` | ✅ |
| 17 | GET | `/api/staff-skills/bad-id` | invalid UUID | Get skill with invalid UUID format → bad request | `code:400` | `code:400` | ✅ |
| 18 | GET | `/api/staff-skills/staff/00000000-...` | non-existent UUID | List skills for non-existent staff → empty result | `code:200, data:[]` | `code:200, data:[]` | ✅ |
| 19 | GET | `/api/staff-skills/staff/bad-id` | invalid UUID | List skills with invalid staff UUID format → bad request | `code:400` | `code:400` | ✅ |
| 20 | GET | `/api/staff-skills/treatment/00000000-...` | non-existent UUID | List staff for non-existent treatment → empty result | `code:200, data:[]` | `code:200, data:[]` | ✅ |
| 21 | GET | `/api/staff-skills/treatment/bad-id` | invalid UUID | List staff with invalid treatment UUID format → bad request | `code:400` | `code:400` | ✅ |

**Result: 21/21 passed.**

### Notes
- All response DTOs include joined `user_full_name`, `user_phone_number`, `treatment_name`, `treatment_price`
- Duplicate `(user_id, treatment_id)` detection at application layer returns 409
- No soft delete — junction table uses hard delete only
- Search uses JOIN + ILIKE on user's full_name/phone_number

---

## Summary

| Module | Tests | Passed | Failed |
|--------|-------|--------|--------|
| User / Auth | 39 | 39 | 0 |
| Treatment CRUD | 21 | 21 | 0 |
| StaffSkill CRUD | 21 | 21 | 0 |
| **Total** | **81** | **81** | **0** |
