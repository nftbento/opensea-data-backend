/*

 */

package main

import (
	"log"

	"github.com/NFTActions/opensea-data-backend/server"
)

func main() {
	log.Fatal(server.Start())
}
