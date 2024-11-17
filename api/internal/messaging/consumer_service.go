package messaging

import (
	"api/internal/infra/rabbitmq"
	"api/internal/usecase"
	"encoding/json"
	"log"
)

type ConsumerService struct {
	rabbitMQ             *rabbitmq.RabbitMQ
	createProductUseCase *usecase.CreateProductUseCase
}

func NewConsumerService(rabbitMQ *rabbitmq.RabbitMQ, createProductUseCase *usecase.CreateProductUseCase) *ConsumerService {
	return &ConsumerService{
		rabbitMQ:             rabbitMQ,
		createProductUseCase: createProductUseCase,
	}
}

func (s *ConsumerService) StartConsumer(queueName string) error {
	return rabbitmq.StartConsumer(s.rabbitMQ, queueName, func(body []byte) error {
		var msg ProductMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			log.Printf("Erro ao desserializar mensagem: %v", err)
			return err
		}

		_, err := s.createProductUseCase.Execute(usecase.CreateProductInputDto{
			Name:  msg.Name,
			Price: msg.Price,
		})
		if err != nil {
			log.Printf("Erro ao criar produto: %v", err)
			return err
		}

		log.Printf("Produto criado com sucesso: %s", msg.Name)
		return nil
	})
}
