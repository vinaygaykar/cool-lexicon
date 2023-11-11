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


### 0. Setup & basic scenarios

#### 0.1 Setup the DB or validate existing setup
As a user, during the first run, it is possible that necessary database and tables are not setup, to do that use the `-check` flag which will perform necessary
operation if required. Any failures here would need manual intervention.

Usage
```console
  ./lxc -check -ex नमस्कार
```

#### 0.2 File based & CLI inputs

Inputs can be given from following sources for all provided operations.

1. **CLI** : (Default) Using the word directly from the terminal. 
In the following example the operation exists uses `नमस्कार` and operation add uses `धन्यवाद` as input word.
```console
  ./lxc -ex नमस्कार -ad धन्यवाद
```
The limitation to this is that only one word can be given as input to one operation.


2. **File** : Using the words from the provided text file path. Use the `-if` flag to indicate this option. 
In the following example the operation exists uses the file `./path-to/file1.txt` and add uses `./file2.txt` as input, given the file exists.
```console
  ./lxc -if -ex ./path-to/file.txt -ad ./file2.txt
```
There are some requirements, 
  - File should exists at given place and is a valid text file
  - File should have required access to be read by the program
  - Words are space delimited and a line should not be more than 64K characters long
  - Once the `-if` flag is used input for all operations is streamed from file

#### 0.3 File base & CLI output

Output can be streamed to either of the places for _all the operations_.

1. **CLI** : (Default) Using terminal to display result of every operation.

2. **File** : Writing output of every sepcified operation to individual files at provided location. Use the `-of` flag and provided expected location of the output.
In the following example, output of both the operations will be written to `./output-path` location under different file for each operation.
```console
  ./lxc -of ./output-path -ex नमस्कार -ad धन्यवाद
```
There are some requirements,
  - Output folder should exists
  - If file exists at the output location with name of the operation then it will be overwritten
  - Program should have access to the output location
  - Once the `-of` flag is used output for all operations is streamed to file


### 1. Check if a word exists

As a user, you can check if a given word exists in the lexicon by using the `-ex` operation. 

Usage
```console
  ./lxc -ex नमस्कार
```


### 2. Search words that start with a substring

As a user, you can retrieve a list of words that start with a given substring using the `-ss` operation. The result will be a sorted list of words that match the provided substring.

Usage
```console
  ./lxc -ss नम
```


### 3. Search words that end with a substring

As a user, you can find words that end with a specific substring, use the `-se` operation. It will return a sorted list of words that match the provided substring.

Usage
```console
  ./lxc -se कार
```


### 4. Add words to the lexicon

As a user, you can add new words to the lexicon using the `-ad` operation. 

Usage
```console
  ./lxc -ad धन्यवाद
```



## Getting Started

To get started with the Lexicon project, follow these steps:

- Clone the project repository to your local machine
  `git clone https://github.com/vinaygaykar/cool-lexicon.git`

- Install the required dependencies
  `go get`

- Set up the MySQL database server, make sure you have MySQL installed and running (Using docker will be less hassle)

- Configure the database connection in the `config.json` file, make sure the file is present at root level of the project

- (Optinal) Test
`got test --timeout 5m ./...`

- Run the application.
  `go run main.go`

To create and run a binary:

- Run `go build ./cmd/main.go -o lxc`, this will create a executable named `lxc` depending upon your os & arch.
- Make sure MySQL is running
- Make sure the config file `config.json` is present at same level as that of the executable and has valid & working db connection values
- Execute the binary, check [User Scenarios Supported](https://github.com/vinaygaykar/cool-lexicon/edit/tech/docs/README.md#user-scenarios-supported) for supported operations



## Troubleshooting

- If you encounter any issues with the database setup, ensure that the database is correctly configured and that you have the necessary privileges.

- To provide different configs use the `-cfg` flag and pass new config file location. Usage

  ```console
    ./lxc -cfg location/to/diffent/config.json
  ```



## Contributing

If you'd like to contribute to the Lexicon project, please open an issue or submit a pull request. We welcome any improvements or feature enhancements.
License

This project is licensed under the Apache 2.0 License.
