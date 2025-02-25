package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hemagome/hefesto/config"
	"github.com/hemagome/hefesto/logging"
	"github.com/hemagome/hefesto/storage"
	"github.com/hemagome/hefesto/tui"

	"go.uber.org/zap"
)

const tokenFile = "tokens.json"

func main() {
	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Error loading configuration: %v", err))
	}

	// Initialize logger
	logging.InitLogger()
	logging.Logger.Info("🚀 Iniciando Hefesto...",
		zap.String("environment", cfg.Environment))

	// Permitir cambiar el nivel de logs desde una variable de entorno
	logLevel := os.Getenv("HEFESTO_LOG_LEVEL")
	if logLevel != "" {
		logging.SetLogLevel(logLevel)
	}

	err = ensureTokenFile(tokenFile)
	if err != nil {
		logging.Logger.Fatal("❌ Error asegurando el archivo de tokens", zap.Error(err))
	}

	// Initialize and run TUI
	p := tea.NewProgram(tui.NewTokenInput())
	if _, err := p.Run(); err != nil {
		logging.Logger.Fatal("❌ Error running TUI", zap.Error(err))
	}

	// Verify token was stored
	retrievedToken, err := storage.GetToken("github")
	if err != nil {
		logging.Logger.Error("❌ Error obteniendo token", zap.Error(err))
	} else {
		logging.Logger.Info("✅ Token guardado exitosamente",
			zap.String("proveedor", "github"),
			zap.String("token_length", fmt.Sprintf("%d chars", len(retrievedToken))))
	}
}

func ensureTokenFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		logging.Logger.Warn("📂 Archivo de tokens no encontrado, creando uno nuevo...")
		return os.WriteFile(filename, []byte("{}"), 0600)
	}
	return err
}
