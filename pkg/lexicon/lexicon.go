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

	// Lookup checks existence of the given words and returns an array of bool indicating their
	// presence in the lexicon.
	// If any problem occurs during lookup then non nil error is returned.
	// If and error occurs after checking some number of words then the response will have existence
	// proof of all the checked words and error together. In this scenarion size of returned bool array
	// is less than the number of words.
	Lookup(words ...string) ([]bool, error)

	// GetAllStartingWith returns an array of words which all starts with toSearch string.
	// If any problem occurs during lookup then non nil error is returned.
	GetAllStartingWith(toSearch string) ([]string, error)

	// SearchForStartingWith will search given 'substrings' strings and return an array of all the words that start with the string.
	// Return value is a map where key is the 'substrings' string and value is array of matching words.
	// If any problem occurs during lookup then non nil error is returned.
	// If an error occurs in between while some substrings are searched and others are pending then the
	// return map will have only succesfully searched substrings aong with non nil error.
	SearchForStartingWith(substrings ...string) (map[string][]string, error)

	// GetAllEndingWith returns an array of words which all ends with toSearch string.
	// If any problem occurs during lookup then non nil error is returned.
	GetAllEndingWith(toSearch string) ([]string, error)

	// SearchForEndingWith will search given 'substrings' strings and return an array of all the words that end with the string.
	// Return value is a map where key is the 'substrings' string and value is array of matching words.
	// If any problem occurs during lookup then non nil error is returned.
	// If an error occurs in between while some substrings are searched and others are pending then the
	// return map will have only succesfully searched substrings aong with non nil error.
	SearchForEndingWith(substrings ...string) (map[string][]string, error)

	// Add adds the given array of words/string to current lexicon.
	// If any problem occurs during lookup then non nil error is returned.
	Add(words ...string) error

	// Close will close the lexicon.
	// Just like a book which is closed after usage.
	Close()
}
