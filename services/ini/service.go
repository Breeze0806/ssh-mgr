package ini

import (
	"github.com/Breeze0806/ssh-mgr/dao"
)

type ConnMapper interface {
	Write(conn dao.SshConnInfo) (err error)
}

type Service struct {
	ConnMapper ConnMapper
}

func NewService(connMapper ConnMapper) *Service {
	return &Service{
		ConnMapper: connMapper,
	}
}

func (s *Service) Add(conn dao.SshConnInfo) (err error) {
	return s.ConnMapper.Write(conn)
}
