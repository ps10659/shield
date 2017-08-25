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

type announcement struct {
	StartTime int `yaml:"startTime"`
	EndTime int `yaml:"endTime"`
	Title string `yaml:"title"`
	Description string `yaml:"description"`
	Pic string `yaml:"pic"`	
}

var (
	announcementURL string = "https://raw.githubusercontent.com/ps10659/shield/master/announcement.yaml"
	// announcementURL string = "file:///Users/jackychen/workspace/go/src/github.com/ps10659/shield/announcement.yaml"

	ErrNoAnnouncement = errors.New("Announcement file not found")

	// announcements = map[string] announcement{}	
)

func getAnnouncement(c *gin.Context, announcementURL string) (map[string] announcement, error) {
	resp, err := http.Get(announcementURL)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		// handle err
		//return nil, ErrNoAnnouncement
	}
	


	body, _ := ioutil.ReadAll(resp.Body)
    c.String(http.StatusOK, string(body))


    announcements := map[string] announcement{}	
	err = yaml.Unmarshal(body, &announcements)
	if err != nil {
        // return unmarshall error
    }
    // c.String(http.StatusOK, "\nannouncement: %#v\n", announcement)
    // c.String(http.StatusOK, "\nannouncement: %+v\n", announcements["EN"])


	return announcements, nil
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
		announcements, err := getAnnouncement(c, announcementURL) 
		if err != nil {
			// return http status 204
		}

		c.String(http.StatusOK, "\nannouncement: %+v\n", announcements["EN"])
		c.String(http.StatusOK, "\nannouncement: %+v\n", announcements["TW"])
		// c.JSON(http.StatusOK, announcement)

		return
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

