package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/sprt/zoossh"
)

type Args struct {
	Filename string
}

type Zoossh int

type Result []*zoossh.RouterStatus

func (t *Zoossh) ConsensusRouters(r *http.Request, args *Args, result *Result) error {
	log.Println("Parsing", args.Filename)

	consensus, err := zoossh.ParseConsensusFile(args.Filename)
	if err != nil {
		return err
	}

	routerStatuses := make([]*zoossh.RouterStatus, 0, len(consensus.RouterStatuses))
	for _, rs := range consensus.RouterStatuses {
		routerStatuses = append(routerStatuses, rs())
	}

	*result = routerStatuses
	return nil
}

func main() {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

	arith := new(Zoossh)
	s.RegisterService(arith, "")

	r := mux.NewRouter()
	r.Handle("/rpc", s)

	log.Println("Listening...")
	http.ListenAndServe(":1234", r)
}
