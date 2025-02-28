package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/app"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/app/state"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

func main() {
	serverAddr := "localhost:50051" // Адрес gRPC сервера
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения к gRPC серверу: %v", err)
	}
	defer conn.Close()

	// Создаем gRPC-клиент
	userServiceClient := pb.NewUserServiceClient(conn)

	// Create initial state
	initState := state.NewInitState(state.Choices, 0)
	//initState := state.InitState{Choices: state.Choices}

	// Create app context with all dependencies
	appContext := app.NewApp(initState, userServiceClient)

	p := tea.NewProgram(appContext)

	//p := tea.NewProgram(state.InitState{Choices: state.Choices})

	_, err = p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v", err)
		os.Exit(1)
	}
}
