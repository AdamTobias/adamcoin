package server

import (
	"blockchain/block"
	"blockchain/controller"
	"net/http"
)

type Server struct {
	cont controller.Controller
}

//func (s Server) PostTransaction(t transaction) error {
// verify transaction
// if valid, add to openTransactions
//	return nil
//}

func (s Server) InitializePeers(pl []string) error {
	// for each peer in list, establish connection (get current block?)
	return nil
}

func (s Server) PostBlock(b block.Block) {
	// receive incoming block.  verify it then add to chain
}

func (s *Server) AddPeer(p string) error {
	// add peer to peer list (any verification required?
	s.cont.AddPeer(p)
	return nil
}

func (s *Server) InitializeServer(peer string) error {
	// pass in initial peer(s?)
	// get chain from peer
	c := controller.NewController(peer)
	// init chain
	// keep chain locally?
	s.cont = c
	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("you got me!"))
	}
	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(":9421", nil)
	return nil
}
