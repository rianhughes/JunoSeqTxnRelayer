package main

import (
	"context"
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/rpc"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

var (
	fullNodeBase = ""
	seqBase      = "http://localhost:8545"
)

func main() {

	c, err := ethrpc.DialContext(context.Background(), fullNodeBase)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the full node: %v", err))
	}
	fullNode := rpc.NewProvider(c)
	log.Default().Printf("Connected to the full node")

	// connect to sequencer
	c, err = ethrpc.DialContext(context.Background(), seqBase)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the sequencer: %v", err))
	}
	seq := rpc.NewProvider(c)
	log.Default().Printf("Connected to the sequencer")

	for blockN := uint64(0); blockN <= 1; blockN++ {

		log.Default().Printf("Querying full node for block %d\n", blockN)

		// Get txns from full node
		block, err := fullNode.BlockWithTxs(context.Background(), rpc.BlockID{Number: &blockN})
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to get the block: %v", err))
		}
		blockCast := block.(*rpc.Block)

		log.Default().Printf("Got block %d with %d transactions\n", blockN, len(blockCast.Transactions))
		log.Default().Printf("-- Sending transactions to sequencer\n")

		// Send txns to sequencer
		for _, txn := range blockCast.Transactions {

			var err error

			switch t := txn.(type) {
			case *rpc.BlockInvokeTxnV0:
				_, err = seq.AddInvokeTransaction(context.Background(), t.InvokeTxnV0)
			case *rpc.BlockInvokeTxnV1:
				_, err = seq.AddInvokeTransaction(context.Background(), t.InvokeTxnV1)
			case *rpc.BlockDeclareTxnV0:
				_, err = seq.AddDeclareTransaction(context.Background(), t.DeclareTxnV0)
			case *rpc.BlockDeclareTxnV1:
				_, err = seq.AddDeclareTransaction(context.Background(), t.DeclareTxnV1)
			case *rpc.BlockDeclareTxnV2:
				_, err = seq.AddDeclareTransaction(context.Background(), t.DeclareTxnV2)
			case *rpc.BlockDeployAccountTxn:
				_, err = seq.AddDeployAccountTransaction(context.Background(), t.DeployAccountTxn)
			case *rpc.BlockDeployTxn:
				log.Default().Printf("-- Skipping deploy transaction since not supported in latest rpc version\n")
			}
			if err != nil {
				panic(fmt.Sprintf("-- Failed to send the transaction: %v", err))
			}
		}
		log.Default().Printf("-- Sent %d transactions to sequencer\n", len(blockCast.Transactions))
	}

}
