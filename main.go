package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"
	// "strconv"
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

type status struct {
	Shielding_all bool                    `json:"shielding_all"`
	Announcements map[string]announcement `json:"announcements"`
}

var (
	// announcements is a global variable holding the information sourced from announcementURL
	announcements          = map[string]announcement{}
	announcementURL string = "https://raw.githubusercontent.com/ps10659/shield/master/announcement.yaml"

	// shielding_all determines whether to send announcements response to the clients
	// can also use shielding_iOS, shielding_android, shielding_web to block specific platform users
	shielding_all = false

	// update announcement from github at most once an updateDuration time
	updateDuration  = time.Second * 5
	lastUpadateTime time.Time

	// err
	ErrNoAnnouncement = errors.New("Announcement file not found")
)

func main() {
	r := gin.Default()

	// authorized user only
	r.GET("/shield/set/:action", setShield)
	r.GET("/shield/info", getShieldInfo)
	r.POST("/webhook", webhook)

	// for all users
	r.GET("/announcement", getAnnouncement)

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

func setShield(c *gin.Context) {
	action := c.Param("action")

	switch action {
	case "on":
		shielding_all = true
		c.JSON(http.StatusOK, gin.H{"shielding_all": shielding_all})
	case "off":
		shielding_all = false
		c.JSON(http.StatusOK, gin.H{"shielding_all": shielding_all})
	default:
		// Should this case return 4XX ???
		c.Status(http.StatusNoContent)
	}
}

func getShieldInfo(c *gin.Context) {
	err := updateAnnouncement(c, announcementURL)
	if err != nil {
		// err
	}

	info := status{Shielding_all: shielding_all, Announcements: announcements}
	c.JSON(http.StatusOK, info)
}

func webhook(c *gin.Context) {
	err := updateAnnouncement(c, announcementURL)
	if err != nil && err == ErrNoAnnouncement {
		// c.Status(http.StatusNoContent)
		c.Status(http.StatusNoContent)
		return
	}
	return
}

func getAnnouncement(c *gin.Context) {
	if shielding_all == false {
		c.Status(http.StatusNoContent)
		return
	}

	language := c.Request.Header.Get("language") // MustGet fail why???
	duration := time.Now().Sub(lastUpadateTime)  // Will this line cause server bad performance?
	if duration > updateDuration {
		err := updateAnnouncement(c, announcementURL)
		if err != nil && err == ErrNoAnnouncement {
			// c.Status(http.StatusNoContent)
			c.Status(http.StatusNoContent)
			return
		}
		// if other err, do something
	}

	switch language {
	case "TW":
		c.JSON(http.StatusOK, announcements["TW"])
	default:
		c.JSON(http.StatusOK, announcements["EN"])
	}
	return
}

func updateAnnouncement(c *gin.Context, announcementURL string) error { //Still need to pass gin.Context?
	lastUpadateTime = time.Now()
	resp, err := http.Get(announcementURL)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil { // & err == ErrFileNotFound
		return ErrNoAnnouncement
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
