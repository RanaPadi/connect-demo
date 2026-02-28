import time
import json
from kafka import KafkaProducer


def main():

    IP = "IP:PORT"
    TOPIC = "YOUR TOPIC HERE"

    # Producer configuration with UTF-8 encoding
    producer = KafkaProducer(
        bootstrap_servers=IP,
        value_serializer=lambda v:json.dumps(v).encode('utf-8')  # Serialize data to bytes
    )

    MESSAGE = {
        "message": "Hello World"
    }

    producer.send(TOPIC, MESSAGE)
    print(MESSAGE)



if __name__ == "__main__":
    main()

