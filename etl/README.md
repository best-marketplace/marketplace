# ETL Service Documentation

## Overview
ETL (Extract, Transform, Load) сервис предназначен для обработки данных о продуктах из Kafka и их индексации в Elasticsearch. Сервис обеспечивает непрерывную обработку событий о создании продуктов и их сохранение в поисковом индексе.

## Структура данных

### Входные данные (Kafka)
Сервис ожидает события в формате:
```json
{
  "action": "product_created",
  "product_id": "uuid",
  "title": "Название продукта",
  "description": "Описание продукта"
}
```

### Выходные данные (Elasticsearch)
Данные индексируются в Elasticsearch со следующей структурой:
```json
{
  "action": "product_created",
  "product_id": "uuid",
  "title": "Название продукта",
  "description": "Описание продукта"
}
```

## Конфигурация

### Kafka
Настройки подключения к Kafka задаются через переменные окружения:
- `KAFKA_BOOTSTRAP_SERVERS` - адреса брокеров Kafka (по умолчанию: kafka:9092)
- `KAFKA_TOPIC` - топик для чтения событий (по умолчанию: product-events)
- `group_id` - идентификатор группы потребителей (по умолчанию: product-etl-group)

### Elasticsearch
Настройки подключения к Elasticsearch:
- Хост: http://elasticsearch:9200
- Индекс: products

## Маппинг Elasticsearch
Индекс использует следующие поля:
- `action` (keyword) - тип действия
- `product_id` (keyword) - уникальный идентификатор продукта
- `title` (text) - название продукта с русской морфологией
- `description` (text) - описание продукта с русской морфологией

## Запуск сервиса

### Требования
- Python 3.8+
- Docker
- Kafka
- Elasticsearch

### Запуск через Docker
```bash
docker-compose up -d product-etl
```

### Проверка работоспособности
1. Проверка логов:
```bash
docker logs product-etl
```

2. Проверка индексации в Elasticsearch:
```bash
curl -X GET "http://localhost:9200/products/_search" -H "Content-Type: application/json" -d '{
  "query": {
    "match_all": {}
  }
}'
```

## Обработка ошибок
Сервис включает механизмы обработки ошибок:
- Логирование всех операций
- Повторные попытки при сбоях подключения
- Graceful shutdown при получении сигналов завершения

## Мониторинг
Сервис логирует следующие события:
- Инициализация компонентов
- Получение сообщений из Kafka
- Трансформация данных
- Индексация в Elasticsearch
- Ошибки и исключения

## Структура проекта
```
etl/
├── main.py              # Основной код ETL сервиса
├── config.py            # Конфигурация сервиса
├── services/
│   ├── kafka_service.py      # Сервис для работы с Kafka
│   └── elasticsearch_service.py  # Сервис для работы с Elasticsearch
└── README.md            # Документация
``` 