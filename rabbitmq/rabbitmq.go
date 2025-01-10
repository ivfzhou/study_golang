package main

import (
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	samqp "github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

func main() {}

func SendAndConsume() {
	amqp.SetLogger(log.Default())
	connection, err := amqp.Dial("amqp://guest:guest@ivfzhou-debian:5672/")
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	notifyClose := channel.NotifyClose(make(chan *amqp.Error))
	go func() {
		for v := range notifyClose {
			if v != nil {
				fmt.Println(v)
			}
		}
	}()

	// name: 为空串意味着生成随机队列名
	// durable: 当服务重启后，队列是否保留，持久化消息是否保留。只能绑定对应是否持久化交换机
	// autoDelete: 当最后一个消费者关闭连接时，是否删除队列
	// exclusive: 只允许声明队列的连接操作队列，当连接关闭时，队列删除
	// noWait: 不等待服务器返回结果，失败了没返回 err
	queue, err := channel.QueueDeclare("SendAndConsume", false, true, false, false, nil)
	if err != nil {
		panic(err)
	}

	// 发送方
	go func() {
		// 发送消息失败，接收通知
		notifyReturn := channel.NotifyReturn(make(chan amqp.Return))
		go func() {
			for v := range notifyReturn {
				fmt.Printf("unpublish message: %v\n", v)
			}
		}()

		// 每个队列都会和默认交换机绑定，路由键是队列名。
		// 当 mandatory 为 true，如果队列不存在将通知 NotifyReturn
		// 当 immediate 为 true，如果消息没有消费者接收，将通知 NotifyReturn
		err = channel.Publish(amqp.DefaultExchange, queue.Name, true, false, amqp.Publishing{
			Body: []byte("Hello World"),
		})
		if err != nil {
			panic(err)
		}
	}()

	// 接收方
	// 同其他消费者轮流消费消息
	// autoAck: 是否自动确认消息已消费。默认超过 30 分钟没确认将自动确认。如果消费者 connection/channel/tcp 关闭，消息将会再次入列
	// autoAck: 服务器在发送消息前，就确认消息被消费了
	// exclusive: 只让一个消费者消费
	// noWait: 不等待服务器确认就开始等待消息消费，如果有错误发生，channel 将会关闭
	consume, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	msg := <-consume
	fmt.Printf("Received a message: %s\n", msg.Body)
	time.Sleep(10 * time.Second)
}

func Confirm() {
	amqp.SetLogger(log.Default())
	connection, err := amqp.Dial("amqp://guest:guest@ivfzhou-debian:5672/")
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()
	queue, err := channel.QueueDeclare("", false, true, false, false, nil)
	if err != nil {
		panic(err)
	}

	go func() {
		err = channel.Confirm(false)
		if err != nil {
			panic(err)
		}

		// channel 关闭则 publish 关闭
		publish := channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		go func() {
			for v := range publish {
				fmt.Printf("publish message: %v\n", v)
			}
		}()

		err = channel.Publish(amqp.DefaultExchange, queue.Name, false, false, amqp.Publishing{
			Body: []byte("Hello World"),
		})
		if err != nil {
			panic(err)
		}
		cf, err := channel.PublishWithDeferredConfirm(amqp.DefaultExchange, queue.Name, false, false, amqp.Publishing{
			Body: []byte("Hello World"),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(cf.Acked())
	}()

	consume, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	msg := <-consume
	fmt.Printf("Received a message: %s\n", msg.Body)
	err = msg.Ack(false)
	if err != nil {
		panic(err)
	}
	msg = <-consume
	fmt.Printf("Received a message: %s\n", msg.Body)
	err = msg.Ack(false)
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
}

func DeadQueue() {
	amqp.SetLogger(log.Default())
	connection, err := amqp.Dial("amqp://guest:guest@ivfzhou-debian:5672/")
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()
	queue, err := channel.QueueDeclare("my-queue", false, true, false, false, amqp.Table{
		"x-dead-letter-routing-key": "dead-letter-queue",
		"x-dead-letter-exchange":    amqp.DefaultExchange,
		"x-message-ttl":             1000 * 10,
	})
	if err != nil {
		panic(err)
	}
	dqueue, err := channel.QueueDeclare("dead-letter-queue", false, true, false, false, nil)
	if err != nil {
		panic(err)
	}

	go func() {
		err = channel.Publish(amqp.DefaultExchange, queue.Name, true, false, amqp.Publishing{
			Body: []byte("Hello World"),
		})
		if err != nil {
			panic(err)
		}
	}()

	consume, err := channel.Consume(dqueue.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	msg := <-consume
	fmt.Printf("Received a message: %s\n", msg.Body)
	time.Sleep(10 * time.Second)
}

func TX() {
	amqp.SetLogger(log.Default())
	connection, err := amqp.Dial("amqp://guest:guest@ivfzhou-debian:5672/")
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()
	queue, err := channel.QueueDeclare("", false, true, false, false, nil)
	if err != nil {
		panic(err)
	}

	go func() {
		err = channel.Tx()
		if err != nil {
			panic(err)
		}

		err = channel.Publish(amqp.DefaultExchange, queue.Name, true, false, amqp.Publishing{
			Body: []byte("Hello World"),
		})
		if err != nil {
			panic(err)
		}

		err = channel.TxCommit()
		if err != nil {
			panic(err)
		}
	}()

	consume, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	msg := <-consume
	fmt.Printf("Received a message: %s\n", msg.Body)
	time.Sleep(10 * time.Second)
}

func Ack() {
	amqp.SetLogger(log.Default())
	connection, err := amqp.Dial("amqp://guest:guest@ivfzhou-debian:5672/")
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()
	queue, err := channel.QueueDeclare("Ack", false, true, false, false, nil)
	if err != nil {
		panic(err)
	}

	// 发送方
	go func() {
		err = channel.Publish("", queue.Name, false, false, amqp.Publishing{
			Body: []byte("Hello World"),
		})
		if err != nil {
			panic(err)
		}
	}()

	// 接收方
	// prefetchCount: 在消费者 ack 前，允许 MQ 推送过来的消息数
	// prefetchSize: MQ 将推送这么多字节没有确认的消息给消费者
	// global: 只应用于当前 channel，还是当前 connection 下的所有 channel
	err = channel.Qos(1, 0, false)
	if err != nil {
		panic(err)
	}

	// 只能在同一个接收消息的 channel 下 ack 消息
	consume, err := channel.Consume(queue.Name, "my-consumer", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	msg := <-consume
	fmt.Printf("Received a message: %s\n", msg.Body)

	// 没有 ack 的消息将被放回队列尾部
	err = msg.Ack(false)
	if err != nil {
		panic(err)
	}

	// 停止消费。要确保 consume 全部读出，close 了
	channel.Cancel(msg.ConsumerTag, false)
}

func Fanout() {
	amqp.SetLogger(log.Default())
	connection, err := amqp.Dial("amqp://guest:guest@ivfzhou-debian:5672/")
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	// type: fanout, direct, topic, headers
	// 交换机名以 amq. 开头已被使用
	// durable 服务重启后，交换器还存在
	// autoDelete: 没有绑定后删除交换机
	err = channel.ExchangeDeclare("FanoutName", "fanout", false, true, false, false, nil)
	if err != nil {
		panic(err)
	}
	queue, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		panic(err)
	}
	err = channel.QueueBind(queue.Name, "", "FanoutName", false, nil)
	if err != nil {
		panic(err)
	}

	// 发送方
	go func() {
		err = channel.Publish("FanoutName", "", false, false, amqp.Publishing{Body: []byte("Hello World")})
		if err != nil {
			panic(err)
		}
	}()

	// 接收方
	consume, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	msg := <-consume
	fmt.Printf("Received a message: %s\n", msg.Body)
}

func Direct() {
	amqp.SetLogger(log.Default())
	connection, err := amqp.Dial("amqp://guest:guest@ivfzhou-debian:5672/")
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()
	err = channel.ExchangeDeclare("DirectName", "direct", false, true, false, false, nil)
	if err != nil {
		panic(err)
	}
	queue, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		panic(err)
	}
	err = channel.QueueBind(queue.Name, "key", "DirectName", false, nil)
	if err != nil {
		panic(err)
	}

	// 发送方
	go func() {
		err = channel.Publish("DirectName", "key", false, false, amqp.Publishing{Body: []byte("Hello World")})
		if err != nil {
			panic(err)
		}
	}()

	// 接收方
	consume, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	msg := <-consume
	fmt.Printf("Received a message: %s\n", msg.Body)
}

func Topic() {
	amqp.SetLogger(log.Default())
	connection, err := amqp.Dial("amqp://guest:guest@ivfzhou-debian:5672/")
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()
	err = channel.ExchangeDeclare("TopicName", "topic", false, true, false, false, nil)
	if err != nil {
		panic(err)
	}
	queue, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		panic(err)
	}

	// routing key 最长 255 字节，有点号分隔
	// * 号表示通配一个单词
	// # 号表示通配零到多个单词
	// 同一个队列被多个路由键匹配将只被递交一次消息
	err = channel.QueueBind(queue.Name, "key.*", "TopicName", false, nil)
	if err != nil {
		panic(err)
	}

	// 发送方
	go func() {
		err = channel.Publish("TopicName", "key.one", false, false, amqp.Publishing{Body: []byte("Hello World")})
		if err != nil {
			panic(err)
		}
	}()

	// 接收方
	consume, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	msg := <-consume
	fmt.Printf("Received a message: %s\n", msg.Body)
}

func RPC() {
	amqp.SetLogger(log.Default())
	connection, err := amqp.Dial("amqp://guest:guest@ivfzhou-debian:5672/")
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	// 服务器侧
	queue, err := channel.QueueDeclare("rpc", false, true, false, false, nil)
	if err != nil {
		panic(err)
	}
	err = channel.Qos(1, 0, false)
	if err != nil {
		panic(err)
	}
	consume, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	go func() {
		for msg := range consume {
			fmt.Printf("Received a message: %s\n", msg.Body)
			channel.Publish("", msg.ReplyTo, false, false, amqp.Publishing{
				CorrelationId: msg.CorrelationId,
				Body:          []byte("World"),
			})
			err = msg.Ack(false)
			if err != nil {
				panic(err)
			}
		}
	}()

	// 客户端侧
	queue, err = channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		panic(err)
	}
	consume1, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	err = channel.Publish("", "rpc", false, false, amqp.Publishing{
		CorrelationId: "1",
		ReplyTo:       queue.Name,
		Body:          []byte("Hello"),
	})
	if err != nil {
		panic(err)
	}
	msg := <-consume1
	fmt.Printf("Received a message: %s\n", msg.Body)
}

func Stream() {
	env, err := stream.NewEnvironment(stream.NewEnvironmentOptions().SetHost("ivfzhou-debian").SetUser("guest").SetPassword("guest").SetPort(5552))
	if err != nil {
		panic(err)
	}
	err = env.DeclareStream("streamName", &stream.StreamOptions{MaxLengthBytes: stream.ByteCapacity{}.GB(2)})
	if err != nil {
		panic(err)
	}

	// 发送方
	producer, err := env.NewProducer("streamName", stream.NewProducerOptions())
	if err != nil {
		panic(err)
	}
	err = producer.Send(samqp.NewMessage([]byte("Hello world")))
	if err != nil {
		panic(err)
	}
	err = producer.Close()
	if err != nil {
		panic(err)
	}

	// 接收方
	messagesHandler := func(consumerContext stream.ConsumerContext, message *samqp.Message) {
		fmt.Printf("Stream: %s - Received message: %s\n", consumerContext.Consumer.GetStreamName(), message.Data)
	}
	consumer, err := env.NewConsumer("streamName", messagesHandler, stream.NewConsumerOptions().SetOffset(stream.OffsetSpecification{}.First()))
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
	err = consumer.Close()
	if err != nil {
		panic(err)
	}
}

func TrackStream() {
	env, err := stream.NewEnvironment(stream.NewEnvironmentOptions().SetHost("ivfzhou-debian").SetUser("guest").SetPassword("guest").SetPort(5552))
	if err != nil {
		panic(err)
	}
	err = env.DeclareStream("streamName", &stream.StreamOptions{MaxLengthBytes: stream.ByteCapacity{}.GB(2)})
	if err != nil {
		panic(err)
	}

	go func() {
		producer, err := env.NewProducer("streamName", stream.NewProducerOptions())
		if err != nil {
			panic(err)
		}
		messageCount := 10
		ch := make(chan bool)
		chPublishConfirm := producer.NotifyPublishConfirmation()
		go func() {
			confirmedCount := 0
			for confirmed := range chPublishConfirm {
				for _, msg := range confirmed {
					if msg.IsConfirmed() {
						confirmedCount++
						if confirmedCount == messageCount {
							ch <- true
						}
					}
				}
			}
		}()

		fmt.Printf("Publishing %d messages\n", messageCount)
		for i := 0; i < messageCount; i++ {
			var body string
			if i == messageCount-1 {
				body = "marker"
			} else {
				body = "hello"
			}
			producer.Send(samqp.NewMessage([]byte(body)))
		}
		_ = <-ch
		fmt.Println("Messages confirmed")
		producer.Close()
	}()

	var firstOffset int64 = -1
	var messageCount int64 = -1
	var lastOffset atomic.Int64
	ch := make(chan bool)
	messagesHandler := func(consumerContext stream.ConsumerContext, message *samqp.Message) {
		if atomic.CompareAndSwapInt64(&firstOffset, -1, consumerContext.Consumer.GetOffset()) {
			fmt.Println("First message received.")
		}
		if atomic.AddInt64(&messageCount, 1)%10 == 0 {
			consumerContext.Consumer.StoreOffset() // 保存位置
		}
		if string(message.GetData()) == "marker" {
			lastOffset.Store(consumerContext.Consumer.GetOffset())
			consumerContext.Consumer.StoreOffset()
			_ = consumerContext.Consumer.Close()
			ch <- true
		}
	}

	offsetSpecification := stream.OffsetSpecification{}.First() // stream.OffsetSpecification{}.Offset(42) stream.OffsetSpecification{}.Next()
	consumerName := "offset-tracking-tutorial"
	storedOffset, err := env.QueryOffset(consumerName, "streamName")
	if errors.Is(err, stream.OffsetNotFoundError) {
		offsetSpecification = stream.OffsetSpecification{}.First()
	} else {
		offsetSpecification = stream.OffsetSpecification{}.Offset(storedOffset + 1)
	}
	_, _ = env.NewConsumer("streamName", messagesHandler, stream.NewConsumerOptions().
		SetOffset(offsetSpecification).
		SetManualCommit().             // 手动位置追踪
		SetConsumerName(consumerName), // 设置消费者名称
	)
	fmt.Println("Started consuming...")
	_ = <-ch
	fmt.Printf("Done consuming, first offset %d, last offset %d.\n", firstOffset, lastOffset.Load())
}
