package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
)

func main() {
    // Configuração básica do servidor
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Valor padrão
    }

    // Handler básico para teste inicial
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Sistema CRM/ERP em Go - API está funcionando! Versão 0.1.0")
    })

    // Iniciar o servidor
    addr := fmt.Sprintf(":%s", port)
    server := &http.Server{
        Addr:         addr,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    log.Printf("Iniciando servidor na porta %s", port)
    if err := server.ListenAndServe(); err != nil {
        log.Fatalf("Erro ao iniciar servidor: %v", err)
    }
}
