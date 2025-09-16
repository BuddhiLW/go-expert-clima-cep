#!/bin/bash

# Script para rebuild e restart dos serviços

set -e

echo "🔄 Rebuild e restart dos serviços..."

# Parar containers
echo "🛑 Parando containers..."
docker-compose down

# Rebuild das imagens
echo "🏗️  Fazendo rebuild das imagens..."
docker-compose build --no-cache

# Iniciar containers
echo "🚀 Iniciando containers..."
docker-compose up -d

# Aguardar serviços ficarem prontos
echo "⏳ Aguardando serviços ficarem prontos..."
sleep 10

# Verificar status
echo "📊 Verificando status dos containers..."
docker-compose ps

# Testar health checks
echo "🔍 Testando health checks..."
./scripts/test-head-requests.sh

echo "✅ Rebuild e restart concluídos!"
echo ""
echo "💡 Para ver logs:"
echo "   docker-compose logs -f"
echo ""
echo "💡 Para testar endpoints:"
echo "   ./scripts/test-services.sh"
