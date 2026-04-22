//go:build ignore

package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/Cellul4r/go-quiz-app/backend/domain"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(&domain.Profile{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
