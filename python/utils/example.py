import os
import sys

from python.utils.models import TASMessage
from python.utils.utils import load_config, setup_logger, TopicBoundConsumer, TopicBoundProducer, \
    poll_and_extract_kafka_message

# Adjust the system path to include the parent directories
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..')))

############################################

def main(producer: TopicBoundProducer, consumer:TopicBoundConsumer) -> None:

    # how to receive messages from Kafka
    while True:
        # Poll Kafka messages
        data = poll_and_extract_kafka_message(consumer=consumer)
        if data is None:
            continue
        print(data)

######################

    # how to create TAS Message
    message = TASMessage(
        messageType="EXAMPLE TYPE",
        message={"EXAMPLE KEY": "EXAMPLE VALUE"}
    )

    # how to send message
    producer.send_message(message)

############################################

if __name__ == "__main__":
    # Load configuration and set up logging
    config = load_config()
    logger = setup_logger("EXAMPLE TOPIC")

    # Create instances of Kafka consumer and producer
    consumer = TopicBoundConsumer(server="remote_server", topic="EXAMPLE TOPIC", logger=logger)
    producer = TopicBoundProducer(server="remote_server", topic="EXAMPLE TOPIC", logger=logger)