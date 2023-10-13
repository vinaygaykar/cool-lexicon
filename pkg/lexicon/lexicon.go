package lexicon

type Lexicon interface {
	CheckIfExists(word string) bool

	GetAllStartingWith(toSearch string) []string

	GetAllEndingWith(toSearch string) []string

	AddAll(words []string)

	Close()
}
