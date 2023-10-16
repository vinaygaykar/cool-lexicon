package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vinaygaykar/cool-lexicon/configs"
	"github.com/vinaygaykar/cool-lexicon/pkg/lexicon"
)

var (
	checkSetup                                                                 bool
	cfgFileLoc, opExistsWord, opSearchStartWord, opSearchEndWord, opAddAllFile string
)

func init() {
	flag.BoolVar(&checkSetup, "check", false, "Setup all necessary configs if required. This is optional, if the all configs are already setup correctly this operation will have no effect")
	flag.StringVar(&cfgFileLoc, "cfg", "cool-lexicon-cfg.json", "Config file location")
	flag.StringVar(&opExistsWord, "ex", "", "Check if the given word exist")
	flag.StringVar(&opSearchStartWord, "ss", "", "Search the lexicon to find words that start with given substring")
	flag.StringVar(&opSearchEndWord, "se", "", "Search the lexicon to find words that end with given substring")
	flag.StringVar(&opAddAllFile, "ad", "", "Add words present in given file location to lexicon")
}

func main() {
	flag.Parse()

	sanitizeInputs()
	validateInputs()

	lxc := configs.GetLexicon(cfgFileLoc, checkSetup)
	defer lxc.Close()

	tryOperateExists(lxc)
	tryOperateGetAllStartingWith(lxc)
	tryOperateGetAllEndingWith(lxc)
	tryOperateAddAll(lxc)
}

func sanitizeInputs() {
	// remove any whitespaces
	cfgFileLoc = strings.TrimSpace(cfgFileLoc)
	opExistsWord = strings.TrimSpace(opExistsWord)
	opSearchStartWord = strings.TrimSpace(opSearchStartWord)
	opSearchEndWord = strings.TrimSpace(opSearchEndWord)
	opAddAllFile = strings.TrimSpace(opAddAllFile)
}

func validateInputs() {
	if len(cfgFileLoc) == 0 {
		// config file location string must be present; default file location string is provided to `flag`
		log.Panic("config file location not provided")
	}

	// config file location string is there but is the location valid
	if _, err := os.Stat(cfgFileLoc); err != nil {
		log.Panic(err.Error())
	}

	if len(opExistsWord) == 0 && len(opSearchStartWord) == 0 && len(opSearchEndWord) == 0 && len(opAddAllFile) == 0 {
		flag.PrintDefaults()
		log.Panic("no operation provided")
	}

	if len(opAddAllFile) != 0 { // Check if the given file exists
		if _, err := os.Stat(opAddAllFile); err != nil {
			log.Panic(err.Error())
		}
	}
}

func tryOperateExists(lxc lexicon.Lexicon) {
	if len(opExistsWord) == 0 {
		return
	}

	if exists, err := lxc.Lookup(opExistsWord); err != nil {
		log.Fatalf("could not perform 'exists' for the word (%s), error: %s\n", opExistsWord, err.Error())
	} else {
		fmt.Printf("exists (%s) : %t\n", opExistsWord, exists)
	}
}

func tryOperateGetAllStartingWith(lxc lexicon.Lexicon) {
	if len(opSearchStartWord) == 0 {
		return
	}

	if words, err := lxc.GetAllWordsStartingWith(opSearchStartWord); err != nil {
		log.Fatalf("could not perform 'starts with' for the word (%s), error: %s\n", opExistsWord, err.Error())
	} else {
		fmt.Printf("starts with (%s) : %v\n", opSearchStartWord, words)
	}
}

func tryOperateGetAllEndingWith(lxc lexicon.Lexicon) {
	if len(opSearchEndWord) == 0 {
		return
	}

	if words, err := lxc.GetAllWordsEndingWith(opSearchEndWord); err != nil {
		log.Fatalf("could not perform 'ends with' for the word (%s), error: %s\n", opExistsWord, err.Error())
	} else {
		fmt.Printf("ends with (%s) : %v\n", opSearchEndWord, words)
	}
}

func tryOperateAddAll(lxc lexicon.Lexicon) {
	if len(opAddAllFile) == 0 {
		return
	}

	// open file containing words to add
	file, err := os.Open(opAddAllFile)
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
		log.Fatalf("could not perform 'add words' from file (%s), error: %s\n", opAddAllFile, err.Error())
	} else {
		fmt.Printf("added words from the file (%s)\n", opAddAllFile)
	}
}
