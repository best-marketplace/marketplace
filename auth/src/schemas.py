from pydantic import BaseModel, EmailStr, constr, ConfigDict
from typing import Optional
from datetime import datetime
from .models import UserRole
import uuid

class UserBase(BaseModel):
    email: EmailStr
    username: constr(min_length=3, max_length=255)
    role: UserRole

    model_config = ConfigDict(from_attributes=True)

class UserCreate(UserBase):
    password: constr(min_length=8, max_length=255)

class UserLogin(BaseModel):
    email: EmailStr
    password: str

    class Config:
        json_schema_extra = {
            "example": {
                "email": "user@example.com",
                "password": "password123"
            }
        }

class UserResponse(UserBase):
    id: str
    is_active: bool
    created_at: datetime
    updated_at: Optional[datetime] = None

    model_config = ConfigDict(from_attributes=True)

    @classmethod
    def from_orm(cls, obj):
        print(f"Converting UUID to string: {obj.id} -> {str(obj.id)}")
        data = {
            'id': str(obj.id),
            'email': obj.email,
            'username': obj.username,
            'role': obj.role,
            'is_active': obj.is_active,
            'created_at': obj.created_at,
            'updated_at': obj.updated_at
        }
        return cls(**data)

class Token(BaseModel):
    access_token: str
    token_type: str

class TokenData(BaseModel):
    user_id: Optional[str] = None
    email: Optional[str] = None 