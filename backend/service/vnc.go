package service

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type VNCInfo struct {
	SourceAddr string
	TargetAddr string
	upgrader   *websocket.Upgrader
}

func (v *VNCInfo) forwardTcp(wsConn *websocket.Conn, conn net.Conn) {
	var tcpBuffer [1024]byte
	defer func() {
		if conn != nil {
			err := conn.Close()
			if err != nil {
				log.Printf("%s: closing connection: %s", time.Now().Format(time.Stamp), err)
			}
		}
		if wsConn != nil {
			err := wsConn.Close()
			if err != nil {
				log.Printf("%s: closing connection: %s", time.Now().Format(time.Stamp), err)
			}
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

func (v *VNCInfo) forwardWeb(wsConn *websocket.Conn, conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%s: reading from WS failed: %s", time.Now().Format(time.Stamp), err)
		}
		if conn != nil {
			err := conn.Close()
			if err != nil {
				log.Printf("%s: closing connection: %s", time.Now().Format(time.Stamp), err)
			}
		}
		if wsConn != nil {
			err := wsConn.Close()
			if err != nil {
				log.Printf("%s: closing connection: %s", time.Now().Format(time.Stamp), err)
			}
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

func (v *VNCInfo) ServeWs(c *gin.Context) {
	if v.upgrader == nil {
		v.upgrader = &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
	}
	ws, err := v.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("%s: failed to upgrade to WS: %s", time.Now().Format(time.Stamp), err)
		return
	}

	vnc, err := net.Dial("tcp", v.TargetAddr)
	if err != nil {
		log.Printf("%s: failed to bind to the VNC Server: %s", time.Now().Format(time.Stamp), err)
	}

	go v.forwardTcp(ws, vnc)
	go v.forwardWeb(ws, vnc)
}
