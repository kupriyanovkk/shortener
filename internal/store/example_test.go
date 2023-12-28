package store

import (
	"context"
	"fmt"

	"github.com/kupriyanovkk/shortener/internal/models"
	inmemory "github.com/kupriyanovkk/shortener/internal/store/in_memory"
)

func Example() {
	store := inmemory.NewStore()

	shortURL, _ := store.AddValue(context.Background(), models.AddValueOptions{
		Original: "https://example.com",
		Short:    "abc",
		BaseURL:  "https://short.ly",
	})

	originalURL, _ := store.GetOriginalURL(context.Background(), shortURL)

	fmt.Println(originalURL)
}
