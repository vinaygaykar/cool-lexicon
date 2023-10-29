package io

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"

	lx "github.com/vinaygaykar/cool-lexicon/pkg/lexicon"
)

// A ConsumeOutput provides ability to consume output of any operation supported by the lexicon.
// It is a generic type interface to match with pkg/lexicon/OperationResult.
type ConsumeOutput[V bool | []string] interface {
	// Consume will consume the given map of string and OperationResult.
	// If consumeErrors is set then error message will also be consumed.
	// Check implementations for concrete details.
	Consume(output *map[string]lx.OperationResult[V], consumeErrors bool)
}

// A ConsumeOutputToLog is one of the implementation of ConsumeOutput which forwards the output
// to the log. If logging is not configured, which is generally the case, then output is print
// to the screen.
type ConsumeOutputToLog[V bool | []string] struct{}

func (co *ConsumeOutputToLog[V]) Consume(output *map[string]lx.OperationResult[V], consumeErrors bool) {
	for k, v := range *output {
		if v.Err != nil && consumeErrors {
			log.Printf("%s : error : %s\n", k, v.Err.Error())
		} else {
			if json, err := json.Marshal(v.Value); err != nil && consumeErrors {
				log.Printf("%s : error printing result : %s\n", k, err.Error())
			} else {
				log.Printf("%s : %v\n", k, string(json))
			}
		}
	}
}

// A ConsumeOutputToFile is one of the implementation of ConsumeOutput which forwards the output
// to the provided file.
type ConsumeOutputToFile[V bool | []string] struct {
	filePath string
}

func (co *ConsumeOutputToFile[V]) Consume(output *map[string]lx.OperationResult[V], consumeErrors bool) {
	data := make([]byte, 0)

	for k, v := range *output {
		if v.Err != nil && consumeErrors {
			data = append(data, fmt.Sprintf("%s : error : %s\n", k, v.Err.Error())...)
		} else {
			if json, err := json.Marshal(v.Value); err != nil && consumeErrors {
				data = append(data, fmt.Sprintf("%s : error printing result : %s\n", k, err.Error())...)
			} else {
				data = append(data, fmt.Sprintf("%s : %v\n", k, string(json))...)
			}
		}
	}

	os.WriteFile(co.filePath, data, fs.FileMode(os.O_WRONLY))
}
