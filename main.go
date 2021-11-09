package main

import (
	"encoding/json"
	"net/http"

	"github.com/ShwetaT0mar/FoodDeliveryCalculator/service"
)

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
