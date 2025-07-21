package perplexity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const Endpoint = "https://api.perplexity.ai/chat/completions"

type Service struct {
	apiKey string
}

func NewService(
	apiKey string,
) *Service {
	return &Service{
		apiKey: apiKey,
	}
}

type Request struct {
	Model    string `json:"model"`
	Messages []Msg  `json:"messages"`
}

type Msg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (s *Service) Ask(
	question string,
) (string, error) {
	reqBody, _ := json.Marshal(
		Request{
			Model: "sonar",
			Messages: []Msg{
				{Role: "user", Content: FormatQuestionForTg(question)},
			},
		},
	)

	req, err := http.NewRequest(
		http.MethodPost,
		Endpoint,
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Println(string(bodyBytes))

	var result Response
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return "", err
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}

	return "No answer received.", nil
}

var DefaultPrompt = `This is a question asked via Perplexity API which will only be displayed inside a Telegram chat.
Thus, answer it by following these rules strictly:
- Including all sources as a list at the end of the message.
- Format it accordingly to be pretty inside the telegram chat.
- Answer by using the language the actual question is asked.

The actual question is: "%s"`

func FormatQuestionForTg(
	question string,
) string {
	return fmt.Sprintf(DefaultPrompt, question)
}
