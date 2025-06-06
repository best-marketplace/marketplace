from fastapi import FastAPI
from app.handlers.auth_handler import router as auth_router
from app.config import settings
from app.database import engine
from app.models import Base, User, RefreshToken 

app = FastAPI(
    title=settings.PROJECT_NAME,
    version=settings.VERSION,
)

Base.metadata.create_all(bind=engine)

app.include_router(
    auth_router,
    prefix="/auth",
    tags=["auth"]
)