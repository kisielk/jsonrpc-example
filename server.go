package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type JSONRPCServer struct {
	*rpc.Server
}

func NewJSONRPCServer() *JSONRPCServer {
	return &JSONRPCServer{rpc.NewServer()}
}

func (s *JSONRPCServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("got a request")
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Print("rpc hijacking ", req.RemoteAddr, ": ", err.Error())
		return
	}
	io.WriteString(conn, "HTTP/1.0 200 Connected to Go JSON-RPC\n\n")
	codec := jsonrpc.NewServerCodec(conn)
	log.Println("ServeCodec")
	s.Server.ServeCodec(codec)
	log.Println("finished serving request")
}

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func main() {
	s := NewJSONRPCServer()
	arith := new(Arith)
	s.Register(arith)
	http.Handle("/rpc", s)
	http.ListenAndServe(":8080", nil)
}
