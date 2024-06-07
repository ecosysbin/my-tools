package main

import (
	"context"
	"fmt"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func Product() {
	// Kafka broker的地址。更复杂的场景会有多个broker。
	brokerAddress := "10.220.9.21:9092"
	// 主题名称，根据你的Kafka设置定义
	topic := "test"

	// 初始化一个kafka的writer实例，用于发送消息
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{brokerAddress},
		Topic:        topic,
		BatchTimeout: 200 * time.Millisecond,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		// 根据需要调整其他配置，例如批量大小、超时时间等
	})

	// 初始化消息
	msg := kafka.Message{
		Key:   []byte("Key-s"),
		Value: []byte("Hello Kafka2d2!"),
		Time:  time.Now(),
	}

	// 发送消息
	err := writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Fatalf("failed to write messages: %s", err)
	}

	// 输出结果并关闭writer
	fmt.Println("Message sent successfully")
	writer.Close()
}
