package main

import (
	"context"
	"os"

	cont "shuv1wolf/skillmatch/core/containers"
)

func main() {
	proc := cont.NewCoreProcess()
	proc.SetConfigPath("../config/config.yml")
	proc.Run(context.Background(), os.Args)
}
