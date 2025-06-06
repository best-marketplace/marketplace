import logging
import sys
import time
import signal
import json
import uuid
from datetime import datetime
from services.kafka_service import KafkaService
from services.elasticsearch_service import ElasticsearchService


logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    stream=sys.stdout
)
logger = logging.getLogger(__name__)

class ETLService:
    def __init__(self):
        self.kafka_service = None
        self.es_service = None
        self.running = False

    def initialize(self):
        try:
            logger.info("Starting ETL service initialization...")
            
            logger.info("Initializing Elasticsearch service...")
            self.es_service = ElasticsearchService()
            
            logger.info("Initializing Kafka service...")
            self.kafka_service = KafkaService()
            
            logger.info("ETL service initialized successfully")
            return True
        except Exception as e:
            logger.error(f"Failed to initialize ETL service: {str(e)}")
            return False

    def transform_product_data(self, data):
        try:
            logger.info(f"Transforming product data: {json.dumps(data, ensure_ascii=False)}")
            
            transformed_data = {
                'product_id': data.get('product_id'),
                'title': data.get('title', data.get('name', '')),
                'seller_name': data.get('seller_name', '')
            }
            
            logger.info(f"Transformed product data: {json.dumps(transformed_data, ensure_ascii=False)}")
            return transformed_data
        except Exception as e:
            logger.error(f"Error transforming product data: {str(e)}")
            return None

    def process_product(self, data):
        try:
            logger.info(f"RAW incoming data: {json.dumps(data, ensure_ascii=False)}")
            logger.info(f"Starting to process product data: {json.dumps(data, ensure_ascii=False)}")
            
            transformed_data = self.transform_product_data(data)
            if not transformed_data:
                logger.error("Failed to transform product data")
                return
            
            success = self.es_service.index_product(transformed_data)
            if success:
                logger.info(f"Successfully processed product {transformed_data['product_id']}")
            else:
                logger.error(f"Failed to process product {transformed_data['product_id']}")
        except Exception as e:
            logger.error(f"Error processing product: {str(e)}")

    def start(self):
        if not self.initialize():
            logger.error("Failed to initialize ETL service. Exiting...")
            sys.exit(1)

        self.running = True
        
        signal.signal(signal.SIGINT, self.handle_shutdown)
        signal.signal(signal.SIGTERM, self.handle_shutdown)

        try:
            logger.info("Starting ETL service...")
            self.kafka_service.consume_messages(self.process_product)
        except Exception as e:
            logger.error(f"Error in ETL service: {str(e)}")
        finally:
            self.shutdown()

    def handle_shutdown(self, signum, frame):
        logger.info(f"Received signal {signum}. Initiating graceful shutdown...")
        self.running = False
        self.shutdown()

    def shutdown(self):
        logger.info("Shutting down ETL service...")
        
        try:
            if self.kafka_service:
                self.kafka_service.close()
            
            if self.es_service:
                self.es_service.close()
                
            logger.info("ETL service shut down successfully")
        except Exception as e:
            logger.error(f"Error during shutdown: {str(e)}")
        finally:
            sys.exit(0)

def main():
    etl_service = ETLService()
    etl_service.start()

if __name__ == "__main__":
    main() 