package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	wg := sync.WaitGroup{}
	wg.Add(2)

	const (
		host     = "localhost"
		port     = 5432
		user     = "docker"
		password = "docker"
		dbname   = "my_database"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	createTable := `
		CREATE TABLE IF NOT EXISTS albums (
			id SERIAL PRIMARY KEY,
			title TEXT,
			artist TEXT,
			price FLOAT
		);
	`

	_, err = db.Exec(createTable)
	if err != nil {
		panic(err)
	}

	insertData := `
		INSERT INTO albums (title, artist, price)
		VALUES ('Imagine', 'John Lennon', 14.99), ('R U Mine?', 'Arctic Monkeys', 12.50), ('Firestarter', 'The Prodigy', 19.75);
	`

	_, err = db.Exec(insertData)
	if err != nil {
		panic(err)
	}

	go func() {
		defer wg.Done()
		if err := router.Run(":8080"); err != nil {
			fmt.Println("Error starting server: ", err)
		}
	}()

	// Test data to POST /albums
	go func() {
		defer wg.Done()

		file, err := os.Open("MOCK_DATA.json")
		if err != nil {
			fmt.Println(errorMessage{Status: "Error", Message: err.Error()})
			return
		}
		defer file.Close()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			displayError(err)
			return
		}

		var payload []album
		err = json.Unmarshal(data, &payload)
		if err != nil {
			displayError(err)
			return
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			displayError(err)
			return
		}

		resp, err := http.Post("http://0.0.0.0:8080/albums", "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			displayError(err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			displayError(err)
			return
		}
		fmt.Println("Response status code: ", resp.StatusCode)
		fmt.Println("Response body: ", string(body))
	}()

	wg.Wait()
}

type album struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type errorMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func displayError(e error) {
	fmt.Println(errorMessage{Status: "Error", Message: e.Error()})
}

var albums []album

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, errorMessage{Status: "Error", Message: err.Error()})
		return
	}

	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func postAlbums(c *gin.Context) {
	var newAlbums []album

	if err := c.BindJSON(&newAlbums); err != nil {
		c.IndentedJSON(http.StatusBadRequest, errorMessage{Status: "Error", Message: err.Error()})
		return
	}

	albums = append(albums, newAlbums...)
	c.IndentedJSON(http.StatusCreated, newAlbums)
}
