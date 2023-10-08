package pkg

type Word string

type Lexicon interface {
	CheckIfExists(word Word) bool

	GetAllStartingWith(toSearch string) []Word

	GetAllEndingWith(toSearch string) []Word

	Add(word Word)

	AddAll(words []Word)
}
