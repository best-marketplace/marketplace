# Работа с Elasticsearch

## Запуск

Elasticsearch запускается автоматически вместе с остальными сервисами при выполнении команды:
```bash
docker-compose up --build
```

Сервис доступен по адресу: `http://localhost:9200`

## Проверка работоспособности

Чтобы проверить, что Elasticsearch работает корректно, выполните:
```bash
curl http://localhost:9200
```

Вы должны получить ответ с информацией о версии и статусе кластера.

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

### 3. Поиск документов

#### Простой поиск по названию
```bash
curl -X GET "localhost:9200/products/_search" -H "Content-Type: application/json" -d'
{
  "query": {
    "match": {
      "name": "тестовый"
    }
  }
}'
```

#### Поиск с фильтрацией по цене
```bash
curl -X GET "localhost:9200/products/_search" -H "Content-Type: application/json" -d'
{
  "query": {
    "bool": {
      "must": [
        { "match": { "name": "тестовый" } }
      ],
      "filter": [
        { "range": { "price": { "gte": 50, "lte": 100 } } }
      ]
    }
  }
}'
```

#### Поиск по категории
```bash
curl -X GET "localhost:9200/products/_search" -H "Content-Type: application/json" -d'
{
  "query": {
    "term": {
      "category": "electronics"
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

## Мониторинг

### Проверка здоровья кластера
```bash
curl http://localhost:9200/_cluster/health
```

### Получение списка индексов
```bash
curl http://localhost:9200/_cat/indices?v
```

### Получение статистики
```bash
curl http://localhost:9200/_stats
```

## Интеграция с приложением

Для интеграции Elasticsearch с вашим Go-приложением:

1. Добавьте зависимость в `go.mod`:
```go
require (
    github.com/olivere/elastic/v7 v7.0.32
)
```

2. Пример кода для подключения:
```go
import (
    "github.com/olivere/elastic/v7"
)

func NewElasticsearchClient() (*elastic.Client, error) {
    return elastic.NewClient(
        elastic.SetURL("http://elasticsearch:9200"),
        elastic.SetSniff(false),
    )
}
```

3. Пример использования в коде:
```go
type Product struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    Category    string    `json:"category"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

func IndexProduct(client *elastic.Client, product *Product) error {
    _, err := client.Index().
        Index("products").
        Id(product.ID).
        BodyJson(product).
        Do(context.Background())
    return err
}

func SearchProducts(client *elastic.Client, query string) ([]Product, error) {
    searchResult, err := client.Search().
        Index("products").
        Query(elastic.NewMatchQuery("name", query)).
        Do(context.Background())
    
    if err != nil {
        return nil, err
    }

    var products []Product
    for _, hit := range searchResult.Hits.Hits {
        var product Product
        err := json.Unmarshal(hit.Source, &product)
        if err != nil {
            return nil, err
        }
        products = append(products, product)
    }
    
    return products, nil
}
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