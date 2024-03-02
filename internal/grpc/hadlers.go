package grpc

import (
	"context"
	"errors"
	"net/url"

	"github.com/kupriyanovkk/shortener/internal/failure"
	"github.com/kupriyanovkk/shortener/internal/generator"
	pb "github.com/kupriyanovkk/shortener/internal/grpc/proto"
	storeInterface "github.com/kupriyanovkk/shortener/internal/store/interface"
	"github.com/kupriyanovkk/shortener/internal/userid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetShortURL retrieves a short URL for the given original URL and user ID.
//
// ctx context.Context, request *pb.GenerateShortURLRequest
// *pb.GenerateShortURLResponse, error
func (s *ShortenerServer) GetShortURL(ctx context.Context, request *pb.GetShortURLRequest) (*pb.GetShortURLResponse, error) {
	var response pb.GetShortURLResponse

	baseURL := s.app.Flags.BaseURL
	userID := userid.Get(ctx)
	parsedURL, err := url.ParseRequestURI(request.Url)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error parsing URL")
	}
	id, _ := generator.GetRandomStr(10)
	short, saveErr := s.app.Store.AddValue(ctx, storeInterface.AddValueOptions{
		Original: parsedURL.String(),
		BaseURL:  baseURL,
		Short:    id,
		UserID:   userID,
	})

	if errors.Is(saveErr, failure.ErrConflict) {
		return nil, status.Error(codes.AlreadyExists, failure.ErrConflict.Error())
	} else {
		response.Result = short
		return &response, nil
	}
}

// DeleteAPIUserURLs deletes the user's URLs in the ShortenerServer.
//
// ctx context.Context, request *pb.DeleteAPIUserURLsRequest
// *pb.DeleteAPIUserURLsResponse, error
func (s *ShortenerServer) GetOriginalURLByShort(ctx context.Context, request *pb.GetOriginalURLByShortRequest) (*pb.GetOriginalURLByShortResponse, error) {
	var response pb.GetOriginalURLByShortResponse

	origURL, err := s.app.Store.GetOriginalURL(ctx, request.Short)
	if err != nil {
		if err.Error() == "URL is deleted" {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	response.FullUrl = origURL
	return &response, nil
}

// GetAPIUserURLs retrieves the URLs for a given user.
//
// ctx context.Context, request *pb.GetAPIUserURLsRequest -> *pb.GetAPIUserURLsResponse, error
func (s *ShortenerServer) GetAPIUserURLs(ctx context.Context, request *pb.GetAPIUserURLsRequest) (*pb.GetAPIUserURLsResponse, error) {
	var response pb.GetAPIUserURLsResponse

	userID := userid.Get(ctx)
	URLs, err := s.app.Store.GetUserURLs(ctx, storeInterface.GetUserURLsOptions{
		UserID:  userID,
		BaseURL: s.app.Flags.BaseURL,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, v := range URLs {
		response.Urls = append(response.Urls, &pb.URL{Original: v.Original, Short: v.Short})
	}

	return &response, nil
}

// GetInternalStats retrieves internal statistics.
//
// Context, request.
// GetInternalStatsResponse, error.
func (s *ShortenerServer) GetInternalStats(ctx context.Context, in *emptypb.Empty) (*pb.GetInternalStatsResponse, error) {
	if s.app.Flags.TrustedSubnet == "" {
		return nil, status.Error(codes.PermissionDenied, "Trusted subnet is not set")
	}

	var response pb.GetInternalStatsResponse

	stats, err := s.app.Store.GetInternalStats(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response.Urls = int32(stats.URLs)
	response.Users = int32(stats.Users)

	return &response, nil
}

// DeleteAPIUserURLs deletes the user's URLs in the ShortenerServer.
//
// ctx context.Context, request *pb.DeleteAPIUserURLsRequest
// *pb.DeleteAPIUserURLsResponse, error
func (s *ShortenerServer) DeleteAPIUserURLs(ctx context.Context, request *pb.DeleteAPIUserURLsRequest) (*pb.DeleteAPIUserURLsResponse, error) {
	var response pb.DeleteAPIUserURLsResponse

	userID := userid.Get(ctx)

	if len(request.Urls) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	s.app.URLChan <- storeInterface.DeletedURLs{
		UserID: userID,
		URLs:   request.Urls,
	}

	return &response, nil
}
