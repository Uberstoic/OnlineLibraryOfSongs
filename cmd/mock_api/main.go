package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	r.GET("/info", func(c *gin.Context) {
		group := c.Query("group")
		song := c.Query("song")
		if group == "" || song == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "group and song parameters are required"})
			return
		}

		// Возвращаем данные в точном соответствии со спецификацией Swagger
		songDetail := SongDetail{
			ReleaseDate: "16.07.2006",
			Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
			Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		}

		c.JSON(http.StatusOK, songDetail)
	})

	port := "8081"
	log.Printf("Mock API starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start mock API:", err)
	}
}
