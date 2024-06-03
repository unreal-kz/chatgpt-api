package tts

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

const (
	apiEndpoint = "https://api.openai.com/v1/chat/completions"
)

type TTS struct {
	Text string
}

func (t *TTS) getChatGPT() ([]byte, error) {
	godotenv.Load()
	client := resty.New()
	resp, err := client.R().
		SetAuthToken(os.Getenv("GPT_PROJ_KEY")).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"messages":   []interface{}{map[string]interface{}{"role": "system", "content": t.Text}},
			"model":      "gpt-3.5-turbo",
			"max_tokens": 200,
		}).
		Post(apiEndpoint)

	if err != nil {
		return nil, err
	}

	return resp.Body(), nil
}

func (t *TTS) GPTResponce() (string, error) { //
	body, err := t.getChatGPT()
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error while decoding JSON response:", err)
		return "", nil
	}
	// Extract the content from the JSON response
	// content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"]
	// fmt.Printf("%[1]T, %[1]v\n", data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"])
	return data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string), nil
}
