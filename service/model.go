package service

type OrderStruct struct {
	OrderId  int
	Items    []string
	Distance float32
}

type deliveryResponse struct {
	Placed  bool
	Message string
}
