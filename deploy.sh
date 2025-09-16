#!/bin/bash

# Script para deploy no Google Cloud Run
# Uso: ./deploy.sh [PROJECT_ID]

set -e

PROJECT_ID=${1:-"your-project-id"}
SERVICE_NAME="cep-temperatura"
REGION="southamerica-east1"
WEATHER_API_KEY=${WEATHER_API_KEY:-"b5d4..."}

echo "🚀 Iniciando deploy no Google Cloud Run"
echo "📋 Projeto: $PROJECT_ID"
echo "🌍 Região: $REGION"
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

# Carregar variáveis de ambiente
if [ -f ".env" ]; then
    echo "📖 Carregando variáveis de ambiente do arquivo .env..."
    export $(grep -v '^#' .env | xargs)
else
    echo "⚠️  Arquivo .env não encontrado, usando valores padrão"
fi

# Build da imagem com variáveis de ambiente
echo "🏗️  Fazendo build da imagem Docker..."
sudo docker build \
  --build-arg WEATHER_API_KEY="$WEATHER_API_KEY" \
  --build-arg PORT="${PORT:-8080}" \
  --build-arg HOST="${HOST:-0.0.0.0}" \
  -t gcr.io/$PROJECT_ID/$SERVICE_NAME:latest .

# Push da imagem
echo "📤 Enviando imagem para Google Container Registry..."
sudo docker push gcr.io/$PROJECT_ID/$SERVICE_NAME:latest

# Deploy no Cloud Run
echo "🚀 Fazendo deploy no Cloud Run..."
gcloud run deploy $SERVICE_NAME \
  --image gcr.io/$PROJECT_ID/$SERVICE_NAME:latest \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=$WEATHER_API_KEY,PORT=8080,HOST=0.0.0.0 \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 10 \
  --port 8080

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
