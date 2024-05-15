# My Virtual AI-Dietitian

## Table of Contents
- [About](#about)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## About
This project is a chatbot designed to provide accurate and helpful information about diabetes. 
Leveraging phi3 for its natural language processing capabilities, this chatbot uses Retrieval-Augmented Generation (RAG) to deliver precise and contextual responses by accessing diabetes-related documents.


## Installation
To get a local copy up and running, follow these simple steps.

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

3. Install required Python packages for rag.py
   ```sh
   pip install llama_index

4. Pull phi3 in ollama
   ```sh
   ollama pull phi3

6. Run the go server
   ```sh
   go run main.go

  Usage
To interact with the chatbot, navigate to http://localhost:3000 in your web browser after starting the Go server. 
