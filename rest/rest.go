package rest

import (
	"database/sql"
	"net/http"

	"github.com/LukeMcAuleyDublin/web-service-gin/models"
	"github.com/gin-gonic/gin"
)

func getAlbums(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM albums")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer rows.Close()

		albums := []models.Album{}
		for rows.Next() {
			var a models.Album
			err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			albums = append(albums, a)
		}

		if err := rows.Err(); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, albums)
	}
}

func getAlbumByID(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		row := db.QueryRow("SELECT * FROM albums WHERE id=$1", id)

		var a models.Album
		err := row.Scan(&a.ID, &a.Title, &a.Artist, &a.Price)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		c.JSON(http.StatusOK, a)
	}
}

func postAlbum(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var albums struct {
			Albums []models.Album `json:"albums"`
		}
		if err := c.BindJSON(&albums); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		for _, a := range albums.Albums {
			var id int
			err := db.QueryRow("INSERT INTO albums (title, artist, price) VALUES ($1, $2, $3) RETURNING id", a.Title, a.Artist, a.Price).Scan(&id)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			a.ID = id

			c.JSON(http.StatusOK, a)
		}
	}
}

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/albums", getAlbums(db))
	router.GET("/albums/:id", getAlbumByID(db))
	router.POST("/albums", postAlbum(db))
}
