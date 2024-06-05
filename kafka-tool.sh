#!/bin/bash

# deploy kafka
wget https://downloads.apache.org/kafka/3.2.3/kafka_2.13-3.2.3.tgz
tar -zxvf kafka_2.13-3.2.3.tgz

# 进入kafka目录
cd kafka_2.12-2.5.0

# 修改配置文件
# vim kafka_2.12-2.5.0/config/server.properties
# 注:若需要外部访问  一定需要配置listeners   默认为本机名，跨机器无法解析主机名则无法访问，主机外访问配置ip或域名, 默认是主机名跨主机访问会解析出错。配置域名时需要外部也能正常解析
# listeners=PLAINTEXT://10.220.9.91:9092

# 启动zookeeper
./bin/zookeeper-server-start.sh -daemon ./config/zookeeper.properties
# 启动kafka
./bin/kafka-server-start.sh -daemon ./config/server.properties


# 创建topic
./bin/kafka-topics.sh --create --partitions 1 --replication-factor 1 --topic test1 --bootstrap-server 10.220.9.91:9092

# topic列表
./bin/kafka-topics.sh --list --zookeeper 10.220.9.91:2181

# topic 详情
./bin/kafka-topics.sh --describe --zookeeper 10.220.9.91:2181 --topic test1

# 生产
./bin/kafka-console-producer.sh --broker-list 10.220.9.91:9092 --topic test1

# 消费
./bin/kafka-console-consumer.sh --bootstrap-server 10.220.9.91:9092 --topic test1 --from-beginning