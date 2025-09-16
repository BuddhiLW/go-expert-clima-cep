#!/bin/bash

# Script para configurar ambiente de desenvolvimento

set -e

echo "🚀 Configurando ambiente de desenvolvimento..."

# Verificar se o .env existe
if [ ! -f ".env" ]; then
    echo "📝 Criando arquivo .env..."
    cp env.example .env
    echo "⚠️  Configure sua WEATHER_API_KEY no arquivo .env"
else
    echo "✅ Arquivo .env já existe"
fi

# Verificar se o Docker está rodando
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker não está rodando. Inicie o Docker e tente novamente."
    exit 1
fi

# Verificar se o docker-compose está disponível
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose não está instalado. Instale e tente novamente."
    exit 1
fi

echo "✅ Ambiente configurado com sucesso!"
echo ""
echo "Para iniciar os serviços:"
echo "  docker-compose up --build"
echo ""
echo "Para testar:"
echo "  curl -X POST http://localhost:8080/cep -H 'Content-Type: application/json' -d '{\"cep\":\"01310100\"}'"
echo ""
echo "Para ver os traces:"
echo "  http://localhost:9411"
