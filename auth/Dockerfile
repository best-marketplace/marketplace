FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .
COPY entrypoint.sh .

ENV PYTHONPATH=/app

CMD ["./entrypoint.sh"] 