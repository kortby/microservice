package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/godev/tolls/types"
	"github.com/gorilla/websocket"
)

var sendInterval = time.Second

const wsEndpoint = "ws://127.0.0.1:30000/ws"

func genLatLng() (float64, float64) {
	return genCoord(), genCoord()
}

func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func main() {
	obuIDS := generateOBUIDS(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := 0; i < len(obuIDS); i++ {
			lat, lng := genLatLng()
			data := types.OBUData{
				OBUID: obuIDS[i],
				Lat: lat,
				Lng: lng,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
			// fmt.Printf("%+v\n", data)
		}
		time.Sleep(sendInterval)
	}
}

func generateOBUIDS(n int) []int {
	ids := make([]int,n)
	for i:= 0; i < n; i++ {
		ids[i] = rand.Intn(999999)
	}
	return ids
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
