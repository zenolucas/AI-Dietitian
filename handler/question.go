package handler

import (
	"AI-Dietitian/view/QandA"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

func HandleQuestionIndex(w http.ResponseWriter, r *http.Request) error {
	return question.Index().Render(r.Context(), w)
}

func HandleQuestionCreate(w http.ResponseWriter, r *http.Request) error {
	prompt := r.FormValue("prompt")

		print(prompt)
		cmd := exec.Command("python", "rag_QA.py", prompt)
		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			fmt.Println("Error executing command:", err)
			return err
		}

		// Extract the output from the buffer
		raggedPrompt := out.String()

		// some preprocessing, removal of irrelevant CMD output from running rag.py
		raggedPrompt = strings.ReplaceAll(raggedPrompt, "LLM is explicitly disabled. Using MockLLM.", "")
		print(raggedPrompt)

		// Replace newline characters with spaces so that raggedPrompt is 1 single line.
		raggedPrompt = strings.ReplaceAll(raggedPrompt, "\n", " ")

		// to send a request to ollama hosted locally
		url := "http://localhost:11434/api/generate"
		headers := `{"Content-Type": "application/json"}`
		payload := []byte(`{
						"model": "phi3",
						"prompt": "` + raggedPrompt + `",
						"stream": false 
						}`)

		response, err := http.Post(url, headers, bytes.NewBuffer(payload))
		if err != nil {
			fmt.Println("something went wrong:", err)
			http.Error(w, "Error sending request to Ollama", http.StatusInternalServerError)
			return err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			fmt.Println("ERROR code: ", response.StatusCode)
			return nil // Or return an error if you want to handle it differently
		}

		// Read the response body
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		type OllamaResponse struct {
			Response string `json:"response"`
			// Add more fields if needed
		}

		var ollamaResp OllamaResponse
		if err := json.Unmarshal(body, &ollamaResp); err != nil {
			return err
		}

		chat_params := question.ChatParams{
			Answer: ollamaResp.Response,
		}

		errors := question.ChatErrors{}

		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return render(r, w, question.ChatForm(chat_params, errors))
		}

		return render(r, w, question.ChatForm(chat_params, errors))
	}
