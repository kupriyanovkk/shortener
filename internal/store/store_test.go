package store

import (
	"context"
	"testing"

	infile "github.com/kupriyanovkk/shortener/internal/store/in_file"
	inmemory "github.com/kupriyanovkk/shortener/internal/store/in_memory"
	storeInterface "github.com/kupriyanovkk/shortener/internal/store/interface"
)

func BenchmarkAddValue(b *testing.B) {
	memStore := inmemory.NewStore()
	fileStore := infile.NewStore("/tmp/short-url-db.json")

	b.Run("in memory", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			memStore.AddValue(context.Background(), storeInterface.AddValueOptions{
				Original: "https://example.com",
				Short:    "abc",
				BaseURL:  "https://short.ly",
			})
		}
	})

	b.Run("in file", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fileStore.AddValue(context.Background(), storeInterface.AddValueOptions{
				Original: "https://example.com",
				Short:    "abc",
				BaseURL:  "https://short.ly",
			})
		}
	})
}
