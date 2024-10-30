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
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Natoque proin id commodo potenti aliquet etiam phasellus duis. Etiam eu viverra ante quis faucibus volutpat urna commodo.",
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Habitant hac gravida gravida sed suscipit ligula interdum mi. Lectus cras quis adipiscing taciti risus mauris condimentum sagittis. Consequat vitae urna feugiat morbi conubia duis mattis porta. Ligula nullam justo blandit dui turpis conubia praesent quis.",
		"Consequat quisque mauris platea taciti blandit sociosqu sociosqu magna purus. Fusce inceptos nam id nulla congue vivamus fermentum turpis nullam. Luctus facilisi metus non vehicula sit ipsum donec orci quisque. Velit consequat risus suspendisse id curabitur cras congue ac magnis.",
		"Gravida adipiscing consectetur ad per dictumst est montes vulputate suscipit eget blandit aliquam urna. Tincidunt ligula urna mattis venenatis conubia feugiat suspendisse sem eget urna per vehicula mollis. Dignissim congue faucibus cum vestibulum malesuada rutrum lectus risus suspendisse lorem sollicitudin non laoreet.",
		"Potenti neque vehicula integer pellentesque curabitur praesent potenti pellentesque bibendum hac torquent nisi taciti. In mattis nunc ultricies vestibulum nisl sagittis amet risus sit rhoncus class commodo ornare. Cubilia maecenas et donec odio rhoncus cubilia fringilla a tempus mus turpis arcu penatibus. Laoreet dapibus maecenas primis senectus donec natoque consectetur litora at habitasse ridiculus sagittis cras.",
		"Feugiat nullam potenti netus pellentesque id consequat phasellus hendrerit. Cras class pretium lorem vehicula purus risus sociosqu nostra. Dolor iaculis hac natoque in condimentum scelerisque rhoncus erat. Platea ultricies convallis platea porttitor scelerisque egestas euismod orci.",
		"Sed litora vestibulum dolor inceptos non nam magna porttitor integer neque per etiam dictumst hendrerit. Nulla dictumst felis vitae viverra aptent sociis etiam pretium torquent massa montes odio porta leo. Sodales ultricies curae feugiat nulla metus cum ornare cubilia duis cum litora convallis porttitor sit. Euismod inceptos luctus ad nam rutrum ornare magna nam donec lectus odio lobortis magna sociis.",
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
