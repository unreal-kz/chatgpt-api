package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/unreal-kz/chatgpt-api/tts"
	"google.golang.org/api/option"
	texttospeech "google.golang.org/api/texttospeech/v1"
)

// Function to convert text to speech
func textToSpeech(text string) ([]byte, error) {
	godotenv.Load()

	ctx := context.Background()
	client, err := texttospeech.NewService(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
	if err != nil {
		return nil, err
	}

	resp, err := client.Text.Synthesize(&texttospeech.SynthesizeSpeechRequest{
		Input: &texttospeech.SynthesisInput{
			Text: text,
		},
		Voice: &texttospeech.VoiceSelectionParams{
			LanguageCode: "en-US",
			SsmlGender:   "NEUTRAL",
		},
		AudioConfig: &texttospeech.AudioConfig{
			AudioEncoding: "MP3",
		},
	}).Do()

	if err != nil {
		return nil, err
	}

	// Decode base64 string to []byte
	audioContent, err := base64.StdEncoding.DecodeString(resp.AudioContent)
	if err != nil {
		return nil, err
	}

	return audioContent, nil
}

func main() {
	my_text := strings.Join(os.Args[1:], " ")
	log.Println(my_text)
	t1 := tts.TTS{Text: my_text}
	responseText, err := t1.GPTResponce()
	if err != nil {
		log.Fatalf("Error fetching response: %v", err)
	}

	audioData, err := textToSpeech(responseText)
	if err != nil {
		log.Fatalf("Error converting text to speech: %v", err)
	}

	if err := os.WriteFile("output.mp3", audioData, 0644); err != nil {
		log.Fatalf("Failed to save audio file: %v", err)
	}

	fmt.Println("Audio file has been saved as output.mp3")
}
