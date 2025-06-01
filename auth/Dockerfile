FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY . .

RUN apt-get update && \
    apt-get install -y netcat-traditional postgresql-client curl && \
    rm -rf /var/lib/apt/lists/*

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:5000/health || exit 1

ENV PYTHONPATH=/app
ENV DATABASE_URL=postgresql://marketplace_user:marketplace_password@postgres-auth:5432/marketplace_auth
ENV JWT_SECRET_KEY=your-secret-key-here

EXPOSE 5000

CMD ["uvicorn", "src.main:app", "--host", "0.0.0.0", "--port", "5000", "--reload"] 