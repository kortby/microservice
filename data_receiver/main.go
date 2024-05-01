package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/godev/tolls/types"
	"github.com/gorilla/websocket"
)

func main() {
	recv := newDataReceiver()
	http.HandleFunc("/ws", recv.handlerWS)
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	msgch chan types.OBUData
	conn *websocket.Conn
}

func (dr *DataReceiver) handlerWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}

	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("new OBU client connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			continue
		}
		fmt.Printf("received obu data from [%d] : lat %.2f, lng, %.2f \n", data.OBUID, data.Lat, data.Lng)
		dr.msgch <- data
	}
}

func newDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
	}
}