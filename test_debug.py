import subprocess, json, time, sys

BASE = "http://localhost:8080/api"

def req(method, path, body=None, token=None):
    url = BASE + path
    cmd = ["curl", "-s", "-w", "\n%{http_code}", "-X", method, url, "-H", "Content-Type: application/json"]
    if token:
        cmd += ["-H", "Authorization: Bearer *** + token]
    if body:
        cmd += ["-d", json.dumps(body)]
    print(f"  DEBUG: running: curl -s -w ... {method} {url}", file=sys.stderr)
    r = subprocess.run(cmd, capture_output=True, text=True, timeout=10)
    print(f"  DEBUG: stdout len={len(r.stdout)}", file=sys.stderr)
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

# Test one request
print("Testing single request...", file=sys.stderr)
r, c = req("POST", "/api/auth/login", {"phone": "6281234567890"})
print(f"Result: {r}, code: {c}", file=sys.stderr)
print(json.dumps({"result": r, "code": c}))
