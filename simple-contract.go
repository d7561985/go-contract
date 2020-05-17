package main

import (
	"github.com/d7561985/go-contract/contracts/leveldb"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	simpleContract := new(leveldb.SimpleQueueContract)

	cc, err := contractapi.NewChaincode(simpleContract)

	if err != nil {
		panic(err.Error())
	}

	if err := cc.Start(); err != nil {
		panic(err.Error())
	}
}
