package main

import (
	"Alarm/internal/config"
	"Alarm/internal/web/controllers"
	"Alarm/internal/web/models"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Config Init
	globalConfig, err := config.NewConfig(".", "config")
	if err != nil {
		log.Fatal(err)
	}

	// Mysql Init
	db, err := models.NewDatabase(globalConfig.Gin.Mysql)
	if err != nil {
		log.Fatal(err)
	}
	cache, err := models.NewCache(globalConfig.Gin.Redis)
	if err != nil {
		log.Fatal(err)
	}

	// Auth Init
	privateKey, err := readPrivateKeyFromFile("./privateKey.pem")
	if err != nil {
		log.Fatal(err)
	}

	publicKey, err := readPublicKeyFromFile("./publicKey.pem")
	if err != nil {
		log.Fatal(err)
	}

	// Gin Init
	engine := gin.Default()

	// --Controller Init
	ctrlConfig := map[string]interface{}{
		"db":    db,
		"cache": cache,
	}
	authConfig := map[string]interface{}{
		"db":                db,
		"cache":             cache,
		"privateKey":        privateKey,
		"publicKey":         publicKey,
		"tokenValidSeconds": globalConfig.Gin.Token.ValidSeconds,
	}
	AccountCtrl := controllers.NewAccount(ctrlConfig)
	AuthCtrl := controllers.NewAuthController(authConfig)

	// --Router Init
	group := engine.Group("")
	{
		group.POST("/register", AccountCtrl.CreateUser)
		group.GET("/users", AccountCtrl.AllUser)
		group.GET("/users/:id", AccountCtrl.GetUserByID)
		group.PATCH("/users/:id", AccountCtrl.UpdateUserByID)

		group.GET("/authtest", AuthCtrl.LoginMiddleware, AuthCtrl.Test)
		group.POST("/login", AuthCtrl.Login)
		group.POST("/logout", AuthCtrl.Logout)
	}
	engine.Run(globalConfig.Gin.Port)
}

func readPrivateKeyFromFile(filepath string) (interface{}, error) {
	keyFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyFile)
	if block == nil {
		return nil, fmt.Errorf("decode private key error")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func readPublicKeyFromFile(filepath string) (interface{}, error) {
	keyFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyFile)
	if block == nil {
		return nil, fmt.Errorf("decode public key error")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}
