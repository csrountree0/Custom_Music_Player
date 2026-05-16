package services

import (
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"musicapp/backend/internal/models"
	"musicapp/backend/internal/repository"
)

type SearchResults struct {
	Artists []models.Artist `json:"artists"`
	Albums  []models.Album  `json:"albums"`
	Tracks  []models.Track  `json:"tracks"`
}

type SearchService struct {
	artistRepo repository.ArtistRepository
	albumRepo  repository.AlbumRepository
	trackRepo  repository.TrackRepository
}

func NewSearchService(
	artistRepo repository.ArtistRepository,
	albumRepo repository.AlbumRepository,
	trackRepo repository.TrackRepository,
) *SearchService {
	return &SearchService{artistRepo: artistRepo, albumRepo: albumRepo, trackRepo: trackRepo}
}

func (s *SearchService) Search(ownerID uuid.UUID, query string) (*SearchResults, error) {
	var results SearchResults

	g := new(errgroup.Group)

	g.Go(func() error {
		artists, err := s.artistRepo.Search(ownerID, query)
		if err != nil {
			return err
		}
		results.Artists = artists
		return nil
	})

	g.Go(func() error {
		albums, err := s.albumRepo.Search(ownerID, query)
		if err != nil {
			return err
		}
		results.Albums = albums
		return nil
	})

	g.Go(func() error {
		tracks, err := s.trackRepo.Search(ownerID, query)
		if err != nil {
			return err
		}
		results.Tracks = tracks
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if results.Artists == nil {
		results.Artists = []models.Artist{}
	}
	if results.Albums == nil {
		results.Albums = []models.Album{}
	}
	if results.Tracks == nil {
		results.Tracks = []models.Track{}
	}

	return &results, nil
}
