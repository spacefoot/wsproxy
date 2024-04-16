package static

import _ "embed"

//go:embed index.html
var Index string

type IndexData struct {
	Debug          bool
	SimulateSerial bool
}
