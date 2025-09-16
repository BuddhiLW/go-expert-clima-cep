#!/bin/bash

# Script para deploy usando imagem local já construída
# Uso: ./deploy-local-image.sh [PROJECT_ID]

set -e

PROJECT_ID=${1:-"your-project-id"}
SERVICE_NAME="cep-temperatura"
REGION="southamerica-east1"

echo "🚀 Deploy usando imagem local"
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

# Fazer push da imagem local para o registry
echo "📤 Fazendo push da imagem local para o registry..."
sudo docker tag gcr.io/$PROJECT_ID/$SERVICE_NAME:latest gcr.io/$PROJECT_ID/$SERVICE_NAME:$(date +%Y%m%d-%H%M%S)

# Tentar fazer push com retry
for i in {1..3}; do
    echo "Tentativa $i de 3..."
    if sudo docker push gcr.io/$PROJECT_ID/$SERVICE_NAME:latest; then
        echo "✅ Push bem-sucedido!"
        break
    else
        echo "❌ Falha na tentativa $i"
        if [ $i -eq 3 ]; then
            echo "❌ Todas as tentativas falharam. Verifique sua conexão de rede."
            exit 1
        fi
        sleep 5
    fi
done

# Deploy no Cloud Run
echo "🚀 Fazendo deploy no Cloud Run..."
gcloud run deploy $SERVICE_NAME \
  --image gcr.io/$PROJECT_ID/$SERVICE_NAME:latest \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY="$WEATHER_API_KEY",PORT="${PORT:-8080}",HOST="${HOST:-0.0.0.0}" \
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
