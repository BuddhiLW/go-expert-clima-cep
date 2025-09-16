#!/bin/bash

# Script para testar a imagem Docker localmente
# Uso: ./test-docker.sh [tag]

set -e

TAG=${1:-"cep-temperatura:latest"}
ENV_FILE=${2:-".env"}

echo "🧪 Testando imagem Docker"
echo "📋 Tag: $TAG"
echo "📄 Arquivo de ambiente: $ENV_FILE"
echo ""

# Verificar se o arquivo .env existe
if [ ! -f "$ENV_FILE" ]; then
    echo "❌ Arquivo $ENV_FILE não encontrado"
    echo "💡 Execute: ./setup-env.sh"
    exit 1
fi

# Carregar variáveis de ambiente
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
docker build \
  --build-arg WEATHER_API_KEY="$WEATHER_API_KEY" \
  --build-arg PORT="${PORT:-8080}" \
  --build-arg HOST="${HOST:-0.0.0.0}" \
  -t "$TAG" .

if [ $? -ne 0 ]; then
    echo "❌ Erro na construção da imagem"
    exit 1
fi

echo "✅ Imagem construída com sucesso"
echo ""

# Executar container
echo "🚀 Executando container..."
docker run -d --name cep-temperatura-test -p 8080:8080 "$TAG"

if [ $? -ne 0 ]; then
    echo "❌ Erro ao executar container"
    exit 1
fi

echo "✅ Container executando"
echo ""

# Aguardar aplicação iniciar
echo "⏳ Aguardando aplicação iniciar..."
sleep 5

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

# Parar e remover container
echo "🛑 Parando e removendo container..."
docker stop cep-temperatura-test
docker rm cep-temperatura-test

echo ""
echo "✅ Teste concluído com sucesso!"
echo "🎉 Imagem pronta para deploy!"
