# Migration App

## How to start

1. start the virtual environment (if not done already): `cd connect-demo` `source ./kafka_python_venv/bin/activate`
2. to start the migration app:
   1. `cd connect-demo/python/migration_app/`
   2. `python3 migration_app.py`

## About the Migration App

The migration_app.py is the main file.
It follows the schema of the kafka story line and uses multiple handlers to handle the different messages.

The TASMessage is the main `pydantic` model to model the messages.
Using custom functions (TopicBoundProducer and TopicBoundConsumer), the messages are send to the correct topic.
They are parsed from and to json using the pydantic `model_dump()` function.
If you want to know more about the custom utils functions, please refer to `/utils/utils.py`

Currently, the migration takes place when the calculated trust value/projected probability is below the RTL threshold.
If there are any errors in the migration, it is most likely to be found in the `tas_notify_handler`.
