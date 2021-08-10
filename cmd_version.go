package main

import (
	"context"
	"fmt"
)

const version = "v0.0.0"

func runVersion(ctx context.Context) error {
	fmt.Printf("dbumper %s\n", version)
	return nil
}
