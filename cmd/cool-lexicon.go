package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vinaygaykar/cool-lexicon/configs"
	"github.com/vinaygaykar/cool-lexicon/configs/io"
	"github.com/vinaygaykar/cool-lexicon/pkg/lexicon"
)

// A ProgramInput holds all the input values provided to the program.
type ProgramArgs struct {
	configFile               string // Location of the config file
	shouldPerformSetupChecks bool   // true if setup checks should be performed
	isFileBasedInput         bool   // true if the input should be read from the given file instead of the command line
	outputFolderPath         string   // true if the output should be printed to file instead of the command line

	opLookup             string // value of the LOOKUP operation, if `isFileBasedInput` is true then this is file location else this is a word to operate on
	opSearchStartingWith string // value of the SEARCH START WITH operation, if `isFileBasedInput` is true then this is file location else this is a word to operate on
	opSearchEndingWith   string // value of the SEARCH END WITH operation, if `isFileBasedInput` is true then this is file location else this is a word to operate on
	opAdd                string // value of the ADD operation, if `isFileBasedInput` is true then this is file location else this is a word to operate on
}

var (
	args          ProgramArgs
	wordSupplier  io.SupplyInput
	outputPrinter io.ConsumeOutput
)

func init() {
	flag.BoolVar(&args.shouldPerformSetupChecks, "check", false, "Setup all necessary configs if required. This is optional, if the all configs are already setup correctly this operation will have no effect")
	flag.BoolVar(&args.isFileBasedInput, "if", false, "This flag indicates that input words to every operation should be taken from the file passed as value to individual operation")
	flag.StringVar(&args.outputFolderPath, "of", "", "This flag indicates that output to every operation should be printed to files (created for every operation) at given path")
	flag.StringVar(&args.configFile, "cfg", "cool-lexicon-cfg.json", "Config file location")
	flag.StringVar(&args.opLookup, "ex", "", "Check if the given word exist")
	flag.StringVar(&args.opSearchStartingWith, "ss", "", "Search the lexicon to find words that start with given substring")
	flag.StringVar(&args.opSearchEndingWith, "se", "", "Search the lexicon to find words that end with given substring")
	flag.StringVar(&args.opAdd, "ad", "", "Add words present in given file location to lexicon")
}

func main() {
	flag.Parse()

	sanitizeInputs()
	validateInputs()

	if args.isFileBasedInput {
		wordSupplier = &io.SupplyWordsFromFile{}
	} else {
		wordSupplier = &io.SupplyWordsFromCLI{}
	}

	if len(args.outputFolderPath) == 0 {
		outputPrinter = &io.ConsumeOutputToLog{}
	} else {
		outputPrinter = &io.ConsumeOutputToFile{ OutputFolderPath: args.outputFolderPath }
	}

	lxc := configs.GetLexicon(args.configFile, args.shouldPerformSetupChecks)
	defer lxc.Close()

	tryOperateExists(lxc)
	tryOperateGetAllStartingWith(lxc)
	tryOperateGetAllEndingWith(lxc)
	tryOperateAddAll(lxc)
}

func sanitizeInputs() {
	// remove any whitespaces
	args.configFile = strings.TrimSpace(args.configFile)
	args.opLookup = strings.TrimSpace(args.opLookup)
	args.opSearchStartingWith = strings.TrimSpace(args.opSearchStartingWith)
	args.opSearchEndingWith = strings.TrimSpace(args.opSearchEndingWith)
	args.opAdd = strings.TrimSpace(args.opAdd)
	args.outputFolderPath = strings.TrimSpace(args.outputFolderPath)
}

func validateInputs() {
	if len(args.configFile) == 0 {
		// config file location string must be present; default file location string is provided to `flag`
		log.Panic("config file location not provided")
	}

	// config file location string is there but is the location valid
	if _, err := os.Stat(args.configFile); err != nil {
		log.Panic(err.Error())
	}

	if len(args.opLookup) == 0 && len(args.opSearchStartingWith) == 0 && len(args.opSearchEndingWith) == 0 && len(args.opAdd) == 0 {
		flag.PrintDefaults()
		log.Panic("no operation provided")
	}
}

func tryOperateExists(lxc lexicon.Lexicon) {
	words, err := wordSupplier.Get(args.opLookup)
	if len(words) == 0 || errors.Is(io.ErrNoInputValue, err) {
		return // this operation was not selected
	} else if err != nil {
		log.Printf("could not perform 'exists' for input (%s), error: %s\n", args.opLookup, err.Error())
	}

	if response, err := lxc.Lookup(words...); err == nil {
		outputPrinter.ConsumeWords("lookup", response)
	} else {
		log.Printf("could not perform 'exists' for input (%s), error: %s\n", args.opLookup, err.Error())
	}
}

func tryOperateGetAllStartingWith(lxc lexicon.Lexicon) {
	words, err := wordSupplier.Get(args.opSearchStartingWith)
	if len(words) == 0 || errors.Is(io.ErrNoInputValue, err) {
		return // this operation was not selected
	} else if err != nil {
		log.Printf("could not perform 'search starts with' for input (%s), error: %s\n", args.opSearchStartingWith, err.Error())
	}

	if searches, err := lxc.GetAllWordsStartingWith(words...); err == nil {
		outputPrinter.ConsumeMapOfWords("search starts with", searches)
	} else {
		log.Fatalf("could not perform 'search starts with' for input (%s), error: %s\n", args.opSearchStartingWith, err.Error())
	}
}

func tryOperateGetAllEndingWith(lxc lexicon.Lexicon) {
	words, err := wordSupplier.Get(args.opSearchEndingWith)
	if len(words) == 0 || errors.Is(io.ErrNoInputValue, err) {
		return // this operation was not selected
	} else if err != nil {
		log.Printf("could not perform 'search ends with' for input (%s), error: %s\n", args.opSearchEndingWith, err.Error())
	}


	if searches, err := lxc.GetAllWordsEndingWith(words...); err == nil {
		outputPrinter.ConsumeMapOfWords("search ends with", searches)
	} else {
		log.Fatalf("could not perform 'search ends with' for input (%s), error: %s\n", args.opSearchEndingWith, err.Error())
	}
}

func tryOperateAddAll(lxc lexicon.Lexicon) {
	words, err := wordSupplier.Get(args.opAdd)
	if len(words) == 0 || errors.Is(io.ErrNoInputValue, err) {
		return // this operation was not selected
	} else if err != nil {
		log.Printf("could not perform 'add' for input (%s), error: %s\n", args.opAdd, err.Error())
	}

	if err = lxc.Add(words...); err != nil {
		log.Fatalf("could not perform 'add' from file (%s), error: %s\n", args.opAdd, err.Error())
	} else {
		fmt.Println("add operation completed")
	}
}
