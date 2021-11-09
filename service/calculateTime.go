package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
)

//AllSlots stores the slot numbers and their waiting times
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
	resp := deliveryResponse{}
	if time == -1 {
		resp = deliveryResponse{false, "Sorry bruh it's taking too long"}

	} else {
		resp = deliveryResponse{true, "Your order will be delivered in " + fmt.Sprintf("%v", time) + " minutes"}
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
		if maxSlots > 7 {
			return -1
		}
	}

	routeTime := o.Distance * 8
	return calculateTotalTime(o.Items, routeTime)
}

func calculateTotalTime(meals []string, dTime float32) float32 {

	waitPeriod := []int{}
	slotsRequired := 0
	expectedTime := 0
	lastMeal := ""

	for _, meal := range meals {
		slotsRequired++
		if meal == "M" {
			slotsRequired += 1
		}
		for i := 1; i <= slotsRequired; i++ {
			s := expectedSlot()

			if meal == "M" {
				if !GetPrepTime(s, &expectedTime, 29) {
					return -1
				}
			} else if !GetPrepTime(s, &expectedTime, 17) {
				return -1
			}

			//check if meals can be prepared simultaneously
			if len(waitPeriod) > 0 {
				delay := waitPeriod[len(waitPeriod)-1] - GetSlotTime(s)
				if (lastMeal == "M" || meal == "M") && delay > 29 {
					return -1
				} else if delay > 17 {
					return -1
				}
			}

			expectedTime = 0
			waitPeriod = append(waitPeriod, GetSlotTime(s))
		}

		slotsRequired = 0
		lastMeal = meal
	}

	sort.Ints(waitPeriod)
	prepTime := waitPeriod[len(waitPeriod)-1]
	totalTime := float32(prepTime) + dTime

	if totalTime > 150 {
		return -1
	}
	return totalTime
}

//returns the slot with smallest waiting time
func expectedSlot() int {
	s := -1
	for slot, time := range SlotStruct.AllSlots {
		if s == -1 {
			s = slot
		} else if GetSlotTime(s) > time {
			s = slot
		}
	}
	return s
}

//returns preparation time for the meal
func GetPrepTime(s int, eTime *int, prepTime int) bool {
	*eTime = GetSlotTime(s) + prepTime
	if *eTime > 150 {
		return false
	}
	SlotStruct.Lock()
	SlotStruct.AllSlots[s] += prepTime
	SlotStruct.Unlock()
	return true

}

func GetSlotTime(s int) int {
	SlotStruct.RLock()
	t := SlotStruct.AllSlots[s]
	SlotStruct.RUnlock()
	return t
}
