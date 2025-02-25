package main

import (
	"context"
)

func main() {
	ctx := context.Background()

	app := NewApp(ctx)

	app.Run()
}
