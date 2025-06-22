package main

import (
	"fmt"

	"employee-app/configs"
)

func init() {
	configs.ConnectDB(&configs.AppConfig)
}

func main() {
	configs.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	//configs.DB.AutoMigrate(&models.User{}, &models.Post{})
	fmt.Println("👍 Migration complete")
}