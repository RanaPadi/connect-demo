# CONNECT Demo Repository Overview

## What It Does

CONNECT-DEMO is a comprehensive simulation and trust-based application migration platform for cooperative autonomous vehicles. It demonstrates a V2X (Vehicle-to-Everything) enabled Connected and Automated Mobility (CCAM) system that safely migrates running applications (like cooperative adaptive cruise control) between vehicle computers based on real-time trust assessments.

The system consists of:

- **SUMO Simulation**: Traffic and vehicle movement simulation
- **CACC (Cooperative Adaptive Cruise Control)**: A cooperative driving control application that runs on vehicles
- **Trust Assessment Framework (TAF)**: A Go-based component that continuously evaluates the trustworthiness of execution environments
- **AIV (Automotive Integrity Verification)**: C-based component that gathers integrity evidence from vehicles
- **Migration App**: Orchestrates the safe migration of CACC between vehicle computers when trust drops below acceptable thresholds

The demo runs on a network with one main machine and two vehicle computers (VCs/ECUs) communicating via Kafka, simulating a realistic autonomous vehicle fleet scenario.

## Technologies Used

### Languages & Runtime

- **Python 3**: Core application logic (CACC, migration app, SUMO integration)
- **C++**: High-performance simulation backend via Libtraci bindings
- **Go 1.22.1**: Trust Assessment Framework (TAF) - the core trust evaluation engine
- **C**: Cryptographic utilities (AIV - integrity verification)
- **Bash**: System setup and orchestration scripts

### Messaging & Communication

- **Apache Kafka**: Event-driven message broker for inter-component communication (with SSL/TLS support)
- **Confluent Kafka**: Python client for Kafka integration

### Simulation

- **SUMO (Simulation of Urban Mobility)**: Accelerated simulation using **libtraci** (C++ implementation) for high-frequency V2X control loops and GUI-enabled visualization

### Data & Configuration

- **Pydantic**: Type-safe Python data modeling for message schemas
- **YAML/JSON**: Configuration files for system setup and component settings
- **JSON Schema**: Message validation and structure definition

### Cryptography & Security

- **OpenSSL**: Cryptographic operations
- **RSA/ASNI1**: For message signing and verification
- **Secure boot** and certificate-based authentication
- **TTL (Time-to-Live)** based V2X node tracking

### Additional Tools

- **Gin Framework**: Web interface for TAF debugging
- **MkDocs**: Documentation generation
- **PyYAML/pydantic**: Configuration management with environment variable support

## Methodologies Used

### Architecture & Design Patterns

- **Microservices Architecture**: Loosely coupled components (TAF, AIV, CACC, Migration App) communicating via Kafka
- **Event-Driven Architecture**: Asynchronous message-based communication with topic-based pub/sub
- **Handler-Based Message Processing**: Multiple specialized handlers for different message types (TASInit, TASNotify, TASSubscribe, TASTear, TASUnsubscribe)

### Trust & Security

- **Evidence-Based Trust Assessment**: TAF collects integrity evidence from AIV and other sources to compute trust values
- **Subjective Logic**: Trust computation using probabilistic belief models (via go-subjectivelogic library)
- **Threshold-Based Migration Triggers**: Application migration triggered when trust metrics fall below RTL/ATL thresholds
- **Cryptographic Verification**: Message signing and verification for secure inter-component communication

### System Design

- **Multi-Machine Deployment**: Hybrid deployment using high-performance C++ shared libraries for local simulation control and Kafka for cross-machine trust orchestration
- **Stateful Message Flows**: Subscription-based evidence collection with session and request management
- **TTL-Based Node Management**: V2X dynamics with configurable node lifetime and periodic health checks
- **Configuration-Driven Behavior**: Externalized configuration for easy tuning and adaptation

### Operational Patterns

- **Automated Setup**: Bash scripts for environment configuration, Kafka setup, and Python virtual environments
- **Logging & Debugging**: Structured logging (JSON, PRETTY, PLAIN formats) and debug mode support
- **Schema-Driven Development**: Auto-generated Go structs from JSON schemas using quicktype

## Summary

This is a high-performance research platform for trust-aware automated vehicle systems, utilizing C++ accelerated simulation (libtraci) to achieve the low-latency control required for CACC and real-time migration.
