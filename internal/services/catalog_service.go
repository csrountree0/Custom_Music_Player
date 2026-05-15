package services

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musicapp/backend/internal/models"
	"musicapp/backend/internal/repository"
	"musicapp/backend/pkg/apierror"
	"musicapp/backend/pkg/pagination"
)

type CreateArtistInput struct {
	Name     string `json:"name"      binding:"required"`
	Bio      string `json:"bio"`
	ImageURL string `json:"image_url"`
}

type UpdateArtistInput struct {
	Name     string `json:"name"`
	Bio      string `json:"bio"`
	ImageURL string `json:"image_url"`
}

type CreateAlbumInput struct {
	Title        string     `json:"title"      binding:"required"`
	ArtistID     *uuid.UUID `json:"artist_id"`
	Year         int        `json:"year"`
	CoverArtPath string     `json:"cover_art_path"`
}

type UpdateAlbumInput struct {
	Title        string     `json:"title"`
	ArtistID     *uuid.UUID `json:"artist_id"`
	Year         int        `json:"year"`
	CoverArtPath string     `json:"cover_art_path"`
}

type UpdateTrackInput struct {
	Title       string     `json:"title"`
	Visibility  string     `json:"visibility" binding:"omitempty,oneof=private public"`
	ArtistID    *uuid.UUID `json:"artist_id"`
	AlbumID     *uuid.UUID `json:"album_id"`
	TrackNumber int        `json:"track_number"`
	DiscNumber  int        `json:"disc_number"`
	Year        int        `json:"year"`
}

type CatalogService struct {
	artistRepo repository.ArtistRepository
	albumRepo  repository.AlbumRepository
	trackRepo  repository.TrackRepository
}

func NewCatalogService(
	artistRepo repository.ArtistRepository,
	albumRepo repository.AlbumRepository,
	trackRepo repository.TrackRepository,
) *CatalogService {
	return &CatalogService{artistRepo: artistRepo, albumRepo: albumRepo, trackRepo: trackRepo}
}

func (s *CatalogService) ListArtists(ownerID uuid.UUID, p pagination.Params) ([]models.Artist, int64, error) {
	return s.artistRepo.List(ownerID, p)
}

func (s *CatalogService) CreateArtist(ownerID uuid.UUID, input CreateArtistInput) (*models.Artist, error) {
	artist := &models.Artist{
		OwnerID:  ownerID,
		Name:     input.Name,
		Bio:      input.Bio,
		ImageURL: input.ImageURL,
	}
	if err := s.artistRepo.Create(artist); err != nil {
		return nil, fmt.Errorf("creating artist: %w", err)
	}
	return artist, nil
}

func (s *CatalogService) GetArtist(id, ownerID uuid.UUID) (*models.Artist, error) {
	artist, err := s.artistRepo.FindByID(id, ownerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.New(404, "ARTIST_NOT_FOUND", "artist not found")
		}
		return nil, err
	}
	return artist, nil
}

func (s *CatalogService) UpdateArtist(id, ownerID uuid.UUID, input UpdateArtistInput) (*models.Artist, error) {
	artist, err := s.GetArtist(id, ownerID)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		artist.Name = input.Name
	}
	if input.Bio != "" {
		artist.Bio = input.Bio
	}
	if input.ImageURL != "" {
		artist.ImageURL = input.ImageURL
	}
	if err := s.artistRepo.Update(artist); err != nil {
		return nil, fmt.Errorf("updating artist: %w", err)
	}
	return artist, nil
}

func (s *CatalogService) DeleteArtist(id, ownerID uuid.UUID) error {
	if _, err := s.GetArtist(id, ownerID); err != nil {
		return err
	}
	return s.artistRepo.Delete(id, ownerID)
}

func (s *CatalogService) ListAlbums(ownerID uuid.UUID, p pagination.Params) ([]models.Album, int64, error) {
	return s.albumRepo.List(ownerID, p)
}

func (s *CatalogService) CreateAlbum(ownerID uuid.UUID, input CreateAlbumInput) (*models.Album, error) {
	album := &models.Album{
		OwnerID:      ownerID,
		ArtistID:     input.ArtistID,
		Title:        input.Title,
		Year:         input.Year,
		CoverArtPath: input.CoverArtPath,
	}
	if err := s.albumRepo.Create(album); err != nil {
		return nil, fmt.Errorf("creating album: %w", err)
	}
	return album, nil
}

func (s *CatalogService) GetAlbum(id, ownerID uuid.UUID) (*models.Album, error) {
	album, err := s.albumRepo.FindByID(id, ownerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.New(404, "ALBUM_NOT_FOUND", "album not found")
		}
		return nil, err
	}
	return album, nil
}

func (s *CatalogService) UpdateAlbum(id, ownerID uuid.UUID, input UpdateAlbumInput) (*models.Album, error) {
	album, err := s.GetAlbum(id, ownerID)
	if err != nil {
		return nil, err
	}
	if input.Title != "" {
		album.Title = input.Title
	}
	if input.ArtistID != nil {
		album.ArtistID = input.ArtistID
	}
	if input.Year != 0 {
		album.Year = input.Year
	}
	if input.CoverArtPath != "" {
		album.CoverArtPath = input.CoverArtPath
	}
	if err := s.albumRepo.Update(album); err != nil {
		return nil, fmt.Errorf("updating album: %w", err)
	}
	return album, nil
}

func (s *CatalogService) DeleteAlbum(id, ownerID uuid.UUID) error {
	if _, err := s.GetAlbum(id, ownerID); err != nil {
		return err
	}
	return s.albumRepo.Delete(id, ownerID)
}

func (s *CatalogService) ListTracks(ownerID uuid.UUID, p pagination.Params) ([]models.Track, int64, error) {
	return s.trackRepo.List(ownerID, p)
}

func (s *CatalogService) GetTrack(id, ownerID uuid.UUID) (*models.Track, error) {
	track, err := s.trackRepo.FindByID(id, ownerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.New(404, "TRACK_NOT_FOUND", "track not found")
		}
		return nil, err
	}
	return track, nil
}

func (s *CatalogService) UpdateTrack(id, ownerID uuid.UUID, input UpdateTrackInput) (*models.Track, error) {
	track, err := s.GetTrack(id, ownerID)
	if err != nil {
		return nil, err
	}
	if input.Title != "" {
		track.Title = input.Title
	}
	if input.Visibility != "" {
		track.Visibility = input.Visibility
	}
	if input.ArtistID != nil {
		track.ArtistID = input.ArtistID
	}
	if input.AlbumID != nil {
		track.AlbumID = input.AlbumID
	}
	if input.TrackNumber != 0 {
		track.TrackNumber = input.TrackNumber
	}
	if input.DiscNumber != 0 {
		track.DiscNumber = input.DiscNumber
	}
	if input.Year != 0 {
		track.Year = input.Year
	}
	if err := s.trackRepo.Update(track); err != nil {
		return nil, fmt.Errorf("updating track: %w", err)
	}
	return track, nil
}

func (s *CatalogService) DeleteTrack(id, ownerID uuid.UUID) error {
	track, err := s.GetTrack(id, ownerID)
	if err != nil {
		return err
	}
	if err := s.trackRepo.Delete(id, ownerID); err != nil {
		return fmt.Errorf("deleting track: %w", err)
	}
	if track.FilePath != "" {
		_ = os.Remove(track.FilePath)
	}
	return nil
}

func (s *CatalogService) findOrCreateArtist(ownerID uuid.UUID, name string) (*uuid.UUID, error) {
	if name == "" {
		return nil, nil
	}
	artist, err := s.artistRepo.FindByNameAndOwner(name, ownerID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("looking up artist: %w", err)
		}
		artist = &models.Artist{OwnerID: ownerID, Name: name}
		if err := s.artistRepo.Create(artist); err != nil {
			return nil, fmt.Errorf("creating artist: %w", err)
		}
	}
	return &artist.ID, nil
}

func (s *CatalogService) findOrCreateAlbum(ownerID uuid.UUID, title string, artistID *uuid.UUID, year int) (*uuid.UUID, error) {
	if title == "" {
		return nil, nil
	}
	album, err := s.albumRepo.FindByTitleArtistOwner(title, artistID, ownerID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("looking up album: %w", err)
		}
		album = &models.Album{OwnerID: ownerID, ArtistID: artistID, Title: title, Year: year}
		if err := s.albumRepo.Create(album); err != nil {
			return nil, fmt.Errorf("creating album: %w", err)
		}
	}
	return &album.ID, nil
}
