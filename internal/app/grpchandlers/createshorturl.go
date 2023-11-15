package grpchandlers

import (
	context "context"
	"errors"

	"github.com/VladKvetkin/shortener/internal/app/entities"
	"github.com/VladKvetkin/shortener/internal/app/shortener"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (s *GRPCServer) CreateShortURL(ctx context.Context, in *CreateShortURLRequest) (*CreateShortURLResponse, error) {
	url := in.Url

	if len(url) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Empty URL")
	}

	userID := in.UserId
	if len(userID) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	id, err := s.createAndAddID(ctx, url, userID)
	if err != nil {
		if errors.Is(err, ErrOriginalURLAlreadyExists) {
			return &CreateShortURLResponse{
				ShortUrl: s.formatShortURL(id),
			}, nil
		}

		return nil, status.Errorf(codes.Internal, "Internal error")
	}

	return &CreateShortURLResponse{
		ShortUrl: s.formatShortURL(id),
	}, nil
}

func (s *GRPCServer) CreateShortURLBatch(ctx context.Context, in *CreateShortURLBatchRequest) (*CreateShortURLBatchResponse, error) {
	userID := in.UserId
	if len(userID) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	urls := make([]entities.URL, 0, len(in.Urls))
	responseModel := make([]*ShortURLs, 0, len(in.Urls))

	for _, batchData := range in.Urls {
		shortURL, err := shortener.CreateID(batchData.OriginalUrl)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Internal error")
		}

		urls = append(
			urls,
			entities.URL{
				OriginalURL: batchData.OriginalUrl,
				ShortURL:    shortURL,
				UserID:      userID,
			},
		)

		responseModel = append(
			responseModel,
			&ShortURLs{
				CorrelationId: batchData.CorrelationId,
				ShortUrl:      s.formatShortURL(shortURL),
			},
		)
	}

	err := s.storage.AddBatch(ctx, urls)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error")
	}

	return &CreateShortURLBatchResponse{
		Urls: responseModel,
	}, nil
}
