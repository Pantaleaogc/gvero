#!/bin/bash
if [ -f app.pid ]; then
    kill $(cat app.pid) 2>/dev/null
    rm app.pid
    echo "Aplicativo parado"
else
    echo "Aplicativo não está em execução"
fi
