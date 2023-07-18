package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Account struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	r := gin.Default()

	r.POST("/account", func(c *gin.Context) {
		account := &Account{}
		if err := c.ShouldBindJSON(account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(account)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while marshalling struct to bytes"})
		}
		res, err := http.Post("http://localhost:8888/account", "application/json", &buf)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
		}
		c.JSON(http.StatusOK, res)
	})

	r.GET("/account", func(c *gin.Context) {
		res, err := http.Get("http://localhost:8888/account")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
		}
		c.JSON(http.StatusOK, res)
	})

	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
