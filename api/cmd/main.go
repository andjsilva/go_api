package main

import (
	"api/internal/config"
	"api/internal/di"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Carregar configurações
	cfg := config.LoadConfig()

	// Inicializar dependências
	container, err := di.InitializeDependencies(cfg)
	if err != nil {
		log.Fatalf("Erro ao inicializar dependências: %v", err)
	}
	defer container.Close()

	// Capturar sinais do sistema para desligamento seguro
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Printf("Recebido sinal: %s. Encerrando...", sig)
		container.Close()
		os.Exit(0)
	}()

	// Iniciar consumidores do RabbitMQ
	go func() {
		err := container.ConsumerService.StartConsumer(cfg.QueueName)
		if err != nil {
			log.Fatalf("Erro ao iniciar consumidor RabbitMQ: %v", err)
		}
	}()

	// Configurar rotas HTTP
	router := di.SetupRouter(container)

	fmt.Println("Servidor rodando na porta :8000")
	http.ListenAndServe(":8000", router)
}
