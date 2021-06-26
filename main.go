package main

import (
	"fmt"

	"github.com/pkg/errors"
)

func main() {
	errors.Wrap()
	fmt.Println("Hello, World!")
}
