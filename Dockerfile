FROM python:3.11-slim

RUN apt-get update && apt-get install -y git build-essential \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /comfyui
RUN git clone https://github.com/comfyanonymous/ComfyUI.git . \
    && pip install --no-cache-dir torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cpu \
    && pip install --no-cache-dir -r requirements.txt

WORKDIR /comfyui/custom_nodes
RUN git clone https://github.com/city96/ComfyUI-GGUF.git \
    && pip install --no-cache-dir -r ComfyUI-GGUF/requirements.txt

WORKDIR /comfyui
CMD python main.py $CLI_ARGS
