# 🚀 Guia de Deploy - Google Cloud Run

Este guia explica como fazer o deploy da aplicação CEP Temperatura no Google Cloud Run.

## 📋 Pré-requisitos

1. **Google Cloud Account** com billing habilitado
2. **Google Cloud SDK** instalado
3. **Docker** instalado
4. **Projeto no Google Cloud** criado

## 🔧 Configuração Inicial

### 1. Instalar Google Cloud SDK

```bash
# Ubuntu/Debian
curl https://sdk.cloud.google.com | bash
exec -l $SHELL

# macOS
brew install google-cloud-sdk

# Windows
# Baixe o instalador em: https://cloud.google.com/sdk/docs/install
```

### 2. Configurar Autenticação

```bash
# Login
gcloud auth login

# Configurar projeto
gcloud config set project SEU_PROJECT_ID

# Verificar configuração
gcloud config list
```

### 3. Habilitar APIs Necessárias

```bash
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable containerregistry.googleapis.com
```

## 🚀 Deploy Automático

### Opção 1: Script de Deploy

```bash
# Tornar o script executável
chmod +x deploy.sh

# Executar deploy
./deploy.sh SEU_PROJECT_ID
```

### Opção 2: Cloud Build

```bash
# Build e deploy com Cloud Build
gcloud builds submit --config cloudbuild.yaml
```

### Opção 3: Deploy Manual

```bash
# 1. Build da imagem
docker build -t gcr.io/SEU_PROJECT_ID/cep-temperatura .

# 2. Push para Container Registry
docker push gcr.io/SEU_PROJECT_ID/cep-temperatura

# 3. Deploy no Cloud Run
gcloud run deploy cep-temperatura \
  --image gcr.io/SEU_PROJECT_ID/cep-temperatura \
  --platform managed \
  --region southamerica-east1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=b5d4215a52bf4e2da2f144209251609
```

## 🔄 Deploy com GitHub Actions

### 1. Configurar Secrets

No GitHub, vá para Settings > Secrets and variables > Actions e adicione:

- `GCP_PROJECT_ID`: ID do seu projeto no Google Cloud
- `GCP_SA_KEY`: Chave da service account (JSON)
- `WEATHER_API_KEY`: Chave da WeatherAPI

### 2. Criar Service Account

```bash
# Criar service account
gcloud iam service-accounts create github-actions \
  --display-name="GitHub Actions" \
  --description="Service account for GitHub Actions"

# Dar permissões necessárias
gcloud projects add-iam-policy-binding SEU_PROJECT_ID \
  --member="serviceAccount:github-actions@SEU_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/run.admin"

gcloud projects add-iam-policy-binding SEU_PROJECT_ID \
  --member="serviceAccount:github-actions@SEU_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"

gcloud projects add-iam-policy-binding SEU_PROJECT_ID \
  --member="serviceAccount:github-actions@SEU_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/storage.admin"

# Criar e baixar chave
gcloud iam service-accounts keys create key.json \
  --iam-account=github-actions@SEU_PROJECT_ID.iam.gserviceaccount.com
```

### 3. Configurar Workflow

O arquivo `.github/workflows/deploy.yml` já está configurado. Apenas configure os secrets no GitHub.

## 🧪 Testando o Deploy

### 1. Verificar Status

```bash
gcloud run services describe cep-temperatura \
  --platform managed \
  --region southamerica-east1
```

### 2. Testar Endpoints

```bash
# Obter URL do serviço
SERVICE_URL=$(gcloud run services describe cep-temperatura \
  --platform managed \
  --region southamerica-east1 \
  --format 'value(status.url)')

# Testar health check
curl $SERVICE_URL/health

# Testar CEP
curl $SERVICE_URL/temperature/01310100
```

### 3. Usar Script de Teste

```bash
# Testar com a URL do Cloud Run
./test-api.sh https://cep-temperatura-xxxxx-uc.a.run.app
```

## 📊 Monitoramento

### 1. Logs

```bash
# Ver logs em tempo real
gcloud logs tail --service=cep-temperatura

# Ver logs específicos
gcloud logs read --service=cep-temperatura --limit=50
```

### 2. Métricas

- Acesse o [Google Cloud Console](https://console.cloud.google.com/)
- Vá para Cloud Run > cep-temperatura
- Visualize métricas de CPU, memória, requisições, etc.

### 3. Health Check

```bash
# Verificar saúde do serviço
curl https://sua-url.run.app/health
```

## 🔧 Configurações Avançadas

### 1. Variáveis de Ambiente

```bash
gcloud run services update cep-temperatura \
  --set-env-vars "WEATHER_API_KEY=nova_chave,PORT=8080"
```

### 2. Recursos

```bash
gcloud run services update cep-temperatura \
  --memory 1Gi \
  --cpu 2 \
  --max-instances 20
```

### 3. Domínio Customizado

```bash
# Mapear domínio
gcloud run domain-mappings create \
  --service cep-temperatura \
  --domain api.seudominio.com \
  --region southamerica-east1
```

## 🚨 Troubleshooting

### Problemas Comuns

1. **Erro de permissão**
   ```bash
   gcloud auth login
   gcloud config set project SEU_PROJECT_ID
   ```

2. **Erro de billing**
   - Verifique se o billing está habilitado no projeto

3. **Erro de API não habilitada**
   ```bash
   gcloud services enable run.googleapis.com
   ```

4. **Erro de imagem não encontrada**
   - Verifique se a imagem foi enviada corretamente
   - Verifique as permissões do Container Registry

### Logs de Debug

```bash
# Logs detalhados
gcloud run services describe cep-temperatura \
  --platform managed \
  --region southamerica-east1 \
  --format="export" > service-config.yaml
```

## 💰 Custos

### Estimativa de Custos (Região: southamerica-east1)

- **CPU**: $0.00002400 por vCPU-segundo
- **Memória**: $0.00000250 por GB-segundo
- **Requisições**: $0.40 por milhão de requisições
- **Tráfego de rede**: $0.12 por GB

### Otimizações

1. **Configurar min-instances = 0** para economizar
2. **Usar CPU limitada** quando possível
3. **Configurar timeout** adequado
4. **Monitorar uso** regularmente

## 🔄 Atualizações

### Deploy de Nova Versão

```bash
# Build nova versão
docker build -t gcr.io/SEU_PROJECT_ID/cep-temperatura:v2 .

# Push
docker push gcr.io/SEU_PROJECT_ID/cep-temperatura:v2

# Deploy
gcloud run deploy cep-temperatura \
  --image gcr.io/SEU_PROJECT_ID/cep-temperatura:v2
```

### Rollback

```bash
# Listar revisões
gcloud run revisions list --service=cep-temperatura

# Fazer rollback
gcloud run services update-traffic cep-temperatura \
  --to-revisions=REVISION_NAME=100
```

## 📚 Recursos Adicionais

- [Documentação do Cloud Run](https://cloud.google.com/run/docs)
- [Preços do Cloud Run](https://cloud.google.com/run/pricing)
- [Melhores Práticas](https://cloud.google.com/run/docs/tips)
- [Troubleshooting](https://cloud.google.com/run/docs/troubleshooting)

---

🎉 **Deploy concluído com sucesso!** Sua API está rodando no Google Cloud Run.
