package main

import (
	"fmt"
	"os"

	"github.com/arabian9ts/builder/pkg/fileoperator"
)

func genBuilder(targetPkg string) {
	err := fileoperator.CleanBuilder(targetPkg)
	if err != nil {
		panic(err)
	}

	err = fileoperator.CreateBuilder(targetPkg)
	if err != nil {
		panic(err)
	}
}

func main() {
	buildTarget := os.Args[1:]

	if len(buildTarget) <= 0 {
		fmt.Println("ackage is not specified")
		fmt.Println("[USAGE]: builder <Package Name>")
		os.Exit(1)
	}

	for i := range buildTarget {
		genBuilder(buildTarget[i])
	}
}
