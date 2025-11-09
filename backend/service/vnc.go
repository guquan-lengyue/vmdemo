// A Go version WebSocket to TCP socket proxy
// Copyright 2021 Michael.liu
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// from https://github.com/novnc/websockify-other/tree/master/golang

package service

import (
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
	"vmdemo/kvm"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	targetAddr string
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func forwardTcp(wsConn *websocket.Conn, conn net.Conn) {
	var tcpBuffer [1024]byte
	defer func() {
		if conn != nil {
			conn.Close()
		}
		if wsConn != nil {
			wsConn.Close()
		}
	}()
	for {
		if (conn == nil) || (wsConn == nil) {
			return
		}
		n, err := conn.Read(tcpBuffer[0:])
		if err != nil {
			log.Printf("%s: reading from TCP failed: %s", time.Now().Format(time.Stamp), err)
			return
		} else {
			if err := wsConn.WriteMessage(websocket.BinaryMessage, tcpBuffer[0:n]); err != nil {
				log.Printf("%s: writing to WS failed: %s", time.Now().Format(time.Stamp), err)
			}
		}
	}
}

func forwardWeb(wsConn *websocket.Conn, conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%s: reading from WS failed: %s", time.Now().Format(time.Stamp), err)
		}
		if conn != nil {
			conn.Close()
		}
		if wsConn != nil {
			wsConn.Close()
		}
	}()
	for {
		if (conn == nil) || (wsConn == nil) {
			return
		}

		_, buffer, err := wsConn.ReadMessage()
		if err == nil {
			if _, err := conn.Write(buffer); err != nil {
				log.Printf("%s: writing to TCP failed: %s", time.Now().Format(time.Stamp), err)
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("%s: failed to upgrade to WS: %s", time.Now().Format(time.Stamp), err)
		return
	}

	vnc, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("%s: failed to bind to the VNC Server: %s", time.Now().Format(time.Stamp), err)
	}

	go forwardTcp(ws, vnc)
	go forwardWeb(ws, vnc)
}

func CrateVncWs(r *gin.RouterGroup) {
	r.GET("/vnc/:vmName", func(c *gin.Context) {
		vmName := c.Param("vmName")
		info, err := kvm.GetVMInfo(vmName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		vncInfo := kvm.Domain{}
		err = xml.Unmarshal([]byte(info), &vncInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// 设置目标地址
		address := strings.Replace(vncInfo.Graphics[0].Listen, "0.0.0.0", "127.0.0.1", 1)
		targetAddr = fmt.Sprintf("%s:%s", address, vncInfo.Graphics[0].Port)
		serveWs(c.Writer, c.Request)
	})
}
