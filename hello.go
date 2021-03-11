package main

import (
	"fmt"
	"net/http"
	"time"

	usbdrivedetector "github.com/deepakjois/gousbdrivedetector"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

var usbDrive string

func pollforUSBDrives(m *melody.Melody) {
	for {
		time.Sleep(2 * time.Second)
		go func() {
			if drives, err := usbdrivedetector.Detect(); err == nil {
				fmt.Printf("%d USB Devices Found\n", len(drives))
				if len(drives) > 0 {
					usbDrive = drives[0]
					msg := fmt.Sprintf(`{ "type": "drive",  "data" : "%s"}`, usbDrive)
					m.Broadcast([]byte(msg))
				} else {
					m.Broadcast([]byte(`{}`))
				}
			} else {
				fmt.Println(err)
			}
		}()
	}
}

func main() {

	// WEB SERVER

	r := gin.Default()
	m := melody.New()

	// DRIVE POLLER

	go pollforUSBDrives(m)

	r.Static("./static", "./static")

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	// FILE UPLOADS
	r.POST("/", func(c *gin.Context) {
		// single file
		file, _ := c.FormFile("file")
		c.SaveUploadedFile(file, usbDrive+"/"+file.Filename)
		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	// WEB SOCKETS
	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		fmt.Println((msg))
		m.Broadcast(msg)
	})

	r.Run(":8765")
}
