package main

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/LukeMcAuleyDublin/web-service-gin/models"
	"github.com/LukeMcAuleyDublin/web-service-gin/rest"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	wg := sync.WaitGroup{}
	wg.Add(1)

	var conf = models.DbConfig{Host: "localhost", Port: 5432, User: "docker", Password: "docker", DatabaseName: "my_database"}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Host, conf.Port, conf.User, conf.Password, conf.DatabaseName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS albums (
			ID SERIAL PRIMARY KEY,
			Title TEXT NOT NULL,
			Artist TEXT NOT NULL,
			Price FLOAT NOT NULL
		)
	`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	rest.RegisterRoutes(router, db)
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
}
