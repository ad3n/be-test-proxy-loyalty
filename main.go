package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	r := gin.Default()

	proxy := createReverseProxy(os.Getenv("APP_SCHEMA"), os.Getenv("APP_HOST"))

	proxyRoutes := r.Group("/demo/v1/api")
	{
		proxyRoutes.POST("/manager", proxy)
		proxyRoutes.PUT("/manager/:id", proxy)
		proxyRoutes.DELETE("/manager/:id", proxy)
		proxyRoutes.GET("/manager/:id", proxy)
		proxyRoutes.GET("/manager", proxy)
	}

	r.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
}

func createReverseProxy(scheme string, host string) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := &http.Client{}

		request, err := http.NewRequest(c.Request.Method, fmt.Sprintf("%s://%s%s", scheme, host, strings.ReplaceAll(c.Request.RequestURI, "manager", "valuta")), nil)
		if err != nil {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusNotFound)

			return
		}

		request.Header = c.Request.Header
		response, err := client.Do(request)

		if err != nil {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusNotFound)

			return
		}

		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusNotFound)

			return
		}

		jsonBody := make(map[string]interface{})

		err = json.Unmarshal(body, &jsonBody)
		if err != nil {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusNotFound)

			return
		}

		c.JSON(response.StatusCode, jsonBody)
	}
}
