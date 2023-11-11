# cool-lexicon

The Lexicon project is a Golang application designed to create and manage a lexicon of Devnagri words.
Current version of the project relies upon MySQL to store all the words because -- simplicity.

## Table of Contents
- [Context](https://github.com/vinaygaykar/cool-lexicon/edit/tech/docs/README.md#context)
- [User Scenarios Supported](https://github.com/vinaygaykar/cool-lexicon/edit/tech/docs/README.md#user-scenarios-supported)
- [Getting Started](https://github.com/vinaygaykar/cool-lexicon/edit/tech/docs/README.md#context)
- [Troubleshooting](https://github.com/vinaygaykar/cool-lexicon/edit/tech/docs/README.md#context)
- [Contributing](https://github.com/vinaygaykar/cool-lexicon/edit/tech/docs/README.md#context)
- [License](https://github.com/vinaygaykar/cool-lexicon/edit/tech/docs/README.md#context)

## Context

The Lexicon project aims to provide a user-friendly CLI interface to work with a lexicon of Devnagri words. It offers various operations to search and manage words efficiently.

## User Scenarios Supported

### 0. Setup the DB or validate existing setup
As a user, during the first run, it is possible that necessary database and tables are not setup, to do that use the `-check` flag which will perform necessary
operation if required. Any failures here would need manual intervention.

Usage
```console
  ./cool-lexicon -check -ex word-to-search
```

#### 0.1 File based & CLI inputs

Inputs can be given from following sources for all provided operations.

1. **CLI** : (Default) Using the word directly from the terminal. 
In the following example the operation exists uses `नमस्कार` and operation add uses `धन्यवाद` as input word.
```console
  ./cool-lexicon -ex नमस्कार -ad धन्यवाद
```
The limitation to this is that only one word can be given as input to one operation.


2. **File** : Using the words from the provided text file path. Use the `-if` flag to indicate this option. 
In the following example the operation exists uses the file `./path-to/file1.txt` and add uses `./file2.txt` as input, given the file exists.
```console
  ./cool-lexicon -ex ./path-to/file.txt -ad ./file2.txt
```
There are some requirements, 
  a. File should exists at given place and is a valid text file
  b. File should have required access to be read by the program
  c. Words are space delimited and a line should not be more than 64K characters long
  d. Once the `-if` flag is used input for all operations is streamed from file

#### 0.2 File base & CLI output

Output can be streamed to either of the places for _all the operations_.

1. **CLI** : (Default) Using terminal to display result of every operation.

2. **File** : Writing output of every sepcified operation to individual files at provided location. Use the `-of` flag and provided expected location of the output.
In the following example, output of both the operations will be written to `./output-path` location under different file for each operation.
```console
  ./cool-lexicon -of ./output-path -ex नमस्कार -ad धन्यवाद
```
There are some requirements,
  a. Output folder should exists
  b. If file exists at the output location with name of the operation then it will be overwritten
  c. Program should have access to the output location
  d. Once the `-of` flag is used output for all operations is streamed to file



### 1. Check if a word exists

As a user, you can check if a given word exists in the lexicon by using the `-ex` operation. 

Usage
```console
  ./cool-lexicon -ex word-to-search
```

### 2. Search words that start with a substring

As a user, you can retrieve a list of words that start with a given substring using the `-ss` operation. The result will be a sorted list of words that match the provided substring.

Usage
```console
  ./cool-lexicon -ss substring-to-check
```

### 3. Search words that end with a substring

As a user, you can find words that end with a specific substring, use the `-se` operation. It will return a sorted list of words that match the provided substring.

Usage
```console
  ./cool-lexicon -se substring-to-check
```

### 4. Add words to the lexicon

As a user, you can add a batch of new words to the lexicon by providing a text file with each word on a new line and using the `-ad` operation. 

Usage
```console
  ./cool-lexicon -ad location/of/file/that/contains-words.txt
```

## Getting Started

To get started with the Lexicon project, follow these steps:

- Clone the project repository to your local machine
  `git clone https://github.com/vinaygaykar/cool-lexicon.git`

- Install the required dependencies
  `go get`

- Set up the MySQL database server, make sure you have MySQL installed and running (Using docker will be less hassle)

- Configure the database connection in the `config.json` file, make sure the file is present at root level of the project

- Run the application.
  `go run main.go`

To create and run a binary:

- Run `go build ./cmd/cool-lexicon.go`, this will create a executable named `cool-lexicon` depending upon your os & arch.
- Make sure MySQL is running
- Make sure the config file `config.json` is present at same level as that of the executable and has valid & working db connection values
- Execute the binary, check [User Scenarios Supported](https://github.com/vinaygaykar/cool-lexicon/edit/tech/docs/README.md#user-scenarios-supported) for supported operations

## Troubleshooting

- If you encounter any issues with the database setup, ensure that the database is correctly configured and that you have the necessary privileges.

- To provide different configs use the `-cfg` flag and pass new config file location. Usage

  ```console
    ./cool-lexicon -cfg location/to/diffent/config.json
  ```

## Contributing

If you'd like to contribute to the Lexicon project, please open an issue or submit a pull request. We welcome any improvements or feature enhancements.
License

This project is licensed under the Apache 2.0 License.
