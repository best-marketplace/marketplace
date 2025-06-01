from datetime import timedelta
from typing import Annotated
from fastapi import FastAPI, Depends, HTTPException, status, Request, Form
from fastapi.security import OAuth2PasswordRequestForm, OAuth2PasswordBearer
from fastapi.openapi.docs import get_swagger_ui_html
from fastapi.openapi.utils import get_openapi
from sqlalchemy.orm import Session
from . import models, schemas, auth
from .database import engine, get_db
import jwt
from jose import JWTError
from .settings import settings

models.Base.metadata.create_all(bind=engine)

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token", auto_error=False)

app = FastAPI(
    title="Authentication Service",
    description="Сервис аутентификации для маркетплейса",
    version="1.0.0",
    docs_url="/auth/docs",
    redoc_url="/auth/redoc",
    openapi_url="/auth/openapi.json"
)

def custom_openapi():
    if app.openapi_schema:
        return app.openapi_schema
    
    openapi_schema = get_openapi(
        title="Authentication Service API",
        version="1.0.0",
        description="Сервис аутентификации для маркетплейса",
        routes=app.routes,
    )
    
    # Определяем схему для UserRole
    user_role_schema = {
        "type": "string",
        "enum": ["customer", "seller"],
        "description": "User role in the system"
    }
    
    # Схема для ошибок валидации
    validation_error = {
        "type": "object",
        "properties": {
            "detail": {
                "type": "array",
                "items": {
                    "type": "object",
                    "properties": {
                        "loc": {"type": "array", "items": {"type": "string"}},
                        "msg": {"type": "string"},
                        "type": {"type": "string"}
                    }
                }
            }
        }
    }
    
    # Базовые схемы
    schemas_dict = {
        "UserRole": user_role_schema,
        "UserBase": {
            "type": "object",
            "properties": {
                "email": {"type": "string", "format": "email"},
                "username": {"type": "string", "minLength": 3, "maxLength": 255},
                "role": {"$ref": "#/components/schemas/UserRole"}
            },
            "required": ["email", "username", "role"]
        },
        "UserCreate": {
            "type": "object",
            "properties": {
                "email": {"type": "string", "format": "email"},
                "username": {"type": "string", "minLength": 3, "maxLength": 255},
                "password": {"type": "string", "minLength": 8, "maxLength": 255},
                "role": {"$ref": "#/components/schemas/UserRole"}
            },
            "required": ["email", "username", "password", "role"]
        },
        "UserResponse": {
            "type": "object",
            "properties": {
                "id": {"type": "string", "format": "uuid"},
                "email": {"type": "string", "format": "email"},
                "username": {"type": "string"},
                "role": {"$ref": "#/components/schemas/UserRole"},
                "is_active": {"type": "boolean"},
                "created_at": {"type": "string", "format": "date-time"},
                "updated_at": {"type": ["string", "null"], "format": "date-time"}
            },
            "required": ["id", "email", "username", "role", "is_active", "created_at"]
        },
        "Token": {
            "type": "object",
            "properties": {
                "access_token": {"type": "string"},
                "token_type": {"type": "string"}
            },
            "required": ["access_token", "token_type"]
        },
        "TokenData": {
            "type": "object",
            "properties": {
                "user_id": {"type": ["string", "null"]},
                "email": {"type": ["string", "null"]}
            }
        }
    }
    
    # Обновляем компоненты OpenAPI схемы
    openapi_schema["components"] = {
        "securitySchemes": {
            "Bearer": {
                "type": "http",
                "scheme": "bearer",
                "bearerFormat": "JWT",
            }
        },
        "schemas": {
            **schemas_dict,
            "HTTPValidationError": validation_error
        }
    }
    
    # Добавляем требование Bearer аутентификации для всех операций
    for path in openapi_schema["paths"].values():
        for operation in path.values():
            operation["security"] = [{"Bearer": []}]
            
    # Вручную добавляем описание requestBody для эндпоинта auth/token
    if "/auth/token" in openapi_schema["paths"]:
        if "post" in openapi_schema["paths"]["/auth/token"]:
            openapi_schema["paths"]["/auth/token"]["post"]["requestBody"] = {
                "content": {
                    "application/x-www-form-urlencoded": {
                        "schema": {
                            "type": "object",
                            "properties": {
                                "username": {
                                    "type": "string",
                                    "description": "Your email address",
                                    "example": "user@example.com"
                                },
                                "password": {
                                    "type": "string",
                                    "description": "Your password",
                                    "example": "password123"
                                }
                            },
                            "required": [
                                "username",
                                "password"
                            ]
                        }
                    }
                },
                "required": True
            }

    app.openapi_schema = openapi_schema
    return app.openapi_schema

app.openapi = custom_openapi

@app.get("/auth/docs", include_in_schema=False)
async def custom_swagger_ui_html():
    return get_swagger_ui_html(
        openapi_url="/auth/openapi.json",
        title="Authentication Service API",
        oauth2_redirect_url="/auth/docs/oauth2-redirect",
        swagger_js_url="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.9.0/swagger-ui-bundle.js",
        swagger_css_url="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.9.0/swagger-ui.css",
        swagger_favicon_url="https://fastapi.tiangolo.com/img/favicon.png",
    )

@app.get("/auth/openapi.json", include_in_schema=False)
async def get_open_api_endpoint():
    return app.openapi()

@app.get("/auth/health")
def health_check():
    return {"status": "healthy"}

@app.post("/auth/register", response_model=schemas.UserResponse)
def register_user(
    user: schemas.UserCreate,
    db: Session = Depends(get_db)
):
    db_user = db.query(models.User).filter(models.User.email == user.email).first()
    if db_user:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Email already registered"
        )
    
    hashed_password = auth.get_password_hash(user.password)
    db_user = models.User(
        email=user.email,
        username=user.username,
        password_hash=hashed_password,
        role=user.role
    )
    
    db.add(db_user)
    db.commit()
    db.refresh(db_user)
    
    return schemas.UserResponse.from_orm(db_user)

@app.post("/auth/token", response_model=schemas.Token)
def login_for_access_token(
    username: str = Form(..., description="Your email address", example="user@example.com"),
    password: str = Form(..., description="Your password", example="password123"),
    db: Session = Depends(get_db)
):
    """
    Login endpoint to get access token.
    
    - **username**: Your email address
    - **password**: Your password
    
    Returns a JWT token that can be used for authenticated requests.
    """
    user = auth.authenticate_user(db, username, password)
    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect email or password",
            headers={"WWW-Authenticate": "Bearer"},
        )
    
    access_token_expires = timedelta(minutes=settings.ACCESS_TOKEN_EXPIRE_MINUTES)
    access_token = auth.create_access_token(
        data={"sub": str(user.id)},
        expires_delta=access_token_expires
    )
    
    return {"access_token": access_token, "token_type": "bearer"}

@app.get("/auth/profile", response_model=schemas.UserResponse)
def read_user_profile(
    current_user: Annotated[models.User, Depends(auth.get_current_user)]
):
    return schemas.UserResponse.from_orm(current_user)

@app.post("/auth/logout")
def logout(
    current_user: Annotated[models.User, Depends(auth.get_current_user)]
):
    return {"message": "Successfully logged out"}

@app.post("/auth/verify")
async def verify_token(
    request: Request,
    db: Session = Depends(get_db)
):
    auth_header = request.headers.get("Authorization")
    if not auth_header or not auth_header.startswith("Bearer "):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid authentication credentials",
            headers={"WWW-Authenticate": "Bearer"},
        )
    
    token = auth_header.split(" ")[1]
    try:
        payload = jwt.decode(token, settings.JWT_SECRET_KEY, algorithms=[settings.JWT_ALGORITHM])
        user_id: str = payload.get("sub")
        if user_id is None:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Invalid token",
                headers={"WWW-Authenticate": "Bearer"},
            )
        
        user = db.query(models.User).filter(models.User.id == user_id).first()
        if user is None:
             raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="User not found",
                headers={"WWW-Authenticate": "Bearer"},
            )

        return {"status": "valid", "user_id": user_id, "role": user.role}
    except JWTError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid token",
            headers={"WWW-Authenticate": "Bearer"},
        )

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=5000)

