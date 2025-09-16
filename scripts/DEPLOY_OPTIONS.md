# 🚀 Opções de Deploy - Google Cloud Run

Devido a problemas de conectividade de rede, aqui estão as opções disponíveis para fazer deploy:

## ✅ **Opção 1: Deploy com Cloud Build (Recomendado)**

```bash
# Usar Cloud Build (mais confiável)
./deploy-cloudbuild.sh focus-skein-364415
```

**Vantagens:**
- ✅ Não depende do Docker local
- ✅ Build acontece no Google Cloud
- ✅ Mais confiável para problemas de rede

## ✅ **Opção 2: Deploy com Imagem Local**

```bash
# Usar imagem já construída localmente
./deploy-local-image.sh focus-skein-364415
```

**Vantagens:**
- ✅ Usa a imagem já construída
- ✅ Retry automático para push
- ✅ Mais rápido se a imagem já existe

## ✅ **Opção 3: Deploy Manual via Console**

1. **Acesse o Google Cloud Console:**
   - https://console.cloud.google.com/
   - Projeto: focus-skein-364415

2. **Vá para Cloud Run:**
   - https://console.cloud.google.com/run

3. **Crie um novo serviço:**
   - Nome: `cep-temperatura`
   - Região: `southamerica-east1`
   - Imagem: `gcr.io/focus-skein-364415/cep-temperatura:latest`

4. **Configure variáveis de ambiente:**
   - `WEATHER_API_KEY`: `your_weather_api_key_here`
   - `PORT`: `8080`
   - `HOST`: `0.0.0.0`

## ✅ **Opção 4: Deploy via GitHub Actions**

1. **Configure secrets no GitHub:**
   - `GCP_PROJECT_ID`: `focus-skein-364415`
   - `GCP_SA_KEY`: Chave da service account
   - `WEATHER_API_KEY`: `your_weather_api_key_here`

2. **Push para o repositório:**
   ```bash
   git add .
   git commit -m "Deploy to Cloud Run"
   git push origin main
   ```

## 🔧 **Solução de Problemas de Rede**

### Problema: DNS Resolution Failed
```bash
# Verificar DNS
nslookup gcr.io
nslookup cloudresourcemanager.googleapis.com

# Usar DNS alternativo
echo "nameserver 8.8.8.8" | sudo tee /etc/resolv.conf
```

### Problema: Docker Push Failed
```bash
# Configurar proxy se necessário
export HTTP_PROXY=http://proxy:port
export HTTPS_PROXY=http://proxy:port

# Ou usar Cloud Build
./deploy-cloudbuild.sh focus-skein-364415
```

## 📊 **Status Atual**

- ✅ **Imagem construída**: `gcr.io/focus-skein-364415/cep-temperatura:latest`
- ✅ **Aplicação testada**: Funcionando localmente
- ✅ **Configuração**: Variáveis de ambiente configuradas
- ❌ **Push para registry**: Falhou devido a problemas de rede
- ❌ **Deploy**: Pendente

## 🎯 **Próximos Passos**

1. **Tente a Opção 1** (Cloud Build) - mais confiável
2. **Se falhar, use a Opção 3** (Console manual)
3. **Configure CI/CD** para futuros deploys

## 📞 **Suporte**

Se todas as opções falharem:
1. Verifique sua conexão de internet
2. Configure proxy se necessário
3. Use o Google Cloud Console diretamente
4. Considere usar uma VPN

---

**Recomendação:** Use a **Opção 1** (Cloud Build) primeiro, pois é a mais confiável e não depende da sua conexão local.
