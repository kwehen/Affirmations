package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	// This makes it so each ip can only make 5 requests per hour
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Hour,
		Limit: 10,
	})
	mw := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})
	router.POST("/chat", mw, chat)
	// router.OPTIONS("/chat", mw, chat)
	err := router.Run("0.0.0.0:8080")
	if err != nil {
		log.Println("Error starting webserver", err)
	}
}

func chat(c *gin.Context) {
	affirmations := []string{
		"Write a unqiue one sentance daily affirmation. Again, the must only be one sentance. Do not use the word embrace. The affirmation must be different every time. Be creative and original.",
		"Write an affirmation that stresses refocus, hard work, hopefulness, and positivity. Make sure this is only one sentance. Be creative and original.",
		"Write a one sentance unique daily affirmation that focuses on gratitude, a positive outlook, and more. This must only be one sentance. Be creative and original.",
		"Write a one sentance affirmation that affirms the user of their self-worth and value, adds a little bit on open-mindedness, and gives a positive jolt for the rest of the day. This must only be one sentance. Be creative and original.",
	}

	api_key := os.Getenv("OPENAI_API_KEY")
	w := openai.NewClient(api_key)
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT4o,
		MaxTokens: 35,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: affirmations[rand.Intn(len(affirmations))],
			},
		},

		// Prompt:    "Lorem ipsum
	}
	resp, err := w.CreateChatCompletion(ctx, req)
	if err != nil {
		log.Printf("Completion error: %v\n", err)
		return
	}
	response := (resp.Choices[0].Message.Content)
	c.JSON(http.StatusOK, gin.H{"message": response})
}
