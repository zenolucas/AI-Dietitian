package handler

import (
	"AI-Dietitian/types"
	"AI-Dietitian/view/chat"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func HandleChatIndex(w http.ResponseWriter, r *http.Request) error {
	return chat.Index().Render(r.Context(), w)
}

func HandleChatCreate(w http.ResponseWriter, r *http.Request) error {
	params := chat.ChatParams{
		Prompt: r.FormValue("prompt"),
	}
	errors := chat.ChatErrors{}
	prompt := params.Prompt

	// below is where we implement, for what .sh to execute.
	if strings.Contains(prompt, "recommend") && strings.Contains(prompt, "meal") {
		// then we assume it's a meal recommendation task, lol

		// to be change into llama-2-7b-instruction.sh
		cmd := exec.Command("llama-2-7b-chat.sh", prompt)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error executing script:", err)
			return err
		}
		fmt.Print(string(output))

		// let's now extract response from LLM
		text, err := os.ReadFile("./output_recommend.txt")
		if err != nil {
			fmt.Print(err)
		}

		// index of opening and closing brackets of JSON response from LLM
		index1 := strings.LastIndex(string(text), "{")
		index2 := strings.LastIndex(string(text), "}")

		extracted := text[index1 : index2+1]

		fmt.Print(string(extracted))

		var meal types.Meal
		// unmarshaling JSON
		errs := json.Unmarshal([]byte(extracted), &meal)
		if errs != nil {
			fmt.Println(errs)
		}

		// Print extracted information
		fmt.Println("Meal Name:", meal.MealName)
		fmt.Println("Image File Path:", meal.ImageFilePath)
		fmt.Println("Ingredients:", meal.Ingredients)
		fmt.Println("Procedure:", meal.Procedure)
		fmt.Println("Friendly Comments:", meal.FriendlyComments)

		rec_params := chat.ChatParams{
			Prompt:   r.FormValue("prompt"),
			MealName: meal.MealName,
			FileName: meal.ImageFilePath,
		}

		return render(r, w, chat.ChatForm(rec_params, errors))

	} else {
		// else let's just get the answer from SLM (Smol Language Model)
		cmd := exec.Command("RAG.sh", prompt)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error executing script:", err)
			return err
		}
		fmt.Print(string(output))

		// let's now extract response from LLM
		text, err := os.ReadFile("./output_chat.txt")
		if err != nil {
			fmt.Print(err)
		}

		// trim response to only output the LLM prompt
		// Find the index where "[/INST]" starts
		startIndex := strings.Index(string(text), "[/INST]")

		if startIndex != -1 {
			text = text[startIndex:]
		} else {
			fmt.Println("Substring '[/INST]' not found")
		}

		chat_params := chat.ChatParams{
			Prompt: r.FormValue("prompt"),
			Answer: string(text),
		}

		return render(r, w, chat.ChatForm(chat_params, errors))
	}

}
