package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/danilopolani/gocialite"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Define our gocialite instance
var gocial = gocialite.NewDispatcher()

func main() {
	router := gin.Default()

	err := godotenv.Load("./.env")
	if err != nil {
		panic(err)
	}
	router.GET("/", indexHandler)
	router.GET("/auth/:provider", redirectHandler)
	router.GET("/auth/:provider/callback", callbackHandler)

	router.Run("127.0.0.1:8000")
}

func indexHandler(c *gin.Context) {
	c.Writer.Write([]byte("<html><head><title>Gocialite example</title></head><body>" +
		"<a href='/auth/github'><button>Login with GitHub</button></a></br>" +
		"<a href='/auth/google'><button>Login with Google</button></a><br>" +
		"</body></html>"))
}

func redirectHandler(c *gin.Context) {
	provider := c.Param("provider")
	providerSecrets := map[string]map[string]string{
		"github": {
			"clientID":     os.Getenv("GITHUB_CLIENT_ID"),
			"clientSecret": os.Getenv("GITHUB_CLIENT_SECRET"),
			"redirectURL":  "http://localhost:8000/auth/github/callback",
		},
		"google": {
			"clientID":     os.Getenv("GOOGLE_CLIENT_ID"),
			"clientSecret": os.Getenv("GOOGLE_CLIENT_SECRET"),
			"redirectURL":  "http://localhost:8000/auth/google/callback",
		},
	}

	providerScopes := map[string][]string{
		"github": []string{"public_repo"},
		"google": []string{},
	}
	providerData := providerSecrets[provider]
	actualScopes := providerScopes[provider]
	authURL, err := gocial.New().Driver(provider).Scopes(actualScopes).Redirect(
		providerData["clientID"],
		providerData["clientSecret"],
		providerData["redirectURL"],
	)

	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	c.Redirect(http.StatusFound, authURL)
}

func callbackHandler(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")
	provider := c.Param("provider")

	user, token, err := gocial.Handle(state, code)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	fmt.Printf("%#v", token)

	fmt.Printf("%#v", user)

	c.Writer.Write([]byte("Hi, " + user.FullName))
	c.Writer.Write([]byte("Provider: " + provider))
}
