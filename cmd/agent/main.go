package main

import (
	"log"
	"os"
	"strconv"

	"github.com/wifi538/CalOnlineParallel/internal/agent"
)

func main() {
	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		log.Fatalf("Invalid COMPUTING_POWER: %v", err)
	}

	for i := 0; i < computingPower; i++ {
		go agent.Worker()
	}

	select {}
}
