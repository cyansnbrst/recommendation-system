#/bin/bash

/opt/bitnami/kafka/bin/kafka-topics.sh --create --topic user_updates --bootstrap-server kafka:9092
echo "topic user_updates was created"

/opt/bitnami/kafka/bin/kafka-topics.sh --create --topic product_updates --bootstrap-server kafka:9092
echo "topic product_updates was created"