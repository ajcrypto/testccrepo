/**
 *
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 *
 */
package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
	"unicode"
	"strings"
	"example.com/fffffefe/lib/util/validators"

	"github.com/creasty/defaults"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// Stub stores the ChaincodeStub for the current chaincode
var Stub shim.ChaincodeStubInterface
var ChaincodeName string

// Get ChaincodeStubInterface
func GetStub() shim.ChaincodeStubInterface {
	return Stub
}

// CreateModel constructs the struct object from the given jsonString
func CreateModel(obj interface{}, inputString string) error {
	// fmt.Printf("CreateModel inputString: %s", inputString)
	err := json.Unmarshal([]byte(inputString), &obj)
	if err != nil {
		return fmt.Errorf("Error in creating asset %s", err.Error())
	}
	if err = defaults.Set(obj); err != nil {
		return fmt.Errorf("Failure in default setting %s ", err.Error())
	}
	err = validators.ValidateStruct(obj)
	if err != nil {
		fmt.Println("Validation Failed")
		return err
	}
	return nil
}

// SetField sets the structField of a struct object with the given value
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)
	if !structFieldValue.IsValid() {
		return fmt.Errorf("Error in setting field: No such field: %s in obj", name)
	}
	if !structFieldValue.CanSet() {
		return fmt.Errorf("Error in setting field: Cannot set %s field value", name)
	}
	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New(fmt.Sprintln("Error in setting field: Provided value type didn't match obj field type. Field-", structFieldValue, "Value-", val))
	}
	structFieldValue.Set(val)
	return nil
}

// SetAssetType to Assign Asset Types
func SetAssetType(inputStruct interface{}) error {

	// TypeOf returns type of
	// interface value passed to it
	typ := reflect.TypeOf(inputStruct).Elem()
	name := "AssetType"
	structValue := reflect.ValueOf(inputStruct).Elem()
	structFieldValue := structValue.FieldByName(name)

	// fmt.Println("first fieldByName is successful")
	f, _ := typ.FieldByName(name)
	value := f.Tag.Get("final")

	if !structFieldValue.IsValid() {
		return fmt.Errorf("Error in setting field: No such field: %s in obj", name)
	}
	if !structFieldValue.CanSet() {
		return fmt.Errorf("Error in setting field: Cannot set %s field value", name)
	}
	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New(fmt.Sprintln("Error in setting field: Provided value type didn't match obj field type. Field-", structFieldValue, "Value-", val))
	}
	structFieldValue.Set(val)
	return nil
}

// ConvertMapToStructBasic is an utility function to construct a struct from a given map[string]interface{}.
// This does not work for complex types.
func ConvertMapToStructBasic(inputMap map[string](interface{}), resultStruct interface{}) error {
	for key, value := range inputMap {
		err := SetField(resultStruct, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// ConvertMapToStruct is an another utility function to construct a struct from a given map[string]interface{}.
// This can handle complex types and custom types.
func ConvertMapToStruct(inputMap map[string](interface{}), resultStruct interface{}) error {
	mapBytes, err := json.Marshal(inputMap)
	if err != nil {
		return fmt.Errorf("Error in marshalling map %s", err.Error())
	}
	err = json.Unmarshal(mapBytes, resultStruct)
	if err != nil {
		return fmt.Errorf("Error in unmarshalling map %s", err.Error())
	}
	return nil
}

func makeFirstLetterLowerCaps(input string) string {
	runes := []rune(input)
	if len(runes) > 0 {
		runes[0] = unicode.ToLower(runes[0])
	}
	return string(runes)
}

func convert(argKind reflect.Kind, arg string, argType reflect.Type) (reflect.Value, error) {
	switch argKind {
	case reflect.Bool:
		val, err := strconv.ParseBool(arg)
		if err == nil {
			return reflect.ValueOf(val).Convert(argType), err
		}
		return reflect.ValueOf((interface{})(nil)), err
	case reflect.Int:
		val, err := strconv.ParseInt(arg, 10, 64)
		if err == nil {
			return reflect.ValueOf(int(val)).Convert(argType), err
		}
		return reflect.ValueOf((interface{})(nil)), err
	case reflect.Int8:
		val, err := strconv.ParseInt(arg, 10, 8)
		if err == nil {
			return reflect.ValueOf(int8(val)).Convert(argType), err
		}
		return reflect.ValueOf((interface{})(nil)), err
	case reflect.Int16:
		val, err := strconv.ParseInt(arg, 10, 16)
		if err == nil {
			return reflect.ValueOf(int16(val)).Convert(argType), err
		}
		return reflect.ValueOf((interface{})(nil)), err
	case reflect.Int32:
		val, err := strconv.ParseInt(arg, 10, 32)
		if err == nil {
			return reflect.ValueOf(int32(val)).Convert(argType), err
		}
		return reflect.ValueOf((interface{})(nil)), err
	case reflect.Int64:
		val, err := time.ParseDuration(arg)
		if err == nil {
			return reflect.ValueOf(val).Convert(argType), err
		} else if val, err := strconv.ParseInt(arg, 10, 64); err == nil {
			return reflect.ValueOf(val).Convert(argType), err
		}
		return reflect.ValueOf((interface{})(nil)), err
	case reflect.Uint:
		val, err := strconv.ParseUint(arg, 10, 64)
		if err == nil {
			return reflect.ValueOf(uint(val)).Convert(argType), err
		}
		return reflect.ValueOf((interface{})(nil)), err
	case reflect.Uint8:
		val, err := strconv.ParseUint(arg, 10, 8)
		if err == nil {
			return reflect.ValueOf(uint8(val)).Convert(argType), err
		}
		return reflect.ValueOf((interface{})(nil)), err
	case reflect.Uint16:
		val, err := strconv.ParseUint(arg, 10, 16)
		if err == nil {
			return reflect.ValueOf(uint16(val)).Convert(argType), err
		}
		return reflect.ValueOf((interface{})(nil)), err
	case reflect.Uint32:
		val, err := strconv.ParseUint(arg, 10, 32)
		if err == nil {
			return reflect.ValueOf(uint32(val)).Convert(argType), err
		}
		return reflect.ValueOf((interface{})(nil)), err
	case reflect.Uint64:
		val, err := strconv.ParseUint(arg, 10, 64)
		if err == nil {
			return reflect.ValueOf(val).Convert(argType), err
		}
		return reflect.ValueOf((interface{})(nil)), err
	case reflect.Uintptr:
		val, err := strconv.ParseUint(arg, 10, 64)
		if err == nil {
			return reflect.ValueOf(uintptr(val)).Convert(argType), err
		}
		return reflect.ValueOf((interface{})(nil)), err
	case reflect.Float32:
		val, err := strconv.ParseFloat(arg, 32)
		if err == nil {
			return reflect.ValueOf(float32(val)).Convert(argType), err
		}
		return reflect.ValueOf((interface{})(nil)), err
	case reflect.Float64:
		val, err := strconv.ParseFloat(arg, 64)
		if err == nil {
			return reflect.ValueOf(val).Convert(argType), err
		}
		return reflect.ValueOf((interface{})(nil)), err
	case reflect.String:
		return reflect.ValueOf(arg).Convert(argType), nil
	case reflect.Slice:
		ref := reflect.New(argType)
		ref.Elem().Set(reflect.MakeSlice(argType, 0, 0))
		if err := json.Unmarshal([]byte(arg), ref.Interface()); err != nil {
			return reflect.ValueOf((interface{})(nil)), err
		}
		return ref.Elem().Convert(argType), nil
	case reflect.Map:
		ref := reflect.New(argType)
		ref.Elem().Set(reflect.MakeMap(argType))
		if err := json.Unmarshal([]byte(arg), ref.Interface()); err != nil {
			return reflect.ValueOf((interface{})(nil)), err
		}
		return ref.Elem().Convert(argType), nil
	case reflect.Struct:
		var obj interface{}
		if err := json.Unmarshal([]byte(arg), &obj); err != nil {
			return reflect.ValueOf((interface{})(nil)), err
		}
		inputArgMap := obj.(map[string]interface{})
		for i := 0; i < argType.NumField(); i++ {
			mandatoryTagValue := argType.Field(i).Tag.Get("mandatory")
			if mandatoryTagValue == "true" {
				_, ok := inputArgMap[argType.Field(i).Name]
				if !ok {
					_, ok2 := inputArgMap[makeFirstLetterLowerCaps(argType.Field(i).Name)]
					if !ok2 {
						return reflect.ValueOf((interface{})(nil)), fmt.Errorf("Mandatory field %s for asset %s is not present in the input", argType.Field(i).Name, strings.Split(argType.String(), ".")[1])
					}
				}
			}
		}
		ref := reflect.New(argType)
		if err := defaults.Set(ref.Interface()); err != nil {
			return reflect.ValueOf((interface{})(nil)), err
		}
		if err := json.Unmarshal([]byte(arg), ref.Interface()); err != nil {
			return reflect.ValueOf((interface{})(nil)), err
		}
		err := validators.ValidateStruct(ref.Interface())
		if err != nil {
			return reflect.ValueOf((interface{})(nil)), err
		}
		val := ref.Elem().Convert(argType)
		return val, nil
	case reflect.Ptr:
		ref := reflect.New(argType.Elem())
		return ref.Elem().Convert(argType), nil
	}
	return reflect.ValueOf((interface{})(nil)), fmt.Errorf(("Argument Parsing/Validation failed: argument kind does not match supported kinds"))
}

func processArgs(inputArgTypes reflect.Type, args []string, functionName string) ([]reflect.Value, error) {
	result := make([]reflect.Value, inputArgTypes.NumIn())

	if inputArgTypes.NumIn() != len(args) {
		if (functionName == "Init" && len(args) == 1 && args[0] == "") {
			dummyresult := make([]reflect.Value, 0)
			return dummyresult, nil 
		}
		return nil, fmt.Errorf("Number of input arguments required by the function %s are %d, which did not match the number arguments passed i.e %d", functionName, inputArgTypes.NumIn(), len(args))
	}

	for i := 0; i < inputArgTypes.NumIn(); i++ {
		// fmt.Println(inputArgTypes.In(i).Kind(), args[i])
		response, err := convert(inputArgTypes.In(i).Kind(), args[i], inputArgTypes.In(i))
		// fmt.Println("Response", response)
		if err == nil {
			result[i] = response
		} else {
			return nil, err
		}
	}
	// fmt.Println("processArgs", result)
	return result, nil
}

// ExecuteMethod calls a method with the given name on the provided reciever
func ExecuteMethod(obj interface{}, function string, stub shim.ChaincodeStubInterface, args []string) peer.Response {
	Stub = stub
	methodValue := reflect.ValueOf(obj).MethodByName(function)
	if methodValue.IsValid() != true {
		return shim.Error(fmt.Sprintf("ExecuteMethod: No method found by given name - %s", function))
	}
	// fmt.Println("args", args)
	// fmt.Println("len args ", len(args))
	// if function == "Init" && len(args) == 1 && args[0] == "" {
	// 	fmt.Println("inside here")
	// 	convertedArgs := make([]reflect.Value, 1)
	// 	convertedArgsobj := make([]string, 1)
	// 	convertedArgsobj[0] = ""
	// 	convertedArgs[0] = reflect.ValueOf(convertedArgsobj)
	// 	result := methodValue.Call(convertedArgs)
	// 	resultError := result[1].Interface()
	// 	if resultError != nil {
	// 		return shim.Error(fmt.Sprintf("ExecuteMethod: Error: %s", resultError.(error).Error()))
	// 	}
	// 	returnObj := result[0].Interface()
	// 	returnBytes, errMarshal := json.Marshal(returnObj)
	// 	if errMarshal != nil {
	// 		return shim.Error(fmt.Sprintf("ExecuteMethod: Marshalling response Error: %s", errMarshal.Error()))
	// 	}
	// 	return shim.Success(returnBytes)
	// }
	convertedArgs, err := processArgs(methodValue.Type(), args, function)
	// fmt.Println(convertedArgs)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error in argument parsing and validation Detailed Error : %s", err.Error()))
	}
	result := methodValue.Call(convertedArgs)
	resultError := result[1].Interface()
	if resultError != nil {
		return shim.Error(fmt.Sprintf("ExecuteMethod: Error: %s", resultError.(error).Error()))
	}
	returnObj := result[0].Interface()
	returnBytes, errMarshal := json.Marshal(returnObj)
	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("ExecuteMethod: Marshalling response Error: %s", errMarshal.Error()))
	}
	return shim.Success(returnBytes)
}
