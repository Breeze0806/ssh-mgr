package mapper

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Breeze0806/ssh-mgr/dao"
	"github.com/tjfoc/gmsm/sm4"
)

type ConnSM4CryptoMappar struct {
	*ConnMapper
	key []byte
}

func NewConnSM4CryptoMappar(sourcePath, key string) *ConnSM4CryptoMappar {
	hash := md5.New().Sum([]byte(key))
	data := make([]byte, 16)
	copy(data, hash)
	return &ConnSM4CryptoMappar{
		ConnMapper: NewConnMapper(sourcePath),
		key:        data,
	}
}

func (c *ConnSM4CryptoMappar) Read(group, name string) (conn dao.SshConnInfo, err error) {
	conn, err = c.ConnMapper.Read(group, name)
	if err != nil {
		return
	}
	var data []byte
	data, err = hex.DecodeString(conn.User)
	if err != nil {
		return
	}
	data, err = sm4.Sm4Ecb(c.key, data, false)
	if err != nil {
		return
	}
	conn.User = string(data)
	data, err = hex.DecodeString(conn.Password)
	if err != nil {
		return
	}
	data, err = sm4.Sm4Ecb(c.key, data, false)
	if err != nil {
		return
	}
	conn.Password = string(data)
	return
}

func (c *ConnSM4CryptoMappar) Write(conn dao.SshConnInfo) (err error) {
	var data []byte

	data, err = sm4.Sm4Ecb(c.key, []byte(conn.User), true)
	if err != nil {
		return
	}
	conn.User = hex.EncodeToString(data)
	data, err = sm4.Sm4Ecb(c.key, []byte(conn.Password), true)
	if err != nil {
		return
	}
	conn.Password = hex.EncodeToString(data)
	return c.ConnMapper.Write(conn)
}

type ConnMapper struct {
	sourcePath string
}

func NewConnMapper(sourcePath string) *ConnMapper {
	return &ConnMapper{
		sourcePath: sourcePath,
	}
}

func (c *ConnMapper) Read(group, name string) (conn dao.SshConnInfo, err error) {
	filename := filepath.Join(c.sourcePath, group, name+".json")
	if group == "" {
		filename = filepath.Join(c.sourcePath, name+".json")
	}
	var data []byte
	data, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &conn)
	if err != nil {
		return
	}

	conn.Group = group
	conn.Name = name
	return
}

func (c *ConnMapper) Write(conn dao.SshConnInfo) (err error) {
	var data []byte

	data, err = json.MarshalIndent(conn, "", "    ")
	if err != nil {
		return
	}

	filename := filepath.Join(c.sourcePath, conn.Group, conn.Name+".json")
	if conn.Group == "" {
		filename = filepath.Join(c.sourcePath, conn.Name+".json")
	} else {
		os.Mkdir(filepath.Join(c.sourcePath, conn.Group), os.ModePerm)
	}

	return ioutil.WriteFile(filename, data, os.ModePerm)
}
