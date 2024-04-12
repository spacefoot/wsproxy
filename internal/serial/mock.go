package serial

type Mock struct {
	status chan<- Status
}

func NewMock(status chan<- Status) *Mock {
	return &Mock{
		status: status,
	}
}

func (s *Mock) Run() {}

func (s *Mock) RequestStatus() {
	s.status <- Status{
		Open: true,
	}
}
