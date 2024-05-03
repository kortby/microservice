package types


type OBUData struct {
	OBUID	int 	`json:"obuID"`
	Lat		float64 `json:"lat"`
	Lng		float64 `json:"lng"`
}

type Distance struct {
	Value float64 `json:"value"`
	OBUID int 	  `json:"obuID`
	Unix int64	  `json:unix"`
}


type Invoice struct {
	OBUID		  int 		`json:"obuID"`
	TotalDistance float64 	`json:"totalDistance"`
	TotalAmount	  float64	`json:"totalAmount"`
}