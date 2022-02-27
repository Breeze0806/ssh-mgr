package cmdline

import (
	"strings"

	"github.com/Breeze0806/ssh-mgr/dao"
	"github.com/c-bata/go-prompt"
)

type ShowService interface {
	Add(conn dao.SshConnInfo) error
	Groups() (suggests []prompt.Suggest)
	Names(group string) (suggests []prompt.Suggest)
	Addrs() (suggests []prompt.Suggest)
	NamesWithPasswd(group string) (suggests []prompt.Suggest)
}

type Completer struct {
	cmdlines []prompt.Suggest
	service  ShowService
}

func NewCompleter(service ShowService) *Completer {
	return &Completer{
		cmdlines: []prompt.Suggest{
			{Text: "ssh", Description: "ssh group name"},
			{Text: "sftp", Description: "sftp group name"},
			{Text: "add", Description: "add group name"},
			{Text: "showAddr", Description: "show address"},
			{Text: "show", Description: "show group name"},
			{Text: "exit", Description: "exit ssh-mgr.prompt"},
		},
		service: service,
	}
}

func (c *Completer) Complete(d prompt.Document) (suggests []prompt.Suggest) {
	w := d.TextBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}

	block := strings.Split(w, " ")
	if len(block) == 1 {
		return prompt.FilterHasPrefix(c.cmdlines, block[0], true)
	}

	if len(block) == 2 {
		switch block[0] {
		case "ssh", "sftp", "add", "show":
			return prompt.FilterHasPrefix(c.service.Groups(), block[1], true)
		case "showAddr":
			return prompt.FilterHasPrefix(c.service.Addrs(), block[1], true)

		default:
			return []prompt.Suggest{}
		}
	}

	if len(block) == 3 {
		switch block[0] {
		case "ssh", "sftp", "add":
			return prompt.FilterHasPrefix(c.service.Names(block[1]), block[2], true)
		case "show":
			return prompt.FilterHasPrefix(c.service.NamesWithPasswd(block[1]), block[2], true)
		default:
			return []prompt.Suggest{}
		}
	}
	return []prompt.Suggest{}
}
