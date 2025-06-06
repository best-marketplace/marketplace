import logging
import sys
import time
import signal
import json
from datetime import datetime
from services.kafka_service import KafkaService
from services.elasticsearch_service import ElasticsearchService

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    stream=sys.stdout
)
logger = logging.getLogger(__name__)

class CommentETLService:
    def __init__(self):
        self.kafka_service = None
        self.es_service = None
        self.running = False

    def initialize(self):
        try:
            logger.info("Starting Comment ETL service initialization...")
            
            logger.info("Initializing Elasticsearch service...")
            self.es_service = ElasticsearchService()
            
            logger.info("Initializing Kafka service...")
            self.kafka_service = KafkaService(topic='comment-events')
            
            logger.info("Comment ETL service initialized successfully")
            return True
        except Exception as e:
            logger.error(f"Failed to initialize Comment ETL service: {str(e)}")
            return False

    def transform_comment_data(self, data):
        try:
            logger.info(f"Transforming comment data: {json.dumps(data, ensure_ascii=False)}")
            
            transformed_data = {
                'comment_id': data.get('comment_id'),
                'comment': data.get('comment', '')
            }
            
            logger.info(f"Transformed comment data: {json.dumps(transformed_data, ensure_ascii=False)}")
            return transformed_data
        except Exception as e:
            logger.error(f"Error transforming comment data: {str(e)}")
            return None

    def process_comment(self, data):
        try:
            logger.info(f"RAW incoming data: {json.dumps(data, ensure_ascii=False)}")
            logger.info(f"Starting to process comment data: {json.dumps(data, ensure_ascii=False)}")
            
            transformed_data = self.transform_comment_data(data)
            if not transformed_data:
                logger.error("Failed to transform comment data")
                return
            
            success = self.es_service.index_comment(transformed_data)
            if success:
                logger.info(f"Successfully processed comment {transformed_data['comment_id']}")
            else:
                logger.error(f"Failed to process comment {transformed_data['comment_id']}")
        except Exception as e:
            logger.error(f"Error processing comment: {str(e)}")

    def start(self):
        if not self.initialize():
            logger.error("Failed to initialize Comment ETL service. Exiting...")
            sys.exit(1)

        self.running = True
        
        signal.signal(signal.SIGINT, self.handle_shutdown)
        signal.signal(signal.SIGTERM, self.handle_shutdown)

        try:
            logger.info("Starting Comment ETL service...")
            self.kafka_service.consume_messages(self.process_comment)
        except Exception as e:
            logger.error(f"Error in Comment ETL service: {str(e)}")
        finally:
            self.shutdown()

    def handle_shutdown(self, signum, frame):
        logger.info(f"Received signal {signum}. Initiating graceful shutdown...")
        self.running = False
        self.shutdown()

    def shutdown(self):
        logger.info("Shutting down Comment ETL service...")
        
        try:
            if self.kafka_service:
                self.kafka_service.close()
            
            if self.es_service:
                self.es_service.close()
                
            logger.info("Comment ETL service shut down successfully")
        except Exception as e:
            logger.error(f"Error during shutdown: {str(e)}")
        finally:
            sys.exit(0)

def main():
    etl_service = CommentETLService()
    etl_service.start()

if __name__ == "__main__":
    main() 