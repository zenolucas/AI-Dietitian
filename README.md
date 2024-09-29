# My Virtual AI-Dietitian

## About
My Virtual AI-Dietitian is a chatbot powered by Phi3 and Retrieval-Augmented Generation (RAG) to provide contextual responses based on diabetes-related documents. Users can also ask for meal recommendations, and the chatbot will respond with a suggestion, including an image, recipe, and step-by-step preparation instructions.

![Screenshot from 2024-09-29 14-20-32](https://github.com/user-attachments/assets/f4cb5eda-3f10-43c8-a815-8012cda3d2e0)

![Screenshot from 2024-05-22 00-27-00](https://github.com/user-attachments/assets/dbd6fc01-9d8b-42a1-ae66-7ee1ae7bafe6)

## Installation
To get a local copy up and running, follow the steps below.

## Disclaimer
This prototype was developed in a Unix environment. Running the program in Windows may introduce some issues.

### Prerequisites
- Python 3.8 or higher
- **Go Lang**: [Download and install Go](https://golang.org/dl/)
- **Ollama**: [Download and install Ollama](https://ollama.com/)

### Installation Steps
1. Clone the repository:
   ```sh
   git clone https://github.com/zenolucas/AI-Dietitian

2. Navigate into project directory
   ```sh
   cd AI-Dietitian

3. Install required Python packages for rag_QA.py
   ```sh
   pip install llama_index
   pip install llama_index.embeddings.huggingface

4. Pull phi3 in ollama
   ```sh
   ollama pull phi3

5. Install required Go packages
   ```sh
   go mod tidy

6. Run the go server
   ```sh
   go run main.go

  Usage
To interact with the chatbot, navigate to http://localhost:3000 in your web browser after starting the Go server. 

