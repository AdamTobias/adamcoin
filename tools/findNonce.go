package main

import (
	"blockchain/block"
	"crypto"
	_ "crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	mrand "math/rand"
	"strconv"
	"time"
)

const target = "141bc330"

func main() {
	pb := block.Block{}
	blk := block.Block{
		Header: block.Header{
			PrevBlk: &pb,
			Hash:    []byte{},
		},
		Body: block.Body{
			Txns: []block.Transaction{},
		},
	}

	exp, err1 := strconv.Atoi(target[:2])
	if err1 != nil {
		fmt.Println("err1 is", err1)
	}
	val, err2 := strconv.ParseUint(target[2:], 16, 32)
	if err2 != nil {
		fmt.Println("err2 is", err2)
	}

	valj32 := uint64(val)
	final := uint64(math.Pow10(exp-3)) * valj32

	times := 0
	valx := final
	var hashed []byte
	var r1 string
	for times < 10000 && valx >= final {
		time.Sleep(2 * time.Millisecond)
		seed := time.Now().UnixNano()
		mrand.Seed(seed)
		r1 = strconv.FormatInt(mrand.Int63(), 10)
		blk.Header.Nonce = r1
		msg, _ := json.Marshal(blk)

		newHash := crypto.SHA256
		pssh := newHash.New()
		pssh.Write(msg)
		hashed = pssh.Sum(nil)

		valx = binary.BigEndian.Uint64(hashed)
		times++
	}
	if valx < final {
		fmt.Printf("Solved it! %v < %v!\n", valx, final)
		fmt.Println("nonce is", r1)
	} else {
		fmt.Println("Didnt solve it T.T")
	}
	fmt.Println("times ran = ", times)
}
