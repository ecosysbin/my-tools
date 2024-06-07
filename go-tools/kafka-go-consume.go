package main

import (
	"context"
	"fmt"
	"log"

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
