package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/app/state"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/logger"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
	"google.golang.org/grpc"
	"log"
	"os"
)

func main() {
	serverAddr := "localhost:50051" // Адрес gRPC сервера
	conn, err := grpc.NewClient(serverAddr)
	if err != nil {
		log.Fatalf("Ошибка подключения к gRPC серверу: %v", err)
	}
	defer conn.Close()

	// Создаем gRPC-клиент
	userServiceClient := pb.NewUserServiceClient(conn)

	// Также нужно иметь какую-то структуру с глобальным состоянием, где будет храниться ключ для шифрования на основе пароля + токен
	// и клиент туда-же можно пихнуть

	//p := tea.NewProgram(app.App{State: state.InitState{Choices: state.Choices}})
	p := tea.NewProgram(state.InitState{Choices: state.Choices})
	_, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v", err)
		os.Exit(1)
	}
}
