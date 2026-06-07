package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/Breeze0806/ssh-mgr/api/cmdline"
	"github.com/Breeze0806/ssh-mgr/dao"
	"github.com/Breeze0806/ssh-mgr/dao/mapper"
	"github.com/Breeze0806/ssh-mgr/services/ini"
	"github.com/Breeze0806/ssh-mgr/services/pass"
	"github.com/Breeze0806/ssh-mgr/services/sftp"
	"github.com/Breeze0806/ssh-mgr/services/show"
	"github.com/Breeze0806/ssh-mgr/services/ssh"
	"github.com/c-bata/go-prompt"
	"golang.org/x/term"
)

type ConnMapper interface {
	Read(group, name string) (conn dao.SshConnInfo, err error)
	Write(conn dao.SshConnInfo) (err error)
}

type Environment struct {
	err      error
	conf     *Config
	password string

	connMapper ConnMapper
	passMapper *mapper.PassMapper

	sftpService *sftp.Service
	sshService  *ssh.Service
	initService *ini.Service
	passService *pass.Service
	showService *show.Service

	executor  *cmdline.Executor
	completer *cmdline.Completer
	prompt    *prompt.Prompt
}

func NewEnvironment(filename string) (e *Environment, err error) {
	e = &Environment{}
	e.conf, err = NewConfig(filename)
	if err != nil {
		return
	}
	return
}

func (e *Environment) Build() error {
	return e.InitPassword().InitMappers().InitServices().InitApis().err
}

func (e *Environment) InitMappers() *Environment {
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

func (e *Environment) InitPassword() *Environment {
	if e.conf.IsEncrypted {
		e.passMapper = mapper.NewPassMapper(e.conf.PasswordFile())
		e.passService = pass.NewService(e.passMapper)
		for {
			if err := e.tryPassword(); err != nil {
				fmt.Println(err)
				continue
			}
			return e
		}
	}
	return e
}

func (e *Environment) tryPassword() (err error) {
	// Write prompts to stderr so they don't interfere with go-prompt's
	// terminal renderer. x/term.ReadPassword disables echo on the terminal
	// and restores it on return — unlike the unmaintained gopass library
	// (which uses pkg/term) this is the modern, supported path.
	var data1 []byte
	data1, err = readPasswordPrompt("please input password:")
	if err != nil {
		err = fmt.Errorf("input password fail error: %v", err)
		return
	}
	if !e.passMapper.HasPassword() {
		var data2 []byte
		data2, err = readPasswordPrompt("please confirm password:")
		if err != nil {
			err = fmt.Errorf("confirm password fail error: %v", err)
			return
		}
		if !reflect.DeepEqual(data1, data2) {
			err = fmt.Errorf("password does not match")
			return
		}
	}
	e.password = string(data1)

	if e.passMapper.HasPassword() {
		err = e.passMapper.Match(e.password)
	} else {
		err = e.passMapper.Write(e.password)
	}

	return
}

// readPasswordPrompt prints `prompt` to stderr (so the go-prompt renderer
// that owns stdout is not disturbed) and reads a password from stdin
// without echoing. The terminal is restored to its original state on
// return, including on error paths, by virtue of term.ReadPassword's
// internal save/restore of the term state.
func readPasswordPrompt(prompt string) ([]byte, error) {
	if _, err := fmt.Fprint(os.Stderr, prompt); err != nil {
		return nil, err
	}
	defer fmt.Fprintln(os.Stderr)
	return term.ReadPassword(int(os.Stdin.Fd()))
}

func (e *Environment) InitServices() *Environment {
	if e.err != nil {
		return e
	}

	if e.err != nil {
		return e
	}
	e.sftpService = sftp.NewService(e.conf.Sftp, e.connMapper)
	e.sshService = ssh.NewService(e.conf.Ssh, e.connMapper)
	e.initService = ini.NewService(e.connMapper)
	e.showService = show.NewService(e.connMapper)
	e.err = e.showService.Init(e.conf.Source)
	return e
}

func (e *Environment) InitApis() *Environment {
	if e.err != nil {
		return e
	}
	e.completer = cmdline.NewCompleter(e.showService)
	e.executor = cmdline.NewExecutor(e.initService, e.showService, e.sshService, e.sftpService)
	e.prompt = prompt.New(e.executor.Execute,
		e.completer.Complete,
		prompt.OptionTitle("ssh-mgr-prompt"),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
	)
	return e
}

func (e *Environment) Run() error {
	e.prompt.Run()
	return nil
}
