package yarserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/renxiaotu/phppack"
	"net"
	"reflect"
	"strconv"
)

type Request struct {
	header header
	body   body
}

func (r Request) GetParam(p interface{}) error {
	value := reflect.ValueOf(p)
	for value.Kind() == reflect.Ptr {
		next := value.Elem().Kind()
		if next == reflect.Struct || next == reflect.Ptr {
			value = value.Elem()
		} else {
			break
		}
	}
	plan := 0
	actual := len(r.body.Param)
	fieldList := make([]int, 0)

	for i := 0; i < value.NumField(); i++ {
		if value.Field(i).CanSet() {
			fieldList = append(fieldList, i)
			plan++
		}
	}

	if plan != actual {
		return errors.New(packageName + ":inconsistent number of parameters,the plan is " + strconv.Itoa(plan) + ", the actual is " + strconv.Itoa(actual))
	}

	for i := 0; i < plan; i++ {
		v := value.Field(i)
		v.Set(reflect.ValueOf(r.body.Param[i]).Convert(v.Type()))
	}
	return nil
}

//解包请求
func parseRequest(conn net.Conn) (Request, error) {
	r := Request{}
	err := getHeader(conn, &r.header)
	if err != nil {
		return r, err
	}

	err = getBody(conn, r.header, &r.body)
	if err != nil {
		return r, err
	}
	return r, nil
}

func getHeader(conn net.Conn, h *header) error {
	hand := header{}
	handByte := make([]byte, headerSize)
	_, err := conn.Read(handByte)
	if err != nil {
		return errors.New(packageName + ":" + err.Error())
	}

	err = phppack.UnpackByStruct(&hand, handByte)
	if err != nil {
		return err
	}

	fmt.Println(handByte)
	fmt.Println(hand)

	if hand.MagicNum != magicNum {
		return errors.New(packageName + ": illegal Yar RPC request")
	}

	hand.BodyLen -= 8
	*h = hand
	return nil
}

func getBody(conn net.Conn, h header, b *body) error {
	data := body{}
	bodyByte := make([]byte, h.BodyLen)
	_, err := conn.Read(bodyByte)
	if err != nil {
		return errors.New(packageName + ": " + err.Error())
	}
	fmt.Println(string(bodyByte))
	switch h.Packager {
	case `JSON`:
		err = json.Unmarshal(bodyByte, &data)
	default:
		err = errors.New(`unsupported packager "` + h.Packager + `"`)
	}
	if err != nil {
		return err
	}

	*b = data
	return nil
}
