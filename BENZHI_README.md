# expr-eval — Go 数学表达式解析、编译与求值 HTTP 服务（含 VM 字节码执行）

读入算术/比较/逻辑表达式与变量绑定，先编成栈机字节码再求值，支持数值、布尔与字符串。除零、未绑定变量或类型不能强制转换时返回错误；同一表达式树遍历与字节码 VM 必须给出相同结果，不得只列出入口。

## 构建 / 运行 / 测试

```text
go build ./...     # 编译
go run .
go test ./...      # 测试
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
