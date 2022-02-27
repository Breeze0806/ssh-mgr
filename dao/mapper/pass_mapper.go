package mapper

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/tjfoc/gmsm/sm3"
)

type PassMapper struct {
	filename string
}

func NewPassMapper(filename string) *PassMapper {
	return &PassMapper{
		filename: filename,
	}
}

func (p *PassMapper) HasPassword() bool {
	if f, err := os.Open(p.filename); err == nil {
		f.Close()
		return true
	}
	return false
}

func (p *PassMapper) Write(password string) (err error) {
	data := sm3.Sm3Sum([]byte(password))
	return ioutil.WriteFile(p.filename, data, 0664)
}

func (p *PassMapper) Match(password string) (err error) {
	var data []byte
	data, err = ioutil.ReadFile(p.filename)
	if err != nil {
		return err
	}
	data1 := sm3.Sm3Sum([]byte(password))
	if reflect.DeepEqual(data, data1) {
		return nil
	}
	return fmt.Errorf("password does not match")
}
