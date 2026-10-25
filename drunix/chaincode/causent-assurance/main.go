package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func main() {
	chaincode, err := contractapi.NewChaincode(&AssuranceContract{})
	if err != nil {
		log.Panicf("error creating causent-assurance chaincode: %v", err)
	}
	if err := chaincode.Start(); err != nil {
		log.Panicf("error starting causent-assurance chaincode: %v", err)
	}
}
