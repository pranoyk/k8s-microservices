package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "postgresdb"
	port     = 5432
	user     = "testUser"
	password = "testPassword"
	dbname   = "testDB"
)

func main() {
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

	res, err := db.Query("CREATE TABLE IF NOT EXISTS \"users\" (id serial PRIMARY KEY, name VARCHAR(100) NOT NULL, email VARCHAR(100) NOT NULL, password VARCHAR(100) NOT NULL)")
	CheckError(err)

	fmt.Println(res)

	fmt.Println("Connected!")
	Routes()
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func Routes() {
	r := gin.Default()
	r.GET("/db", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	r.Run(":8888")
}
