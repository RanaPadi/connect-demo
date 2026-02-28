import time
import json
import threading
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
import argparse
import logging
import subprocess
import sys
import paramiko
logging.basicConfig(filename='aiv_eval.log',
                            filemode='a',
                            format='%(asctime)s,%(msecs)d %(name)s %(levelname)s %(message)s',
                            datefmt='%H:%M:%S',
                            level=logging.DEBUG)

path = os.getcwd()
aiv = ctypes.cdll.LoadLibrary(path + "/libaiv_crypto.so")
thread_map = {}
AttestattionMode = True
TrustModeCheck = False

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
                "trusteeReports": self.trusteeReport,
                "aivEvidence": self.aivEvidence
            }
        }

class AIVNotify:
    def __init__(self, sender, serviceType, messageType, responseId, trusteeReports, aivEvidence, subscriptionId):
        self.sender = sender
        self.serviceType = serviceType
        self.messageType = messageType
        self.responseId = responseId
        self.trusteeReports = trusteeReports
        self.aivEvidence = aivEvidence
        self.subscriptionId = subscriptionId

    def to_dict(self):
        return {
            "sender": self.sender,
            "serviceType": self.serviceType,
            "messageType": self.messageType,
            "responseId": self.responseId,
            "message": {
                "subscriptionId": self.subscriptionId,
                "trusteeReports": self.trusteeReports,
                "aivEvidence": self.aivEvidence
            }
        }

class AIVSubscribeRequest:
    def __init__(self, subscribe, checkInterval, evidence, attestationCertificate):
        self.subscribe = subscribe
        self.checkInterval = checkInterval
        self.evidence = evidence
        self.attestationCertificate = attestationCertificate

    def to_dict(self):
        return {
            "subscribe": self.subscribe,
            "checkInterval": self.checkInterval,
            "evidence": self.evidence,
            "attestationCertificate": self.attestationCertificate
        }

class AIVSubscribeResponse:
    def __init__(self, sender, serviceType, messageType, requestId, success=None, error=None, subscriptionId=None):
        self.sender = sender
        self.serviceType = serviceType
        self.messageType = messageType
        self.requestId = requestId
        self.success = success
        self.error = error
        self.subscriptionId = subscriptionId

    def to_dict(self):
        return {
            "sender": self.sender,
            "serviceType": self.serviceType,
            "messageType": self.messageType,
            "responseId": self.requestId,
            "message": {
                "success": self.success,
                "subscriptionId": self.subscriptionId
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

    def to_dict(self):
        return {
            "timestamp": self.timestamp,
            "nonce": self.nonce,
            "signatureAlgorithmType": self.signatureAlgorithmType,
            "signature": self.signature,
            "keyRef": self.keyRef
        }

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

    to_be_signed = nonceBytestream + messageBytestream

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

        message = json.loads(msg.value().decode())


        # Determine message type and handle accordingly
        message_type = message.get("messageType")
        print(f"Message Type Received from TAF is: {message_type}")
        if message_type == "AIV_REQUEST":
            handle_aiv_request(message["message"], message["requestId"], producer, topic_p)
        elif message_type == "AIV_SUBSCRIBE_REQUEST":
            handle_aiv_subscribe_request(message["message"], message["requestId"], producer, topic_p)
        elif message_type == "AIV_UNSUBSCRIBE_REQUEST":
            handle_aiv_unsubscribe_request(message["message"], message["requestId"], producer, topic_p)
        else:
            print(f"Unknown message type: {message_type}")

def handle_aiv_request(message, requestID, producer, topic_p):
    # Verify TAF signature
    start = time.time()
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

    queries = message.get('query', [])
    uniqueDevices = list(set(query['TrusteeID'] for query in queries))

    trustee_reports = []
    if AttestattionMode == True:
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
            responseId=requestID,
            trusteeReport=trustee_reports,  # assuming we just take the first for simplicity
            aivEvidence=aiv_evidence.to_dict()
        )

        json_response = json.dumps(response.to_dict())
        producer.produce(topic_p, json_response)
        producer.flush()
        print("AIV_RESPONSE message produced successfully")
        end = time.time()
        logging.info("[*]TIMING ATTESTATION REPORT CONSTRUCTION (SYCHRONOUS)\t {0}".format(end-start))
    else:
        with open("aiv_taf.json", "r") as f:
                    claims = json.load(f)
                    f.close()
                
        published_report = ""
        counter = 0 
        for subscription in message["query"]:
            timestamp = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
            claims_VC = claims[subscription["TrusteeID"]]["requestedClaims"]
            trust_properties = claims_VC.keys()
            for i, trust_property in enumerate(trust_properties):
                trustee_report = TrusteeReport(
                    trusteeID=subscription["TrusteeID"],
                    claim=trust_property,
                    timestamp=timestamp,
                    appraisal=claims_VC[trust_property][counter%len(claims_VC[trust_property])]
                )
                if( i == 0):
                    
                    published_report = trustee_report.to_dict()
                else:
                    published_report["attestationReport"].append(trustee_report.to_dict()["attestationReport"][0])
            trustee_reports.append(published_report)

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
                responseId=requestID,
                trusteeReport=trustee_reports,  # assuming we just take the first for simplicity
                aivEvidence=aiv_evidence.to_dict()
            )
            json_response = json.dumps(response.to_dict())
            producer.produce(topic_p, json_response)
            producer.flush()
            print("AIV_RESPONSE message produced successfully")
            counter+=1
            end = time.time()
            logging.info("[*]TIMING ATTESTATION REPORT CONSTRUCTION (SYCHRONOUS)\t {0}".format(end-start))
        


def handle_aiv_unsubscribe_request(message, requestID, producer, topic_p):
    
    

    print("AIV_UNSUBSCRIBE_REQUEST is sent by TAF and it is authentic.")
    
    notify_thread, notify_stop_event, _ = thread_map[requestID]
    notify_stop_event.set()

    notify_thread.join()
    del thread_map[requestID]
    unsubscribe_success = True
        
    

    if unsubscribe_success:
        response_message = {
            "success": "Unsubscribed successfully"
        }
    else:
        response_message = {
            "error": "Failed to unsubscribe"
        }
    

    generic_response = {
        "sender": "a77b29bac8f1-aiv",
        "serviceType": "ECI",
        "messageType": "AIV_UNSUBSCRIBE_RESPONSE",
        "responseId": requestID,
        "message": response_message
    }

    json_response = json.dumps(generic_response)
    producer.produce(topic_p, json_response)
    producer.flush()
    print("AIV_UNSUBSCRIBE_RESPONSE message produced successfully")

def handle_aiv_subscribe_request(message, requestID, producer, topic_p):
    # Verify AIV_SUBSCRIBE_REQUEST signature
    pubKey = load_public_key_from_pem(message["evidence"]["keyRef"] + ".pem")

    nonceBytestream = bytes.fromhex(message["evidence"]["nonce"])
    signature = bytes.fromhex(message["evidence"]["signature"])

    try:
        pubKey.verify(
            signature,
            nonceBytestream,
            ec.ECDSA(hashes.SHA256())
        )
        print("AIV_SUBSCRIBE_REQUEST is sent by TAF and it is authentic.")

    except InvalidSignature:
        print("Signature is invalid.")
        # sys.exit("Exiting: Signature is invalid.")

    
    sender = "a77b29bac8f1-aiv"
    serviceType = "ECI"
    messageType = "AIV_SUBSCRIBE_RESPONSE"
    requestId = requestID  # Unique request ID that was used in the subscription request

    success_message = AIVSubscribeResponse(
        sender=sender,
        serviceType=serviceType,
        messageType=messageType,
        requestId=requestId,
        success="Subscription successful",
        subscriptionId=requestId
    )

    json_response = json.dumps(success_message.to_dict())
    producer.produce(topic_p, json_response)
    producer.flush()
    print("AIV_SUBSCRIBE_RESPONSE message produced successfully")

    # Start a separate thread for periodic AIV_NOTIFY messages
    notify_stop_event = threading.Event()
    thread_map_lock = threading.Lock()
    notify_thread = threading.Thread(
        target=periodic_notify,
        args=(message, requestID, producer, topic_p, notify_stop_event, thread_map_lock)
    )
    thread_map.update({requestID:[notify_thread, notify_stop_event, thread_map_lock]})
    notify_thread.start()

def periodic_notify(message, requestID, producer, topic_p, stop_event, thread_map_lock):
    counter_notify = 0
    while not stop_event.is_set():
        start = time.time()
        with thread_map_lock:
            _, check_event, _ = thread_map.get(requestID, (None, None, None))

            if not check_event:
                break
            
            if check_event.is_set():
                break

            trustee_reports = []
            if AttestattionMode == True:
                for subscription in message["subscribe"]:
                    appraisal = request_for_evidence(subscription["TrusteeID"])
                    # appraisal = 1
                    timestamp = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())

                    trustee_report = TrusteeReport(
                        trusteeID=subscription["TrusteeID"],
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

                notify = AIVNotify(
                    sender="a77b29bac8f1-aiv",
                    serviceType="ECI",
                    messageType="AIV_NOTIFY",
                    responseId=requestID,
                    subscriptionId=requestID,
                    trusteeReports=trustee_reports,
                    aivEvidence=aiv_evidence.to_dict()
                )
            elif AttestattionMode == False and TrustModeCheck == True:
                username = "user"
                password = "Password123"
                bash_script_path = "../connect-demo/bash_scripts/check_system_trustworthy.sh"
                published_report = ""

                print("Trust Mode Check is activated")

                for subscription in message["subscribe"]:
                    timestamp = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
                    requested_vc_id = subscription["TrusteeID"]
                    print("vc_id", requested_vc_id)
                    if requested_vc_id == "VC1" :
                        hostname = "192.168.10.102"
                    elif requested_vc_id == "VC2" :
                        hostname = "192.168.10.103"
                    else:
                        print("Invalid VC ID")
                        continue
                    print("hostname:", hostname)

                    requested_claims = subscription["requestedClaims"]
                    print("requested claims", requested_claims)

                    # Create an SSH client and connect to the host
                    print(f"connecting to host: {hostname} via ssh")
                    ssh = paramiko.SSHClient()
                    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
                    try:
                        ssh.connect(hostname, username=username, password=password)
                        print(f"SSH connection to {hostname} established")
                    except paramiko.SSHException as e:
                        print(f"SSH connection to {hostname} failed: {e}")
                        continue
                    i=0
                    for trust_property in requested_claims:
                    #for i, trust_property in enumerate(trust_properties):
                        print('trust_property/turst_source: ', trust_property)

                        print(f"Executing bash script on {hostname} with flag: {trust_property}")

                        # Execute the bash script with the 'trust_property' flag
                        stdin, stdout, stderr = ssh.exec_command(f"bash {bash_script_path} {trust_property}")

                        # Wait for command completion and get the exit status
                        exit_status = stdout.channel.recv_exit_status()

                        # Read and decode the output from the bash script
                        #output = stdout.read().decode().strip()
                        print(f"Bash script on {hostname} with flag: {trust_property} exist status: {exit_status}")

                        trustee_report = TrusteeReport(
                             trusteeID=subscription["TrusteeID"],
                             claim=trust_property,
                             timestamp=timestamp,
                             appraisal=exit_status
                         )
                        # if len(requested_claims) == 0:
                        if( i == 0):
                             published_report = trustee_report.to_dict()
                        else:
                             published_report["attestationReport"].append(
                                 trustee_report.to_dict()["attestationReport"][0]
                             )
                        i = i + 1

                    # Close the SSH connection
                    print(f"closing ssh connection to the host: {hostname}")
                    ssh.close()

                    print("published report:", published_report)
                    trustee_reports.append(published_report)

                    aiv_signature = sign_message(json.dumps(trustee_reports, separators=(',', ':')), message["evidence"]["nonce"])

                    aiv_evidence = AIVEvidence(
                        timestamp=timestamp,
                        nonce=message["evidence"]["nonce"],
                        signature=aiv_signature,
                        signatureAlgorithmType="ECDSA-SHA256",
                        keyRef="aiv_public_key"
                    )

                    notify = AIVNotify(
                        sender="a77b29bac8f1-aiv",
                        serviceType="ECI",
                        messageType="AIV_NOTIFY",
                        responseId=requestID,
                        subscriptionId=requestID,
                        trusteeReports=trustee_reports,
                        aivEvidence=aiv_evidence.to_dict()
                    )
                    print("AIV_NOTIFY <--mode=trust_check--> message produced successfully")
                    
            else:
                with open("aiv_taf.json", "r") as f:
                    claims = json.load(f)
                    f.close()
                
                published_report = ""
                for subscription in message["subscribe"]:
                    timestamp = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
                    claims_VC = claims[subscription["TrusteeID"]]["requestedClaims"]
                    trust_properties = claims_VC.keys()
                    for i, trust_property in enumerate(trust_properties):
                        trustee_report = TrusteeReport(
                            trusteeID=subscription["TrusteeID"],
                            claim=trust_property,
                            timestamp=timestamp,
                            appraisal=claims_VC[trust_property][counter_notify%len(claims_VC[trust_property])]
                        )
                        if( i == 0):
                            
                            published_report = trustee_report.to_dict()
                        else:
                            published_report["attestationReport"].append(trustee_report.to_dict()["attestationReport"][0])
                        trustee_reports.append(published_report)
                aiv_signature = sign_message(json.dumps(trustee_reports, separators=(',', ':')), message["evidence"]["nonce"])

                aiv_evidence = AIVEvidence(
                    timestamp=timestamp,
                    nonce=message["evidence"]["nonce"],
                    signature=aiv_signature,
                    signatureAlgorithmType="ECDSA-SHA256",
                    keyRef="aiv_public_key"
                )

                notify = AIVNotify(
                    sender="a77b29bac8f1-aiv",
                    serviceType="ECI",
                    messageType="AIV_NOTIFY",
                    responseId=requestID,
                    subscriptionId=requestID,
                    trusteeReports=trustee_reports,
                    aivEvidence=aiv_evidence.to_dict()
                )
                print("AIV_NOTIFY <--mode=mutable--> message produced successfully")

            json_notify = json.dumps(notify.to_dict())
            producer.produce(topic_p, json_notify)
            producer.flush()
            counter_notify+=1
            end = time.time()
            logging.info("[*]TIMING ATTESTATION REPORT CONSTRUCTION (ASYCHRONOUS)\t {0}".format(end-start))

            #print("AIV_NOTIFY message produced successfully")

            time.sleep(message["checkInterval"] / 1000)

def request_for_evidence(nonce):
    try:
        nonce = secrets.token_hex(32)
        nonce_hex = nonce
        payload = {
            "Nonce": nonce_hex
        }
        json_string = json.dumps(payload)
        url = 'http://localhost:8000/api/post/RequestForEvidence'
        headers = {'Content-Type': 'application/json'}
        response = requests.post(url, headers=headers, data=json_string)
        response.raise_for_status()
        response_data = response.json()
        sig_r_bb = bytes.fromhex(response_data.get('SigR'))
        sig_s_bb = bytes.fromhex(response_data.get('SigS'))
        public_bb = bytes.fromhex(response_data.get('Public'))
        attestation = aiv.VerifySignature(sig_r_bb, sig_s_bb, bytes.fromhex(nonce), 32, public_bb)
        return attestation
    except requests.exceptions.RequestException as e:
        print(f"Error sending request: {e}")
        return -1

def main(broker):
    # broker = "127.0.0.1:9092"
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
    parser = argparse.ArgumentParser(description='aiv script')
    parser.add_argument('--broker_ip', required=True, help='IP address of the broker')
    parser.add_argument('--mode', required=True, help='Mode of testing model')
    args = parser.parse_args()
    if(args.mode == "mutable"):
        AttestattionMode = False
    elif(args.mode == "trust_check"):
        TrustModeCheck = True
        AttestattionMode = False

    main(args.broker_ip)

