# Migration App

The migration app is the one taking care of the migration of the `CACC` instance for `v2` from VC2 to VC1 and vice versa.
It is a standalone component that can be run independently from other modules in this project, however it needs the Greek Kafka server to be running.

It has one main file, the `migration_app.py` which is the main file for the migration app and six handlers to handle the different messages.
Most notably, the `tas_notify_handler` is the one that takes care of the migration by invoking the migration script.

The messages sent are the custom `TASMessage` which is using the main `pydantic` model to model the messages.
Using custom functions (`TopicBoundProducer` and `TopicBoundConsumer`), the messages are send to the correct topic.
They are parsed from and to json using the pydantic `model_dump()` function.
If you want to know more about the custom utils functions, please refer to `/utils/utils.py`

Currently, the migration takes place when the calculated trust value/projected probability is below the RTL threshold.
If there are any errors in the migration, it is most likely to be found in the `tas_notify_handler`.

## How to start

1. start the virtual environment (if not done already): `cd connect-demo` `source ./kafka_python_venv/bin/activate`
2. to start the migration app:
    1. `cd connect-demo/python/migration_app/`
    2. `python3 migration_app.py`
