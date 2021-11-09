# FoodDeliveryCalculator


Restaurant has 7 slots, these slots are stored in a map with their waiting time(time after which the slots is available). 
Go routines decrement the waiting time per passing minute.

(To avoid concurrent read and write on maps synchronization is used)

Steps:

1) For every meal slot with the least wait time is selected.
2) Preperation time for the meal is calculated 
    if less than 150 minutes 
      then the slot is booked for the meal and waiting time for alot is incremented
3)To check if the orders will be prepared simulatenously, wait period between any two slots booked for an order has to be less then 17/29 depending on the meal type.
3) Max preparation time of all meals + Time on Route if less than 150 minutes is returned
