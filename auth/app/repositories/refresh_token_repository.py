from sqlalchemy.orm import Session
from datetime import datetime
from typing import Optional
from app.models.refresh_token import RefreshToken
import uuid

class RefreshTokenRepository:
    def __init__(self, db: Session):
        self.db = db

    def create_token(self, token: str, user_id: uuid.UUID, expires_at: datetime) -> RefreshToken:
        db_token = RefreshToken(
            token=token,
            user_id=user_id,
            expires_at=expires_at
        )
        self.db.add(db_token)
        self.db.commit()
        self.db.refresh(db_token)
        return db_token

    def get_valid_token(self, token: str) -> Optional[RefreshToken]:
        return self.db.query(RefreshToken).filter(
            RefreshToken.token == token,
            RefreshToken.is_valid == True,
            RefreshToken.expires_at > datetime.utcnow()
        ).first()

    def invalidate_user_tokens(self, user_id: uuid.UUID) -> None:
        self.db.query(RefreshToken).filter(
            RefreshToken.user_id == user_id,
            RefreshToken.is_valid == True
        ).update({"is_valid": False})
        self.db.commit()

    def invalidate_token(self, token: str) -> None:
        self.db.query(RefreshToken).filter(
            RefreshToken.token == token
        ).update({"is_valid": False})
        self.db.commit() 