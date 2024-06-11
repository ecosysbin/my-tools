package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	//    "time"

	"github.com/segmentio/kafka-go"
)

func Consumer() {
	// Kafka broker的地址
	brokerAddress := "10.220.9.21:9092"
	// 主题名称
	topic := "test"

	// 初始化reader实例，连接到Kafka主题
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		GroupID:  "",   // 空字符串表示我们不使用consumer group
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer reader.Close()

	fmt.Println("开始消费消息...")

	// 死循环，持续读取并处理消息
	for {
		// 读取消息
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("无法读取消息: %v", err)
			continue
		}
		fmt.Printf("接收到的消息: %s = %s\n", string(m.Key), string(m.Value))
	}

	// 实际应用中，你可能需要有条件地打破循环，或者处理连接断开等情况
}

func DescribeTopic() {
	// Kafka broker的地址
	brokerAddress := "10.220.9.21:9092"
	// 主题名称
	topic := "test"

	conn, err := kafka.Dial("tcp", brokerAddress)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	// topicConfigs := []kafka.TopicConfig{
	// 	{
	// 		Topic:             topic,
	// 		NumPartitions:     1,
	// 		ReplicationFactor: 1,
	// 	},
	// }

	// 创建主题
	// err = controllerConn.CreateTopics(topicConfigs...)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// 问题：1.k8s集群外访问kafka生产或消费消息时，需要手动配置/etc/hosts域名, 不然会直接报错域名无法解析。
	//      2.k8s集群外（slurm login节点）生成消息时有几率（大约2/3）失败。报错 Kafka write errors (1/1), errors: [[6] Not Leader For Partition: the client attempted to send messages to a replica that is not the leader for some partition, the client's metadata are likely out of date]
	// 原因：1.客户端端先查询topic的partition列表，轮转选择一个partition, 根据partion的leader Host信息发送消息。kafka 集群返回的host是k8s内的域名，客户端无法解析报错。
	//      2.在生产消息时，使用k8s的service轮转发送消息，不能保证topic的partition的leader正好是发送的消息的目标broker。则消息发送几率性失败。
	// 解决方案：
	// 1.kafka hostnetwork方式部署，配置全局dns解析（能够覆盖slurm login节点）或者login节点配置/etc/hosts,两者都需要考虑单实例故障配置自动更新。
	// 2.slurm login节点部署在k8s集群内。客户端则可以直接使用k8s coredns域名解析即可。
	// 3.kafka单实例部署，存在单点故障
	partitions, err := controllerConn.ReadPartitions(topic)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("partitions: %+v\n", partitions)
	// partitions: [{Topic:test ID:0 Leader:{Host:kafka-controller-2.kafka-controller-headless.gcp.svc.cluster.local Port:9092 ID:2 Rack:} Replicas:[{Host:kafka-controller-2.kafka-controller-headless.gcp.svc.cluster.local Port:9092 ID:2 Rack:}] Isr:[{Host:kafka-controller-2.kafka-controller-headless.gcp.svc.cluster.local Port:9092 ID:2 Rack:}] OfflineReplicas:[] Error:<nil>}]
}
