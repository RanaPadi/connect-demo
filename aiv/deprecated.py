import time
import json
from confluent_kafka import Consumer, Producer
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.backends import default_backend
from cryptography.exceptions import InvalidSignature
from cryptography.hazmat.primitives.serialization import Encoding, PrivateFormat, NoEncryption
import requests
import secrets
import os
import ctypes
import sys

path = os.getcwd()
aiv = ctypes.cdll.LoadLibrary(path + "/libaiv_crypto.so")

class AIVResponse:
    def __init__(self, sender, serviceType, messageType, responseId, trusteeReport, aivEvidence):
        self.sender = sender
        self.serviceType = serviceType
        self.messageType = messageType
        self.responseId = responseId
        self.trusteeReport = trusteeReport
        self.aivEvidence = aivEvidence

    def to_dict(self):
        return {
            "sender": self.sender,
            "serviceType": self.serviceType,
            "messageType": self.messageType,
            "responseId": self.responseId,
            "message": {
                "trusteeReports":[
                    {
                        "attestationReport": [
                            {
                                "appraisal": self.trusteeReport["appraisal"],
                                "claim": self.trusteeReport["claim"],
                                "timestamp": self.trusteeReport["timestamp"]
                            }
                        ],
                        "trusteeID": self.trusteeReport["trusteeID"]
                    }
                ],
                "aivEvidence": {
                    "timestamp": self.aivEvidence["timestamp"],
                    "nonce": self.aivEvidence["nonce"],
                    "signatureAlgorithmType": self.aivEvidence["signatureAlgorithmType"],
                    "signature": self.aivEvidence["signature"],
                    "keyRef": self.aivEvidence["keyRef"]
                }
            }
        }

class TrusteeReport:
    def __init__(self, trusteeID, claim, timestamp, appraisal):
        self.trusteeID = trusteeID
        self.claim = claim
        self.timestamp = timestamp
        self.appraisal = appraisal

    def to_dict(self):
        return {
            "attestationReport": [
                {
                    "appraisal": self.appraisal,
                    "claim": self.claim,
                    "timestamp": self.timestamp
                }
            ],
            "trusteeID": self.trusteeID
        }

class AIVEvidence:
    def __init__(self, timestamp, nonce, signature, signatureAlgorithmType, keyRef):
        self.timestamp = timestamp
        self.nonce = nonce
        self.signature = signature
        self.signatureAlgorithmType = signatureAlgorithmType
        self.keyRef = keyRef

def save_private_key_to_pem(private_key, filename, password):
    pem = private_key.private_bytes(
        encoding=Encoding.PEM,
        format=PrivateFormat.TraditionalOpenSSL,
        encryption_algorithm=serialization.BestAvailableEncryption(password)
    )
    with open(filename, 'wb') as f:
        f.write(pem)

def save_public_key_to_pem(public_key, filename):
    pem = public_key.public_bytes(
        encoding=Encoding.PEM,
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

    to_be_signed =  nonceBytestream + messageBytestream

    signature = private_key.sign(
        to_be_signed,
        ec.ECDSA(hashes.SHA256())
    )
    return signature.hex()

def consume_and_process(consumer, producer, topic_p):
    while True:
        msg = consumer.poll(1.0)
        if msg is None:
            continue
        if msg.error():
            print(f"Consumer error: {msg.error()}")
            continue

        message = json.loads(msg.value().decode())["message"]
        # Verify TAF signature
        pubKey = load_public_key_from_pem(message["evidence"]["keyRef"] + ".pem")

        nonceBytestream = bytes.fromhex(message["evidence"]["nonce"])

        signature = bytes.fromhex(message["evidence"]["signature"])

        try:
            pubKey.verify(
                signature,
                nonceBytestream,
                ec.ECDSA(hashes.SHA256())
            )
            print("AIV_REQUEST is sent by TAF and it is authentic.")
        except InvalidSignature:
            print("Signature is invalid.")
            sys.exit("Exiting: Signature is invalid.")

        queries = message.get('query', [])

        uniqueDevices = list(set(query['TrusteeID'] for query in queries))

        trustee_reports = []

        for device in uniqueDevices:
            appraisal = request_for_evidence(device)

            timestamp = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())

            trustee_report = TrusteeReport(
                trusteeID=device,
                claim="runtime-integrity",
                timestamp=timestamp,
                appraisal=appraisal
            )

            trustee_reports.append(trustee_report.to_dict())

        aiv_signature = sign_message(json.dumps(trustee_reports, separators=(',', ':')), message["evidence"]["nonce"])

        aiv_evidence = AIVEvidence(
            timestamp=timestamp,
            nonce=message["evidence"]["nonce"],
            signature=aiv_signature,
            signatureAlgorithmType="ECDSA-SHA256",
            keyRef="aiv_public_key"
        )

        response = AIVResponse(
            sender="a77b29bac8f1-aiv",
            serviceType="ECI",
            messageType="AIV_RESPONSE",
            responseId="4c54a50f8e43",
            trusteeReport=trustee_report.__dict__,
            aivEvidence=aiv_evidence.__dict__
        )

        json_response = json.dumps(response.to_dict())
        producer.produce(topic_p, json_response)
        producer.flush()
        print("Message produced successfully")

def request_for_evidence(deviceName):
    print("Attesting device: " + deviceName)
    nonce = secrets.token_hex(32)
    nonce_hex = nonce
    payload = {
        "Nonce": nonce_hex
    }
    json_string = json.dumps(payload)
    try:
        url = 'http://localhost:8000/api/post/RequestForEvidence'
        headers = {'Content-Type': 'application/json'}
        response = requests.post(url, headers=headers, data=json_string)
        response.raise_for_status()
        response_data = response.json()
        sig_r_bb = bytes.fromhex(response_data.get('SigR'))
        sig_s_bb = bytes.fromhex(response_data.get('SigS'))
        public_bb = bytes.fromhex(response_data.get('Public'))
        attestation = aiv.VerifySignature(sig_r_bb, sig_s_bb, bytes.fromhex(nonce), 32,public_bb)
        return attestation
    except requests.exceptions.RequestException as e:
        print(f"Error sending request: {e}")
        return -1

def main():
    broker = "127.0.0.1:9092"
    topic_rfe = "aiv"
    topic_p = "taf"
    consumer = Consumer({
        'bootstrap.servers': broker,
        'group.id': 'AIV',
        'auto.offset.reset': 'earliest'
    })
    consumer.subscribe([topic_rfe])
    producer = Producer({
        'bootstrap.servers': broker
    })
    policy = bytes([0x05] * 32)
    create_aiv_key(policy)
    try:
        consume_and_process(consumer, producer, topic_p)
    except KeyboardInterrupt:
        pass
    finally:
        consumer.close()

if __name__ == "__main__":
    main()