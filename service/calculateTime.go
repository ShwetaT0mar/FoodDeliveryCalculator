package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

//AllSlots stores the slot number and their waiting times
var SlotStruct = struct {
	sync.RWMutex
	AllSlots map[int]int
}{AllSlots: map[int]int{
	1: 0,
	2: 0,
	3: 0,
	4: 0,
	5: 0,
	6: 0,
	7: 0,
}}

func GetDeliveryTime(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error Reading body", http.StatusInternalServerError)
	}

	order := OrderStruct{}
	err = json.Unmarshal(body, &order)
	if err != nil {
		http.Error(w, "Error Unmarshaling body", http.StatusInternalServerError)
	}

	time := calculateTime(order)
	resp := responseDelivery{}
	fmt.Println("===========================time is======================================", time)
	if time == -1 {
		resp = responseDelivery{false, "Taking too long"}

	} else {
		resp = responseDelivery{true, "Your order will be delivered in " + fmt.Sprintf("%v", time) + " minutes"}
	}
	json.NewEncoder(w).Encode(resp)
}

func calculateTime(o OrderStruct) float32 {
	maxSlots := 0
	for _, i := range o.Items {
		maxSlots++
		if i == "M" {
			maxSlots++
		}
	}
	if maxSlots > 7 {

		fmt.Println("No of slots required", maxSlots)
		return -1
	}
	fmt.Println()
	routeTime := o.Distance * 8
	return calculateTotalTime(o.Items, routeTime)
}

func calculateTotalTime(meals []string, dTime float32) float32 {

	return 1
}

func expectedSlot() int {
	return 1
}

func GetPrepTime(s int, eTime *int, prepTime int) bool {
	return true

}
