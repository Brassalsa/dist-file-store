# Distributed File Storage in Go

This project is a distributed file storage system implemented in Go. It allows for the storage and retrieval of files in a distributed manner across multiple nodes, enhancing both the reliability and accessibility of your data.

## Project Description

This distributed file storage system is designed with the aim of providing a robust and efficient method of storing files across multiple nodes. This ensures that your data is not only secure but also highly available, even in the event of individual node failures.

The system is implemented in Go, taking advantage of its efficiency and performance benefits, particularly in a networked context.

## Code Along Project

This project was developed with the help of a code along video tutorial. You can follow the same to understand the underlying concepts and see the code in action.

Code Along Video Tutorial [Here](https://www.youtube.com/watch?v=bymQakvTY40)

## Features

- **Distributed Storage**: Files are stored across multiple nodes, ensuring high availability and reliability.
- **Fault Tolerance**: The system is designed to handle individual node failures without any loss of data.
- **Efficient Retrieval**: Files can be retrieved from the nearest node, ensuring efficient access.

## Getting Started

Here you can add instructions on how to get a copy of the project up and running on the local machine for development and testing purposes.

## Prerequisites

- Go lang 1.22 or above (required)
- Makefile tools (optional)

## Installing

Clone the repo and run with makefile

```bash
make run
```

OR

```bash
go build -o bin/fs
./bin/fs
```

## Running the tests

For running test with makefile

```bash
make test
```

OR

```bash
go test ./... -v
```

## Built With

- Go - The programming language used.

## Acknowledgments

- Anthony GG youtube channel [@anthonygg\_](https://www.youtube.com/@anthonygg_)
