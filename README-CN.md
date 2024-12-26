# 下载器项目

该项目是一个用于下载向量检索数据集的工具，支持下载sift1m、sift1b、deep1b等数据集。

## 安装

1. 确保已安装Go语言环境。
2. 克隆该项目到本地：
   ```
   git clone <项目地址>
   ```
3. 进入项目目录：
   ```
   cd downloader
   ```
4. 安装依赖：
   ```
   go mod tidy
   ```

## 使用

### 启动工具

运行以下命令启动下载器：
```
go run main.go
```

### 下载数据集

使用以下命令下载特定的数据集：
```
go run main.go download --dataset <数据集名称>
```

可用的数据集名称包括：
- sift1m
- sift1b
- deep1b

## 贡献

欢迎提交问题和拉取请求，帮助我们改进该项目。

## 许可证

该项目遵循MIT许可证。