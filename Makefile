.PHONY: all clean build-go build-c build-python build-jni test

GO_MOD = ee-sdk
BUILD_DIR = build
LIB_DIR = $(BUILD_DIR)/lib
INCLUDE_DIR = $(BUILD_DIR)/include

all: build-go build-c build-python build-jni

build-c:
	@echo "Building C bindings..."
	mkdir -p $(LIB_DIR) $(INCLUDE_DIR)
	CGO_ENABLED=1 go build -buildmode=c-shared -o $(LIB_DIR)/libee_sdk.so ./bindings/c
	CGO_ENABLED=1 go build -buildmode=c-archive -o $(LIB_DIR)/libee_sdk.a ./bindings/c
	mv $(LIB_DIR)/libee_sdk.h $(INCLUDE_DIR)/

build-jni: build-c
	@echo "Building JNI bindings..."
	mkdir -p $(LIB_DIR) $(INCLUDE_DIR)
	# 生成JNI头文件
	javac bindings/jni/src/main/java/com/turingq/eesdk/*.java && javac -h bindings/jni/src/main/c bindings/jni/src/main/java/com/turingq/eesdk/*.java
	# 编译JNI共享库
	gcc -shared -fPIC \
		-I"$(JAVA_HOME)/include" \
		-I"$(JAVA_HOME)/include/linux" \
		-I"$(INCLUDE_DIR)" \
		-L$(LIB_DIR) \
		-o $(LIB_DIR)/libee_sdk_jni.so \
		bindings/jni/src/main/c/*.c $(LIB_DIR)/libee_sdk.a \
		-lpthread -lm -ldl

	# 构建Java包
	cd bindings/jni && mvn clean package
	cp bindings/jni/target/eesdk*.jar $(BUILD_DIR)/

build-python: build-c
	@echo "Building Python bindings..."
	cp bindings/python/ee_sdk.py $(BUILD_DIR)/
	cp $(LIB_DIR)/libee_sdk.so $(BUILD_DIR)/

examples: build-go build-c build-python build-jni
	@echo "Running examples..."
	@echo "=== Go示例 ==="
	cd examples/golang && go run main.go
	@echo "\n=== C示例 ==="
	cd examples/c && \
		gcc -I../../build/include -L../../build/lib -lee_sdk main.c -o main && \
		LD_LIBRARY_PATH=../../build/lib ./main
	@echo "\n=== Python示例 ==="
	cd examples/python && PYTHONPATH=../../build:$$PYTHONPATH python3 main.py
	@echo "\n=== Java示例 ==="
	$(MAKE) example-java

example-go: build-go
	@echo "Running Go example..."
	cd examples/golang && go run main.go

example-c: build-c
	@echo "Running C example..."
	cd examples/c && \
		gcc -I../../build/include -L../../build/lib main.c  -lee_sdk -o ../../build/main && \
		LD_LIBRARY_PATH=../../build/lib ../../build/main

example-python: build-python
	@echo "Running Python example..."
	cd examples/python && PYTHONPATH=../../build:$$PYTHONPATH python3 main.py

example-java: build-jni
	@echo "Running Java example..."
	cd examples/java && \
		javac -cp ../../build/eesdk*.jar Main.java && \
		java -Djava.library.path=../../build/lib -cp "../../build/*:." Main

clean-examples:
	rm -f examples/c/main
	rm -f examples/java/*.class
	rm -f bindings/jni/src/main/java/com/turingq/eesdk/*.class
	rm -f bindings/jni/src/main/java/com_turingq_eesdk_EESdk.h

clean: clean-examples
	rm -rf $(BUILD_DIR)
	cd bindings/eesdk-java && mvn clean


