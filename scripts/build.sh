#!/bin/bash

# Script para build da imagem Docker com variáveis de ambiente
# Uso: ./build.sh [tag] [env_file]

set -e

TAG=${1:-"cep-temperatura:latest"}
ENV_FILE=${2:-".env"}

echo "🐳 Construindo imagem Docker"
echo "📋 Tag: $TAG"
echo "📄 Arquivo de ambiente: $ENV_FILE"
echo ""

# Verificar se o arquivo .env existe
if [ ! -f "$ENV_FILE" ]; then
    echo "❌ Arquivo $ENV_FILE não encontrado"
    echo "💡 Execute: ./setup-env.sh"
    exit 1
fi

# Carregar variáveis do arquivo .env
echo "📖 Carregando variáveis de ambiente..."
export $(grep -v '^#' "$ENV_FILE" | xargs)

# Verificar se as variáveis necessárias estão definidas
if [ -z "$WEATHER_API_KEY" ]; then
    echo "❌ WEATHER_API_KEY não definida no arquivo $ENV_FILE"
    exit 1
fi

echo "✅ Variáveis carregadas:"
echo "   PORT: ${PORT:-8080}"
echo "   HOST: ${HOST:-0.0.0.0}"
echo "   WEATHER_API_KEY: ${WEATHER_API_KEY:0:8}..."
echo ""

# Build da imagem
echo "🏗️  Construindo imagem..."
sudo docker build \
  --build-arg PORT="${PORT:-8080}" \
  --build-arg HOST="${HOST:-0.0.0.0}" \
  -t "$TAG" .

if [ $? -eq 0 ]; then
    echo "✅ Imagem construída com sucesso: $TAG"
    echo ""
    echo "Para executar localmente:"
    echo "  sudo docker run -p 8080:8080 -e WEATHER_API_KEY=\"$WEATHER_API_KEY\" $TAG"
    echo ""
    echo "Para testar:"
    echo "  curl http://localhost:8080/health"
else
    echo "❌ Erro na construção da imagem"
    exit 1
fi
