package dao

import (
	"net/url"
)

type SshConnInfo struct {
	Name     string `json:"-"`
	Group    string `json:"-"`
	Address  string `json:"address"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// SftpURL returns an sftp:// URL for the saved connection.
//
// The userinfo is escaped using a strict percent-encoding that only leaves
// the URI unreserved set untouched (A-Z / a-z / 0-9 / "-" / "." / "_" /
// "~"). Go's stdlib helpers are not aggressive enough for this purpose:
//
//   - url.UserPassword follows RFC 3986 and leaves sub-delims such as
//     '$', '!', '+', '&' un-escaped inside userinfo. WinSCP tolerates
//     this, but FileZilla on Linux interprets some of those bytes
//     specially (notably '$' as a shell variable) and refuses to connect.
//   - url.PathEscape is even worse for userinfo: it deliberately leaves
//     '@', ':', ';' un-escaped, which produces a URL whose userinfo/host
//     boundary is ambiguous as soon as the password contains '@'.
//
// The result of escapeUserInfo is still a syntactically valid URL: any
// conforming SFTP client percent-decodes the userinfo back to its
// original bytes before using it.
func (c *SshConnInfo) SftpURL() string {
	return "sftp://" + escapeUserInfo(c.User) + ":" + escapeUserInfo(c.Password) + "@" + c.Address
}

func (c *SshConnInfo) SshURL() string {
	URL := &url.URL{
		Scheme: "",
		User:   url.User(c.User),
		Host:   c.Address,
	}

	return URL.String()[2:]
}

// escapeUserInfo percent-encodes every byte in s that is not in the URI
// unreserved set (A-Z / a-z / 0-9 / "-" / "." / "_" / "~"). The result is
// safe to drop between "://" and "@" in a URL.
func escapeUserInfo(s string) string {
	const upperhex = "0123456789ABCDEF"
	needs := 0
	for i := 0; i < len(s); i++ {
		if !isUnreserved(s[i]) {
			needs++
		}
	}
	if needs == 0 {
		return s
	}
	b := make([]byte, 0, len(s)+2*needs)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if isUnreserved(c) {
			b = append(b, c)
		} else {
			b = append(b, '%', upperhex[c>>4], upperhex[c&15])
		}
	}
	return string(b)
}

func isUnreserved(c byte) bool {
	return ('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z') ||
		('0' <= c && c <= '9') ||
		c == '-' || c == '.' || c == '_' || c == '~'
}
