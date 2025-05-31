import os
from dotenv import load_dotenv

load_dotenv()

KAFKA_CONFIG = {
    'bootstrap_servers': os.getenv('KAFKA_BOOTSTRAP_SERVERS', 'kafka:9092'),
    'topic': os.getenv('KAFKA_TOPIC', 'product-events'),
    'group_id': 'product-etl-group'
}

ELASTICSEARCH_CONFIG = {
    'hosts': ['http://elasticsearch:9200'],
    'index': 'products'
}

# Index mapping for Elasticsearch
INDEX_MAPPING = {
    'settings': {
        'analysis': {
            'analyzer': {
                'russian_analyzer': {
                    'type': 'custom',
                    'tokenizer': 'standard',
                    'filter': [
                        'lowercase',
                        'russian_stop',
                        'russian_stemmer'
                    ]
                }
            }
        }
    },
    'mappings': {
        'properties': {
            'product_id': {'type': 'keyword'},
            'name': {
                'type': 'text',
                'analyzer': 'russian_analyzer',
                'fields': {
                    'keyword': {
                        'type': 'keyword',
                        'ignore_above': 256
                    }
                }
            },
            'description': {
                'type': 'text',
                'analyzer': 'russian_analyzer'
            },
            'price': {'type': 'float'},
            'categoryName': {
                'type': 'text',
                'analyzer': 'russian_analyzer',
                'fields': {
                    'keyword': {
                        'type': 'keyword',
                        'ignore_above': 256
                    }
                }
            },
            'sellerName': {
                'type': 'text',
                'analyzer': 'russian_analyzer',
                'fields': {
                    'keyword': {
                        'type': 'keyword',
                        'ignore_above': 256
                    }
                }
            },
            'created_at': {'type': 'date'},
            'updated_at': {'type': 'date'}
        }
    }
} 