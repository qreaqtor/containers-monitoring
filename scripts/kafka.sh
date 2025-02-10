set -e

kafka-topics --bootstrap-server kafka:9092 --list

echo -e 'Creating kafka topics'
kafka-topics --bootstrap-server kafka:9092 --topic containers --create --partitions 3 --replication-factor 1
echo -e 'Successfully created the following topics:'

kafka-topics --bootstrap-server kafka:9092 --list