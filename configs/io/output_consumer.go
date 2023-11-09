package io

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	lx "github.com/vinaygaykar/cool-lexicon/pkg/lexicon"
)

// A ConsumeOutput provides ability to consume output of any operation supported by the lexicon.
// It provides a method for every type of output supported by lx.OperationResult
type ConsumeOutput interface {
	// ConsumeBool will consume the given map of string and OperationResult[bool].
	// If consumeErrors is set then error message will also be consumed.
	ConsumeBool(operation string, output *map[string]lx.OperationResult[bool], consumeErrors bool)

	// Consume will consume the given map of string and OperationResult[[]string].
	// If consumeErrors is set then error message will also be consumed.
	ConsumeStringArray(operation string, output *map[string]lx.OperationResult[[]string], consumeErrors bool)
}

// A ConsumeOutputToLog is one of the implementation of ConsumeOutput which forwards the output
// to the log. If logging is not configured, which is generally the case, then output is print
// to the screen.
type ConsumeOutputToLog struct{}

func (co *ConsumeOutputToLog) ConsumeBool(operation string, output *map[string]lx.OperationResult[bool], consumeErrors bool) {
	commonConsumerToLog(output, consumeErrors)
}

func (co *ConsumeOutputToLog) ConsumeStringArray(operation string, output *map[string]lx.OperationResult[[]string], consumeErrors bool) {
	commonConsumerToLog(output, consumeErrors)
}

func commonConsumerToLog[T bool | []string](output *map[string]lx.OperationResult[T], consumeErrors bool) {
	for k, v := range *output {
		if v.Err != nil && consumeErrors {
			log.Printf("%s : error : %s\n", k, v.Err.Error())
		} else {
			if json, err := json.Marshal(v.Value); err != nil && consumeErrors {
				log.Printf("%s : error : printing result : %s\n", k, err.Error())
			} else {
				log.Printf("%s : %v\n", k, string(json))
			}
		}
	}
}

// A ConsumeOutputToFile is one of the implementation of ConsumeOutput which forwards the output
// to the provided file.
type ConsumeOutputToFile struct {
	OutputFolderPath string
}

func (co *ConsumeOutputToFile) ConsumeBool(operation string, output *map[string]lx.OperationResult[bool], consumeErrors bool) {
	commonConsumerToFile(operation, co.OutputFolderPath, output, consumeErrors)
}

func (co *ConsumeOutputToFile) ConsumeStringArray(operation string, output *map[string]lx.OperationResult[[]string], consumeErrors bool) {
	commonConsumerToFile(operation, co.OutputFolderPath, output, consumeErrors)
}

func commonConsumerToFile[T bool | []string](operation, filePath string, output *map[string]lx.OperationResult[T], consumeErrors bool) {
	data := make([]byte, 0)

	for k, v := range *output {
		if v.Err != nil && consumeErrors {
			data = append(data, fmt.Sprintf("%s : error : %s\n", k, v.Err.Error())...)
		} else {
			if json, err := json.Marshal(v.Value); err != nil && consumeErrors {
				data = append(data, fmt.Sprintf("%s : error : printing result : %s\n", k, err.Error())...)
			} else {
				data = append(data, fmt.Sprintf("%s : %v\n", k, string(json))...)
			}
		}
	}

	os.WriteFile(filePath+"/"+operation+".txt", data, 0644)
}
