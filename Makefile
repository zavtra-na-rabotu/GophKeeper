# Имя бинарных файлов
CLIENT_BIN=gophkeeper-client
SERVER_BIN=gophkeeper-server

# Путь к исходникам
CLIENT_SRC=cmd/client/main.go
SERVER_SRC=cmd/server/main.go

# Каталог сборки
BUILD_DIR=bin

# Флаги компиляции
GO_FLAGS=-ldflags="-s -w"

.PHONY: all client server clean run-client run-server

# Сборка всего проекта
all: client server

# Сборка клиента
client:
	@echo "Building client..."
	@mkdir -p $(BUILD_DIR)
	@go build $(GO_FLAGS) -o $(BUILD_DIR)/$(CLIENT_BIN) $(CLIENT_SRC)

# Сборка сервера
server:
	@echo "Building server..."
	@mkdir -p $(BUILD_DIR)
	@go build $(GO_FLAGS) -o $(BUILD_DIR)/$(SERVER_BIN) $(SERVER_SRC)

# Очистка сборки
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

