package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

// Take the data form the block
// create a counter (nonce) which stats at 0
// create a hash of the data plus the counter
// check the hash to see if it meets a set of requirements
// Requirements :
// 1. The First few bytes must contain 0s

const Difficulty = 12 // in real life implementation it will incrementing by period of time when the blockchain is getting larger

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	return bytes.Join([][]byte{pow.Block.PrevHash, pow.Block.Data, ToHex(int64(nonce)), ToHex(int64(Difficulty))}, []byte{})
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var (
		intHash big.Int
		hash    [32]byte
	)

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		}
		nonce++
	}

	fmt.Println()
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}
