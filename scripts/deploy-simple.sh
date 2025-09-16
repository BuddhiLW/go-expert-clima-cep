#!/bin/bash

# Script simplificado para deploy no Cloud Run
# Usa uma imagem já construída e faz deploy direto

set -e

PROJECT_ID=${1:-"focus-skein-364415"}
REGION="southamerica-east1"
SERVICE_NAME="cep-temperatura"

echo "🚀 Deploy Simplificado no Google Cloud Run"
echo "📋 Projeto: $PROJECT_ID"
echo "🌍 Região: $REGION"

# Carregar variáveis de ambiente
if [ -f ".env" ]; then
    echo "📖 Carregando variáveis de ambiente do arquivo .env..."
    export $(grep -v '^#' .env | xargs)
else
    echo "⚠️  Arquivo .env não encontrado, usando valores padrão"
    export WEATHER_API_KEY="your_weather_api_key_here"
fi

echo "✅ Variáveis carregadas:"
echo "   WEATHER_API_KEY: ${WEATHER_API_KEY:0:8}..."

# Configurar projeto
echo "⚙️  Configurando projeto..."
/home/ramanujan/google-cloud-sdk/bin/gcloud config set project $PROJECT_ID

# Habilitar APIs necessárias
echo "🔧 Habilitando APIs necessárias..."
/home/ramanujan/google-cloud-sdk/bin/gcloud services enable run.googleapis.com cloudbuild.googleapis.com

# Fazer deploy usando uma imagem base do Cloud Run
echo "🚀 Fazendo deploy no Cloud Run usando imagem base..."

# Usar uma imagem Go simples do Cloud Run
/home/ramanujan/google-cloud-sdk/bin/gcloud run deploy $SERVICE_NAME \
  --image gcr.io/cloudrun/hello \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=$WEATHER_API_KEY \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 10 \
  --port 8080 \
  --source . \
  --execution-environment gen2

echo ""
echo "✅ Deploy concluído!"
echo "🌐 URL do serviço:"
/home/ramanujan/google-cloud-sdk/bin/gcloud run services describe $SERVICE_NAME --platform managed --region $REGION --format 'value(status.url)'
