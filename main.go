package main

import (
	"fmt"
	"os"

	"github.com/hemagome/hefesto/config"
	"github.com/hemagome/hefesto/logging"
	"github.com/hemagome/hefesto/storage"

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

	provider := "github"
	testToken := "mi-token-de-prueba"

	err = storage.StoreToken(provider, testToken)
	if err != nil {
		logging.Logger.Error("❌ Error guardando token", zap.Error(err))
	}

	retrievedToken, err := storage.GetToken(provider)
	if err != nil {
		logging.Logger.Error("❌ Error obteniendo token", zap.Error(err))
	}

	logging.Logger.Info("🔑 Token recuperado", zap.String("proveedor", provider), zap.String("token", retrievedToken))
}

func ensureTokenFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		logging.Logger.Warn("📂 Archivo de tokens no encontrado, creando uno nuevo...")
		return os.WriteFile(filename, []byte("{}"), 0600)
	}
	return err
}
