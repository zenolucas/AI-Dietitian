package handler

import (
	"AI-Dietitian/types"
	"AI-Dietitian/view/chat"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func HandleChatIndex(w http.ResponseWriter, r *http.Request) error {
	return chat.Index().Render(r.Context(), w)
}

func HandleChatCreate(w http.ResponseWriter, r *http.Request) error {
	prompt := r.FormValue("prompt")

	if strings.Contains(prompt, "recommend") && strings.Contains(prompt, "meal") {
		// we assume it's a recommendation task

		print("this is executed")

		// do RAG first given prompt
		cmd := exec.Command("python", "rag_recommend.py", prompt)
		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			fmt.Println("Error executing command:", err)
			return err
		}

		// Extract the output from the buffer
		cmdOutput := out.String()

		// some preprocessing, removal of irrelevant CMD output from running rag.py
		cmdSplit := strings.Split(cmdOutput, "[INST]")
		raggedPrompt := cmdSplit[1]
		print(raggedPrompt)

		// Replace newline characters with spaces so that raggedPrompt is 1 single line.
		// raggedPrompt = strings.ReplaceAll(raggedPrompt, "\n", " ")

		// Define a struct for your payload
		type Payload struct {
			Model  string `json:"model"`
			Prompt string `json:"prompt"`
			Stream bool   `json:"stream"`
		}

		// Create an instance of Payload with the desired data
		data := Payload{
			Model:  "phi3",
			Prompt: raggedPrompt,
			Stream: false,
		}

		// // to send a request to ollama hosted locally
		url := "http://localhost:11434/api/generate"
		headers := `{"Content-Type": "application/json"}`
		// payload := []byte(`{
		// 				"model": "phi3",
		// 				"prompt": "` + raggedPrompt + `",
		// 				"stream": false
		// 				}`)

		// Marshal the payload to JSON
		payloadBytes, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return err
		}

		response, err := http.Post(url, headers, bytes.NewBuffer(payloadBytes))
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

		// amen! so far so good. next is.

		// ollama Response is a filepath to our recipe.pdf
		// e.g. /home/zenolucas/Projects/webDev/AI-Dietitian/Documents/MealRec/Fruited Chicken Salad.pdf

		// so given this
		recipeFilePath := ollamaResp.Response

		recipe, err := readPdf(recipeFilePath) // Read local pdf file
		if err != nil {
			panic(err)
		}

		recipe = strings.ReplaceAll(recipe, "\n", "")
		recipe = strings.ReplaceAll(recipe, "\f", "")

	
		var meal types.Meal
		// next, unmarshaling JSON
		errs := json.Unmarshal([]byte(recipe), &meal)
		if errs != nil {
			fmt.Println(errs)
		}

		fmt.Println("Meal Name:", meal.MealName)
		fmt.Println("Image File Path:", meal.ImageFilePath)
		fmt.Println("Ingredients:", meal.Ingredients)
		fmt.Println("Procedure:", meal.Directions)

	
		chat_params := chat.ChatParams{
			MealName: meal.MealName,
			Answer: ollamaResp.Response,
			FileName: meal.ImageFilePath,
			Ingredients: meal.Ingredients,
			Procedure: meal.Directions,
		}

		errors := chat.ChatErrors{}

		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
		}

		return render(r, w, chat.ChatForm(chat_params, errors))

		// var meal types.Meal
		// // unmarshaling JSON
		// errs := json.Unmarshal([]byte(raggedPrompt), &meal)
		// if errs != nil {
		// 	fmt.Println(errs)
		// }

		// // Print extracted information
		// // fmt.Println("Meal Name:", meal.MealName)
		// // fmt.Println("Image File Path:", meal.ImageFilePath)
		// // fmt.Println("Ingredients:", meal.Ingredients)
		// // fmt.Println("Procedure:", meal.Directions)
		// // fmt.Println("Friendly Comments:", meal.FriendlyComments)

		// rec_params := chat.ChatParams{
		// 	Prompt:   r.FormValue("prompt"),
		// 	FileName: meal.ImageFilePath,
		// }

		// return render(r, w, chat.ChatForm(rec_params, chat.ChatErrors{}))

		// print("recommend a meal using RAG.")

	} else if strings.Contains(prompt, "What is") || strings.Contains(prompt, "what is") || strings.Contains(prompt, "?") {
		// else if this is a QandA task, do RAG first then feed raggedPrompt into model

		cmd := exec.Command("python", "rag_.py", prompt)
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

		chat_params := chat.ChatParams{
			Answer: ollamaResp.Response,
		}

		errors := chat.ChatErrors{}

		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return render(r, w, chat.ChatForm(chat_params, errors))
		}

		return render(r, w, chat.ChatForm(chat_params, errors))
	} else {
		// else chat with LLM

		print(prompt)
		// to send a request to ollama hosted locally
		url := "http://localhost:11434/api/generate"
		headers := `{"Content-Type": "application/json"}`
		payload := []byte(`{
						"model": "phi3",
						"prompt": "` + prompt + `",
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
		}

		var ollamaResp OllamaResponse
		if err := json.Unmarshal(body, &ollamaResp); err != nil {
			return err
		}

		chat_params := chat.ChatParams{
			Answer: ollamaResp.Response,
		}

		errors := chat.ChatErrors{}

		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return render(r, w, chat.ChatForm(chat_params, errors))
		}

		return render(r, w, chat.ChatForm(chat_params, errors))

	}

	return nil
}

func readPdf(path string) (string, error) {
	// so given this

	cmd := exec.Command("pdf2txt", "-o", "recipe.txt", path)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return "error happened", err
	}

	b, err := os.ReadFile("recipe.txt") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	str := string(b) // convert content to a 'string'

	fmt.Println(str) // print the content as a 'string'

	// fmt.Println(str) // print the content as a 'string'
	return str, err
}
