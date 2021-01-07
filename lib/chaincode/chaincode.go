/**
 * 
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * 
 */
package chaincode

import (
	"fmt"

	"example.com/fffffefe/lib/util"
	"example.com/fffffefe/src"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type ChainCode struct {
}

//Init Function Executes only once while initializing or upgrading chaincode
func (t *ChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	_, args := stub.GetFunctionAndParameters()
	chaincodeController := new(src.Controller)
	return util.ExecuteMethod(chaincodeController, "Init", stub, args)
}

// Invoke Function Executes everytime except on initializing or on updating the chain code
func (t *ChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoking " + function)
	chaincodeController := new(src.Controller)
	return util.ExecuteMethod(chaincodeController, function, stub, args)
}
