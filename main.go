package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ShwetaT0mar/FoodDeliveryCalculator/service"
)

func init() {
	Clock()
}

func main() {

	http.HandleFunc("/deliveryTime", service.GetDeliveryTime)
	http.HandleFunc("/isActive", active)

	http.ListenAndServe(":8080", nil)

	select {}

}

func active(w http.ResponseWriter, r *http.Request) {
	a := struct {
		Status string
	}{"Active"}
	json.NewEncoder(w).Encode(a)
}

func Clock() {
	go Timer(1)
	go Timer(2)
	go Timer(3)
	go Timer(4)
	go Timer(5)
	go Timer(6)
	go Timer(7)
}

func Timer(slot int) {

	for {
		t := service.GetSlotTime(slot)
		if t != 0 {
			time.Sleep(1 * time.Minute)
			service.SlotStruct.Lock()
			service.SlotStruct.AllSlots[slot]--
			service.SlotStruct.Unlock()
		}
	}
}
