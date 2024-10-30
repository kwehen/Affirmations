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
	// Each IP can only make 5 requests per hour.
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Hour,
		Limit: 5,
	})
	mw := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})
	router.POST("/chat", mw, chat)
	err := router.Run("0.0.0.0:8080")
	if err != nil {
		log.Println("Error starting webserver", err)
	}
}

func chat(c *gin.Context) {
	affirmations := []string{
		AFFIRMATION_QUERY,
		AFFIRMATION_QUERY_1,
		AFFIRMATION_QUERY_2,
		AFFIRMATION_QUERY_3,
		AFFIRMATION_QUERY_4,
		AFFIRMATION_QUERY_5,
		AFFIRMATION_QUERY_6,
	}

	api_key := os.Getenv("OPENAI_API_KEY")
	w := openai.NewClient(api_key)
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT4o,
		MaxTokens: 45,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: affirmations[rand.Intn(len(affirmations))],
			},
		},
	}
	resp, err := w.CreateChatCompletion(ctx, req)
	if err != nil {
		log.Printf("Completion error: %v\n", err)
		return
	}
	response := (resp.Choices[0].Message.Content)
	c.JSON(http.StatusOK, gin.H{"message": response})
}
