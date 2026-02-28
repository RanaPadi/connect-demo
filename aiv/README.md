# Steps

rm CMakeCache.txt
cmake CMakeLists.txt
make

pip install cryptography==3.4.8
pip install confluent_kafka==2.4.0

python3 aiv.py