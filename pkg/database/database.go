package database

import (
    "database/sql"
    "fmt"
    "os"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

// DB é a conexão global com o banco de dados
var DB *sql.DB

// InitDB inicializa a conexão com o banco de dados
func InitDB() (*sql.DB, error) {
    dbDriver := os.Getenv("DB_DRIVER")
    if dbDriver == "" {
        dbDriver = "mysql" // Valor padrão
    }

    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASS")
    dbName := os.Getenv("DB_NAME")

    var dsn string
    if dbDriver == "mysql" {
        dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
            dbUser, dbPass, dbHost, dbPort, dbName)
    } else {
        return nil, fmt.Errorf("driver de banco de dados não suportado: %s", dbDriver)
    }

    db, err := sql.Open(dbDriver, dsn)
    if err != nil {
        return nil, err
    }

    // Configuração do pool de conexões
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(time.Minute * 5)

    // Verificar conexão
    if err = db.Ping(); err != nil {
        return nil, err
    }

    DB = db
    return db, nil
}

// GetDB retorna a conexão com o banco de dados
func GetDB() *sql.DB {
    return DB
}
