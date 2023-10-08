package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/vinaygaykar/cool-lexicon/internal"
	"github.com/vinaygaykar/cool-lexicon/pkg"
)

var (
    opExistsWord, opSearchStartWord, opSearchEndWord, opAddWord, opAddAllFile string
)

func init() {
    flag.StringVar(&opExistsWord, "ex", "", "Check if the given word exist")
    flag.StringVar(&opSearchStartWord, "ss", "", "Search the lexicon to find words that start with given substring")
    flag.StringVar(&opSearchEndWord, "se", "", "Search the lexicon to find words that end with given substring")
    flag.StringVar(&opAddWord, "ad", "", "Add the word to lexicon")
    flag.StringVar(&opAddAllFile, "aa", "", "Add words present in given file location to lexicon")
}

func main() {
    flag.Parse()

    cleanupInput()
    
    lxc := internal.NewLexicon()

    if len(opExistsWord) != 0 {
        lxc.CheckIfExists(pkg.Word(opExistsWord))
    } else if len(opSearchStartWord) != 0 {
        lxc.GetAllStartingWith(opSearchStartWord)
    } else if len(opSearchEndWord) != 0 {
        lxc.GetAllEndingWith(opSearchEndWord)
    } else if len(opAddWord) != 0 {
        lxc.Add(pkg.Word(opAddWord))
    } else if len(opAddAllFile) != 0 {
        _, err := os.Stat(opAddAllFile)
        if err != nil {
            lxc.AddAll(getWordsFromFile(opAddAllFile))
        } else {
            log.Fatalf("File does not exits at the path or is corrupted, %s", opAddAllFile)
        }
    } else {
        flag.PrintDefaults()
    }
}

func getWordsFromFile(opAddAllFile string) []pkg.Word {
	file, err := os.Open(opAddAllFile)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    words := make([]pkg.Word, 0)

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        words = append(words, pkg.Word(strings.TrimSpace(scanner.Text())))
    }

    if err2 := scanner.Err(); err != nil {
        log.Fatal(err2)
    }

    return words
}

func cleanupInput() {
    opExistsWord = strings.TrimSpace(opExistsWord)
    opSearchStartWord = strings.TrimSpace(opSearchStartWord)
    opSearchEndWord = strings.TrimSpace(opSearchEndWord)
    opAddWord = strings.TrimSpace(opAddWord)
    opAddAllFile = strings.TrimSpace(opAddAllFile)
}
