package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/configuration"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/security"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/service"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/client/tui/model"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

func main() {
	config := configuration.Configure()

	conn, err := grpc.NewClient(config.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения к gRPC серверу: %v", err)
	}
	defer conn.Close()

	// Create GRPC client
	userServiceClient := pb.NewUserServiceClient(conn)
	secretServiceClient := pb.NewSecretServiceClient(conn)

	// Create required components
	encryptionService := security.NewEncryptionService()

	// Create services
	userService := service.NewUserService(userServiceClient)
	secretService := service.NewSecretService(secretServiceClient, encryptionService)

	// Create initial model
	initModel := model.NewInitModel()

	// Create tui context with all dependencies
	appContext := tui.NewTUIContext(initModel, userService, secretService)

	p := tea.NewProgram(appContext)

	_, err = p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
