## Architecture Diagram

![alt text](https://github.com/dedihartono801/LlamaTalks/blob/master/simple-architecture-AI-BE.png)

## Description

# 🦙 tinyllama-stream-chat

This repository demonstrates a full integration of a simple AI chatbot using:

- **Frontend** built with HTML, CSS, and JavaScript
- **Backend** built with Go
- **AI Model** using TinyLLaMA via Ollama
- **Real-time response streaming** using Server-Sent Events (SSE)

---

## 🚀 How to Install TinyLLaMA using OLLAMA
- brew install ollama 
*(Note: If you are not a macOS user, please find out for yourself how to install it on another OS hehe)*
## To start ollama:
- brew services start ollama
## Or, if you don't want/need a background service you can just run:
- /opt/homebrew/opt/ollama/bin/ollama serve
## Run TinyLLaMA
- ollama run tinyllama

---

## 🦙 What is TinyLLaMA?

TinyLLaMA is a compact, efficient large language model designed to run locally with minimal resource requirements. It’s ideal for personal and offline AI use cases.

- Lightweight: Can run on CPUs or small GPUs
- Open-source and fast
- Perfect for developers experimenting with LLMs locally

## Model Details:
- Name: TinyLLaMA-1.1B
- Parameters: ~1.1 billion
- Architecture: Based on LLaMA (Large Language Model Meta AI), but scaled down
- Training Dataset: A subset of publicly available text corpora (~3 trillion tokens)
- Vocabulary Size: ~32,000 tokens
- Precision: Most models are released in quantized formats like Q4_0 for fast inference with smaller memory footprint

---

## What is Server-Sent Events (SSE)?

Server-Sent Events (SSE) let the server send real-time updates to the browser over one long HTTP connection.

SSE is:

- Easy to use and supported by modern browsers
- Good for one-way data flow from server to client
- Lightweight and efficient for streaming data

### How SSE works in this project

- The backend gets a streaming response from TinyLLaMA.
- It sends each token to the frontend using SSE.
- The browser shows the tokens as they arrive, creating a smooth typing effect.