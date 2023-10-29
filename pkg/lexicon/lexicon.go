// Package lexicon defines an Lexicon interface.
package lexicon

// A OperationResult represents result of lexicon operation performed on an individual word.
// It is a wrapper for value-error pair where either of the two will be valid.
type OperationResult[V bool | []string] struct {
	// Err specifies if this poeration resulted in an error. If err is non nil then value should be considered useless
	Err error 	`json:"error"`

	// Value is result of the operation performed. Value will be garbage or nil if err is non nil
	Value V		`json:"value"`
}

// A Lexicon is an collection of words.
// Unlike dictionary, lexicon only stores words/string and no value (meaning).
// Like dictionary, various operation such as search or add can be performed on a Lexicon.
// A word is just a string in golang terms.
type Lexicon interface {

	// Lookup checks existence of the given words and returns a pointer to the map of OperationResult indicating their presence in the lexicon.
	// For every word passed, it will have a corresponding entry in the result map with key being the word itself and the value
	// being the an instance of OperationResult.
	// Existence of a word will be represented using boolean value within OperationResult,
	// if an error occurs during lookup then err value will be present.
	// If any critical error occurs then the error is returned and in such cases value of the map should not be trusted.
	// If words are nil then map is empty and error is returned.
	// If words are empty then map is empty and error is nil.
	Lookup(words ...string) (*map[string]OperationResult[bool], error)

	// GetAllWordsStartingWith will search given 'substrings' strings and return an array of all the words that start with the string.
	// Return value is a map where key is the 'substrings' string and value is array of matching words.
	// If any problem occurs during lookup then non nil error is returned.
	// If an error occurs in between while some substrings are searched and others are pending then the
	// return map will have only succesfully searched substrings aong with non nil error.
	GetAllWordsStartingWith(substrings ...string) (map[string][]string, error)

	// GetAllWordsEndingWith will search given 'substrings' strings and return an array of all the words that end with the string.
	// Return value is a map where key is the 'substrings' string and value is array of matching words.
	// If any problem occurs during lookup then non nil error is returned.
	// If an error occurs in between while some substrings are searched and others are pending then the
	// return map will have only succesfully searched substrings aong with non nil error.
	GetAllWordsEndingWith(substrings ...string) (map[string][]string, error)

	// Add adds the given array of words/string to current lexicon.
	// If any problem occurs during lookup then non nil error is returned.
	Add(words ...string) error

	// Close will close the lexicon.
	// Just like a book which is closed after usage.
	Close()
}
