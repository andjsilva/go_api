package rabbitmq

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Connection *amqp091.Connection
	Channel    *amqp091.Channel
}

func NewRabbitMQ(dsn string) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir canal RabbitMQ: %v", err)
	}

	return &RabbitMQ{Connection: conn, Channel: ch}, nil
}

func (r *RabbitMQ) Close() {
	r.Channel.Close()
	r.Connection.Close()
}

func (r *RabbitMQ) ConfigureQos(prefetchCount int) error {
	return r.Channel.Qos(prefetchCount, 0, false)
}

func StartConsumer(r *RabbitMQ, queueName string, handler func([]byte) error) error {
	err := r.ConfigureQos(1)
	if err != nil {
		return fmt.Errorf("erro ao configurar QoS: %v", err)
	}

	msgs, err := r.Channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar consumidor: %v", err)
	}

	go func() {
		for msg := range msgs {
			log.Printf("Mensagem recebida: %s", string(msg.Body))
			if err := handler(msg.Body); err != nil {
				log.Printf("Erro ao processar mensagem: %v", err)
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
				log.Println("Mensagem processada com sucesso.")
			}
		}
	}()
	return nil
}
