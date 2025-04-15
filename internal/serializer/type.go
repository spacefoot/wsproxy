package serializer

type Weight struct {
	Weight float64 `json:"weight"`
	Unit   string  `json:"unit"`
}

type Status struct {
	Open bool `json:"open"`
}

type RequestStatus struct{}
type RequestWeight struct{}
type Unstable struct{}
type Zero struct{}

type Log struct {
	Enabled bool `json:"enabled"`
}

type DebugWeight Weight
type DebugUnstable Unstable
