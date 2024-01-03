package store

import (
	"context"
	"fmt"

	inmemory "github.com/kupriyanovkk/shortener/internal/store/in_memory"
	storeInterface "github.com/kupriyanovkk/shortener/internal/store/interface"
)

func Example() {
	store := inmemory.NewStore()

	shortURL, _ := store.AddValue(context.Background(), storeInterface.AddValueOptions{
		Original: "https://example.com",
		Short:    "abc",
		BaseURL:  "https://short.ly",
	})

	originalURL, _ := store.GetOriginalURL(context.Background(), shortURL)

	fmt.Println(originalURL)
}
