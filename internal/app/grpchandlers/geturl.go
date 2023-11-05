package grpchandlers

import (
	context "context"

	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (s *GRPCServer) CreateSGetOriginalURLhortURL(ctx context.Context, in *GetOriginalURLRequest) (*GetOriginalURLResponse, error) {
	shortURL := in.ShortUrl
	if shortURL == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty ShortURL")
	}

	url, err := s.storage.ReadByID(ctx, shortURL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error")
	}

	if url.DeletedFlag {
		return nil, status.Errorf(codes.Internal, "URL was deleted")
	}

	return &GetOriginalURLResponse{
		OriginalUrl: url.OriginalURL,
	}, nil
}

func (s *GRPCServer) GetUserURLs(ctx context.Context, in *GetUserURLsRequest) (*GetUserURLsResponse, error) {
	userID := in.UserId
	if len(userID) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	userURLs, err := s.storage.GetUserURLs(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error")
	}

	if len(userURLs) == 0 {
		return nil, status.Errorf(codes.NotFound, "No content")
	}

	responseModel := make([]*UserUrls, 0, len(userURLs))
	for _, userURL := range userURLs {
		responseModel = append(
			responseModel,
			&UserUrls{
				ShortUrl:    s.formatShortURL(userURL.ShortURL),
				OriginalUrl: userURL.OriginalURL,
			},
		)
	}

	return &GetUserURLsResponse{
		Urls: responseModel,
	}, nil
}
