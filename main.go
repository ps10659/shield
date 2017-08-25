package main

import (
	"errors"
	"net/http"
	"io/ioutil"
	// "reflect"
	// "fmt"
	// "encoding/json"
	"gopkg.in/yaml.v2"	
	"github.com/gin-gonic/gin"
)

//var DB = make(map[string]string)

type Announcement struct {
	startTime int
	endTime int
	title string
	discription string
	pic string
}

var (
	announcementURL string = "https://raw.githubusercontent.com/ps10659/shield/master/announcement.yaml"
	ErrNoAnnouncement = errors.New("Announcement file not found")
	
)

func getAnnouncement(c *gin.Context, announcementURL string) (Announcement, error) {
	file, err := http.Get(announcementURL)
	if err != nil {
		// handle err
		//return nil, ErrNoAnnouncement
	}
	defer file.Body.Close()

	body, _ := ioutil.ReadAll(file.Body)
    // c.String(http.StatusOK, string(body))

    var announcement Announcement
	err = yaml.Unmarshal(body, &announcement)
	if err != nil {
        // return unmarshall error
    }
    // c.String(http.StatusOK, "\nannouncement: %#v\n", announcement)

	return announcement, nil
}


func main() {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/shield", func(c *gin.Context) {
		announcement, err := getAnnouncement(c, announcementURL) 
		if err != nil {
			// return http status 204
		}

		c.String(http.StatusOK, "\nannouncement: %#v\n", announcement)
		// c.JSON(http.StatusOK, announcement)

		return
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

