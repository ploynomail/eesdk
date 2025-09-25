# EE SDK 使用示例

本目录包含了EE SDK在各种编程语言中的使用示例。

## 目录结构

```
examples/
├── golang/          # Go语言示例
├── c/              # C语言示例
├── python/         # Python示例
├── java/           # Java示例
└── README.md       # 本文档
```

## 运行示例

### Go语言示例

```bash
cd examples/golang
go run main.go
```

### C语言示例

首先确保已经构建了C绑定：

```bash
make build-c
cd examples/c
gcc -I../../build/include -L../../build/lib -lee_sdk main.c -o main
LD_LIBRARY_PATH=../../build/lib ./main
```

### Python示例

首先确保已经构建了Python绑定：

```bash
make build-python
cd examples/python
python3 main.py
```

### Java示例

首先确保已经构建了JNI绑定：

```bash
make build-jni
cd examples/java
javac -cp ../../bindings/jni Main.java
java -Djava.library.path=../../build/lib -cp .:../../bindings/jni Main
```

## 示例功能

所有示例都包含以下功能演示：

1. **客户端初始化** - 如何配置和创建SDK客户端
2. **健康检查** - 检查服务可用性
3. **GET请求** - 获取数据
4. **POST请求** - 创建资源
5. **PUT请求** - 更新资源
6. **错误处理** - 如何处理各种错误情况

## 配置说明

### 基本配置项

- `base_url`: API服务的基础URL
- `api_key`: API访问密钥
- `timeout`: 请求超时时间（秒）
- `retry_count`: 重试次数

### 环境变量

可以通过环境变量配置：

```bash
export EE_SDK_BASE_URL="https://api.example.com"
export EE_SDK_API_KEY="your-api-key"
export EE_SDK_TIMEOUT="30"
export EE_SDK_RETRY_COUNT="3"
```

## 最佳实践

1. **错误处理**: 总是检查和处理可能的错误
2. **资源清理**: 在C/JNI绑定中确保正确释放资源
3. **超时设置**: 根据网络环境合理设置超时时间
4. **重试机制**: 对于临时性错误，合理使用重试
5. **日志记录**: 在生产环境中添加适当的日志记录

## 故障排除

### 常见问题

1. **库文件找不到**: 确保LD_LIBRARY_PATH包含共享库路径
2. **API密钥错误**: 检查API密钥是否正确配置
3. **网络连接问题**: 检查网络连接和防火墙设置
4. **超时错误**: 适当增加超时时间设置

### 调试技巧

1. 启用详细日志输出
2. 使用网络抓包工具检查HTTP请求
3. 检查服务端点的可用性
4. 验证请求格式是否正确
