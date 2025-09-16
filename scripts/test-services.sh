#!/bin/bash

# Script para testar os serviços

set -e

echo "🧪 Testando serviços..."

# Aguardar serviços ficarem prontos
echo "⏳ Aguardando serviços ficarem prontos..."
sleep 10

# Testar health checks
echo "🔍 Testando health checks..."

echo "Testando Serviço A - Health básico..."
curl -s http://localhost:8080/health | jq .

echo -e "\nTestando Serviço A - Health detalhado..."
curl -s http://localhost:8080/health/detailed | jq .

echo -e "\nTestando Serviço A - Readiness..."
curl -s http://localhost:8080/ready | jq .

echo -e "\nTestando Serviço A - Liveness..."
curl -s http://localhost:8080/live | jq .

echo -e "\nTestando Serviço B - Health básico..."
curl -s http://localhost:8081/health | jq .

echo -e "\nTestando Serviço B - Health detalhado..."
curl -s http://localhost:8081/health/detailed | jq .

echo -e "\nTestando Serviço B - Readiness..."
curl -s http://localhost:8081/ready | jq .

echo -e "\nTestando Serviço B - Liveness..."
curl -s http://localhost:8081/live | jq .

# Testar CEP válido
echo -e "\n🌡️  Testando CEP válido..."
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep":"01310100"}' | jq .

# Testar CEP inválido
echo -e "\n❌ Testando CEP inválido..."
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep":"1234567"}' | jq .

# Testar CEP não encontrado
echo -e "\n🔍 Testando CEP não encontrado..."
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep":"99999999"}' | jq .

echo -e "\n✅ Testes concluídos!"
echo "📊 Verifique os traces em: http://localhost:9411"
