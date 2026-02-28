from confluent_kafka import Consumer
from datetime import datetime
import time
import json
import subprocess

def poll_fresh_kafka_message(topic, freshness_secs=10, timeout_secs=30, poll_interval=2):
    consumer = Consumer({
        'bootstrap.servers': 'localhost:9092',
        'group.id': 'aiv-consumer',
        'auto.offset.reset': 'latest',
        'enable.auto.commit': False
    })
    consumer.subscribe([topic])

    start_time = datetime.utcnow()
    print(f"Polling Kafka topic '{topic}' for a fresh message...")

    while (datetime.utcnow() - start_time).total_seconds() < timeout_secs:
        msg = consumer.poll(timeout=5.0)

        if msg is None or msg.error():
            print("No message or Kafka error, retrying...")
            time.sleep(poll_interval)
            continue

        try:
            payload = json.loads(msg.value().decode("utf-8"))
            msg_timestamp = payload.get("timestamp")

            if msg_timestamp and (time.time() - msg_timestamp < freshness_secs):
                consumer.commit()
                consumer.close()
                print("Fresh Kafka message received.")
                return payload
            else:
                print("Stale Kafka message. Waiting for a fresh one...")
                time.sleep(poll_interval)
        except Exception as e:
            print(f"Error decoding Kafka message: {e}")
            time.sleep(poll_interval)

    consumer.close()
    print("Failed to get a fresh Kafka message within timeout.")
    return None

def run_remote_bash_script_with_sshpass(host, user, password, remote_script_path):
    cmd = f"sshpass -p '{password}' ssh -o StrictHostKeyChecking=no {user}@{host} 'bash {remote_script_path}'"
    try:
        result = subprocess.check_output(cmd, shell=True, stderr=subprocess.STDOUT)
        return result.decode().strip()
    except subprocess.CalledProcessError as e:
        print("Error output:", e.output.decode())
        return None