package handlers

import (
	"fmt"
	"github.com/joho/godotenv"
)

func LoadEnv(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		fmt.Println("Load Env: ", err)
		//panic("Failed load env file")
		return err
	}
	return nil
}
