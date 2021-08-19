package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Block struct {
	Index        int
	Timestamp    time.Time
	Proof        int
	PreviousHash string
}

type Blockchain []*Block

func (bc Blockchain) CreateBlock(proof int, previousHash string) *Block {
	block := &Block{
		Index:        len(bc) + 1,
		Timestamp:    time.Now(),
		Proof:        proof,
		PreviousHash: previousHash,
	}

	return block
}

func (bc Blockchain) GetLastBlock() *Block {
	return bc[len(bc)-1]
}

func ProofOfWork(previousProof int) int {
	newProof := 1
	checkProof := false
	for !checkProof {
		n := math.Pow(float64(newProof), 2) - math.Pow(float64(previousProof), 2)
		hashOperation := sha256.Sum256([]byte(fmt.Sprintf("%f", n)))

		trunc := fmt.Sprintf("%x", hashOperation)[0:4]
		if trunc == "0000" {
			checkProof = true

			return newProof
		}

		newProof++
	}

	return newProof
}

func (b *Block) Hash() (string, error) {
	encodedBlock, err := json.Marshal(b)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(encodedBlock))

	return fmt.Sprintf("%x", hash), nil
}

func (bc Blockchain) IsValid() (bool, error) {
	for k, v := range bc {
		hash, err := v.Hash()
		if err != nil {
			return false, err
		}

		if hash != bc[k+1].PreviousHash {
			return false, nil
		}
	}

	return true, nil
}

func main() {
	blockchain := Blockchain{
		{
			Index:        1,
			Timestamp:    time.Now(),
			Proof:        1,
			PreviousHash: "0",
		},
	}

	e := echo.New()
	e.GET("/mine_block", func(c echo.Context) error {
		lastBlock := blockchain.GetLastBlock()
		proof := ProofOfWork(lastBlock.Proof)

		previousHash, err := lastBlock.Hash()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "cannot hash last block")
		}
		newBlock := blockchain.CreateBlock(proof, previousHash)
		blockchain = append(blockchain, newBlock)

		return c.JSON(http.StatusOK, newBlock)
	})

	e.GET("/get_chain", func(c echo.Context) error {
		return c.JSON(http.StatusOK, blockchain)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
