import json
import time
from kafka import KafkaProducer, KafkaConsumer

def load_json_file(file_path):
    """Load and return data from a JSON file."""
    with open(file_path, 'r') as file:
        return json.load(file)

def main():
    # Configuration for the Kafka broker
    IP = "IP:PORT"
    TOPIC = "YOUR TOPIC HERE"

    # Producer configuration with UTF-8 encoding
    producer = KafkaProducer(
        bootstrap_servers=IP,
        value_serializer=lambda v: json.dumps(v).encode('utf-8')  # Serialize data to bytes
    )

    # Consumer configuration with UTF-8 decoding
    consumer = KafkaConsumer(
        TOPIC,  # Topic name
        bootstrap_servers=IP,
        value_deserializer=lambda m: json.loads(m.decode('utf-8')),  # Deserialize data from bytes
        auto_offset_reset='latest'
    )

    # Load JSON data from a file
    MESSAGE = {
        "message": "Hello World"
    }


    try:

        while True:
        # Poll for new messages
            kafka_messages = consumer.poll(timeout_ms=100)
            for tp, messages in kafka_messages.items():
                for message in messages:
                    print(f"Received message: {message.value}")  # Should match the produced message


            # Produce a message
            future = producer.send(TOPIC, MESSAGE)
            # Ensure the message is sent and acknowledged
            try:
                future.get(timeout=10.0)
                print(f"Produced message: {MESSAGE}")
            except Exception as e:
                print(f"Failed to deliver message: {e}")

            # Flush the producer to ensure all messages are sent
            producer.flush()

            time.sleep(1)


    except KeyboardInterrupt:
        print("Consumer interrupted")

    finally:
        # Clean up and close the producer and consumer
        producer.close()
        consumer.close()

if __name__ == "__main__":
    main()