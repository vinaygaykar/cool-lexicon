package io

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	// errors
	ErrNoInputValue = errors.New("raw value is empty or blank")
)

// A SupplyInput defines interface to recieve program inputs
type SupplyInput interface {
	// Get returns array of strings which represent input words for the lexicon operation.
	// `rawValue` is the unprocessed input value as recieved from the user interface/terminal.
	// If rawValue is empty or blank, then error ErrNoInputValue is returned with empty array.
	// If error is encountered while parsing rawValue, nil array with error response is returned.
	// If error is present along with array response, it should indicate that the array response is
	// not to be used and instead error should be dealt with.
	Get(rawValue string) ([]string, error)
}

// A SupplyInputFromWord is one of the implementation of SupplyInput.
// It processes and treat the passed rawValue as a single word itself.
type SupplyInputFromWord struct{}

func (si *SupplyInputFromWord) Get(rawValue string) ([]string, error) {
	value := strings.TrimSpace(rawValue)

	if len(value) == 0 {
		return nil, ErrNoInputValue
	}

	return []string{value}, nil
}

// A SupplyInputFromFile is one of the implementation of SupplyInput.
// It processes and treat the passed rawValue as a file path which contains words to be used as input.
type SupplyInputFromFile struct{}

func (si *SupplyInputFromFile) Get(rawValue string) ([]string, error) {
	path := strings.TrimSpace(rawValue)

	if len(path) == 0 {
		return nil, ErrNoInputValue
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("input: file is corrupt or file does not exist: %w", err)
	}

	words := make([]string, 0)

	// read file contents into `words`
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("input: file contents are invalid or file is corrupt: %w", err)
	}

	return words, nil
}
