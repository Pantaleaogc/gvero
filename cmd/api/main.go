package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Pantaleaogc/pantaleaocrmerp/pkg/database"
	"github.com/Pantaleaogc/pantaleaocrmerp/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Inicializar logger
	logger.Init()
	defer logger.Close()
	
	logger.InfoLogger.Println("Iniciando aplicação CRM/ERP...")

	// Carregar configurações
	if err := godotenv.Load("configs/.env"); err != nil {
		logger.InfoLogger.Printf("Aviso: arquivo .env não encontrado: %v", err)
	}

	// Configurações da porta
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Valor padrão
	}

	// Inicializar banco de dados
	db, err := database.InitDB()
	if err != nil {
		logger.ErrorLogger.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()
	logger.InfoLogger.Println("Conexão com banco de dados estabelecida")

	// Configurar rotas
	router := setupRoutes()

	// Configurar servidor
	addr := fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Canal para lidar com sinais de término
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor em uma goroutine
	go func() {
		logger.InfoLogger.Printf("Servidor iniciado na porta %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.ErrorLogger.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Esperar sinal de término
	go func() {
		<-quit
		logger.InfoLogger.Println("Servidor está sendo encerrado...")

		// Criar contexto com timeout para shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.ErrorLogger.Fatalf("Erro ao encerrar servidor: %v", err)
		}
		
		close(done)
	}()

	<-done
	logger.InfoLogger.Println("Servidor encerrado com sucesso")
}
