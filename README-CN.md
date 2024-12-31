# A simple downloader in Go

![License](https://img.shields.io/badge/License-MIT-blue.svg)

[English](./README.md)

这是一个简单的 Go 语言下载器。它可以从一个给定的 URL 下载文件，并将其保存到指定的文件路径。目前支持 HTTP 和 FTP 协议。

## 安装

你可以通过以下步骤安装该项目：

```sh
git clone <仓库地址>
```

导航到项目目录：

```sh
cd go-downloader
```

构建项目：

```sh
make build
```

## 使用方法

你可以使用以下命令下载文件：

```bash
./go-downloader download [URL] -d [目标路径]
```

例如，要将文件下载到当前目录：

```bash
./go-downloader download http://example.com/file.zip -d .
```

## 选项

-d, --dst : 指定下载文件的目标目录，默认为当前目录。

-w, --workers : 指定用于下载的并发工作线程数量，默认为 16。

## 贡献

欢迎贡献！请提交 Pull Request 或报告问题。

## 许可证

本项目根据 MIT 许可证授权。详情请参阅 LICENSE 文件。