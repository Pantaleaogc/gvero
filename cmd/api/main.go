mkdir -p cmd/api
cat > cmd/api/main.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/joho/godotenv"
)

func main() {
    // Carregar configurações
    if err := godotenv.Load("configs/.env"); err != nil {
        log.Printf("Aviso: arquivo .env não encontrado: %v", err)
    }

    // Configurações da porta
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Valor padrão
    }

    // Handler básico para teste
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Sistema CRM/ERP em Go - API funcionando! Versão 0.1.0")
    })

    // Iniciar servidor
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
EOF
