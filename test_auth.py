import subprocess, json, time, sys

BASE = "http://localhost:8080/api"

def req(method, path, body=None, token=None):
    url = BASE + path
    cmd = ["curl", "-s", "-w", "\n%{http_code}", "-X", method, url, "-H", "Content-Type: application/json"]
    if token:
        bear = "Authorization: Bearer " + token
        cmd += ["-H", bear]
    if body:
        cmd += ["-d", json.dumps(body)]
    r = subprocess.run(cmd, capture_output=True, text=True, timeout=10)
    out = r.stdout.strip()
    if not out:
        return {"error": "empty"}, 0
    parts = out.rsplit("\n", 1)
    code = parts[-1]
    data = parts[0] if len(parts) > 1 else ""
    try:
        return json.loads(data) if data else {}, int(code)
    except:
        return {"raw": data}, int(code) if code.isdigit() else 0

def get_token(phone):
    r, c = req("POST", "/auth/login", {"phone": phone})
    if c != 200:
        return None, "login: " + str(c)
    time.sleep(0.3)
    r2 = subprocess.run(
        ["docker", "exec", "jadwalin_redis", "redis-cli", "GET", "login-otp:" + phone],
        capture_output=True, text=True, timeout=5
    )
    otp = r2.stdout.strip()
    if not otp:
        return None, "no otp"
    r3, c3 = req("POST", "/auth/login-verify", {"phone": phone, "otp": otp})
    if c3 != 200:
        return None, "verify: " + str(c3)
    token = r3.get("data", {}).get("token", "")
    if not token:
        return None, "no token field"
    return token, None

admin_token, err = get_token("6281234567890")
print("Admin token:", "OK" if admin_token else "FAIL: " + (err or "?"))
emp_token, err = get_token("6281122334455")
print("Emp token:  ", "OK" if emp_token else "FAIL: " + (err or "?"))

if not admin_token or not emp_token:
    sys.exit(1)

print()
print("=" * 60)
print("AUTH MIDDLEWARE TESTS")
print("=" * 60)

tests = [
    ("[1] List (auth)",        "GET",    "/user/list?role=", None, admin_token, 200),
    ("[2] List (no auth)",     "GET",    "/user/list",       None, None,        401),
    ("[3] Delete (admin)",     "DELETE", "/user/00000000-0000-0000-0000-000000000000", None, admin_token, 404),
    ("[4] Delete (no auth)",   "DELETE", "/user/00000000-0000-0000-0000-000000000000", None, None,        401),
    ("[5] Delete (staff)",     "DELETE", "/user/00000000-0000-0000-0000-000000000000", None, emp_token,   403),
    ("[6] RegStaff (no auth)", "POST",   "/user/register-staff",
     {"phone_number":"+628***e1","email":"e1@t.com","full_name":"N","address":"a"}, None, 401),
    ("[7] RegStaff (admin)",   "POST",   "/user/register-staff",
     {"phone_number":"+6281111111117","email":"e2@t.com","full_name":"A","address":"a"}, admin_token, 201),
    ("[8] RegStaff (staff)",   "POST",   "/user/register-staff",
     {"phone_number":"+628***e3","email":"e3@t.com","full_name":"E","address":"a"}, emp_token, 403),
    ("[9] Update (staff)",     "PUT",    "/user/00000000-0000-0000-0000-000000000000",
     {"full_name":"x"}, emp_token, 404),
    ("[10] Update (no auth)",  "PUT",    "/user/00000000-0000-0000-0000-000000000000",
     {"full_name":"x"}, None, 401),
]

passed = 0
failed = 0
results = []
for name, method, path, body, token, expected in tests:
    r, c = req(method, path, body, token)
    ok = c == expected
    msg = r.get("message", r.get("raw", ""))[:50]
    results.append((name, method, path, c, expected, ok, msg))
    if ok:
        passed += 1
    else:
        failed += 1
    status = "PASS" if ok else "FAIL"
    print("  [" + status + "] " + name + ": " + str(c) + " (expected " + str(expected) + ")")

print()
print("=" * 60)
print("Result: " + str(passed) + "/" + str(passed+failed) + " passed")
if failed:
    print("FAILURES: " + str(failed))
    for name, method, path, c, expected, ok, msg in results:
        if not ok:
            print("  " + name + ": got " + str(c) + " want " + str(expected) + " - " + msg)
