package serializer

type Weight struct {
	Weight float64 `json:"weight"`
	Unit   string  `json:"unit"`
}

type Status struct {
	Open bool `json:"open"`
}

type RequestStatus struct{}

type Zeroed struct{}
