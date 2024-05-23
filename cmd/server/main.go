package main

import (
	"log"

	server_logic "toloka-metrics/internal/server_logic"
)

func main() {
	err := server_logic.StartServer()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
