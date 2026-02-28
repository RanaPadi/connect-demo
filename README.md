# CONNECT Demo

This quick start guide outlines the necessary steps to run the demo on your device. Ensure you are operating on a Linux machine, such as Ubuntu 20.04 or newer.

## Important Notes

- Install this repo exactly at the same user folder path across all the machines with in same LAN network.
- Be careful of the placeholders or sample values contained within `<>` for any commands mentioned in this README.
- Ensure each python script is executed within the virtual environment with the necessary Python libraries installed.
- Before starting the demo, verify that all prerequisites are installed (see Setup instructions below).

## Table of Contents

- [CONNECT Demo](#connect-demo)
- [Important Notes](#important-notes)
- [1. Machine IP Addresses](#1-machine-ip-addresses)
- [2. Execution Locations](#2-execution-locations)
- [3. Setup: Configuration File](#3-setup-configuration-file)
- [4. Setup: AIV](#4-setup-aiv)
- [5. Setup: TAF](#5-setup-taf)
- [6. Setup: Kafka](#6-setup-kafka)
  - [6.1. On Main Machine](#61-on-main-machine)
  - [6.2. On VC1/VC2](#62-on-vc1vc2)
- [7. Setup: Python](#7-setup-python)
  - [7.1. Automated Method Using Bash Script](#71-automated-method-using-bash-script)
  - [7.2. Manual Setup](#72-manual-setup)
- [8. Run: SUMO](#8-run-sumo)
- [9. Run: CACC](#9-run-cacc)
- [10. Run: TAF, AIV, and Migration App](#10-run-taf-aiv-and-migration-app)
  - [10.1. Automated Method Using Bash Script](#101-automated-method-using-bash-script)
  - [10.2. Manual Setup](#102-manual-setup)
    - [10.2.1. Start TAF](#1021-start-taf)
    - [10.2.2. Start AIV](#1022-start-aiv)
    - [10.2.3. Perform Key Exchange Between TAF and AIV](#1023-perform-key-exchange-between-taf-and-aiv)
    - [10.2.4. Start the Migration App](#1024-start-the-migration-app)

## 1. Machine IP Addresses

- **Main Machine:** `192.168.10.101`
- **Vehicle Computer (VC1 / ECU1):** `192.168.10.102`
- **Vehicle Computer (VC2 / ECU2):** `192.168.10.103`

## 2. Execution Locations

- **On the  Main Machine:** Run Kafka, `demo.py`, and CACC with `v1` and `v3`.
- **On VC2:** Start CACC with `v2`, TAF, AIV, CIV, and the Migration App.

## 3. Setup: Configuration File

Each user must customize their configuration file. The template is available at `./python/utils/template_config.yaml`.

- To generate your configuration file, execute the command: `./bash_scripts/create_config.sh`.

## 4. Setup: AIV

- Setup the TAF using it's [AIV README](./aiv/README.md)

## 5. Setup: TAF

- Setup the TAF using it's [TAF README](./taf/README.md)

## 6. Setup: Kafka

**Important Notes:** This only needs to be run on the  Main machine. If installed for one user it need not be installed for other users beacuse kafka executable installed at `/opt/kafka/` is shared by all the users.`

**Warning:** Install Kafka only using the main user which has SUDO permissions.

## 6.1. On Main Machine

To install and configure Kafka and Zookeeper on the  Main Machine, use the provided setup script:

- Run the Kafka setup script:

  ```bash
  ./bash_scripts/kafka_setup.sh
  ```

- Follow the prompts and select the following options as needed:
  - `1` – Install Kafka binaries.
  - `3` – Start Kafka Server and Zookeeper.
  - (Recommended) `5` – Enable Kafka Server and Zookeeper to auto-start on machine reboot.
  - `7` – Create Kafka topics as specified in `./python/utils/template_config.yaml`.
  - `8` – List already active Kafka topics to verify creation.
  - `9` – Display available Kafka commands for further help.

This process ensures Kafka and Zookeeper are properly installed, started, and configured for the demo.

## 6.2. On VC1/VC2

To verify Kafka connectivity from VC1 or VC2, you can use the Kafka setup script as follows:

- Run the Kafka setup script:

  ```bash
  ./bash_scripts/kafka_setup.sh
  ```

- Follow the prompts and select the following options as needed:
  - `1` – Install Kafka binaries (if not already installed).
  - `3` – List already active Kafka topics to verify that the topics are accessible from this machine.
  - `4` – Display available Kafka commands for further help.

This ensures that VC1/VC2 can communicate with the Kafka server running on the  Main Machine.

## 7. Setup: Python

**Important Notes:**  
Python and its dependencies must be set up on all machines and for all users who will run the demo.

### 7.1. Automated Method Using Bash Script

To ensure all required Python scripts and libraries are installed, run the following command from the root directory of the repository:

```bash
./bash_scripts/python_setup.sh
```

When prompted, select the `ipl` (install python libraries) option.

### 7.2. Manual Setup

If you prefer to set up the Python environment manually, follow these steps:

- Create a Python virtual environment:

  ```bash
  python3 -m venv ./kafka_python_venv
  ```

- Activate the virtual environment:

  ```bash
  source ./kafka_python_venv/bin/activate
  ```

- Install the required Python libraries:

  ```bash
  pip install -r ./requirements.txt
  ```

**Note:**  
If you encounter issues with the `kafka-python` library during installation, you may need to install it manually:

```bash
pip install kafka-python
```

## 8. Run: SUMO

**Important Notes:** SUMO only needs to be installed on the  Main machine. If it is installed for one user, there is no need to install it for other users, since the `sumo` executable installed via `apt` is available system-wide.

To install SUMO, follow the official instructions at [the SUMO website](https://sumo.dlr.de/docs/Installing/index.html).

To start SUMO for the demo:

- Activate the Python virtual environment:

  ```bash
  source ./kafka_python_venv/bin/activate
  ```

- Run the SUMO demo script:

  ```bash
  python3 ./python/sumo/demo.py
  ```

## 9. Run: CACC

- Activate the Python virtual environment by running the following command:

  ```bash
  source ./kafka_python_venv/bin/activate
  ```

  **Note:** Always run the CACC application using the full/absolute path as specified below.

- Execute the following commands to start the CACC application:

  - **On the  Main Machine:**

    ```bash
    python3 ./python/cc/cc.py v1
    python3 ./python/cc/cc.py v3
    ```

  - **On VC2:**

    ```bash
    python3 ./python/cc/cc.py v2
    ```

    The migration process will be initiated with this instance.

## 10. Run: TAF, AIV, and Migration App

The migration app takes care of the migration of CACC from VC2 to VC1 (and vice versa).

### 10.1. Automated Method Using Bash Script

To simplify the process, you can use the provided bash script to start TAF, AIV, and the Migration App. Run the following command and follow the prompts to select the appropriate options:

```bash
./bash_scripts/start_taf_aiv_migration.sh
```

### 10.2. Manual Setup

If you prefer to start each component manually, follow the steps below. Ensure that each command is executed in a separate shell window.

#### 10.2.1. Start TAF

- Set the TAF configuration file path:

  ```bash
  export TAF_CONFIG=$PWD/taf/res/taf.json
  ```

- Navigate to the TAF directory:

  ```bash
  cd ./taf
  ```

- Start the TAF application:

  ```bash
  ./go-taf
  ```

#### 10.2.2. Start AIV

- Navigate to the AIV directory:

  ```bash
  cd ./aiv
  ```

- Start the AIV application using one of the following commands, depending on your Kafka server setup:

  - **For local Kafka server:**

      ```bash
      python3 aiv.py --broker_ip=192.168.10.101:9092 --mode=mutable
      ```

      or

      ```bash
      python3 aiv.py --broker_ip=192.168.10.101:9092 --mode=trust_check
      ```

#### 10.2.3. Perform Key Exchange Between TAF and AIV

- Navigate back to the TAF directory:

  ```bash
  cd ./taf
  ```

- Execute the key exchange script:

  ```bash
  ./copy_keys.sh
  ```

#### 10.2.4. Start the Migration App

- Activate the Python virtual environment:

  ```bash
  source ./kafka_python_venv/bin/activate
  ```

- Navigate to the Migration App directory:

  ```bash
  cd ./python/migration_app/
  ```

- Start the Migration App:

  - For imminent driving scenario:
  
    ```bash
    python3 migration_app.py
    ```

  - For upcoming driving scenario:
  
    ```bash
    python3 migration_app_upcoming.py
    ```
