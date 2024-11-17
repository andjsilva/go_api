package di

import (
	"api/internal/config"
	"api/internal/infra/rabbitmq"
	"api/internal/infra/repository"
	"api/internal/infra/web"
	"api/internal/messaging"
	"api/internal/usecase"
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	DB              *pgxpool.Pool
	RabbitMQ        *rabbitmq.RabbitMQ
	UseCases        UseCases
	Handlers        *web.ProductHandlers
	ConsumerService *messaging.ConsumerService
}

type UseCases struct {
	CreateProductUseCase *usecase.CreateProductUseCase
	ListProductsUseCase  *usecase.ListProductsUseCase
}

func InitializeDependencies(cfg *config.Config) (*Container, error) {
	// Banco de dados
	db, err := pgxpool.New(context.Background(), cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}

	// RabbitMQ
	rabbit, err := rabbitmq.NewRabbitMQ(cfg.RabbitMQDSN)
	if err != nil {
		db.Close()
		return nil, err
	}

	// Repositórios
	productRepo := repository.NewProductRepository(db)

	// Casos de Uso
	createProductUseCase := usecase.NewCreateProductUseCase(productRepo)
	listProductsUseCase := usecase.NewListProductsUseCase(productRepo)

	// Handlers HTTP
	handlers := web.NewProductHandlers(createProductUseCase, listProductsUseCase)

	// Serviço de Consumidor
	consumerService := messaging.NewConsumerService(rabbit, createProductUseCase)

	return &Container{
		DB:              db,
		RabbitMQ:        rabbit,
		UseCases:        UseCases{createProductUseCase, listProductsUseCase},
		Handlers:        handlers,
		ConsumerService: consumerService,
	}, nil
}

func (c *Container) Close() {
	c.DB.Close()
	c.RabbitMQ.Close()
}

func SetupRouter(container *Container) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/products", container.Handlers.CreateProductHandler)
	r.Get("/products", container.Handlers.ListProductsHandler)
	return r
}
