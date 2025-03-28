package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"

	cont "shuv1wolf/skillmatch/core/containers"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	proc := cont.NewCoreProcess()
	proc.SetConfigPath("../config/config.yml")
	proc.Run(context.Background(), os.Args)
}
