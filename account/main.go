package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Account struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func CheckEnvErr(env string, key string) {
	if env == "" {
		panic(fmt.Sprintf("Environment variable %s not set", key))
	}
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	host := os.Getenv("POSTGRES_HOST")
	CheckEnvErr(host, "POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	CheckEnvErr(port, "POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	CheckEnvErr(user, "POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	CheckEnvErr(password, "POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	CheckEnvErr(dbname, "POSTGRES_DB")

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	userRes, err := db.Query("CREATE TABLE IF NOT EXISTS \"accounts\" (id serial PRIMARY KEY, name VARCHAR(100) NOT NULL, email VARCHAR(100) NOT NULL)")
	CheckError(err)

	fmt.Println(userRes)
	fmt.Println("Connected!")

	r := gin.Default()

	r.POST("/account", func(c *gin.Context) {
		account := &Account{}
		if err := c.ShouldBindJSON(account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sqlStatement := `
INSERT INTO accounts (name, email)
VALUES ($1, $2)
RETURNING id`
		err = db.QueryRow(sqlStatement, account.Name, account.Email).Scan(&account.ID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, account.ID)
	})

	r.GET("/account", func(c *gin.Context) {
		account := &Account{}
		rows, err := db.Query("SELECT id, name, email FROM accounts")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&account.ID, &account.Name, &account.Email)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			fmt.Println("\n", account.ID, account.Name, account.Email)
		}
		err = rows.Err()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, account)
	})

	r.Run(":8081") // listen and serve on 0.0.0.0:8081 (for windows "localhost:8081")
}
