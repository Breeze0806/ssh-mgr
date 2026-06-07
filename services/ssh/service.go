package ssh

import (
	"os"
	"os/exec"
	"runtime"

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

// puttyCmd builds the exec.Cmd used to launch PuTTY. On Ubuntu 26.04 +
// Wayland the native Linux PuTTY (apt putty 0.83) is built against GTK3
// and the GDK Wayland backend fails to enumerate fonts (returns "unable
// to load font") and trips a glibc symbol conflict; forcing the X11 GDK
// backend fixes both. The env var is a no-op on Windows, where PuTTY is
// a native Win32 app that does not use GDK, so we only set it when not
// building for Windows.
func (s *Service) puttyCmd(url, password string) *exec.Cmd {
	cmd := exec.Command(s.Pragram, url, "-pw", password, "-ssh")
	if runtime.GOOS != "windows" {
		cmd.Env = append(os.Environ(), "GDK_BACKEND=x11")
	}
	return cmd
}

func (s *Service) StartPragram(group, name string) (err error) {
	var conn dao.SshConnInfo
	conn, err = s.ConnMapper.Read(group, name)
	if err != nil {
		return
	}
	return s.puttyCmd(conn.SshURL(), conn.Password).Start()
}

func (s *Service) RunPragram(group, name string) (err error) {
	var conn dao.SshConnInfo
	conn, err = s.ConnMapper.Read(group, name)
	if err != nil {
		return
	}
	return s.puttyCmd(conn.SshURL(), conn.Password).Run()
}
