#!/bin/bash
# Mata qualquer instância anterior
if [ -f app.pid ]; then
    kill $(cat app.pid) 2>/dev/null
    rm app.pid
fi
# Executa o aplicativo em background
nohup ./crm_erp > app.log 2>&1 &
echo $! > app.pid
echo "Aplicativo iniciado com PID: $(cat app.pid)"
