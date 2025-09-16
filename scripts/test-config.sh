#!/bin/bash

# Script para testar diferentes configurações
# Uso: ./test-config.sh [dev|prod]

ENV=${1:-"dev"}

echo "🧪 Testando configuração: $ENV"
echo ""

# Definir variáveis de ambiente baseadas no ambiente
if [ "$ENV" = "prod" ]; then
    export WEATHER_API_KEY="your_weather_api_key_here"
    export PORT="8080"
    export HOST="0.0.0.0"
    echo "📋 Usando configuração de produção"
else
    export WEATHER_API_KEY="your_weather_api_key_here"
    export PORT="8080"
    export HOST="localhost"
    echo "📋 Usando configuração de desenvolvimento"
fi

echo ""

# Compilar aplicação
echo "🏗️  Compilando aplicação..."
go build -o main cmd/main.go

if [ $? -ne 0 ]; then
    echo "❌ Erro na compilação"
    exit 1
fi

echo "✅ Compilação bem-sucedida"
echo ""

# Executar aplicação em background
echo "🚀 Iniciando aplicação..."
./main &
APP_PID=$!

# Aguardar aplicação iniciar
sleep 3

# Testar endpoints
echo "🧪 Testando endpoints..."
echo ""

# Health check
echo "Testando health check..."
curl -s http://localhost:8080/health | jq . || echo "Resposta: $(curl -s http://localhost:8080/health)"
echo ""

# Teste de CEP válido
echo "Testando CEP válido (01310-100)..."
curl -s http://localhost:8080/temperature/01310100 | jq . || echo "Resposta: $(curl -s http://localhost:8080/temperature/01310100)"
echo ""

# Teste de CEP inválido
echo "Testando CEP inválido (123)..."
curl -s http://localhost:8080/temperature/123 | jq . || echo "Resposta: $(curl -s http://localhost:8080/temperature/123)"
echo ""

# Parar aplicação
echo "🛑 Parando aplicação..."
kill $APP_PID

echo ""
echo "✅ Teste de configuração concluído!"
