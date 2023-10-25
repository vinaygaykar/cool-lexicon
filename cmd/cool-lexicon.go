package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vinaygaykar/cool-lexicon/configs"
	"github.com/vinaygaykar/cool-lexicon/pkg/lexicon"
)

type ProgramInputs struct {
	configFile               string // Location of the config file
	shouldPerformSetupChecks bool   // true if setup checks should be performed
	isFileBasedInput         bool   // true if the input should be read from the given file instead of the command line

	opLookup             string // value of the LOOKUP operation, if `isFileBasedInput` is true then this is file location else this is a word to operate on
	opSearchStartingWith string // value of the SEARCH START WITH operation, if `isFileBasedInput` is true then this is file location else this is a word to operate on
	opSearchEndingWith   string // value of the SEARCH END WITH operation, if `isFileBasedInput` is true then this is file location else this is a word to operate on
	opAdd                string // value of the ADD operation, if `isFileBasedInput` is true then this is file location else this is a word to operate on
}

var (
	// arg inputs
	inputs ProgramInputs

	// errors
	errFileInvalid = errors.New("file contents are invalid/empty or file is corrupt or file does not exist")

	// constant responses
	noWords = []string{}
)

func init() {
	flag.BoolVar(&inputs.shouldPerformSetupChecks, "check", false, "Setup all necessary configs if required. This is optional, if the all configs are already setup correctly this operation will have no effect")
	flag.BoolVar(&inputs.isFileBasedInput, "if", false, "This flag indicates that input words to every operation should be taken from a file present at the given location")
	flag.StringVar(&inputs.configFile, "cfg", "cool-lexicon-cfg.json", "Config file location")
	flag.StringVar(&inputs.opLookup, "ex", "", "Check if the given word exist")
	flag.StringVar(&inputs.opSearchStartingWith, "ss", "", "Search the lexicon to find words that start with given substring")
	flag.StringVar(&inputs.opSearchEndingWith, "se", "", "Search the lexicon to find words that end with given substring")
	flag.StringVar(&inputs.opAdd, "ad", "", "Add words present in given file location to lexicon")
}

func main() {
	flag.Parse()

	sanitizeInputs()
	validateInputs()

	lxc := configs.GetLexicon(inputs.configFile, inputs.shouldPerformSetupChecks)
	defer lxc.Close()

	tryOperateExists(lxc)
	tryOperateGetAllStartingWith(lxc)
	tryOperateGetAllEndingWith(lxc)
	tryOperateAddAll(lxc)
}

func sanitizeInputs() {
	// remove any whitespaces
	inputs.configFile = strings.TrimSpace(inputs.configFile)
	inputs.opLookup = strings.TrimSpace(inputs.opLookup)
	inputs.opSearchStartingWith = strings.TrimSpace(inputs.opSearchStartingWith)
	inputs.opSearchEndingWith = strings.TrimSpace(inputs.opSearchEndingWith)
	inputs.opAdd = strings.TrimSpace(inputs.opAdd)
}

func validateInputs() {
	if len(inputs.configFile) == 0 {
		// config file location string must be present; default file location string is provided to `flag`
		log.Panic("config file location not provided")
	}

	// config file location string is there but is the location valid
	if _, err := os.Stat(inputs.configFile); err != nil {
		log.Panic(err.Error())
	}

	if len(inputs.opLookup) == 0 && len(inputs.opSearchStartingWith) == 0 && len(inputs.opSearchEndingWith) == 0 && len(inputs.opAdd) == 0 {
		flag.PrintDefaults()
		log.Panic("no operation provided")
	}

	if len(inputs.opAdd) != 0 { // Check if the given file exists
		if _, err := os.Stat(inputs.opAdd); err != nil {
			log.Panic(err.Error())
		}
	}
}

func tryOperateExists(lxc lexicon.Lexicon) {
	words, err := getWords(inputs.opLookup)
	if err != nil {
		log.Fatalf("could not perform 'exists' for the word (%s), error: %s\n", inputs.opLookup, err.Error())
	} else if len(words) == 0 {
		return // this operation was not selected
	}

	if exists, err := lxc.Lookup(words...); err != nil {
		log.Fatalf("could not perform 'exists' for the word (%s), error: %s\n", inputs.opLookup, err.Error())
	} else {
		fmt.Printf("exists (%s) : %t\n", inputs.opLookup, exists)
	}
}

func tryOperateGetAllStartingWith(lxc lexicon.Lexicon) {
	words, err := getWords(inputs.opSearchStartingWith)
	if err != nil {
		log.Fatalf("could not perform 'search starts with' for the word (%s), error: %s\n", inputs.opSearchStartingWith, err.Error())
	} else if len(words) == 0 {
		return // this operation was not selected
	}
	
	if words, err := lxc.GetAllWordsStartingWith(inputs.opSearchStartingWith); err != nil {
		log.Fatalf("could not perform 'starts with' for the word (%s), error: %s\n", inputs.opSearchStartingWith, err.Error())
	} else {
		fmt.Printf("starts with (%s) : %v\n", inputs.opSearchStartingWith, words)
	}
}

func tryOperateGetAllEndingWith(lxc lexicon.Lexicon) {
	if len(inputs.opSearchEndingWith) == 0 {
		return
	}

	if words, err := lxc.GetAllWordsEndingWith(inputs.opSearchEndingWith); err != nil {
		log.Fatalf("could not perform 'ends with' for the word (%s), error: %s\n", inputs.opLookup, err.Error())
	} else {
		fmt.Printf("ends with (%s) : %v\n", inputs.opSearchEndingWith, words)
	}
}

func tryOperateAddAll(lxc lexicon.Lexicon) {
	if len(inputs.opAdd) == 0 {
		return
	}

	// open file containing words to add
	file, err := os.Open(inputs.opAdd)
	if err != nil {
		log.Fatalf("could not open file to add words, error: %s\n", err.Error())
		return
	}
	defer file.Close()

	words := make([]string, 0)

	// read file contents into `words`
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, strings.TrimSpace(scanner.Text()))
	}

	if err = scanner.Err(); err != nil {
		log.Fatalf("could not read contents of the file to add words, error: %s\n", err.Error())
		return
	}

	// add `words` to lexicon
	if err = lxc.Add(words...); err != nil {
		log.Fatalf("could not perform 'add words' from file (%s), error: %s\n", inputs.opAdd, err.Error())
	} else {
		fmt.Printf("added words from the file (%s)\n", inputs.opAdd)
	}
}

func getWords(value string) ([]string, error) {
	if len(value) == 0 { // empty input value; return empty result
		return noWords, nil
	}

	if inputs.isFileBasedInput {
		return readFileAt(value)
	} else { // input is a single word
		return []string{value}, nil
	}
}

func readFileAt(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return noWords, errors.Join(errFileInvalid, err)
	}

	words := make([]string, 0)

	// read file contents into `words`
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("could not read contents from file %s, error: %s\n", file.Name(), err.Error())
		return noWords, errors.Join(errFileInvalid, err)
	}

	return words, nil
}
