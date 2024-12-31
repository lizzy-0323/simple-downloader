# A simple downloader in Go

![License](https://img.shields.io/badge/License-MIT-blue.svg)

[简体中文](./README-CN.md)

This is a simple downloader in Go. It downloads a file from a given URL and saves it to a specified file path. Currently, it supports HTTP and FTP protocols.

## Installation

You can install the project by following these steps:

1. Clone the repository:

    ```sh
    git clone <repository-url>
    ```

2. Navigate to the project directory:

    ```sh
    cd go-downloader
    ```

3. Build the project:

    ```sh
    make build
    ```

## Usage

You can download files using the following command:

```bash
./go-downloader download [URL] -d [destination]
```

For example, to download a file to the current directory:

```bash
./go-downloader download http://example.com/file.zip -d .
```

## Options

-d, --dst : Specify the destination directory for the downloaded file, default is the current directory.

-w, --workers : Specify the number of concurrent worker threads for downloading, default is 16.

## Contributing

Contributions are welcome! Please submit a Pull Request or report an issue.

## License

This project is licensed under the MIT License. See the LICENSE file for details. 
