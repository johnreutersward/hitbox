package main

import (
	"fmt"
	"github.com/rojters/hitbox"
	"log"
)

func main() {
	hbc := hitbox.NewClient(nil)

	games, _, err := hbc.Games()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("The top 100 games on hitbox.tv sorted by number of viewers")

	for i, g := range games.Categories {
		fmt.Printf("[%d] %s | current viewers = %d\n", i+1, *g.CategoryName, *g.CategoryViewers)
	}

}
