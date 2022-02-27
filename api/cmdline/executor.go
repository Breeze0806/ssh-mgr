package cmdline

import (
	"fmt"
	"os"
	"strings"

	"github.com/Breeze0806/ssh-mgr/dao"
	"github.com/howeyc/gopass"
)

type CmdServer interface {
	StartPragram(group, name string) (err error)
}

type AddServer interface {
	Add(conn dao.SshConnInfo) error
}

type Executor struct {
	ssh  CmdServer
	sftp CmdServer
	ini  AddServer
	show AddServer
}

func NewExecutor(ini, show AddServer, ssh CmdServer, sftp CmdServer) *Executor {
	return &Executor{
		ini:  ini,
		ssh:  ssh,
		sftp: sftp,
		show: show,
	}
}

func (e *Executor) Execute(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")
	switch blocks[0] {
	case "ssh":
		if len(blocks) != 3 {
			fmt.Println("ssh format is invalid")
			return
		}
		if blocks[1] == "" || blocks[2] == "" {
			fmt.Println("group or name is empty")
			return
		}
		if err := e.ssh.StartPragram(blocks[1], blocks[2]); err != nil {
			fmt.Println("ssh start fail. err:", err)
			return
		}
	case "sftp":
		if len(blocks) != 3 {
			fmt.Println("sftp format is invalid")
			return
		}
		if blocks[1] == "" || blocks[2] == "" {
			fmt.Println("group or name is empty")
			return
		}
		if err := e.sftp.StartPragram(blocks[1], blocks[2]); err != nil {
			fmt.Println("sftp start fail. err:", err)
			return
		}
	case "add":
		if len(blocks) != 3 {
			fmt.Println("add format is invalid, must be 4 blocks")
			return
		}

		if blocks[1] == "" || blocks[2] == "" {
			fmt.Println("group or name is empty")
			return
		}

		conn := dao.SshConnInfo{
			Group: blocks[1],
			Name:  blocks[2],
		}

		fmt.Print("please input ssh address:")
		fmt.Scanln(&conn.Address)
		if len(strings.Split(conn.Address, ":")) == 1 {
			conn.Address += ":22"
		}

		fmt.Print("please input ssh user:")
		fmt.Scanln(&conn.User)

		data, err := gopass.GetPasswdPrompt("please input ssh password:", true, os.Stdin, os.Stdout)
		if err != nil {
			fmt.Println("input password fail, error:", err)
			return
		}
		conn.Password = string(data)
		if err := e.ini.Add(conn); err != nil {
			fmt.Println("save ssh info fail, error:", err)
			return
		}
		e.show.Add(conn)

		fmt.Println("add success！")
	case "show", "showAddr", "":
	case "exit":
		fmt.Println("Bye!")
		Input()
		os.Exit(0)
	default:
		fmt.Println("no such command!")
	}
}

func Input() {
	fmt.Println("press return")
	input := ""
	fmt.Scanf("%s", &input)
}
