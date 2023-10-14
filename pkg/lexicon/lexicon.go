package lexicon

type Lexicon interface {
	CheckIfExists(word string) (bool, error)

	GetAllStartingWith(toSearch string) ([]string, error)
	
	GetAllEndingWith(toSearch string) ([]string, error)

	AddAll(words []string) error

	Close()
}
