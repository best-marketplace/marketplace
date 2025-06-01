from kafka import KafkaConsumer
import json
import logging
import time
from config import KAFKA_CONFIG

logger = logging.getLogger(__name__)

class KafkaService:
    def __init__(self, topic=None):
        self.config = KAFKA_CONFIG
        self.topic = topic or self.config['topic']
        self.consumer = None
        self._initialize_consumer()

    def _initialize_consumer(self):
        """Initialize Kafka consumer with retry logic"""
        max_retries = 5
        retry_delay = 5  # seconds

        for attempt in range(max_retries):
            try:
                logger.info(f"Attempting to initialize Kafka consumer (attempt {attempt + 1}/{max_retries})")
                self.consumer = KafkaConsumer(
                    self.topic,
                    bootstrap_servers=self.config['bootstrap_servers'],
                    group_id=self.config['group_id'],
                    auto_offset_reset='earliest',
                    value_deserializer=lambda x: json.loads(x.decode('utf-8-sig')),
                    enable_auto_commit=True,
                    auto_commit_interval_ms=5000,
                    session_timeout_ms=30000,
                    heartbeat_interval_ms=10000
                )
                logger.info(f"Successfully initialized Kafka consumer for topic {self.topic}")
                return
            except Exception as e:
                if attempt < max_retries - 1:
                    logger.warning(f"Failed to initialize Kafka consumer (attempt {attempt + 1}/{max_retries}): {str(e)}")
                    time.sleep(retry_delay)
                else:
                    logger.error(f"Failed to initialize Kafka consumer after {max_retries} attempts: {str(e)}")
                    raise

    def consume_messages(self, callback):
        """Read messages from Kafka and pass them to callback function"""
        while True:
            try:
                if self.consumer is None:
                    logger.info("Reinitializing Kafka consumer...")
                    self._initialize_consumer()

                logger.info("Starting to consume messages...")
                for message in self.consumer:
                    try:
                        logger.info(f"Received message from topic {message.topic}, partition {message.partition}, offset {message.offset}")
                        
                        # Log raw message for debugging
                        logger.debug(f"Raw message value: {message.value}")
                        
                        if not message.value:
                            logger.warning("Received empty message, skipping...")
                            continue

                        data = message.value
                        logger.info(f"Message content: {json.dumps(data, ensure_ascii=False)}")

                        if not isinstance(data, dict):
                            logger.warning(f"Received invalid message format: {data}")
                            continue

                        # Check if this is a product creation event or comment creation event
                        if data.get('action') in ['product_created', 'comment_created']:
                            logger.info(f"Processing {data.get('action')} event: {json.dumps(data, ensure_ascii=False)}")
                            callback(data)
                        else:
                            logger.debug(f"Ignoring message with action: {data.get('action')}")

                    except json.JSONDecodeError as e:
                        logger.error(f"Failed to decode message: {str(e)}")
                        logger.error(f"Raw message value: {message.value}")
                    except Exception as e:
                        logger.error(f"Error processing message: {str(e)}")
                        logger.error(f"Message details - Topic: {message.topic}, Partition: {message.partition}, Offset: {message.offset}")

            except Exception as e:
                logger.error(f"Error in Kafka consumer: {str(e)}")
                time.sleep(5)  # Wait before reconnecting
                self.consumer = None  # Force reinitialization

    def close(self):
        """Safely close the Kafka consumer"""
        if self.consumer:
            try:
                self.consumer.close()
                logger.info("Kafka consumer closed successfully")
            except Exception as e:
                logger.error(f"Error closing Kafka consumer: {str(e)}") 