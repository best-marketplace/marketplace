# Работа с Elasticsearch

## Запуск

Elasticsearch запускается автоматически вместе с остальными сервисами при выполнении команды:
```bash
docker-compose up --build
```

Сервис доступен по адресу: `http://localhost:9200`

## Проверка работоспособности

1. Проверка статуса кластера:
```bash
curl -X GET "http://localhost:9200/_cluster/health?pretty"
```

2. Проверка списка индексов:
```bash
curl -X GET "http://localhost:9200/_cat/indices?v"
```

3. Проверка маппинга индекса products:
```bash
curl -X GET "http://localhost:9200/products/_mapping?pretty"
```

## Основные операции

### 1. Создание индекса для продуктов
```bash
curl -X PUT "localhost:9200/products" -H "Content-Type: application/json" -d'
{
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "name": { "type": "text", "analyzer": "russian" },
      "description": { "type": "text", "analyzer": "russian" },
      "price": { "type": "float" },
      "category": { "type": "keyword" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}'
```

### 2. Добавление документа
```bash
curl -X POST "localhost:9200/products/_doc" -H "Content-Type: application/json" -d'
{
  "id": "123",
  "name": "Тестовый продукт",
  "description": "Это тестовый продукт для проверки работы Elasticsearch",
  "price": 99.99,
  "category": "electronics",
  "created_at": "2024-03-20T12:00:00",
  "updated_at": "2024-03-20T12:00:00"
}'
```

### 3. Поиск товаров

#### Простой поиск по названию
```bash
curl -X GET "http://localhost:9200/products/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "match": {
      "title": "ноутбук"
    }
  }
}'
```

#### Поиск по названию и описанию
```bash
curl -X GET "http://localhost:9200/products/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "multi_match": {
      "query": "мощный",
      "fields": ["title", "description"]
    }
  }
}'
```

#### Поиск с фильтрацией по дате
```bash
curl -X GET "http://localhost:9200/products/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "bool": {
      "must": [
        { "match": { "title": "ноутбук" } }
      ],
      "filter": [
        {
          "range": {
            "created_at": {
              "gte": "2024-01-01"
            }
          }
        }
      ]
    }
  }
}'
```

#### Поиск по точному совпадению product_id
```bash
curl -X GET "http://localhost:9200/products/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "term": {
      "product_id": "p123"
    }
  }
}'
```

### 4. Получение документа по ID
```bash
curl -X GET "localhost:9200/products/_doc/123"
```

### 5. Обновление документа
```bash
curl -X POST "localhost:9200/products/_update/123" -H "Content-Type: application/json" -d'
{
  "doc": {
    "price": 89.99,
    "updated_at": "2024-03-20T13:00:00"
  }
}'
```

### 6. Удаление документа
```bash
curl -X DELETE "localhost:9200/products/_doc/123"
```

## Агрегации

### 1. Подсчет товаров по категориям
```bash
curl -X GET "http://localhost:9200/products/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "size": 0,
  "aggs": {
    "categories": {
      "terms": {
        "field": "category.keyword"
      }
    }
  }
}'
```

## Управление индексом

### 1. Удаление индекса
```bash
curl -X DELETE "http://localhost:9200/products"
```

### 2. Создание индекса с маппингом
```bash
curl -X PUT "http://localhost:9200/products" -H 'Content-Type: application/json' -d'
{
  "mappings": {
    "properties": {
      "product_id": { "type": "keyword" },
      "title": { "type": "text", "analyzer": "russian" },
      "description": { "type": "text", "analyzer": "russian" },
      "created_at": { "type": "date" }
    }
  }
}'
```

## Мониторинг

### 1. Статистика индекса
```bash
curl -X GET "http://localhost:9200/products/_stats?pretty"
```

### 2. Информация о кластере
```bash
curl -X GET "http://localhost:9200/_cluster/stats?pretty"
```

## Полезные команды для отладки

### 1. Проверка настроек анализатора
```bash
curl -X GET "http://localhost:9200/products/_analyze?pretty" -H 'Content-Type: application/json' -d'
{
  "analyzer": "russian",
  "text": "Мощный ноутбук с 16ГБ RAM"
}'
```

### 2. Получение документа по ID
```bash
curl -X GET "http://localhost:9200/products/_doc/p123?pretty"
```

## Примечания

1. Все запросы к Elasticsearch выполняются на порту 9200
2. Для работы с русским языком используется анализатор "russian"
3. Поле product_id используется как уникальный идентификатор документа
4. Все даты хранятся в формате UTC

## Примеры использования в Python

```python
from elasticsearch import Elasticsearch

# Подключение к Elasticsearch
es = Elasticsearch(['http://localhost:9200'])

# Поиск товаров
def search_products(query):
    response = es.search(
        index="products",
        body={
            "query": {
                "multi_match": {
                    "query": query,
                    "fields": ["title", "description"]
                }
            }
        }
    )
    return response['hits']['hits']

# Получение товара по ID
def get_product(product_id):
    return es.get(index="products", id=product_id)

# Создание нового товара
def create_product(product_data):
    return es.index(index="products", id=product_data['product_id'], body=product_data)
```

## Важные замечания

1. Данные Elasticsearch сохраняются в именованном томе `elasticsearch-data`
2. Отключена встроенная безопасность для упрощения разработки
3. Выделено 512MB памяти для JVM
4. Настроен memory lock для улучшения производительности
5. Настроен healthcheck для автоматического перезапуска при проблемах

## Очистка данных

Для полной очистки данных Elasticsearch:
```bash
docker-compose down -v
```

Это удалит все данные и тома. Используйте с осторожностью!

## Решение проблем

### 1. Если Elasticsearch не запускается
Проверьте логи:
```bash
docker-compose logs elasticsearch
```

### 2. Если не хватает памяти
Увеличьте значение `ES_JAVA_OPTS` в `docker-compose.yaml`:
```yaml
environment:
  - "ES_JAVA_OPTS=-Xms1g -Xmx1g"
```

### 3. Если возникают проблемы с доступом
Проверьте, что порт 9200 не занят другим процессом:
```bash
netstat -an | grep 9200
``` 