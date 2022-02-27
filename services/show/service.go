package show

import (
	"io/fs"
	"path/filepath"
	"sort"
	"sync"

	"github.com/Breeze0806/ssh-mgr/dao"
	"github.com/c-bata/go-prompt"
)

type ConnMapper interface {
	Read(group, name string) (conn dao.SshConnInfo, err error)
}

type Service struct {
	mapper ConnMapper

	sync.RWMutex
	names map[string]map[string]dao.SshConnInfo
	addrs map[string]map[string]dao.SshConnInfo
}

func NewService(mapper ConnMapper) *Service {
	return &Service{
		mapper: mapper,
		names:  make(map[string]map[string]dao.SshConnInfo),
		addrs:  make(map[string]map[string]dao.SshConnInfo),
	}
}

func (s *Service) Init(path string) (err error) {
	group := ""
	if err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			name := info.Name()[0 : len(info.Name())-5]

			var conn dao.SshConnInfo
			if conn, err = s.mapper.Read(group, name); err != nil {
				return err
			}
			s.Add(conn)
			return nil
		}
		group = info.Name()
		return nil
	}); err != nil {
		return
	}
	return
}

func (s *Service) Add(conn dao.SshConnInfo) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.names[conn.Group]; !ok {
		s.names[conn.Group] = make(map[string]dao.SshConnInfo)
	}
	s.names[conn.Group][conn.Name] = conn

	if _, ok := s.addrs[conn.Address]; !ok {
		s.addrs[conn.Address] = make(map[string]dao.SshConnInfo)
	}
	s.addrs[conn.Address][conn.Group+"."+conn.Name] = conn
	return nil
}

func (s *Service) Groups() (suggests []prompt.Suggest) {
	s.RLock()
	for k := range s.names {
		suggests = append(suggests, prompt.Suggest{
			Text:        k,
			Description: k,
		})
	}
	s.RUnlock()
	sort.Sort(Suggests(suggests))
	return
}

func (s *Service) Names(group string) (suggests []prompt.Suggest) {
	s.RLock()
	conns := s.names[group]
	s.RUnlock()
	for _, conn := range conns {
		suggests = append(suggests, prompt.Suggest{
			Text:        conn.Name,
			Description: conn.Group + "." + conn.Name + "(" + conn.User + "@" + conn.Address + ")",
		})
	}
	sort.Sort(Suggests(suggests))
	return
}

func (s *Service) NamesWithPasswd(group string) (suggests []prompt.Suggest) {
	s.RLock()
	conns := s.names[group]
	s.RUnlock()
	for _, conn := range conns {
		suggests = append(suggests, prompt.Suggest{
			Text:        conn.Name,
			Description: conn.Group + "." + conn.Name + "(" + conn.User + ":" + conn.Password + "@" + conn.Address + ")",
		})
	}
	sort.Sort(Suggests(suggests))
	return
}

func (s *Service) Addrs() (suggests []prompt.Suggest) {
	s.RLock()
	for _, conns := range s.addrs {
		for _, conn := range conns {
			suggests = append(suggests, prompt.Suggest{
				Text:        conn.Address,
				Description: conn.Group + "." + conn.Name + "(" + conn.User + ":" + conn.Password + "@" + conn.Address + ")",
			})
		}
	}
	s.RUnlock()
	sort.Sort(Suggests(suggests))
	return
}

func (s *Service) Addr(addr string) (suggests []prompt.Suggest) {
	s.RLock()
	conns := s.addrs[addr]
	s.RUnlock()
	for _, conn := range conns {
		suggests = append(suggests, prompt.Suggest{
			Text:        "show",
			Description: conn.Group + "." + conn.Name + "(" + conn.User + ":" + conn.Password + "@" + conn.Address + ")",
		})
	}
	sort.Sort(Suggests(suggests))
	return
}

func (s *Service) ConnsByAddr(addr string) (conns []dao.SshConnInfo) {
	s.RLock()
	for _, conn := range s.addrs[addr] {
		conns = append(conns, conn)
	}
	s.RUnlock()
	sort.Sort(Conns(conns))
	return
}

type Conns []dao.SshConnInfo

func (c Conns) Len() int {
	return len(c)
}

func (c Conns) Less(i, j int) bool {
	return c[i].Group < c[j].Group || (c[i].Group == c[j].Group && c[i].Name < c[j].Name)
}

func (c Conns) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type Suggests []prompt.Suggest

func (s Suggests) Len() int {
	return len(s)
}

func (s Suggests) Less(i, j int) bool {
	return s[i].Text < s[j].Text || (s[i].Text == s[j].Text && s[i].Description < s[j].Description)
}

func (s Suggests) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
