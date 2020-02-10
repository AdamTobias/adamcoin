package server

import (
	"blockchain/block"
	"blockchain/controller"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
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
	http.HandleFunc("/addtxn", s.AddTxnHandler())
	fmt.Println("Listening on 9421...")
	http.ListenAndServe(":9421", nil)
	return nil
}

func (s *Server) AddTxnHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txn, err := s.ParseAddTxnReq(r)
		fmt.Println("tried to parse txn with err", err)
		s.cont.AddTxn(txn)
	}
}

var InvalidReqErr = errors.New("invalid request")

func (s *Server) ParseAddTxnReq(r *http.Request) (block.Transaction, error) {
	txn := block.Transaction{}
	keys := r.URL.Query()
	var err error
	txn.Sender, err = parseSender(keys)
	if err != nil {
		return txn, err
	}
	txn.Receivers, err = parseReceivers(keys)
	if err != nil {
		return txn, err
	}
	txn.Amount, err = parseAmount(keys)
	if err != nil {
		return txn, err
	}
	txn.Signature, err = parseSignature(keys)
	if err != nil {
		return txn, err
	}

	//	msg, jerr := json.Marshal(txn)
	//	if jerr != nil {
	//		return txn, jerr
	//	}
	return txn, nil
}

func parseSender(keys url.Values) (*rsa.PublicKey, error) {
	sender := rsa.PublicKey{}
	if senderN := keys.Get("senderN"); senderN != "" {
		var j big.Int
		_, ok := j.SetString(senderN, 10)
		if !ok {
			return nil, errors.New("could not parse senderN")
		}
		sender.N = &j
	} else {
		return nil, errors.New("no senderN")
	}
	if senderE := keys.Get("senderE"); senderE != "" {
		k, err := strconv.Atoi(senderE)
		if err != nil {
			return nil, errors.New("could not parse senderE")
		}
		sender.E = k
	} else {
		return nil, errors.New("no senderE")
	}
	return &sender, nil
}

func parseReceivers(keys url.Values) (block.UnspentCoins, error) {
	us := make(block.UnspentCoins)
	recN := keys["recN"]
	recE := keys["recE"]
	recAmt := keys["recAmt"]
	if len(recE) == 0 || len(recN) != len(recE) || len(recN) != len(recAmt) {
		fmt.Println("flag 941")
		return us, errors.New("some receiver data missing")
	}
	for i, val := range recN {
		var n big.Int
		_, ok := n.SetString(val, 10)
		if !ok {
			fmt.Println("flag 169")
			return nil, errors.New("could not parse some receiverN")
		}
		e, err := strconv.Atoi(recE[i])
		if err != nil {
			fmt.Println("flag 182")
			return nil, errors.New("could not convert some receiverN to string")
		}
		pubkey := rsa.PublicKey{E: e, N: &n}
		name, _ := json.Marshal(pubkey)
		nameStr := string(name)
		if _, ok := us[nameStr]; !ok {
			us[nameStr] = 0
		}
		amt, err := strconv.Atoi(recAmt[i])
		if err != nil {
			fmt.Println("flag 984")
			return nil, errors.New("could not parse some receiverAmt")
		}
		us[nameStr] += amt
	}
	return us, nil
}
func parseSignature(keys url.Values) ([]byte, error) {
	sig := keys.Get("sig")
	if sig == "" {
		return []byte{}, errors.New("no sig found")
	}
	return []byte(sig), nil
}

func parseAmount(keys url.Values) (int, error) {
	amt := keys.Get("amt")
	a, err := strconv.Atoi(amt)
	if err != nil {
		return 0, errors.New("could not parse amt")
	}
	return a, nil
}
