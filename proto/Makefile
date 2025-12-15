.PHONY: all clean local

# 定义目录
PROTO_DIR := ./definitions
GO_GEN_DIR := ./gen/go
TS_GEN_DIR := ./gen/ts
PY_GEN_DIR := ./gen/py
LOCAL_GO_GEN_DIR := ./local-gen/go
LOCAL_TS_GEN_DIR := ./local-gen/ts
LOCAL_PY_GEN_DIR := ./local-gen/py

# 自动发现所有模块目录（包括多级目录）
MODULES := $(patsubst $(PROTO_DIR)/%/,%,$(wildcard $(PROTO_DIR)/*/)) 

# 查找所有 .proto 文件
PROTO_FILES := $(shell find $(PROTO_DIR) -name "*.proto")

# 提取所有模块路径（支持多级目录）
MODULE_PATHS := $(sort $(dir $(PROTO_FILES)))

# 转换为相对于 PROTO_DIR 的模块名
NESTED_MODULES := $(patsubst $(PROTO_DIR)/%,%,$(MODULE_PATHS))

# 移除尾部斜杠
NESTED_MODULES := $(patsubst %/,%,$(NESTED_MODULES))

# 默认目标：编译所有 proto 文件（包括嵌套模块）
all: prepare $(NESTED_MODULES)

# 本地开发目标：编译所有 proto 文件但不提交
local: local-prepare local-modules

# 本地模块编译（遍历所有模块）
local-modules:
	@for module in $(NESTED_MODULES); do \
		$(MAKE) local-module MODULE=$$module; \
	done

# 清理生成的代码
clean:
	rm -rf $(GO_GEN_DIR)
	rm -rf $(TS_GEN_DIR)
	rm -rf $(PY_GEN_DIR)
	rm -rf $(LOCAL_GO_GEN_DIR)
	rm -rf $(LOCAL_TS_GEN_DIR)
	rm -rf $(LOCAL_PY_GEN_DIR)

# 准备输出目录
prepare:
	@echo "Found modules: $(NESTED_MODULES)"
	@for module in $(NESTED_MODULES); do \
		mkdir -p $(GO_GEN_DIR)/$$module; \
		mkdir -p $(TS_GEN_DIR)/$$module; \
		mkdir -p $(PY_GEN_DIR)/$$module; \
	done
	@touch $(PY_GEN_DIR)/__init__.py

# 准备本地输出目录
local-prepare:
	@echo "Found modules: $(NESTED_MODULES)"
	@for module in $(NESTED_MODULES); do \
		mkdir -p $(LOCAL_GO_GEN_DIR)/$$module; \
		mkdir -p $(LOCAL_TS_GEN_DIR)/$$module; \
		mkdir -p $(LOCAL_PY_GEN_DIR)/$$module; \
	done
	@touch $(LOCAL_PY_GEN_DIR)/__init__.py

# 添加 local-gen 到 .gitignore
	@if ! grep -q "^/local-gen/" .gitignore; then \
		echo "/local-gen/" >> .gitignore; \
	fi

# 通用规则：编译每个模块的 proto 文件（支持多级目录）
$(NESTED_MODULES): prepare
	@echo "Compiling module: $@"
	@if [ -d "$(PROTO_DIR)/$@" ] && [ -n "`ls -A $(PROTO_DIR)/$@/*.proto 2>/dev/null`" ]; then \
		echo "Generating Go code for $@"; \
		set -e; \
		protoc -I=. \
			--go_out=. \
			--go_opt=paths=source_relative \
			--go-grpc_out=. \
			--go-grpc_opt=paths=source_relative \
			--validate_out="lang=go,paths=source_relative:." \
			$(PROTO_DIR)/$@/*.proto || { echo "Error: protoc failed for Go code generation"; exit 1; }; \
		mkdir -p $(GO_GEN_DIR)/$@; \
		mv $(PROTO_DIR)/$@/*.pb.go $(GO_GEN_DIR)/$@/ 2>/dev/null || true; \
		for f in `find $(PROTO_DIR)/$@ -name "*.validate.pb.go" -o -name "*.validate.go" -o -name "*_validate.pb.go" 2>/dev/null`; do \
			mv $$f $(GO_GEN_DIR)/$@/ 2>/dev/null || true; \
		done; \
		echo "Successfully compiled Go code for module: $@"; \
		echo "Checking for TypeScript plugin"; \
		if command -v protoc-gen-ts > /dev/null 2>&1; then \
			echo "Generating TypeScript code for $@"; \
			protoc -I=. \
				--plugin=protoc-gen-ts=`which protoc-gen-ts` \
				--ts_out="use_proto_field_name:$(TS_GEN_DIR)" \
				$(PROTO_DIR)/$@/*.proto || { echo "Error: protoc failed for TypeScript code generation"; exit 1; }; \
			if [ -d "$(TS_GEN_DIR)/definitions/$@" ]; then \
				mkdir -p $(TS_GEN_DIR)/$@; \
				for f in `find $(TS_GEN_DIR)/definitions/$@ -name "*.ts" 2>/dev/null`; do \
					filename=`basename $$f`; \
					cp $$f $(TS_GEN_DIR)/$@/$$filename; \
				done; \
				rm -rf $(TS_GEN_DIR)/definitions; \
			fi; \
			echo "Successfully compiled TypeScript code for module: $@"; \
		else \
			echo "Warning: protoc-gen-ts not found, skipping TypeScript generation for $@"; \
			echo "To install: npm install --global @protobuf-ts/plugin"; \
		fi; \
		echo "Generating Python code for $@"; \
		protoc -I=. \
			--python_out=$(PY_GEN_DIR) \
			$(PROTO_DIR)/$@/*.proto || { echo "Error: protoc failed for Python code generation"; exit 1; }; \
		if python3 -c "import grpc_tools.protoc" 2>/dev/null; then \
			echo "Generating Python gRPC code for $@"; \
			python3 -m grpc_tools.protoc -I=. \
				--grpc_python_out=$(PY_GEN_DIR) \
				$(PROTO_DIR)/$@/*.proto || { echo "Error: grpc_tools.protoc failed for Python gRPC code generation"; exit 1; }; \
		else \
			echo "Warning: grpcio-tools not found, skipping Python gRPC generation for $@"; \
			echo "To install: pip install grpcio-tools"; \
		fi; \
		if [ -d "$(PY_GEN_DIR)/definitions/$@" ]; then \
			mkdir -p $(PY_GEN_DIR)/$@; \
			for f in `find $(PY_GEN_DIR)/definitions/$@ -name "*.py" 2>/dev/null`; do \
				filename=`basename $$f`; \
				cp $$f $(PY_GEN_DIR)/$@/$$filename; \
			done; \
			rm -rf $(PY_GEN_DIR)/definitions; \
		fi; \
		touch $(PY_GEN_DIR)/$@/__init__.py; \
		python3 fix_py_imports.py $(PY_GEN_DIR); \
		echo "Successfully compiled Python code for module: $@"; \
	else \
		echo "No proto files found for module: $@"; \
	fi

# 本地开发规则：编译单个模块的 proto 文件到 local-gen 目录（通过 MODULE 变量传入）
local-module:
	@echo "Compiling module: $(MODULE)"
	@if [ -d "$(PROTO_DIR)/$(MODULE)" ] && [ -n "`ls -A $(PROTO_DIR)/$(MODULE)/*.proto 2>/dev/null`" ]; then \
		echo "Generating Go code for $(MODULE)"; \
		set -e; \
		protoc -I=. \
			--go_out=. \
			--go_opt=paths=source_relative \
			--go-grpc_out=. \
			--go-grpc_opt=paths=source_relative \
			--validate_out="lang=go,paths=source_relative:." \
			$(PROTO_DIR)/$(MODULE)/*.proto || { echo "Error: protoc failed for Go code generation"; exit 1; }; \
		mkdir -p $(LOCAL_GO_GEN_DIR)/$(MODULE); \
		mv $(PROTO_DIR)/$(MODULE)/*.pb.go $(LOCAL_GO_GEN_DIR)/$(MODULE)/ 2>/dev/null || true; \
		for f in `find $(PROTO_DIR)/$(MODULE) -name "*.validate.pb.go" -o -name "*.validate.go" -o -name "*_validate.pb.go" 2>/dev/null`; do \
			mv $$f $(LOCAL_GO_GEN_DIR)/$(MODULE)/ 2>/dev/null || true; \
		done; \
		echo "Successfully compiled Go code for module: $(MODULE)"; \
		echo "Checking for TypeScript plugin"; \
		if command -v protoc-gen-ts > /dev/null 2>&1; then \
			echo "Generating TypeScript code for $(MODULE)"; \
			protoc -I=. \
				--plugin=protoc-gen-ts=`which protoc-gen-ts` \
				--ts_out="use_proto_field_name:$(LOCAL_TS_GEN_DIR)" \
				$(PROTO_DIR)/$(MODULE)/*.proto || { echo "Error: protoc failed for TypeScript code generation"; exit 1; }; \
			if [ -d "$(LOCAL_TS_GEN_DIR)/definitions/$(MODULE)" ]; then \
				mkdir -p $(LOCAL_TS_GEN_DIR)/$(MODULE); \
				for f in `find $(LOCAL_TS_GEN_DIR)/definitions/$(MODULE) -name "*.ts" 2>/dev/null`; do \
					filename=`basename $$f`; \
					cp $$f $(LOCAL_TS_GEN_DIR)/$(MODULE)/$$filename; \
				done; \
				rm -rf $(LOCAL_TS_GEN_DIR)/definitions; \
			fi; \
			echo "Successfully compiled TypeScript code for module: $(MODULE)"; \
		else \
			echo "Warning: protoc-gen-ts not found, skipping TypeScript generation for $(MODULE)"; \
			echo "To install: npm install --global @protobuf-ts/plugin"; \
		fi; \
		echo "Generating Python code for $(MODULE)"; \
		protoc -I=. \
			--python_out=$(LOCAL_PY_GEN_DIR) \
			$(PROTO_DIR)/$(MODULE)/*.proto || { echo "Error: protoc failed for Python code generation"; exit 1; }; \
		if python3 -c "import grpc_tools.protoc" 2>/dev/null; then \
			echo "Generating Python gRPC code for $(MODULE)"; \
			python3 -m grpc_tools.protoc -I=. \
				--grpc_python_out=$(LOCAL_PY_GEN_DIR) \
				$(PROTO_DIR)/$(MODULE)/*.proto || { echo "Error: grpc_tools.protoc failed for Python gRPC code generation"; exit 1; }; \
		else \
			echo "Warning: grpcio-tools not found, skipping Python gRPC generation for $(MODULE)"; \
			echo "To install: pip install grpcio-tools"; \
		fi; \
		if [ -d "$(LOCAL_PY_GEN_DIR)/definitions/$(MODULE)" ]; then \
			mkdir -p $(LOCAL_PY_GEN_DIR)/$(MODULE); \
			for f in `find $(LOCAL_PY_GEN_DIR)/definitions/$(MODULE) -name "*.py" 2>/dev/null`; do \
				filename=`basename $$f`; \
				cp $$f $(LOCAL_PY_GEN_DIR)/$(MODULE)/$$filename; \
			done; \
			rm -rf $(LOCAL_PY_GEN_DIR)/definitions; \
		fi; \
		touch $(LOCAL_PY_GEN_DIR)/$*/__init__.py; \
		python3 fix_py_imports.py $(LOCAL_PY_GEN_DIR); \
		echo "Successfully compiled Python code for module: $*"; \
	else \
		echo "No proto files found for module: $(MODULE)"; \
	fi

# 显示所有可用的模块
list-modules:
	@echo "Available modules: $(NESTED_MODULES)"

# 显示帮助信息
help:
	@echo "Available targets:"
	@echo "  all          - Compile all proto files (Go, TypeScript, Python)"
	@echo "  local        - Compile all proto files for local development (output to local-gen/)"
	@echo "  clean        - Remove generated code"
	@echo "  list-modules - List all available modules"
	@echo "  <module>     - Compile proto files for a specific module"
	@echo "                 Available modules: $(NESTED_MODULES)"
	@echo ""
	@echo "Output directories:"
	@echo "  Go:         $(GO_GEN_DIR) / $(LOCAL_GO_GEN_DIR)"
	@echo "  TypeScript: $(TS_GEN_DIR) / $(LOCAL_TS_GEN_DIR)"
	@echo "  Python:     $(PY_GEN_DIR) / $(LOCAL_PY_GEN_DIR)"
	@echo ""
	@echo "Python dependencies (optional for gRPC):"
	@echo "  pip install grpcio-tools"