// Package lexicon defines an Lexicon interface.
package lexicon

// A Lexicon is an collection of words.
// Unlike dictionary, lexicon only stores words/string and no value (meaning).
// Like dictionary, various operation such as search or add can be performed on a Lexicon.
// A word is just a string in golang terms.
type Lexicon interface {
	// Lookup checks existence of the given words.
	// It returns array of strings of all the words that exists within the lexicon.
	// If any error occurs then it is returned; nil or empty words will return error.
	Lookup(words ...string) (*[]string, error)

	// GetAllWordsStartingWith will search given 'substrings' strings and return an array of all the words that start with the string.
	// Return value is a map where key is the 'substrings' string and value is array of matching words.
	// If any error occurs then it is returned; nil or empty words will return error.
	GetAllWordsStartingWith(substrings ...string) (*map[string][]string, error)

	// GetAllWordsEndingWith will search given 'substrings' strings and return an array of all the words that end with the string.
	// Return value is a map where key is the 'substrings' string and value is array of matching words.
	// If any error occurs then it is returned; nil or empty words will return error.
	GetAllWordsEndingWith(substrings ...string) (*map[string][]string, error)

	// Add adds the given array of words/string to current lexicon.
	// If failure occurs then error is returned; nil or empty words will return error.
	Add(words ...string) error

	// Close will close the lexicon.
	// Just like a book which is closed after usage.
	Close()
}
