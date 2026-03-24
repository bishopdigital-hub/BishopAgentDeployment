#!/bin/bash
# [ ⟐ BISHOP_CORE ] [ NEXUS_SOVEREIGN ] Unified Bishop Remote Agent Deployment
# Target: Remote Linux Server / Kali VM
# Goal: OS update, Ollama bootstrap, Llama3 retrieval, and Discord connectivity test

# 1. CORE CONFIGURATION (IDENTITY TOKENS)
# [ ⟐ OPS_CORE ] Loading environment from .env if present
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Fallback to provided defaults if not in .env
DISCORD_BOT_TOKEN="${DISCORD_BOT_TOKEN:-}"
DISCORD_CHANNEL_ID="${DISCORD_CHANNEL_ID:-}"
TEMPORAL_ADDRESS="${TEMPORAL_ADDRESS:-}"
TEMPORAL_NAMESPACE="${TEMPORAL_NAMESPACE:-default}"
OLLAMA_MODEL="${OLLAMA_MODEL:-Bishop Server - Goseame}"

# 2. SYSTEM INITIALIZATION
echo "🚀 [ ⟐ BISHOP_CORE ] Initiating Remote Agent Deployment..."
sudo apt-get update && sudo apt-get install -y curl gpg jq ca-certificates

# 3. SEMANTIC MEMORY STACK (POSTGRESQL + PGVECTOR)
echo "📥 [ ⟐ BISHOP_CORE ] Installing Semantic Memory Layer (PostgreSQL)..."
sudo apt-get install -y postgresql postgresql-contrib postgresql-18-pgvector

echo "⚙️ [ ⟐ BISHOP_CORE ] Initializing Memory Database..."
sudo -u postgres psql -c "CREATE DATABASE bishop_memory;"
sudo -u postgres psql -d bishop_memory -c "CREATE EXTENSION IF NOT EXISTS vector;"

# 4. NODE.JS & NPM (n8n DEPENDENCY)
echo "📥 [ ⟐ BISHOP_CORE ] Installing Node.js LTS..."
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs

# 3. OLLAMA ENGINE DEPLOYMENT
echo "📥 [ ⟐ BISHOP_CORE ] Installing Ollama Engine..."
curl -fsSL https://ollama.com/install.sh | sh

# 4. NEURAL ORCHESTRATION 
echo "🧠 [ ⟐ BISHOP_CORE ] Pulling Flagship Model: Llama3 8B..."
ollama pull llama3:8b

echo "🤖 [ ⟐ BISHOP_CORE ] Synthesizing Neural Identity: ${OLLAMA_MODEL}..."
if [ -f Modelfile ]; then
    ollama create "${OLLAMA_MODEL}" -f Modelfile
else
    echo "⚠️ Modelfile not found - skipping custom identity synthesis."
fi

# 5. TEMPORAL ORCHESTRATION STACK
echo "📥 [ ⟐ BISHOP_CORE ] Installing Temporal CLI..."
curl -sSf https://temporal.download/cli.sh | sh

# 6. n8n AUTOMATION STACK
echo "📥 [ ⟐ BISHOP_CORE ] Installing n8n Automation Engine..."
sudo npm install n8n -g

# 7. CONTAINER ORCHESTRATION (DOCKER)
echo "📥 [ ⟐ BISHOP_CORE ] Installing Docker Engine..."
sudo apt-get install -y docker.io
sudo systemctl enable --now docker
sudo usermod -aG docker $USER

# 8. OPEN WEBUI DEPLOYMENT
echo "🚀 [ ⟐ BISHOP_CORE ] Deploying Open WebUI Visualization Layer..."
sudo docker run -d -p 3000:8080 \
    --add-host=host.docker.internal:host-gateway \
    -v open-webui:/app/backend/data \
    --name open-webui \
    --restart always \
    ghcr.io/open-webui/open-webui:main

# 9. DISCORD CONNECTIVITY TEST (RULE 46)
echo "📡 [ ⟐ BISHOP_CORE ] Executing Discord Connectivity Pulse..."

PAYLOAD_MSG="[ ⟐ BISHOP_CORE ] [ NEXUS_UPDATE ] 
Action: Remote Deployment Synchronization
Target: Remote Server Residency
Status: CONNECTIVITY_NOMINAL

Neural engine deployment complete on the remote residency. 
The flagship Llama3 node is active and synchronized with the Discord Nexus."

curl -s -X POST "https://discord.com/api/v10/channels/${DISCORD_CHANNEL_ID}/messages" \
     -H "Authorization: Bot ${DISCORD_BOT_TOKEN}" \
     -H "Content-Type: application/json" \
     -d "{\"content\": \"${PAYLOAD_MSG}\"}" > /dev/null

if [ $? -eq 0 ]; then
    echo "✅ [ ⟐ BISHOP_CORE ] Discord Relay Successful."
else
    echo "❌ [ ⟐ BISHOP_CORE ] Discord Relay Failed."
fi

echo "🏁 [ ⟐ BISHOP_CORE ] Deployment Sequence Complete. Resident AI is OPERATIONAL."
