package miner

import (
	"blockchain/block"
	"crypto"
	_ "crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	mrand "math/rand"
	"strconv"
	"time"
)

type Miner struct {
	Kill chan int
}

func (m Miner) Mine(blk block.Block, kill chan int) chan block.Block {
	response := make(chan block.Block)
	targetInt, _ := targetToInt("111bc330")
	go func() {
		for {
			select {
			case <-kill:
				fmt.Println("killing the miner itself")
				return
			default:
				seed := time.Now().UnixNano()
				mrand.Seed(seed)
				nonce := strconv.FormatInt(mrand.Int63(), 10)
				blk.Header.Nonce = nonce // Is this some bad 'pass by reference' shit?
				msg, _ := json.Marshal(blk)
				newHash := crypto.SHA256
				pssh := newHash.New()
				pssh.Write(msg)
				hashed := pssh.Sum(nil)

				guess := binary.BigEndian.Uint64(hashed)
				fmt.Println("guessing")
				if guess < targetInt {
					response <- blk
					return
				}
				// TODO -- should we do more than one mine before checking if the channel is closed?
			}
		}
	}()
	return response
}

func NewMiner() Miner {
	return Miner{Kill: make(chan int)}
}

func targetToInt(target string) (uint64, error) {
	return ParseTarget(target)
}

func ParseTarget(t string) (uint64, error) {
	exp, err := strconv.Atoi(t[:2])
	if err != nil {
		return 0, errors.New("could not parse target -- bad exponent")
	}
	val, err := strconv.ParseUint(t[2:], 16, 32)
	if err != nil {
		return 0, errors.New("could not parse target -- bad value")
	}
	v := uint64(val)
	res := uint64(math.Pow10(exp-3)) * v
	return res, nil
}
