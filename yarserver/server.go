package yarserver

import (
	"errors"
	"fmt"
	"net"
)

type onHandle func(request Request, response Response)
type errorHandle func(error)

type address struct {
	network string
	address string
}

type YarServer struct {
	run       bool
	addr      address
	ln        net.Listener
	onList    map[string]*onHandle
	errorList *[]*errorHandle
}

func (ys YarServer) Error(h errorHandle) {
	*(ys.errorList) = append(*(ys.errorList), &h)
}

func (ys YarServer) On(methodName string, h onHandle) {
	ys.onList[methodName] = &h
}

func (ys YarServer) error(err error) {
	fmt.Println(ys.errorList)
	for i := 0; i < len(*(ys.errorList)); i++ {
		fmt.Println(err.Error())
		(*(*(ys.errorList))[i])(err)
	}
}

func (ys YarServer) Run() error {
	err := errors.New("")
	ys.ln, err = net.Listen(ys.addr.network, ys.addr.address)
	if err != nil {
		return err
	}
	ys.run = true

	fmt.Println("yarServer:  " + ys.addr.network + "://" + ys.addr.address)

	for ys.run {
		conn, err := ys.ln.Accept()
		if err != nil {
			ys.error(err)
		}
		go handle(conn, ys)
	}
	return nil
}

func (ys YarServer) Close() error {
	ys.run = false
	return ys.ln.Close()
}

func New(network string, addr string) YarServer {
	a := address{network, addr}
	oh := make(map[string]*onHandle)
	eh := make([]*errorHandle, 0)
	var l net.Listener
	return YarServer{false, a, l, oh, &eh}
}

func handle(conn net.Conn, ys YarServer) {
	//自动关闭
	defer func() {
		err := conn.Close()
		if err != nil {
			ys.error(errors.New("net.Conn error: " + err.Error()))
			return
		}
	}()

	request, err := parseRequest(conn)
	if err != nil {
		ys.error(err)
		return
	}

	response := Response{conn}

	f, ok := ys.onList[request.body.Method]
	if !ok {
		ys.error(errors.New(packageName + ":method does not exist"))
	}
	(*f)(request, response)
}
