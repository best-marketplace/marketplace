from fastapi import FastAPI, Request, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse, Response
import httpx
import os
from dotenv import load_dotenv

load_dotenv()

app = FastAPI(title="API Gateway")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

SERVICES = {
    "auth": os.getenv("AUTH_SERVICE_URL", "http://auth:5000"),
    "catalog": os.getenv("CATALOG_SERVICE_URL", "http://catalog:8080"),
}

PUBLIC_ROUTES = {
    ("/auth/token", "POST"),
    ("/auth/register", "POST"),
    ("/auth/verify", "POST"),
    ("/auth/docs", "GET"),
    ("/auth/redoc", "GET"),
    ("/auth/openapi.json", "GET"),
    ("/auth/health", "GET"),
    ("/catalog/products/", "GET"),
    ("/catalog/product/", "GET"),
    ("/catalog/comments/", "GET"),
    ("/catalog/search/", "GET"),
}

PROTECTED_ROUTES = {
    "/catalog/product/": {"POST": ["seller"]},
    "/catalog/comment/": {"POST": ["customer", "seller"]},
}

async def verify_token(request: Request):
    auth_header = request.headers.get("Authorization")
    if not auth_header or not auth_header.startswith("Bearer "):
        return None

    token = auth_header.split(" ")[1]
    async with httpx.AsyncClient() as client:
        try:
            response = await client.post(
                f"{SERVICES['auth']}/auth/verify",
                headers={"Authorization": f"Bearer {token}"}
            )
            if response.status_code == 200:
                data = response.json()
                if data.get("status") == "valid":
                    return {
                        "token": token,
                        "user_id": data.get("user_id"),
                        "role": data.get("role")
                    }
        except Exception:
            return None
    return None

@app.middleware("http")
async def auth_middleware(request: Request, call_next):
    path = request.url.path
    method = request.method

    if (path, method) in PUBLIC_ROUTES:
        return await call_next(request)

    # Check for protected route with required role
    requires_auth = False
    required_roles = []
    for prefix, methods_roles in PROTECTED_ROUTES.items():
        if path.startswith(prefix) and method in methods_roles:
            requires_auth = True
            required_roles = methods_roles[method]
            break

    verification_result = None
    if requires_auth or path.startswith("/catalog/"):
        verification_result = await verify_token(request)
        if not verification_result:
            return JSONResponse(status_code=401, content={"detail": "Unauthorized"})
        if required_roles and verification_result["role"] not in required_roles:
            return JSONResponse(status_code=403, content={"detail": "Forbidden: Insufficient permissions"})

        request.headers.__dict__["_list"].append((b"x-user-id", verification_result["user_id"].encode()))
        request.headers.__dict__["_list"].append((b"x-user-role", verification_result["role"].encode()))
        request.headers.__dict__["_list"].append((b"authorization", f"Bearer {verification_result['token']}".encode()))

    return await call_next(request)

@app.api_route("/{service}/{path:path}", methods=["GET", "POST", "PUT", "DELETE", "PATCH"])
async def proxy_request(service: str, path: str, request: Request):
    if service not in SERVICES:
        raise HTTPException(status_code=404, detail="Service not found")

    if service == "auth":
        target_url = f"{SERVICES[service]}/auth/{path}"
    elif service == "catalog":
        target_url = f"{SERVICES[service]}/api/{path}"
    else:
        target_url = f"{SERVICES[service]}/{path}"

    headers = dict(request.headers)
    headers.pop("host", None)

    try:
        body = await request.body()
        async with httpx.AsyncClient() as client:
            response = await client.request(
                method=request.method,
                url=target_url,
                headers=headers,
                content=body,
                params=request.query_params
            )

        return Response(
            content=response.content,
            status_code=response.status_code,
            media_type=response.headers.get("content-type", "application/json")
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
