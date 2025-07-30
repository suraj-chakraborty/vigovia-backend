package main

import "itinerary/router"

func main() {
	r := router.SetupRouter()
	r.Run(":8080")
}
