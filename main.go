package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	// "reflect"
	// "fmt"
	// "encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

//var DB = make(map[string]string)

type announcement struct {
	StartTime   int    `yaml:"startTime" json:"startTime"`
	EndTime     int    `yaml:"endTime" json:"endTime"`
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description" json:"description"`
	Pic         string `yaml:"pic" json:"pic"`
}

var (
	announcementURL string = "https://raw.githubusercontent.com/ps10659/shield/master/announcement.yaml"

	announcements = map[string]announcement{}

	ErrNoAnnouncement = errors.New("Announcement file not found")
)

func getAnnouncement(c *gin.Context, announcementURL string) error {
	resp, err := http.Get(announcementURL)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		// handle err
		//return nil, ErrNoAnnouncement
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// err
	}
	// c.String(http.StatusOK, string(body))

	err = yaml.Unmarshal(body, &announcements)
	if err != nil {
		// return unmarshall error
	}

	return nil
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
		err := getAnnouncement(c, announcementURL)
		if err != nil {
			if err == ErrNoAnnouncement {
				// return http status 204
			}
			// other err do something?

			return
		}

		c.JSON(http.StatusOK, announcements["EN"])
		// c.JSON(http.StatusOK, announcements["TW"])

		return
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
