package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
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

	userRes, err := db.Query("CREATE TABLE IF NOT EXISTS \"users\" (id serial PRIMARY KEY, username VARCHAR(100) NOT NULL, email VARCHAR(100) NOT NULL, password VARCHAR(100) NOT NULL)")
	CheckError(err)

	fmt.Println(userRes)

	fmt.Println("Connected!")

	r := gin.Default()

	r.POST("/user", func(c *gin.Context) {
		user := &User{}
		if err := c.ShouldBindJSON(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		sqlStatement := `
INSERT INTO users (username, email, password)
VALUES ($1, $2, $3)
RETURNING id`
		err = db.QueryRow(sqlStatement, user.Username, user.Email, user.Password).Scan(&user.Id)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, user.Id)
	})

	r.GET("/user", func(c *gin.Context) {
		user := &User{}
		rows, err := db.Query("SELECT id, usename, email FROM users")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&user.Id, &user.Username, &user.Email)
			if err != nil {
				panic(err)
			}
			fmt.Println("\n", user.Id, user.Username, user.Email)
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, user)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
