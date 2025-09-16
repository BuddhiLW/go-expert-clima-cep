#!/bin/bash

# Script para configurar o ambiente de desenvolvimento
# Uso: ./setup-env.sh

echo "🔧 Configurando ambiente de desenvolvimento"
echo ""

# Verificar se o arquivo .env já existe
if [ -f ".env" ]; then
    echo "⚠️  Arquivo .env já existe"
    read -p "Deseja sobrescrever? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ Operação cancelada"
        exit 1
    fi
fi

# Criar arquivo .env
echo "📝 Criando arquivo .env..."
cat > .env << EOF
# Weather API Configuration
WEATHER_API_KEY=your_weather_api_key_here

# Server Configuration
PORT=8080
HOST=0.0.0.0
EOF

echo "✅ Arquivo .env criado com sucesso"
echo ""

# Verificar se as dependências estão instaladas
echo "📦 Verificando dependências..."
go mod tidy

if [ $? -eq 0 ]; then
    echo "✅ Dependências atualizadas"
else
    echo "❌ Erro ao atualizar dependências"
    exit 1
fi

echo ""

# Compilar aplicação
echo "🏗️  Compilando aplicação..."
go build -o main cmd/main.go

if [ $? -eq 0 ]; then
    echo "✅ Aplicação compilada com sucesso"
else
    echo "❌ Erro na compilação"
    exit 1
fi

echo ""
echo "🎉 Ambiente configurado com sucesso!"
echo ""
echo "Para executar a aplicação:"
echo "  ./main"
echo ""
echo "Para executar com Docker:"
echo "  docker-compose up"
echo ""
echo "Para testar a API:"
echo "  ./test-api.sh"
