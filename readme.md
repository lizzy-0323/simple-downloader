# ANN dataset downloader

![License](https://img.shields.io/badge/License-MIT-blue.svg)

[简体中文](./README-CN.md)

This is a tool for downloading vector retrieval datasets, supporting the download of datasets such as sift1m, sift1b, deep1b, etc.

## Installation

1. Ensure that the Go language environment is installed.
2. Clone the project to your local machine:
   ```
   git clone <project URL>
   ```
3. Enter the project directory:
   ```
   cd downloader
   ```
4. Install dependencies:
   ```
   go mod tidy
   ```

## Usage

### Start the Tool

Run the following command to start the downloader:
```
go run main.go
```

### Download Datasets

Use the following command to download a specific dataset:
```
go run main.go download --dataset <dataset name>
```

Available dataset names include:
- sift1m
- sift1b
- deep1b

## Contributing

We welcome the submission of issues and pull requests to help us improve this project.

## License

This project is licensed under the MIT license.
