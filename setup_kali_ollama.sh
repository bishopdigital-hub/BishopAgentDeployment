#!/bin/bash
# [ ⟐ BISHOP_CORE ] [ NEXUS_SOVEREIGN ] Ollama Deployment Bootstrap (Kali Linux)
# Target: Kali Linux VM on Hyper-V
# Goal: Automated engine installation and service configuration

echo "🚀 Initiating Ollama deployment sequence..."

# 1. Update and Dependencies
sudo apt-get update
sudo apt-get install -y curl gpg ca-certificates

# 2. Ollama Installation
echo "📥 Downloading and installing Ollama..."
curl -fsSL https://ollama.com/install.sh | sh

# 3. Environment Configuration
echo "⚙️ Configuring Ollama environment..."
# Ensure Ollama is listening on all interfaces for potential cross-VM communication
sudo mkdir -p /etc/systemd/system/ollama.service.d
cat <<EOF | sudo tee /etc/systemd/system/ollama.service.d/override.conf
[Service]
Environment="OLLAMA_HOST=0.0.0.0"
EOF

# 4. Reload and Restart
sudo systemctl daemon-reload
sudo systemctl enable ollama
sudo systemctl restart ollama

# 5. Model Inventory Retrieval
echo "🧠 Pulling model inventory..."
ollama pull llama3:8b
ollama pull mistral

# 6. Verification
echo "✅ Verification Pulse..."
ollama list
ollama run llama3:8b "Neural handshake complete. Bishop Core operational on Kali Linux."

echo "🏁 Ollama deployment sequence FINISHED."
