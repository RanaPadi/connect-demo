import os
import sys
import subprocess
import time
from typing import Dict, Any
from python.utils.utils import setup_logger, load_config
from pandas.plotting import table

# Adjust the system path to include the parent directories
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..')))



# Set up loggers for different purposes
migration_logger = setup_logger("migration_app")
benchmark_logger = setup_logger("benchmark")

# Load configuration settings
config = load_config()

# Global Variables
source_vehicle_computer = "VC2"        # Source VC
target_vehicle_computer = "VC1"        # Target VC

# RTL (Required Trust Level) values, hardcoded for now
rtl_bdu_values = {
    "belief": 0.6,
    "disbelief": 0.2,
    "uncertainty": 0.2,
    "base_rate": 0.5,
    "trust_value": 0.7
}

# Initialize global counters for benchmarks
benchmark_3_counter = 0
benchmark_8_counter = 0


def compare_all_trust_values(rtl_bdu_values: Dict[str, float], aggregated_propositions: Dict[str, Dict[str, float]]) -> bool:
    """
    Compares trust-related values in a set of aggregated propositions against reference threshold values.

    The function iterates through the `aggregated_propositions` dictionary, which contains trust-related metrics
    for various propositions, and checks if each metric meets the conditions specified relative to `rtl_bdu_values`.
    The conditions are as follows:
        - The "belief" value must be greater than `rtl_bdu_values["belief"]`.
        - The "disbelief" value must be less than `rtl_bdu_values["disbelief"]`.
        - The "uncertainty" value must be less than `rtl_bdu_values["uncertainty"]`.
        - The "trust_value" must be greater than `rtl_bdu_values["trust_value"]`.

    If all conditions are satisfied for every proposition, the function allows the loop to proceed. Otherwise,
    it continues to the next proposition.

    Args:
        rtl_bdu_values (Dict[str, float]): A dictionary containing reference values for trust metrics.
            Expected keys: "belief", "disbelief", "uncertainty", "trust_value".
        aggregated_propositions (Dict[str, Dict[str, float]]): A dictionary containing trust metrics for multiple propositions.
            Each key represents a proposition, and the value is another dictionary with the same trust metric keys.

    Returns:
        bool: Currently, the function does not return a meaningful result. You may want to modify it
              to return True/False based on the conditions or the number of valid propositions.
    """

    for prop, values in aggregated_propositions.items():
        if values["belief"] <= rtl_bdu_values["belief"]: continue
        if values["disbelief"] >= rtl_bdu_values["disbelief"]: continue
        if values["uncertainty"] >= rtl_bdu_values["uncertainty"]: continue
        if values["trust_value"] <= rtl_bdu_values["trust_value"]: continue





def execute_migration(target_vehicle_computer: str) -> None:
    """
    Execute the migration script for the specified target vehicle computer.

    Args:
        target_vehicle_computer (str): The target vehicle computer for the migration script.

    Returns:
        None
    """
    global benchmark_8_counter
    benchmark_8_counter += 1
    current_count = benchmark_8_counter

    migration_logger.info(f"Executing migration script for {target_vehicle_computer}")

    # Log the start of Benchmark 8 with counter

    benchmark_start_time = time.time()

    try:
        # Path to the migration Bash script
        migration_script_path = config["env_variables"]["migration_script"]

        # Variable to pass to the script
        migration_argument = target_vehicle_computer

        # Run the migration script with the target container as an argument
        result = subprocess.run(
            [migration_script_path, migration_argument],
            capture_output=True,
            text=True,
            check=True  # Raises CalledProcessError if the command exits with a non-zero status
        )
        migration_logger.info(f"Migration Script Output:\n{result.stdout}")
    except subprocess.CalledProcessError as e:
        migration_logger.error(f"Migration Script Errors:\n{e.stderr}")
        migration_logger.error(f"Migration Script Exit Code: {e.returncode}")
    finally:
        benchmark_end_time = time.time()
        benchmark_duration_ms = (benchmark_end_time - benchmark_start_time) * 1000
        benchmark_logger.info(f"Benchmark 8 ({current_count}/200): [{source_vehicle_computer}->{target_vehicle_computer}] {benchmark_duration_ms:.3f}ms.")


def compute_and_evaluate_trust_values(
        propositions_data: Dict[str, Dict[str, float]],
        benchmark_3_start_time: float,
        benchmark_3_count: int
) -> None:
    """
    Computes and evaluates trust values for propositions, checking the trustworthiness of source and target vehicle computers
    based on the computed values and benchmarks. Initiates migration if necessary.

    The trust value for each proposition is calculated using the formula:
        Trust Value = belief + uncertainty * base_rate

    The function also logs benchmark times for evaluation and checks if the source and target vehicles meet the required trust thresholds.

    Args:
        propositions_data (Dict[str, Dict[str, float]]): A dictionary containing proposition IDs as keys and their corresponding trust-related values
            as dictionaries. Each value dictionary must contain keys: "belief", "uncertainty", and "base_rate".
        benchmark_3_start_time (float): The start time of Benchmark 3, used for calculating its duration.
        benchmark_3_count (int): The current count or index of Benchmark 3 for logging purposes.

    Global Variables:
        source_vehicle_computer (str): The source vehicle computer of the migration (e.g., VC2). This is used in the comparison for trustworthiness.
        target_vehicle_computer (str): The target vehicle computer of the migration (e.g., VC1). It is swapped after the first migration to ensure bi-directional migration.

    RTL Thresholds:
        The function uses the `rtl_bdu_values` to compare trust values. These include:
            - belief
            - disbelief
            - uncertainty
            - base_rate
            - trust_value (used for comparison)

    Returns:
        None: The function does not return a value, but it performs side effects such as logging benchmarks and initiating migration based on trust evaluations.

    Side Effects:
        - Logs benchmarks and the decision-making process related to trust evaluation and migration.
        - Executes migration via the `execute_migration` function if the target vehicle computer is deemed trustworthy.

    Notes:
        - The source vehicle is considered trustworthy if its trust value meets or exceeds the threshold in `rtl_bdu_values["trust_value"]`.
        - If the source is not trustworthy, the function checks the target vehicle for trustworthiness and performs a migration if the target is trustworthy.
        - The global variables `source_vehicle_computer` and `target_vehicle_computer` are swapped after the first migration for bi-directional migration.
    """

    global source_vehicle_computer, target_vehicle_computer

    #migration_logger.debug(f"Computing trust values for Source {source_vehicle_computer} and Target {target_vehicle_computer}.")

    # Compute ATL values using dictionary comprehension
    atl_values = {
        prop_id: bdu["belief"] + bdu["uncertainty"] * bdu["base_rate"]
        for prop_id, bdu in propositions_data.items()
    }

    # Aggregate all proposition data with their ATL values
    aggregated_propositions = {
        prop_id: {**bdu, "trust_value": atl}
        for prop_id, (bdu, atl) in zip(propositions_data.keys(), zip(propositions_data.values(), atl_values.values()))
    }

    # Log aggregated data
    migration_logger.info(f"Aggregated Propositions: {aggregated_propositions}")

    # Benchmark 3: Check if source is trustworthy and measure duration
    source_trustworthy = atl_values[source_vehicle_computer] >= rtl_bdu_values["trust_value"]
    benchmark_3_end_time = time.time()
    benchmark_3_duration_ms = (benchmark_3_end_time - benchmark_3_start_time) * 1000
    benchmark_logger.info(f"Benchmark 3 ({benchmark_3_count}/200): {benchmark_3_duration_ms:.3f}ms")

    # Migration requirement based on trustworthiness of the source
    migration_required = not source_trustworthy


    # If source is not trustworthy, proceed to benchmark 7 for target trust check
    if migration_required:
        migration_logger.info(f"Source '{source_vehicle_computer}' trust below RTL ({atl_values[source_vehicle_computer]} < {rtl_bdu_values['trust_value']}). Migration check in progress...")

        # Benchmark 7: Check if target is trustworthy
        benchmark_7_start_time = time.time()  # Start time of target trust check
        # target_trustworthy = atl_values[target_vehicle_computer] >= rtl_bdu_values["trust_value"]
        target_trustworthy = atl_values[target_vehicle_computer] <= rtl_bdu_values["trust_value"] #Added_for_testing

        # If target is trustworthy, initiate migration
        if target_trustworthy:

            benchmark_7_end_time = time.time()
            benchmark_7_duration_ms = (benchmark_7_end_time - benchmark_7_start_time) * 1000

            benchmark_logger.info(f"Benchmark 7: {benchmark_7_duration_ms:.3f}ms")
            migration_logger.info(f"-> Migration required [{source_vehicle_computer}->{target_vehicle_computer}]. Initiating migration...")

            execute_migration(target_vehicle_computer=target_vehicle_computer)

            # Swap source and target for bi-directional migration
            source_vehicle_computer, target_vehicle_computer = target_vehicle_computer, source_vehicle_computer
        else:
            migration_logger.info(f"-> Target [{target_vehicle_computer}] is not trustworthy. Migration not initiated.")
    else:
        migration_logger.info(f"-> Source [{source_vehicle_computer}] is trustworthy. No migration required.")




def handle_tas_notify_response(data: Any) -> None:
    """
    Handles TAS_NOTIFY_RESPONSE messages, extracts the proposition data, and processes it for trust evaluation.

    This function processes the TAS_NOTIFY_RESPONSE message by extracting the relevant proposition data,
    including trust-related metrics like belief, disbelief, uncertainty, and base rate. The data is then passed
    to `compute_and_evaluate_trust_values` for further evaluation and migration handling.

    Args:
        data (Any): The data received in the TAS_NOTIFY_RESPONSE message. This is expected to contain updates with
                    proposition data, including trust metrics such as belief, disbelief, uncertainty, and base rate.

    Returns:
        None: The function does not return any value, but performs side effects such as processing the proposition
              data and invoking `compute_and_evaluate_trust_values`.

    Side Effects:
        - Increments the `benchmark_3_counter` global variable.
        - Logs the reception of the `TAS_NOTIFY_RESPONSE` message and the start of Benchmark 3.
        - Extracts the proposition data and forwards it to the `compute_and_evaluate_trust_values` function for further processing.
        - The benchmark logging related to Benchmark 3 and trust value computation is performed.

    Notes:
        - This function assumes the `data` contains a dictionary with the structure:
            - `message` -> contains a list of `updates` -> each `update` contains `propositions` ->
            - each `proposition` contains `propositionId` and `actualTrustworthinessLevel`.
        - Only `propositions` with a `propositionId` and an `actualTrustworthinessLevel` of type "SUBJECTIVE_LOGIC_OPINION" are considered.
        - The extracted proposition data is passed to `compute_and_evaluate_trust_values`, which handles trust value computation.
    """

    global benchmark_3_counter

    # Increment Benchmark 3 counter
    benchmark_3_counter += 1
    current_count = benchmark_3_counter

    # Log the start of Benchmark 3 with counter
    benchmark_3_start_time = time.time()

    migration_logger.info("Received TAS_NOTIFY_RESPONSE")

    # Extract updates from the message using dictionary comprehension
    extracted_proposition_data = {
        proposition.get("propositionId"): {
            "belief": tl.get("output", {}).get("belief", 0.0),
            "disbelief": tl.get("output", {}).get("disbelief", 0.0),
            "uncertainty": tl.get("output", {}).get("uncertainty", 0.0),
            "base_rate": tl.get("output", {}).get("baseRate", 0.0)
        }
        for update in data.get("message", {}).get("updates", [])
        for proposition in update.get("propositions", [])
        if proposition.get("propositionId") is not None
        for tl in proposition.get("actualTrustworthinessLevel", [])
        if tl.get("type") == "SUBJECTIVE_LOGIC_OPINION"
    }

    # Perform comparisons and calculations
    compute_and_evaluate_trust_values(
        propositions_data=extracted_proposition_data,
        benchmark_3_start_time=benchmark_3_start_time,
        benchmark_3_count=current_count
    )
