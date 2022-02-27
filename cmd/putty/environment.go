package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Breeze0806/ssh-mgr/dao"
	"github.com/Breeze0806/ssh-mgr/dao/mapper"
	"github.com/Breeze0806/ssh-mgr/services/pass"
	"github.com/Breeze0806/ssh-mgr/services/show"
	"github.com/c-bata/go-prompt"
	"github.com/howeyc/gopass"
)

var (
	configFile = flag.String("c", "config.json", "config file")
	sshFlag    = flag.Bool("ssh", false, "config file")
	user       = flag.String("l", "", "user")
	passwd     = flag.String("pw", "", "passwd")
)

type ConnMapper interface {
	Read(group, name string) (conn dao.SshConnInfo, err error)
	Write(conn dao.SshConnInfo) (err error)
}

type Environment struct {
	err        error
	conf       *Config
	password   string
	addr       string
	connMapper ConnMapper
	passMapper *mapper.PassMapper

	passService *pass.Service
	showService *show.Service

	prompt *prompt.Prompt
}

func NewEnvironment(filename string) (e *Environment, err error) {
	e = &Environment{}
	e.conf, err = NewConfig(filename)
	if err != nil {
		return
	}
	start := strings.LastIndex(*passwd, "@")
	s := (*passwd)[start+1:]
	ss := strings.Split(s, ":")
	e.addr = ss[0] + ":" + ss[2]
	fmt.Println(e.addr)
	return
}

func (e *Environment) Build() error {
	return e.initSSH().initPassword().initMappers().initServices().initApis().err
}

func (e *Environment) initMappers() *Environment {
	if e.err != nil {
		return e
	}
	if e.conf.IsEncrypted {
		e.connMapper = mapper.NewConnSM4CryptoMappar(e.conf.Source, e.password)
	} else {
		e.connMapper = mapper.NewConnMapper(e.conf.Source)
	}
	return e
}

func (e *Environment) initPassword() *Environment {
	if e.conf.IsEncrypted {
		e.passMapper = mapper.NewPassMapper(e.conf.PasswordFile())
		var data1 []byte
		data1, e.err = gopass.GetPasswdPrompt("please input password:", false, os.Stdin, os.Stdout)
		if e.err != nil {
			e.err = fmt.Errorf("input password fail error: %v", e.err)
			return e
		}
		e.password = string(data1)
	}
	return e
}

func (e *Environment) initServices() *Environment {
	if e.err != nil {
		return e
	}

	if e.conf.IsEncrypted {
		e.passService = pass.NewService(e.passMapper)
		if e.passMapper.HasPassword() {
			e.err = e.passMapper.Match(e.password)
		}
	}
	if e.err != nil {
		return e
	}
	e.showService = show.NewService(e.connMapper)
	e.err = e.showService.Init(e.conf.Source)
	return e
}

func (e *Environment) initApis() *Environment {
	if e.err != nil {
		return e
	}
	e.prompt = prompt.New(e.Execute, e.Complete,
		prompt.OptionTitle(e.addr),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray))
	return e
}

func (e *Environment) initSSH() *Environment {
	if e.err != nil {
		return e
	}
	cmd := exec.Command(e.conf.Ssh, "-ssh", "-l", *user, "-pw", *passwd, os.Args[len(os.Args)-2], os.Args[len(os.Args)-1])
	e.err = cmd.Start()

	return e
}

func (e *Environment) Run() error {
	e.prompt.Run()
	return nil
}

func (e *Environment) Complete(d prompt.Document) (suggests []prompt.Suggest) {
	w := d.TextBeforeCursor()
	if w == "" {
		return
	}
	suggests = e.showService.Addr(e.addr)
	suggests = append(suggests, prompt.Suggest{
		Text:        "exit",
		Description: "exit putty",
	})
	return prompt.FilterHasPrefix(suggests, w, true)
}

func (e *Environment) Execute(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")
	switch blocks[0] {
	case "":
	case "show":
		conns := e.showService.ConnsByAddr(e.addr)
		for _, v := range conns {
			fmt.Println(v.User)
		}
	case "exit":
		fmt.Println("Bye!")
		os.Exit(0)
	default:
		fmt.Println("no such command!")
	}
}
