package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"templ_workout/internals/app"
	"templ_workout/internals/config"
)

func main() {

	app := app.NewApp(config.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed start app:", err)
	}

}
