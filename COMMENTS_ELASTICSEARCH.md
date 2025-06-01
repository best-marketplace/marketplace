# Работа с комментариями в Elasticsearch

## Структура данных

Комментарии хранятся в индексе `comments` со следующей структурой:
```json
{
    "comment_id": "uuid комментария",
    "comment": "текст комментария"
}
```

## Основные операции

### 1. Проверка существования индекса
```bash
curl -X GET "http://localhost:9200/comments?pretty"
```

### 2. Проверка маппинга индекса
```bash
curl -X GET "http://localhost:9200/comments/_mapping?pretty"
```

### 3. Поиск комментариев

#### Поиск всех комментариев
```bash
curl -X GET "http://localhost:9200/comments/_search?pretty" -H "Content-Type: application/json" -d '{
    "query": {
        "match_all": {}
    }
}'
```

#### Поиск по тексту комментария
```bash
curl -X GET "http://localhost:9200/comments/_search?pretty" -H "Content-Type: application/json" -d '{
    "query": {
        "match": {
            "comment": "хорошая куртка"
        }
    }
}'
```

#### Поиск по ID комментария
```bash
curl -X GET "http://localhost:9200/comments/_search?pretty" -H "Content-Type: application/json" -d '{
    "query": {
        "term": {
            "comment_id": "00000000-0000-0000-0000-000000000006"
        }
    }
}'
```

### 4. Получение конкретного комментария по ID
```bash
curl -X GET "http://localhost:9200/comments/_doc/00000000-0000-0000-0000-000000000006?pretty"
```

### 5. Удаление комментария
```bash
curl -X DELETE "http://localhost:9200/comments/_doc/00000000-0000-0000-0000-000000000006"
```

### 6. Удаление всего индекса комментариев
```bash
curl -X DELETE "http://localhost:9200/comments"
```

## Примеры использования

### 1. Проверка добавления комментария
После создания комментария через API, проверьте его наличие в Elasticsearch:
```bash
curl -X GET "http://localhost:9200/comments/_search?pretty" -H "Content-Type: application/json" -d '{
    "query": {
        "term": {
            "comment_id": "00000000-0000-0000-0000-000000000006"
        }
    }
}'
```

### 2. Поиск комментариев по ключевым словам
```bash
curl -X GET "http://localhost:9200/comments/_search?pretty" -H "Content-Type: application/json" -d '{
    "query": {
        "match": {
            "comment": "хороший"
        }
    }
}'
```

### 3. Получение статистики по индексу
```bash
curl -X GET "http://localhost:9200/comments/_stats?pretty"
```

## Примечания

1. Все запросы к Elasticsearch выполняются на порту 9200
2. Для работы с русским языком используется анализатор "russian"
3. Поле comment_id используется как уникальный идентификатор документа
4. Текст комментария индексируется с поддержкой русского языка

## Отладка

### Проверка анализатора
```bash
curl -X GET "http://localhost:9200/comments/_analyze?pretty" -H "Content-Type: application/json" -d '{
    "analyzer": "russian_analyzer",
    "text": "хорошая куртка, отличное качество"
}'
```

### Проверка настроек индекса
```bash
curl -X GET "http://localhost:9200/comments/_settings?pretty"
``` 