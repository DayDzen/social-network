package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func CheckEnvs() error {
	log.Println("Checking envs file...")
	var err error
	if err = godotenv.Load("prod_config.yaml"); err != nil {
		if err = godotenv.Load("local_config.yaml"); err != nil {
			return fmt.Errorf("checkEnvs err: %w", err)
		} else {
			log.Println("Local envs OK")
		}
	} else {
		log.Println("Prod envs OK")
	}

	return nil
}
