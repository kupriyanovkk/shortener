package main

import (
	"context"
	"log"
	"os"

	"github.com/kupriyanovkk/shortener/internal/config"
	pb "github.com/kupriyanovkk/shortener/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TestShortener tests the Shortener client.
func TestShortener(c pb.ShortenerClient) {
	ctx := context.Background()

	short, err := c.GetShortURL(ctx, &pb.GetShortURLRequest{
		Url: "https://google.com",
	})
	if err != nil {
		log.Fatal(err)
	}

	original, err := c.GetOriginalURLByShort(ctx, &pb.GetOriginalURLByShortRequest{
		Short: short.Result,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Short: %s\nOriginal: %s\n", short.Result, original.FullUrl)
}

func main() {
	flags, _ := config.ParseFlags(os.Args[0], os.Args[1:])
	conn, err := grpc.Dial(flags.GRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	c := pb.NewShortenerClient(conn)

	TestShortener(c)
}
