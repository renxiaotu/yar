package yarserver

import (
	"encoding/json"
	"github.com/renxiaotu/phppack"
	"net"
)

type Response struct {
	conn net.Conn
}

func (r Response) JsonSuccess(param interface{}) error {
	return jsonReply(r, 0, param)
}

func (r Response) JsonFail(param interface{}) error {
	return jsonReply(r, 1, param)
}

func jsonReply(r Response, s int, p interface{}) error {
	res := make(map[string]interface{})
	res[`i`] = 0
	res[`s`] = s
	res[`r`] = p
	b, err := json.Marshal(res)
	if err != nil {
		return err
	}
	h, err := packHeader(uint32(len(b)))
	if err != nil {
		return err
	}
	h = append(h, b...)
	_, err = r.conn.Write(h)
	if err != nil {
		return err
	}
	return nil
}

//打包响应头
func packHeader(len uint32) ([]byte, error) {
	return phppack.PackByStruct(header{
		Id:       0,
		Version:  0,
		MagicNum: magicNum,
		Reserved: 0,
		Provider: "Yar GO TCP Server",
		Token:    "",
		BodyLen:  len + 8,
		Packager: packager,
	})
}
