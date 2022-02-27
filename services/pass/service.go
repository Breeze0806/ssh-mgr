package pass

type PassMapper interface {
	Write(password string) (err error)
	Match(password string) (err error)
}

type Service struct {
	mapper PassMapper
}

func NewService(mapper PassMapper) *Service {
	return &Service{
		mapper: mapper,
	}
}

func (s *Service) Init(password string) error {
	return s.mapper.Write(password)
}

func (s *Service) Match(password string) error {
	return s.mapper.Match(password)
}
