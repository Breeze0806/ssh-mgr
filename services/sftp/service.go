package sftp

import (
	"os/exec"

	"github.com/Breeze0806/ssh-mgr/dao"
)

type ConnMapper interface {
	Read(group, name string) (conn dao.SshConnInfo, err error)
}

type Service struct {
	Pragram    string
	ConnMapper ConnMapper
}

func NewService(pragram string, connMapper ConnMapper) *Service {
	return &Service{
		Pragram:    pragram,
		ConnMapper: connMapper,
	}
}

func (s *Service) StartPragram(group, name string) (err error) {
	var conn dao.SshConnInfo
	conn, err = s.ConnMapper.Read(group, name)
	if err != nil {
		return
	}
	cmd := exec.Command(s.Pragram, conn.SftpURL())
	return cmd.Start()
}

func (s *Service) RunPragram(group, name string) (err error) {
	var conn dao.SshConnInfo
	conn, err = s.ConnMapper.Read(group, name)
	if err != nil {
		return
	}
	cmd := exec.Command(s.Pragram, conn.SftpURL())
	return cmd.Run()
}
