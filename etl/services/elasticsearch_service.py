from elasticsearch import Elasticsearch
from elasticsearch.exceptions import ConnectionError, RequestError
from datetime import datetime
import logging
import time
import json
from config import ELASTICSEARCH_CONFIG, INDEX_MAPPING, COMMENT_INDEX_MAPPING

logger = logging.getLogger(__name__)

class ElasticsearchService:
    def __init__(self):
        self.config = ELASTICSEARCH_CONFIG
        self.es = None
        self.index = self.config['index']
        self.comment_index = 'comments'
        self._initialize_elasticsearch()

    def _initialize_elasticsearch(self):
        max_retries = 5
        retry_delay = 5 

        for attempt in range(max_retries):
            try:
                logger.info(f"Attempting to initialize Elasticsearch connection (attempt {attempt + 1}/{max_retries})")
                self.es = Elasticsearch(
                    hosts=self.config['hosts'],
                    retry_on_timeout=True,
                    max_retries=3,
                    timeout=30
                )

                if not self.es.ping():
                    raise ConnectionError("Failed to connect to Elasticsearch")
                
                logger.info("Successfully connected to Elasticsearch")
                self._create_index_if_not_exists()
                self._create_comment_index_if_not_exists()
                return
            except Exception as e:
                if attempt < max_retries - 1:
                    logger.warning(f"Failed to initialize Elasticsearch (attempt {attempt + 1}/{max_retries}): {str(e)}")
                    time.sleep(retry_delay)
                else:
                    logger.error(f"Failed to initialize Elasticsearch after {max_retries} attempts: {str(e)}")
                    raise

    def _create_index_if_not_exists(self):
        try:
            if not self.es.indices.exists(index=self.index):
                logger.info(f"Creating index {self.index} with mapping")
                self.es.indices.create(index=self.index, body=INDEX_MAPPING)
                logger.info(f"Successfully created index {self.index}")
            else:
                logger.info(f"Index {self.index} already exists")
        except Exception as e:
            logger.error(f"Error creating index {self.index}: {str(e)}")
            raise

    def _create_comment_index_if_not_exists(self):
        try:
            if not self.es.indices.exists(index=self.comment_index):
                logger.info(f"Creating index {self.comment_index} with mapping")
                self.es.indices.create(index=self.comment_index, body=COMMENT_INDEX_MAPPING)
                logger.info(f"Successfully created index {self.comment_index}")
            else:
                logger.info(f"Index {self.comment_index} already exists")
        except Exception as e:
            logger.error(f"Error creating index {self.comment_index}: {str(e)}")
            raise

    def index_product(self, product_data):
        if not self.es:
            logger.info("Elasticsearch connection lost, reinitializing...")
            self._initialize_elasticsearch()

        max_retries = 3
        retry_delay = 2 

        for attempt in range(max_retries):
            try:
                logger.info(f"Processing product data for indexing: {json.dumps(product_data, ensure_ascii=False)}")
                
                product_id = product_data.get('product_id')
                if not product_id:
                    logger.error("Product ID is missing in the data")
                    return False

                document = {
                    'product_id': product_id,
                    'title': product_data.get('title', ''),
                    'seller_name': product_data.get('seller_name', '')
                }

                logger.info(f"Prepared document for indexing: {json.dumps(document, ensure_ascii=False)}")

                response = self.es.index(
                    index=self.index,
                    id=product_id,
                    body=document,
                    refresh=True  
                )

                if response['result'] in ['created', 'updated']:
                    logger.info(f"Successfully indexed product {product_id} with result: {response['result']}")
                    return True
                else:
                    logger.warning(f"Unexpected response when indexing product {product_id}: {response['result']}")
                    return False

            except ConnectionError as e:
                if attempt < max_retries - 1:
                    logger.warning(f"Connection error while indexing product (attempt {attempt + 1}/{max_retries}): {str(e)}")
                    time.sleep(retry_delay)
                    self._initialize_elasticsearch() 
                else:
                    logger.error(f"Failed to index product after {max_retries} attempts: {str(e)}")
                    return False
            except Exception as e:
                logger.error(f"Error indexing product {product_data.get('product_id')}: {str(e)}")
                return False

    def index_comment(self, comment_data):
        """Index comment in Elasticsearch with retry logic"""
        if not self.es:
            logger.info("Elasticsearch connection lost, reinitializing...")
            self._initialize_elasticsearch()

        max_retries = 3
        retry_delay = 2 

        for attempt in range(max_retries):
            try:
                logger.info(f"Processing comment data for indexing: {json.dumps(comment_data, ensure_ascii=False)}")

                comment_id = comment_data.get('comment_id')
                if not comment_id:
                    logger.error("Comment ID is missing in the data")
                    return False

                document = {
                    'comment_id': comment_id,
                    'comment': comment_data.get('comment', '')
                }

                logger.info(f"Prepared document for indexing: {json.dumps(document, ensure_ascii=False)}")

                response = self.es.index(
                    index=self.comment_index,
                    id=comment_id,
                    body=document,
                    refresh=True 
                )

                if response['result'] in ['created', 'updated']:
                    logger.info(f"Successfully indexed comment {comment_id} with result: {response['result']}")
                    return True
                else:
                    logger.warning(f"Unexpected response when indexing comment {comment_id}: {response['result']}")
                    return False

            except ConnectionError as e:
                if attempt < max_retries - 1:
                    logger.warning(f"Connection error while indexing comment (attempt {attempt + 1}/{max_retries}): {str(e)}")
                    time.sleep(retry_delay)
                    self._initialize_elasticsearch() 
                else:
                    logger.error(f"Failed to index comment after {max_retries} attempts: {str(e)}")
                    return False
            except Exception as e:
                logger.error(f"Error indexing comment {comment_data.get('comment_id')}: {str(e)}")
                return False

    def close(self):
        if self.es:
            try:
                self.es.close()
                logger.info("Elasticsearch connection closed successfully")
            except Exception as e:
                logger.error(f"Error closing Elasticsearch connection: {str(e)}") 