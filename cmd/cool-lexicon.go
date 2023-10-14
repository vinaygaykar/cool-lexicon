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

var (
	cfgFileLoc, opExistsWord, opSearchStartWord, opSearchEndWord, opAddAllFile string
)

func init() {
	flag.StringVar(&cfgFileLoc, "cfg", "cool-lexicon-cfg.json", "Config file location")
	flag.StringVar(&opExistsWord, "ex", "", "Check if the given word exist")
	flag.StringVar(&opSearchStartWord, "ss", "", "Search the lexicon to find words that start with given substring")
	flag.StringVar(&opSearchEndWord, "se", "", "Search the lexicon to find words that end with given substring")
	flag.StringVar(&opAddAllFile, "ad", "", "Add words present in given file location to lexicon")
}

func main() {
	flag.Parse()

	sanitizeInput()
	if err := validateInput(); err != nil {
		return
	}

	lxc, err := configs.GetLexicon(cfgFileLoc)
	if err != nil {
		fmt.Printf("Error initilising the program. %s\n", err.Error())
	}
	defer lxc.Close()

	operateExists(lxc)
	operateGetAllStartingWith(lxc)
	operateGetAllEndingWith(lxc)
	operateAddAll(lxc)
}

func sanitizeInput() {
	// remove any whitespaces
	cfgFileLoc = strings.TrimSpace(cfgFileLoc)
	opExistsWord = strings.TrimSpace(opExistsWord)
	opSearchStartWord = strings.TrimSpace(opSearchStartWord)
	opSearchEndWord = strings.TrimSpace(opSearchEndWord)
	opAddAllFile = strings.TrimSpace(opAddAllFile)
}

func validateInput() error {
	if len(cfgFileLoc) == 0 {
		// config file location string must be present; default file location string is provided to `flag`
		return errors.New("config file location not provided")
	}

	// config file location string is there but is the location valid
	if _, err := os.Stat(cfgFileLoc); err != nil {
		log.Printf("%s", cfgFileLoc)
		fmt.Println("config file not present at the provided location")
		return err
	}

	if len(opExistsWord) == 0 && len(opSearchStartWord) == 0 && len(opSearchEndWord) == 0 && len(opAddAllFile) == 0 {
		flag.PrintDefaults()
		return errors.New("no operation flag provided")
	}

	if len(opAddAllFile) != 0 { // Check if the given file exists
		if _, err := os.Stat(opAddAllFile); err != nil {
			fmt.Println("file location provided for add all operation is invalid")
			return err
		}
	}

	return nil
}

func operateExists(lxc lexicon.Lexicon) {
	if len(opExistsWord) == 0 {
		return
	}

	exists := lxc.CheckIfExists(string(opExistsWord))
	fmt.Printf("exists (%s) : %t\n", opExistsWord, exists)
}

func operateGetAllStartingWith(lxc lexicon.Lexicon) {
	if len(opSearchStartWord) == 0 {
		return
	}

	words := lxc.GetAllStartingWith(opSearchStartWord)
	fmt.Printf("starts with (%s) : %v\n", opSearchStartWord, words)
}

func operateGetAllEndingWith(lxc lexicon.Lexicon) {
	if len(opSearchEndWord) == 0 {
		return
	}

	words := lxc.GetAllEndingWith(opSearchEndWord)
	fmt.Printf("ends with (%s) : %v\n", opSearchEndWord, words)
}

func operateAddAll(lxc lexicon.Lexicon) {
	if len(opAddAllFile) == 0 {
		return
	}

	file, err := os.Open(opAddAllFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	words := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, string(strings.TrimSpace(scanner.Text())))
	}

	if err2 := scanner.Err(); err2 != nil {
		log.Fatal(err2)
	}

	fmt.Printf("adding words from file (%s)\n", opAddAllFile)
	lxc.AddAll(words)
}
