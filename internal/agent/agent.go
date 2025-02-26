package agent

import (
	"log"
	"os"
	"strconv"
)

func InitAgent() {
	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		log.Fatalf("Invalid COMPUTING_POWER: %v", err)
	}

	for i := 0; i < computingPower; i++ {
		go Worker()
	}

	select {}
}
