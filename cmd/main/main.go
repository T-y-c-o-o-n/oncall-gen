package main

import (
	"fmt"
	"log"
	"oncall-gen/internal/app"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Println(fmt.Sprintf("Usage: %s <config file> <onCall API uri>", os.Args[0]))
		fmt.Println("<config file>    - is .yaml file with team and users")
		fmt.Println("<onCall API uri> - is uri with onCall API (for example, localhost:8080)")
		os.Exit(0)
	}

	filename := os.Args[1]
	onCallUri := os.Args[2]

	a := app.NewAppImpl(fmt.Sprintf("http://%s", onCallUri))
	err := a.CreateTeams(filename)

	if err != nil {
		log.Printf("Error: %v", err)
	}
}
