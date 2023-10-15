// Package lexicon defines an Lexicon interface.
package lexicon


// A Lexicon is an collection of words.
// Unlike dictionary, lexicon only stores words/string and no value (meaning).
// Like dictionary, various operation such as search or add can be performed on a Lexicon.
// A word is just a string in golang terms.
type Lexicon interface {

	// CheckIfExists returns true if the given word is present in current Lexicon, false otherwise.
	// If any problem occurs during lookup then non nil error is returned.
	CheckIfExists(word string) (bool, error)

	// GetAllStartingWith returns an array of words which all starts with toSearch string.
	// If any problem occurs during lookup then non nil error is returned.
	GetAllStartingWith(toSearch string) ([]string, error)
	
	// GetAllEndingWith returns an array of words which all ends with toSearch string.
	// If any problem occurs during lookup then non nil error is returned.
	GetAllEndingWith(toSearch string) ([]string, error)

	// AddAll adds the given array of words/string to current lexicon.
	// If any problem occurs during lookup then non nil error is returned.
	AddAll(words []string) error

	// Close will close the lexicon.
	// Just like a book which is closed after usage.
	Close()
}
