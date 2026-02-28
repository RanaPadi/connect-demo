import time
import json
from kafka import KafkaConsumer


def main():
    IP = "IP:PORT"
    TOPIC = "YOUR TOPIC HERE"

# Consumer configuration with UTF-8 decoding
    consumer = KafkaConsumer(
        TOPIC,  # Topic name
        bootstrap_servers=IP,
        value_deserializer=lambda m: json.loads(m.decode('utf-8')),  # Deserialize data from bytes
        auto_offset_reset='latest'
    )

    while True:

        data = consumer.poll(timeout_ms=1)
        for tp, messages in data.items():
            for message in messages:
                print(f"Received message: {message.value}")

        time.sleep(10)



if __name__ == "__main__":
    main()