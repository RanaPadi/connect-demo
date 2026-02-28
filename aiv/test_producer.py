import json
import time
import secrets
import os
import ctypes
from confluent_kafka import Producer
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.backends import default_backend
requestID_map = []
path = os.getcwd()
aiv = ctypes.cdll.LoadLibrary(path + "/libaiv_crypto.so")

def save_private_key_to_pem(private_key, filename, password):
    pem = private_key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.TraditionalOpenSSL,
        encryption_algorithm=serialization.BestAvailableEncryption(password)
    )
    with open(filename, 'wb') as f:
        f.write(pem)

def save_public_key_to_pem(public_key, filename):
    pem = public_key.public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo
    )
    with open(filename, 'wb') as f:
        f.write(pem)

def create_aiv_key(policy):
    private_key = ec.generate_private_key(ec.SECP256R1(), default_backend())
    public_key = private_key.public_key()
    save_private_key_to_pem(private_key, "aiv_private_key.pem", policy)
    save_public_key_to_pem(public_key, "aiv_public_key.pem")

def load_private_key_from_pem(filename, password):
    with open(filename, 'rb') as f:
        pem_data = f.read()
    private_key = serialization.load_pem_private_key(
        pem_data,
        password=password if password else None,
        backend=default_backend()
    )
    return private_key

def load_public_key_from_pem(filename):
    with open(filename, 'rb') as f:
        pem_data = f.read()
    public_key = serialization.load_pem_public_key(
        pem_data,
        backend=default_backend()
    )
    return public_key

def sign_message(message, nonce):
    policy = bytes([0x05] * 32)
    private_key = load_private_key_from_pem("aiv_private_key.pem", policy)
    nonceBytestream = bytes.fromhex(nonce)
    messageBytestream = message.encode("utf-8")

    to_be_signed = nonceBytestream + messageBytestream

    signature = private_key.sign(
        to_be_signed,
        ec.ECDSA(hashes.SHA256())
    )
    return signature.hex()

def acked(err, msg):
    if err is not None:
        print(f"Failed to deliver message: {err}")
    else:
        print(f"Message produced: {msg.key()}")

def produce_message(producer, topic, message):
    producer.produce(topic, key=str(time.time()), value=json.dumps(message), callback=acked)
    producer.flush()
    print("Message produced successfully")

def create_generic_subscription_request(message_type, message_content):
    subID = secrets.token_hex(16)
    if message_type == "AIV_UNSUBSCRIBE_REQUEST":
        generic_subscription_request = {
            "sender": "sender123",
            "serviceType": "serviceType123",
            "messageType": message_type,
            "responseTopic": "responseTopic123",
            "requestId": requestID_map[-1],
            "message": message_content
        }
    else:
        generic_subscription_request = {
            "sender": "sender123",
            "serviceType": "serviceType123",
            "messageType": message_type,
            "responseTopic": "responseTopic123",
            "requestId": subID,
            "message": message_content
        }
        requestID_map.append(subID)

    print(requestID_map)
    return generic_subscription_request

def create_aiv_request():
    query = [
        {
            "TrusteeID": "VC1",
            "requestedClaims": ["claim1", "claim2"]
        }
    ]
    evidence = {
        "timestamp": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "nonce": secrets.token_hex(32),
        "signatureAlgorithmType": "ECDSA-SHA256",
        "signature": secrets.token_hex(64),
        "keyRef": "aiv_public_key"
    }
    attestationCertificate = "base64_encoded_certificate"

    message = {
        "query": query,
        "evidence": evidence,
        "attestationCertificate": attestationCertificate
    }

    return create_generic_subscription_request("AIV_REQUEST", message)

def create_aiv_subscribe_request():
    subscribe = [
        {
            "TrusteeID": "VC1",
            "requestedClaims": ["claim1", "claim2"]
        },
        {
            "TrusteeID": "VC2",
            "requestedClaims": ["claim1", "claim2"]
        }
    ]
    checkInterval = 2000  # 6 seconds
    evidence = {
        "timestamp": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "nonce": secrets.token_hex(32),
        "signatureAlgorithmType": "ECDSA-SHA256",
        "signature": secrets.token_hex(64),
        "keyRef": "aiv_public_key"
    }
    attestationCertificate = "base64_encoded_certificate"

    message = {
        "subscribe": subscribe,
        "checkInterval": checkInterval,
        "evidence": evidence,
        "attestationCertificate": attestationCertificate
    }

    return create_generic_subscription_request("AIV_SUBSCRIBE_REQUEST", message)

def create_aiv_unsubscribe_request():
    evidence = {
        "timestamp": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "nonce": secrets.token_hex(32),
        "signatureAlgorithmType": "ECDSA-SHA256",
        "signature": secrets.token_hex(64),
        "keyRef": "aiv_public_key"
    }
    attestationCertificate = "base64_encoded_certificate"
    print(requestID_map)
    message = {
        "subscriptionId": requestID_map[-1],
        "attestationCertificate": attestationCertificate
    }

    return create_generic_subscription_request("AIV_UNSUBSCRIBE_REQUEST", message)

def main():
    broker = "127.0.0.1:9092"
    topic = "aiv"

    conf = {
        'bootstrap.servers': broker,
    }

    producer = Producer(**conf)
    
    policy = bytes([0x05] * 32)
    create_aiv_key(policy)

    # Create and send AIV_REQUEST message
    aiv_request_message = create_aiv_request()
    produce_message(producer, topic, aiv_request_message)
    print(aiv_request_message)

    # Create and send AIV_SUBSCRIBE_REQUEST message
    aiv_subscribe_request_message = create_aiv_subscribe_request()
    produce_message(producer, topic, aiv_subscribe_request_message)

    import time
    time.sleep(50)
    # Create and send AIV_UNSUBSCRIBE_REQUEST message
    aiv_unsubscribe_request_message = create_aiv_unsubscribe_request()
    time.sleep(10)
    produce_message(producer, topic, aiv_unsubscribe_request_message)
    print(aiv_unsubscribe_request_message)

if __name__ == "__main__":
    main()
