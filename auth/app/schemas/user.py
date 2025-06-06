from pydantic import BaseModel, EmailStr, UUID4
from datetime import datetime
from typing import Optional
from app.models.user import UserRole

class UserBase(BaseModel):
    email: EmailStr
    username: str
    role: UserRole

    model_config = {
        "arbitrary_types_allowed": True
    }

class UserCreate(UserBase):
    password: str

class UserUpdate(BaseModel):
    username: Optional[str] = None
    password: Optional[str] = None

class UserInDB(UserBase):
    id: UUID4
    created_at: datetime
    updated_at: datetime

    model_config = {
        "arbitrary_types_allowed": True
    }

class Token(BaseModel):
    access_token: str
    refresh_token: str
    token_type: str = "bearer"

class TokenData(BaseModel):
    user_id: Optional[UUID4] = None
    role: Optional[str] = None 