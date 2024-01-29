package engine

import (
	"case/blockchain"
	"case/utils"
	"fmt"
	"strconv"
)

// Engine responsible for processing command line arguments
type Engine struct {
	// bc *blockchain.Blockchain
}

func (engine *Engine) creatBlock(address string) {
	bc := blockchain.CreateBlockchain(address)
	bc.Db.Close()
	fmt.Println("Done!")
}
func (engine *Engine) getBalance(address string) {
	bc := blockchain.NewBlockchain(address)
	defer bc.Db.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
func (engine *Engine) printChain() {
	bc := blockchain.NewBlockchain("")
	defer bc.Db.Close()

	bci := bc.Iterator()
	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		// fmt.Printf("Transactions: %s\n", block.Transactions)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
func (eng *Engine) send(from, to string, amount int) {
	bc := blockchain.NewBlockchain(from)
	defer bc.Db.Close()

	tx := blockchain.NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!")
}

func Run() {
	eng := Engine{}
	if !utils.DbExists() {
		eng.creatBlock("xiasang")
	}

	eng.getBalance("xiasang")
	eng.send("xiasang", "ikun1", 10)
	eng.getBalance("xiasang")
	eng.send("xiasang", "ikun2", 10)
	eng.getBalance("xiasang")
	eng.getBalance("ikun1")
	eng.getBalance("ikun2")
	eng.printChain()
}
