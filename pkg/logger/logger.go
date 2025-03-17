package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
	logFile     *os.File
)

// Init inicializa os loggers
func Init() {
	// Criar diretório de logs se não existir
	logsDir := "logs"
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		os.Mkdir(logsDir, 0755)
	}

	// Nome do arquivo de log com data
	logFileName := filepath.Join(logsDir, fmt.Sprintf("app_%s.log", time.Now().Format("2006-01-02")))
	
	// Abrir arquivo de log
	var err error
	logFile, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Erro ao abrir arquivo de log:", err)
		// Fallback para saída padrão se não puder abrir o arquivo
		InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		DebugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
		return
	}

	// Criar loggers com output para o arquivo e para o console
	multiWriter := log.MultiWriter(os.Stdout, logFile)
	errorWriter := log.MultiWriter(os.Stderr, logFile)

	InfoLogger = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	
	// Debug logger só vai para o console em ambiente de desenvolvimento
	if os.Getenv("ENV") == "development" {
		DebugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		// Em produção, debug vai apenas para o arquivo
		DebugLogger = log.New(logFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

// Close fecha o arquivo de log
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

// RotateLog realiza a rotação do arquivo de log
func RotateLog() error {
	// Fechar o arquivo atual
	if logFile != nil {
		logFile.Close()
	}

	// Iniciar um novo arquivo
	Init()
	return nil
}
