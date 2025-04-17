package main

import (
	"context"
	"github.com/randhid/sanderee"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

func main() {
	err := realMain()
	if err != nil {
		panic(err)
	}
}

func realMain() error {
	ctx := context.Background()
	logger := logging.NewLogger("cli")

	deps := resource.Dependencies{}
	// can load these from a remote machine if you need

	thing, err := sanderee.NewSander(ctx, deps, resource.Config{Name: "foo"}, logger)
	if err != nil {
		return err
	}
	defer thing.Close(ctx)

	return nil
}
