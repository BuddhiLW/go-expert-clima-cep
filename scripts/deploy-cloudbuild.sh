#!/bin/bash

# Script para deploy no Google Cloud Run usando Cloud Build
# Uso: ./deploy-cloudbuild.sh [PROJECT_ID]

set -e

PROJECT_ID=${1:-"your-project-id"}
SERVICE_NAME="cep-temperatura"
REGION="southamerica-east1"

echo "🚀 Deploy no Google Cloud Run (Cloud Build)"
echo "📋 Projeto: $PROJECT_ID"
echo "🌍 Região: $REGION"
echo ""

# Verificar se o arquivo .env existe
if [ ! -f ".env" ]; then
    echo "❌ Arquivo .env não encontrado"
    echo "💡 Execute: ./setup-env.sh"
    exit 1
fi

# Carregar variáveis de ambiente
echo "📖 Carregando variáveis de ambiente..."
export $(grep -v '^#' .env | xargs)

# Verificar se as variáveis necessárias estão definidas
if [ -z "$WEATHER_API_KEY" ]; then
    echo "❌ WEATHER_API_KEY não definida no arquivo .env"
    exit 1
fi

echo "✅ Variáveis carregadas:"
echo "   PORT: ${PORT:-8080}"
echo "   HOST: ${HOST:-0.0.0.0}"
echo "   WEATHER_API_KEY: ${WEATHER_API_KEY:0:8}..."
echo ""

# Verificar se o gcloud está instalado
if ! command -v gcloud &> /dev/null; then
    echo "❌ Google Cloud SDK não encontrado. Instale em: https://cloud.google.com/sdk/docs/install"
    exit 1
fi

# Verificar se está autenticado
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
    echo "🔐 Fazendo login no Google Cloud..."
    gcloud auth login
fi

# Configurar projeto
echo "⚙️  Configurando projeto..."
gcloud config set project $PROJECT_ID

# Habilitar APIs necessárias
echo "🔧 Habilitando APIs necessárias..."
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable containerregistry.googleapis.com

# Deploy usando Cloud Build
echo "🏗️  Fazendo build e deploy com Cloud Build..."
gcloud builds submit \
  --config cloudbuild.yaml \
  --substitutions=_WEATHER_API_KEY="$WEATHER_API_KEY",_PORT="${PORT:-8080}",_HOST="${HOST:-0.0.0.0}"

# Obter URL do serviço
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --platform managed --region $REGION --format 'value(status.url)')

echo ""
echo "✅ Deploy concluído com sucesso!"
echo "🌐 URL do serviço: $SERVICE_URL"
echo ""
echo "🧪 Testando o serviço..."
echo ""

# Testar o serviço
echo "Testando health check..."
curl -s "$SERVICE_URL/health" | jq . || echo "Resposta: $(curl -s $SERVICE_URL/health)"

echo ""
echo "Testando CEP válido (01310-100)..."
curl -s "$SERVICE_URL/temperature/01310100" | jq . || echo "Resposta: $(curl -s $SERVICE_URL/temperature/01310100)"

echo ""
echo "🎉 Deploy finalizado! Acesse: $SERVICE_URL"