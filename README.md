# README.md for Customer Importer CLI Command

## Overview

Customer Importer is a command-line interface (CLI) application designed to process CSV files containing customer data. It reads customer data from a CSV file, counts the number of customers with email addresses from each domain, and outputs the results either to the terminal or to a file. The application can operate in either sequential or concurrent mode.

## Installation

Before running the application, ensure you have Go installed on your system. You can download and install Go from [the official Go website](https://golang.org/dl/). (Required go modules support)

Clone the repository to your local machine:
```
git clone [repository URL]
cd customerimporter
```

## Usage

Run the application using the `go run` command from the root directory of the project:

```
go run cmd/customerimporter/main.go -path=[path to CSV file] -output=[output directory] -mode=[processing mode]
```


### Flags

- `-path`: Specifies the path to the CSV file to be processed. (Required)
- `-output`: Specifies the path to the output directory. If not provided, the output will be printed to stdout. (Optional)
- `-mode`: Specifies the mode of processing. Can be `sequentially` or `concurrency`. Default is `concurrency`. (Optional)

### Example

```
go run cmd/customerimporter/main.go -path="./data/customers.csv" -output="./data" -mode="concurrency"
```

This command will process the file `customers.csv` in concurrent mode and save the output in the `./data` directory. Output file format is `output-${UnixTimestamp}.csv`

## Running Tests

The application includes unit tests for various components. To run these tests, navigate to the specific package directory and use the `go test` command.

For example, to run tests in the `reader` package:
```
cd pkg/reader
go test
```

Repeat this process for other packages like `email` to execute their respective tests.