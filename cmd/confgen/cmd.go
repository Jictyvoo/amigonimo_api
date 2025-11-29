package main

import "log"

func main() {
	if err := configGenCMD(); err != nil {
		log.Fatalf("Failed to execute command due: `%v`", err)
	}
	log.Println("Configuration generation complete")
}
