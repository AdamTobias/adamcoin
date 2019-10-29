package server

import "net/http"

type Server struct {
	peers []string
}

type ServInterface interface{
	PostTransaction(t transaction)
	EmitTransaction(t transaction)
	InitializePeers(peerList []string) error
	PostBlock(b block)
	EmitBlock(b block)
	AddPeer(p peer)
}

type transaction string
type block string
type peer string

func (s Server) PostTransaction(t transaction) error {
	// verify transaction
	// if valid, add to openTransactions
	return nil
}

func (s Server) EmitTransaction(t transaction) error {
	// emit transaction to all peers in peer list
	return nil
}

func (s Server) InitializePeers(pl []string) error {
	// for each peer in list, establish connection (get current block?)
	return nil
}

func (s Server) PostBlock(b block) {
	// receive incoming block.  verify it then add to chain
}

func (s Server) EmitBlock(b block) {
	// send block to all peers
}

func (s Server) AddPeer(p peer) {
	// add peer to peer list (any verification required?
}

func (s Server) InitializeServer() {
	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("you got me!"))
	}
	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(":9421", nil)
}
