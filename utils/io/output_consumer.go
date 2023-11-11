package io

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

// A ConsumeOutput provides ability to consume output of any operation supported by the lexicon.
// It provides a method for every type of output supported by lx.OperationResult
type ConsumeOutput interface {
	// ConsumeWords will consume the given array of words
	ConsumeWords(operation string, output *[]string)

	// ConsumeMapOfWords will consume the given map where key is a word and value is array of words
	ConsumeMapOfWords(operation string, output *map[string][]string)
}

// A ConsumeOutputToLog is one of the implementation of ConsumeOutput which forwards the output
// to the log. If logging is not configured, which is generally the case, then output is print
// to the screen.
type ConsumeOutputToLog struct{}

func (co *ConsumeOutputToLog) ConsumeWords(operation string, output *[]string) {
	log.Printf("%s result: \n%v\n", operation, *output)
}

func (co *ConsumeOutputToLog) ConsumeMapOfWords(operation string, output *map[string][]string) {
	log.Printf("%s result: \n%v\n", operation, *output)
}

// A ConsumeOutputToFile is one of the implementation of ConsumeOutput which forwards the output
// to the provided file.
type ConsumeOutputToFile struct {
	OutputFolderPath string
}

func (co *ConsumeOutputToFile) ConsumeWords(operation string, output *[]string) {
	var sb strings.Builder
	for _, word := range *output {
		sb.WriteString(word)
	}

	path := co.OutputFolderPath+"/"+operation+".txt"
	log.Printf("result of %s : %s\n", operation, path)
	os.WriteFile(path, []byte(sb.String()), 0644)
}

func (co *ConsumeOutputToFile) ConsumeMapOfWords(operation string, output *map[string][]string) {
	if jsonString, err := json.Marshal(output); err == nil {
		path := co.OutputFolderPath+"/"+operation+".txt"
		log.Printf("result of %s : %s\n", operation, path)
		os.WriteFile(path, jsonString, 0644)
	}
}
