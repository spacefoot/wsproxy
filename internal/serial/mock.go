package serial

type Mock struct {
	read   chan<- []byte
	write  <-chan []byte
	status chan<- Status
}

func NewMock(read chan<- []byte, write <-chan []byte, status chan<- Status) *Mock {
	return &Mock{
		read:   read,
		write:  write,
		status: status,
	}
}

func (s *Mock) Run() {
	for msg := range s.write {
		// TODO: hard coded for courier 5000
		if string(msg) == "Z\r\n" {
			s.read <- []byte("0 g\r\n")
		}
	}
}

func (s *Mock) RequestStatus() {
	s.status <- Status{
		Open: true,
	}
}
