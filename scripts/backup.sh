#!/bin/bash
# Script para backup do banco de dados MySQL/MariaDB

# Carregar variáveis de ambiente
source ../configs/.env

# Diretório para armazenar backups
BACKUP_DIR="../backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="$BACKUP_DIR/db_backup_$TIMESTAMP.sql"

# Criar diretório de backups se não existir
mkdir -p $BACKUP_DIR

# Log function
log_message() {
  echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1" >> "$BACKUP_DIR/backup_log.txt"
  echo "$1"
}

log_message "Iniciando backup do banco de dados: $DB_NAME"

# Executar backup com mysqldump
mysqldump --host=$DB_HOST --port=$DB_PORT --user=$DB_USER --password=$DB_PASS $DB_NAME > $BACKUP_FILE

# Verificar se o backup foi bem-sucedido
if [ $? -eq 0 ]; then
  # Comprimir o arquivo
  gzip $BACKUP_FILE
  log_message "Backup concluído com sucesso: ${BACKUP_FILE}.gz"
  
  # Remover backups antigos (manter os últimos 7 dias)
  find $BACKUP_DIR -name "db_backup_*.sql.gz" -type f -mtime +7 -delete
  log_message "Backups mais antigos que 7 dias foram removidos"
else
  log_message "ERRO: Falha ao realizar o backup do banco de dados"
  exit 1
fi

exit 0
