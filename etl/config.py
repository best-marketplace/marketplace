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
                    'type': 'russian'
                }
            }
        }
    },
    'mappings': {
        'properties': {
            'product_id': {'type': 'keyword'},
            'title': {
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
            }
        }
    }
}

# Index mapping for comments
COMMENT_INDEX_MAPPING = {
    'settings': {
        'analysis': {
            'analyzer': {
                'russian_analyzer': {
                    'type': 'russian'
                }
            }
        }
    },
    'mappings': {
        'properties': {
            'comment_id': {'type': 'keyword'},
            'comment': {
                'type': 'text',
                'analyzer': 'russian_analyzer'
            }
        }
    }
} 