package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
)

//slots stores the slot number and their waiting times

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

	fmt.Println("ORDER COMING IN==============================================", order)

	fmt.Println("slot waiting periods are")
	for i := 1; i <= 7; i++ {
		fmt.Println(i, "  :", GetSlotTime(i))
	}

	fmt.Println()
	fmt.Println()
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

	waitPeriod := []int{}

	slotsRequired := 0
	expectedTime := 0
	lastMeal := ""
	fmt.Println()
	fmt.Println("***calculateTotalTime****")
	for i, meal := range meals {
		fmt.Println(i, "Meal is", meal)
		slotsRequired++
		if meal == "M" {
			slotsRequired += 1
		}
		fmt.Println()
		fmt.Println("Slots needed for", meal, "are", slotsRequired)
		for i := 1; i <= slotsRequired; i++ {
			s := expectedSlot()

			if meal == "M" {
				if !GetPrepTime(s, &expectedTime, 29) {
					return -1
				}

			} else {
				if !GetPrepTime(s, &expectedTime, 17) {
					return -1
				}
			}

			//check if meal can be prepared simultaneously
			if len(waitPeriod) > 0 && (lastMeal == "M" || meal == "M") {
				delay := waitPeriod[len(waitPeriod)-1] - GetSlotTime(s)
				fmt.Println("delay between A & M, M&M", delay)
				if delay > 29 {
					fmt.Println("Delay is too much")
					return -1
				}
			} else if len(waitPeriod) > 0 {
				delay := waitPeriod[len(waitPeriod)-1] - GetSlotTime(s)
				fmt.Println("delay between A & A", delay)
				if delay > 17 {
					fmt.Println("Delay is too much")
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
	fmt.Println()
	fmt.Println("WaitPeriods for meals are     ==", waitPeriod)
	fmt.Println()
	fmt.Println()

	prepTime := waitPeriod[len(waitPeriod)-1]
	totalTime := float32(prepTime) + dTime
	fmt.Println("Total time for meals are", totalTime)
	if totalTime > 150 {
		fmt.Println("Total time is too much", totalTime)
		return -1
	}

	return totalTime
}

func expectedSlot() int {
	s := -1
	fmt.Println("++++++++++++++++++++Finding slots with min waiting period++++++++++++++")
	for slot, time := range SlotStruct.AllSlots {
		fmt.Println("Slot", slot, "is booked for", time, "minutes")
		if s == -1 {
			s = slot
		} else if GetSlotTime(s) > time {
			s = slot
		}
	}
	fmt.Println("+++++++++++++++++++++++++++")
	fmt.Println("Slot picked is", s)
	return s
}

func GetPrepTime(s int, eTime *int, prepTime int) bool {
	fmt.Println("In GetPrepTime")
	*eTime = GetSlotTime(s) + prepTime
	if *eTime > 150 {
		return false
	}
	SlotStruct.Lock()
	SlotStruct.AllSlots[s] += prepTime
	SlotStruct.Unlock()
	fmt.Println("Expected prep time", *eTime)
	//fmt.Println("Slot", s, "booked for", GetSlotTime(s), "minutes")
	return true

}

func GetSlotTime(s int) int {
	SlotStruct.RLock()
	t := SlotStruct.AllSlots[s]
	SlotStruct.RUnlock()
	return t
}
