#!/bin/bash

# Script para executar em ambiente de desenvolvimento

set -e

echo "🚀 Iniciando ambiente de desenvolvimento..."

# Verificar se .env existe
if [ ! -f ".env" ]; then
    echo "📝 Criando arquivo .env..."
    cp env.example .env
    echo "⚠️  Configure sua WEATHER_API_KEY no arquivo .env antes de continuar"
    echo "   Edite o arquivo .env e execute novamente este script"
    exit 1
fi

# Carregar variáveis de ambiente
export $(grep -v '^#' .env | xargs)

# Verificar se WEATHER_API_KEY está configurada
if [ "$WEATHER_API_KEY" = "your_weather_api_key_here" ] || [ -z "$WEATHER_API_KEY" ]; then
    echo "❌ WEATHER_API_KEY não configurada no arquivo .env"
    echo "   Edite o arquivo .env e configure sua chave da WeatherAPI"
    exit 1
fi

echo "✅ Configuração validada"

# Iniciar Zipkin
echo "🔍 Iniciando Zipkin..."
docker run -d --name zipkin -p 9411:9411 openzipkin/zipkin > /dev/null 2>&1 || echo "Zipkin já está rodando"

# Aguardar Zipkin ficar pronto
echo "⏳ Aguardando Zipkin ficar pronto..."
sleep 5

# Iniciar Serviço B
echo "🌡️  Iniciando Serviço B..."
go run cmd/service-b/main.go &
SERVICE_B_PID=$!

# Aguardar Serviço B ficar pronto
echo "⏳ Aguardando Serviço B ficar pronto..."
sleep 3

# Iniciar Serviço A
echo "📝 Iniciando Serviço A..."
go run cmd/service-a/main.go &
SERVICE_A_PID=$!

# Aguardar Serviço A ficar pronto
echo "⏳ Aguardando Serviço A ficar pronto..."
sleep 3

echo ""
echo "✅ Todos os serviços iniciados!"
echo ""
echo "🔗 URLs disponíveis:"
echo "   Serviço A: http://localhost:8080"
echo "   Serviço B: http://localhost:8081"
echo "   Zipkin UI: http://localhost:9411"
echo ""
echo "🧪 Para testar:"
echo "   curl -X POST http://localhost:8080/cep -H 'Content-Type: application/json' -d '{\"cep\":\"01310100\"}'"
echo ""
echo "🛑 Para parar os serviços:"
echo "   kill $SERVICE_A_PID $SERVICE_B_PID"
echo "   docker stop zipkin && docker rm zipkin"
echo ""

# Manter script rodando
wait
