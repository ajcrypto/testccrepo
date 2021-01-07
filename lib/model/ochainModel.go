/**
 * 
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * 
 */
package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"example.com/fffffefe/lib/util"
	"example.com/fffffefe/lib/util/validators"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

func getID(obj interface{}) (string, error) {
	objValue := reflect.ValueOf(obj).Elem()
	objType := objValue.Type()

	for i := 0; i < objType.NumField(); i++ {
		structField := objType.Field(i)
		_, ok := structField.Tag.Lookup("id")
		if ok {
			idValue := objValue.Field(i)
			idValInterface := idValue
			idString := fmt.Sprintf("%v", idValInterface)
			// fmt.Println(idString, reflect.TypeOf(idString))
			return idString, nil
		}
	}
	return "", errors.New("Id tag is not set")
}

// Save writes the asset to the ledger
func Save(args ...interface{}) (interface{}, error) {
	stub := util.Stub
	obj := args[0]

	id, idErr := getID(obj)
	if idErr != nil {
		return nil, fmt.Errorf("Error in getting Id. Id is mandatory. Error %s", idErr.Error())
	}

	_, err := Get(id)
	if err == nil {
		return nil, fmt.Errorf("Error in saving: asset already exist in ledger with Id %s ", id)
	}

	err = util.SetAssetType(obj)
	if err != nil {
		return nil, fmt.Errorf("AssetType is missing or resetting is a problem %s", err.Error())
	}

	errValidation := validators.ValidateStruct(obj)
	if errValidation != nil {
		fmt.Println("Validation Failed")
		return nil, fmt.Errorf("Error in saving: Asset Id %s marshal error %s", id, errValidation.Error())
	}

	if len(args) > 1 {
		ptrToAsset := obj
		assetValue := reflect.ValueOf(ptrToAsset).Elem()
		metadata := args[1]
		metdataField := assetValue.FieldByName("Metadata")
		metdataField.Set(reflect.ValueOf(metadata))
	}

	assetAsBytes, errMarshal := json.Marshal(obj)
	if errMarshal != nil {
		return nil, fmt.Errorf("Error in saving: Asset Id %s marshal error %s", id, errMarshal.Error())
	}

	errPut := stub.PutState(id, assetAsBytes)
	if errPut != nil {
		return nil, fmt.Errorf("Error in saving: Asset Id %s transaction error %s", id, errPut.Error())
	}

	fmt.Println("Success in Initiating Transaction Asset", obj)
	return obj, nil
}

func GenerateCompositeKey(indexName string, attributes []string) (string, error) {
	stub := util.Stub
	if len(attributes) == 0 {
		const errorMessage = "Attributes param is expected to be an array of string"
		return "", fmt.Errorf(errorMessage)
	}

	compositeKey, err := stub.CreateCompositeKey(indexName, attributes)
	if err != nil {
		return "", fmt.Errorf("Failed creating composite Key")
	}
	return compositeKey, nil
}

func GetByCompositeKey(key string, columns []string, index int) (interface{}, error) {

	stub := util.Stub

	resultsIterator, err := stub.GetStateByPartialCompositeKey(key, columns)
	if err != nil {
		return nil, fmt.Errorf("Error in returning iterator: %d", resultsIterator)
	}

	defer resultsIterator.Close()
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("Error in getting GetByCompositeKey: iteration error %s", err.Error())
		}

		_, compositeKeyParts, err := stub.SplitCompositeKey(queryResult.Key)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}

		returnedID := compositeKeyParts[index]

		assetAsBytes, _ := stub.GetState(returnedID)

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedID)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(assetAsBytes))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	var result []interface{}
	unmarshalError := json.Unmarshal(buffer.Bytes(), &result)
	if unmarshalError != nil {
		return nil, fmt.Errorf("Error in getting history by id: unmarshalling error %s", unmarshalError.Error())
	}
	return result, nil
}

func GetTransactionId() string {
	return util.Stub.GetTxID()
}

func GetTransactionTimestamp() (*timestamp.Timestamp, error) {
	return util.Stub.GetTxTimestamp()
}

func GetChannelID() string {
	return util.Stub.GetChannelID()
}

func GetCreator() ([]byte, error) {
	return util.Stub.GetCreator()
}

func GetSignedProposal() (*peer.SignedProposal, error) {
	return util.Stub.GetSignedProposal()
}

func GetArgs() [][]byte {
	return util.Stub.GetArgs()
}

func GetStringArgs() []string {
	return util.Stub.GetStringArgs()
}

func GetNetworkStub() shim.ChaincodeStubInterface {
	return util.GetStub()
}

func Get(Id string, result ...interface{}) (interface{}, error) {
	stub := util.Stub

	assetAsBytes, _ := stub.GetState(Id)
	if assetAsBytes == nil {
		return nil, fmt.Errorf("Error in getting: Asset with Id %s does not exists", Id)
	}

	var genericResult interface{}
	unmarshalError := json.Unmarshal(assetAsBytes, &genericResult)
	if unmarshalError != nil {
		return nil, fmt.Errorf("Error in getting: marshalling error %s", unmarshalError.Error())
	}
	assetTypeFromLedgerString := genericResult.(map[string]interface{})["AssetType"]
	assetTypeFromLedger := strings.Split(assetTypeFromLedgerString.(string), ".")[1]
	if len(result) > 0 {
		inputAssetTypeString := reflect.ValueOf(result[0]).Elem().Type().String()
		inputAssetType := strings.Split(inputAssetTypeString, ".")[1]
		if inputAssetType != assetTypeFromLedger {
			return nil, fmt.Errorf("No Asset %s exist with id %s", inputAssetType, Id)
		}
		unmarshalError := json.Unmarshal(assetAsBytes, result[0])
		if unmarshalError != nil {
			return nil, fmt.Errorf("Error in getting: marshalling error %s", unmarshalError.Error())
		}
		errValidation := validators.ValidateStruct(result[0])
		if errValidation != nil {
			fmt.Println("Validation Failed")
			return nil, fmt.Errorf("Error in retrieving asset: Asset %v error %s", result[0], errValidation.Error())
		}
		return result[0], nil
	}
	return genericResult, nil
}

// Update the asset to the ledger
func Update(args ...interface{}) (interface{}, error) {
	stub := util.Stub

	obj := args[0]
	id, idErr := getID(obj)
	if idErr != nil {
		return nil, errors.New("Id tag is not set in the struct, id is necessary for saving the object")
	}

	assetAsBytes, _ := stub.GetState(id)
	if assetAsBytes == nil {
		return nil, fmt.Errorf("Error in updating: Unable to get the asset from ledger with ID %s", id)
	}

	err := util.SetAssetType(obj)
	if err != nil {
		return nil, fmt.Errorf("AssetType is missing or resetting is a problem %s", err.Error())
	}

	errValidation := validators.ValidateStruct(obj)
	if errValidation != nil {
		fmt.Println("Validation Failed")
		return nil, fmt.Errorf("Error in updating: Asset Id %s marshal error %s", id, errValidation.Error())
	}

	assetAsBytes, errMarshal := json.Marshal(obj)
	if errMarshal != nil {
		return nil, fmt.Errorf("Error in updating: Asset Id %s marshal error %s", id, errMarshal.Error())
	}

	errPut := stub.PutState(id, assetAsBytes)
	if errPut != nil {
		return nil, fmt.Errorf("Error in updating: Asset Id %s marshal error %s", id, errPut.Error())
	}

	fmt.Println("Success in initiating Transaction Asset", obj)
	return obj, nil
}

// Delete deletes the asset from the ledger
func Delete(Id string) (interface{}, error) {
	stub := util.Stub

	assetAsBytes, _ := stub.GetState(Id)
	if assetAsBytes == nil {
		return nil, fmt.Errorf("Error in deleting: could not find asset with Id %s", Id)
	}

	errPut := stub.DelState(Id)
	if errPut != nil {
		return nil, fmt.Errorf("Error in deleting: failed to delete asset with Id %s error %s", Id, errPut.Error())
	}

	var result interface{}
	unmarshalError := json.Unmarshal(assetAsBytes, &result)
	if unmarshalError != nil {
		return nil, fmt.Errorf("Error in deleting: marshalling error %s", unmarshalError.Error())
	}
	return result, nil
}

// Query runs the given transaction on the peer
func Query(queryString string) ([]interface{}, error) {
	stub := util.Stub
	fmt.Printf("Query: queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("Query: iteration error %s", err.Error())
	}

	defer resultsIterator.Close()
	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	var result []interface{}
	unmarshalError := json.Unmarshal(buffer.Bytes(), &result)
	if unmarshalError != nil {
		return nil, fmt.Errorf("Query: unmarshalling result error %s", unmarshalError.Error())
	}
	return result, nil
}

// GetByRange gets all the assets with key between the provided range
func GetByRange(startKey string, endKey string, asset ...interface{}) ([]map[string]interface{}, error) {
	stub := util.Stub
	if len(asset) > 0 {
		resultsIterator, err := stub.GetStateByRange(startKey, endKey)

		// fmt.Println("GetSupplierByRange", reflect.TypeOf(asset[0]))
		// fmt.Println("GetSupplierByRange", reflect.TypeOf(asset[0]).Kind())
		// fmt.Println("GetSupplierByRange", reflect.TypeOf(asset[0]).Elem().Elem())

		inputAssetTypeString := reflect.TypeOf(asset[0]).Elem().Elem().String()
		inputAssetType := strings.Split(inputAssetTypeString, ".")[1]
		//fmt.Println("GetSupplierByRange", reflect.TypeOf(asset[0]))
		//fmt.Println("GetSupplierByRange", reflect.TypeOf(asset[0]).Kind())

		if err != nil {
			return nil, fmt.Errorf("Error in getting by range: %s", err.Error())
		}

		defer resultsIterator.Close()
		var buffer bytes.Buffer
		buffer.WriteString("[")
		bArrayMemberAlreadyWritten := false
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return nil, fmt.Errorf("Error in getting by range: iteration error %s", err.Error())
			}
			// Add a comma before array members, suppress it for the first array member
			if bArrayMemberAlreadyWritten == true {
				buffer.WriteString(",")
			}
			buffer.WriteString("{\"Key\":")
			buffer.WriteString("\"")
			buffer.WriteString(queryResponse.Key)
			buffer.WriteString("\"")
			buffer.WriteString(", \"Record\":")
			// Record is a JSON object, so we write as-is
			buffer.WriteString(string(queryResponse.Value))
			buffer.WriteString("}")
			bArrayMemberAlreadyWritten = true
		}

		buffer.WriteString("]")
		var result []interface{}
		unmarshalError := json.Unmarshal(buffer.Bytes(), &result)
		if unmarshalError != nil {
			return nil, fmt.Errorf("Error in getting by range: unmarshalling error %s", unmarshalError.Error())
		}
		var resultAssets []map[string]interface{}
		// fmt.Println("Length of result", len(result))
		for i := 0; i < len(result); i++ {
			// fmt.Println("here inside")
			entry := result[i].(map[string]interface{})
			//fmt.Println(entry)
			value := entry["Record"]
			mapAsset := value.(map[string]interface{})
			assetTypeString := mapAsset["AssetType"].(string)
			assetTypeSplit := strings.Split(assetTypeString, ".")

			chaincodeName := assetTypeSplit[0]
			// fmt.Println("Printing", chaincodeName, util.ChaincodeName)
			// fmt.Println("Printing", assetTypeSplit[1], inputAssetType)
			if len(assetTypeSplit) > 1 && assetTypeSplit[1] == inputAssetType && chaincodeName == util.ChaincodeName {
				// fmt.Println("here")
				// assetType := assetTypeSplit[1]
				// fmt.Println("inputAssetType", inputAssetType, "AssetType", assetType)
				resultAssets = append(resultAssets, mapAsset)
			}
		}
		mapBytes, err := json.Marshal(resultAssets)
		if err != nil {
			return nil, fmt.Errorf("Error in marshalling map %s", err.Error())
		}
		err = json.Unmarshal(mapBytes, asset[0])
		if err != nil {
			return nil, fmt.Errorf("Error in unmarshalling map %s", err.Error())
		}
		// validation
		//typeToValidate := reflect.TypeOf(asset[0]).Elem()
		// for i := 0; i < len(asset[0]([]interface{})); i++ {
		// 	asset2 := asset[0].(*[]interface{});
		// 	errValidation := validators.ValidateStruct(asset2[i])
		// 	if errValidation != nil {
		// 		fmt.Println("Validation Failed")
		// 		return nil, fmt.Errorf("Error in retrieving asset: Asset %v error %s", asset2[i], errValidation.Error())
		// 	}
		// }
		return resultAssets, nil
	}
	resultsIterator, err := stub.GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, fmt.Errorf("Error in getting by range: %s", err.Error())
	}

	defer resultsIterator.Close()
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("Error in getting by range: iteration error %s", err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")
	var result []map[string]interface{}
	unmarshalError := json.Unmarshal(buffer.Bytes(), &result)
	if unmarshalError != nil {
		return nil, fmt.Errorf("Error in getting by range: unmarshalling error %s", unmarshalError.Error())
	}
	return result, nil
}

// GetHistoryByID gets the history of an asset from the ledger
func GetHistoryByID(Id string) ([]interface{}, error) {
	recordKey := Id
	stub := util.Stub
	// fmt.Printf("- start getHistoryForRecord: %s\n", recordKey)

	resultsIterator, err := stub.GetHistoryForKey(recordKey)
	if err != nil {
		return nil, fmt.Errorf("Error in getting history by id: %s", err.Error())
	}
	defer resultsIterator.Close()
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("Error in getting history by id: iteration error %s", err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		// corresponding value null. Else, we will write the response.Values
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}
		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")
		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	var result []interface{}
	unmarshalError := json.Unmarshal(buffer.Bytes(), &result)
	if unmarshalError != nil {
		return nil, fmt.Errorf("Error in getting history by id: unmarshalling error %s", unmarshalError.Error())
	}
	return result, nil
}
