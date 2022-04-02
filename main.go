/*

 */

package main

import (
	"log"

	"opensea-data-backend/server"
)

func main() {
	log.Fatal(server.Start())
}
