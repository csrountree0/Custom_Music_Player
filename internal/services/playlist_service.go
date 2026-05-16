package services

import (
	"errors"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musicapp/backend/internal/models"
	"musicapp/backend/internal/repository"
	"musicapp/backend/pkg/apierror"
	"musicapp/backend/pkg/pagination"
)

type CreatePlaylistInput struct {
	Name            string `json:"name"             binding:"required"`
	Description     string `json:"description"`
	IsPublic        bool   `json:"is_public"`
	IsCollaborative bool   `json:"is_collaborative"`
	CoverURL        string `json:"cover_url"`
}

type UpdatePlaylistInput struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	IsPublic        *bool  `json:"is_public"`
	IsCollaborative *bool  `json:"is_collaborative"`
	CoverURL        string `json:"cover_url"`
}

type AddTrackInput struct {
	TrackID uuid.UUID `json:"track_id" binding:"required"`
}

type ReorderInput struct {
	Positions []repository.TrackPosition `json:"positions" binding:"required"`
}

type PlaylistService struct {
	playlistRepo repository.PlaylistRepository
	trackRepo    repository.TrackRepository
}

func NewPlaylistService(
	playlistRepo repository.PlaylistRepository,
	trackRepo repository.TrackRepository,
) *PlaylistService {
	return &PlaylistService{playlistRepo: playlistRepo, trackRepo: trackRepo}
}

func (s *PlaylistService) ListPlaylists(ownerID uuid.UUID, p pagination.Params) ([]models.Playlist, int64, error) {
	return s.playlistRepo.List(ownerID, p)
}

func (s *PlaylistService) CreatePlaylist(ownerID uuid.UUID, input CreatePlaylistInput) (*models.Playlist, error) {
	playlist := &models.Playlist{
		OwnerID:         ownerID,
		Name:            input.Name,
		Description:     input.Description,
		IsPublic:        input.IsPublic,
		IsCollaborative: input.IsCollaborative,
		CoverURL:        input.CoverURL,
	}
	if err := s.playlistRepo.Create(playlist); err != nil {
		return nil, fmt.Errorf("creating playlist: %w", err)
	}
	return playlist, nil
}

func (s *PlaylistService) GetPlaylist(id, ownerID uuid.UUID) (*models.Playlist, error) {
	playlist, err := s.playlistRepo.FindByID(id, ownerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.New(404, "PLAYLIST_NOT_FOUND", "playlist not found")
		}
		return nil, err
	}
	return playlist, nil
}

func (s *PlaylistService) UpdatePlaylist(id, ownerID uuid.UUID, input UpdatePlaylistInput) (*models.Playlist, error) {
	playlist, err := s.GetPlaylist(id, ownerID)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		playlist.Name = input.Name
	}
	if input.Description != "" {
		playlist.Description = input.Description
	}
	if input.IsPublic != nil {
		playlist.IsPublic = *input.IsPublic
	}
	if input.IsCollaborative != nil {
		playlist.IsCollaborative = *input.IsCollaborative
	}
	if input.CoverURL != "" {
		playlist.CoverURL = input.CoverURL
	}
	if err := s.playlistRepo.Update(playlist); err != nil {
		return nil, fmt.Errorf("updating playlist: %w", err)
	}
	return playlist, nil
}

func (s *PlaylistService) DeletePlaylist(id, ownerID uuid.UUID) error {
	if _, err := s.GetPlaylist(id, ownerID); err != nil {
		return err
	}
	return s.playlistRepo.Delete(id, ownerID)
}

func (s *PlaylistService) ListPlaylistTracks(playlistID, ownerID uuid.UUID) ([]models.Track, error) {
	if _, err := s.GetPlaylist(playlistID, ownerID); err != nil {
		return nil, err
	}

	rows, err := s.playlistRepo.ListTracks(playlistID)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []models.Track{}, nil
	}

	trackIDs := make([]uuid.UUID, len(rows))
	positionByTrack := make(map[uuid.UUID]int, len(rows))
	for i, r := range rows {
		trackIDs[i] = r.TrackID
		positionByTrack[r.TrackID] = r.Position
	}

	tracks, err := s.trackRepo.FindByIDs(trackIDs, ownerID)
	if err != nil {
		return nil, err
	}

	sort.Slice(tracks, func(i, j int) bool {
		return positionByTrack[tracks[i].ID] < positionByTrack[tracks[j].ID]
	})

	return tracks, nil
}

func (s *PlaylistService) AddTrack(playlistID, ownerID uuid.UUID, input AddTrackInput) error {
	if _, err := s.GetPlaylist(playlistID, ownerID); err != nil {
		return err
	}

	if _, err := s.trackRepo.FindByID(input.TrackID, ownerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierror.New(404, "TRACK_NOT_FOUND", "track not found")
		}
		return err
	}

	maxPos, err := s.playlistRepo.MaxPosition(playlistID)
	if err != nil {
		return err
	}

	return s.playlistRepo.AddTrack(playlistID, input.TrackID, ownerID, maxPos+1)
}

func (s *PlaylistService) RemoveTrack(playlistID, trackID, ownerID uuid.UUID) error {
	if _, err := s.GetPlaylist(playlistID, ownerID); err != nil {
		return err
	}
	return s.playlistRepo.RemoveTrack(playlistID, trackID)
}

func (s *PlaylistService) ReorderTracks(playlistID, ownerID uuid.UUID, input ReorderInput) error {
	if _, err := s.GetPlaylist(playlistID, ownerID); err != nil {
		return err
	}
	return s.playlistRepo.ReorderTracks(playlistID, input.Positions)
}
