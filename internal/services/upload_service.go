package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dhowden/tag"
	"github.com/google/uuid"

	"musicapp/backend/internal/models"
	"musicapp/backend/internal/repository"
)

type UploadService struct {
	catalogSvc *CatalogService
	trackRepo  repository.TrackRepository
	jobRepo    repository.UploadJobRepository
	audioDir   string
}

func NewUploadService(
	catalogSvc *CatalogService,
	trackRepo repository.TrackRepository,
	jobRepo repository.UploadJobRepository,
	audioDir string,
) *UploadService {
	return &UploadService{
		catalogSvc: catalogSvc,
		trackRepo:  trackRepo,
		jobRepo:    jobRepo,
		audioDir:   audioDir,
	}
}

func (s *UploadService) UploadFile(ownerID uuid.UUID, file multipart.File, header *multipart.FileHeader) (*models.UploadJob, error) {
	dir := filepath.Join(s.audioDir, ownerID.String())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating audio directory: %w", err)
	}

	ext := filepath.Ext(header.Filename)
	destName := uuid.New().String() + ext
	destPath := filepath.Join(dir, destName)

	dest, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("creating destination file: %w", err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		_ = os.Remove(destPath)
		return nil, fmt.Errorf("writing file: %w", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		_ = os.Remove(destPath)
		return nil, fmt.Errorf("seeking file: %w", err)
	}

	track, err := s.buildTrackFromFile(ownerID, file, header, destPath)
	if err != nil {
		_ = os.Remove(destPath)
		return nil, err
	}

	if err := s.trackRepo.Create(track); err != nil {
		_ = os.Remove(destPath)
		return nil, fmt.Errorf("saving track: %w", err)
	}

	job := &models.UploadJob{
		OwnerID:  ownerID,
		Status:   "complete",
		Source:   "upload",
		TrackID:  &track.ID,
	}
	if err := s.jobRepo.Create(job); err != nil {
		return nil, fmt.Errorf("creating upload job: %w", err)
	}

	return job, nil
}

func (s *UploadService) ImportYouTube(ownerID uuid.UUID, url string) (*models.UploadJob, error) {
	job := &models.UploadJob{
		OwnerID:   ownerID,
		Status:    "pending",
		Source:    "youtube",
		SourceURL: url,
	}
	if err := s.jobRepo.Create(job); err != nil {
		return nil, fmt.Errorf("creating upload job: %w", err)
	}

	go s.processYouTube(job.ID, ownerID, url)

	return job, nil
}

func (s *UploadService) GetJob(id, ownerID uuid.UUID) (*models.UploadJob, error) {
	return s.jobRepo.FindByID(id, ownerID)
}

func (s *UploadService) processYouTube(jobID, ownerID uuid.UUID, url string) {
	dir := filepath.Join(s.audioDir, ownerID.String())
	if err := os.MkdirAll(dir, 0755); err != nil {
		_ = s.jobRepo.UpdateStatus(jobID, "failed", err.Error())
		return
	}

	_ = s.jobRepo.UpdateStatus(jobID, "processing", "")

	fileID := uuid.New().String()
	rawPath := filepath.Join(dir, fileID+"_raw.mp3")
	outPath := filepath.Join(dir, fileID+".mp3")

	out, err := exec.Command("yt-dlp", "-x", "--audio-format", "mp3", "-o", rawPath, url).CombinedOutput()
	if err != nil {
		msg := truncate(string(out), 1000)
		if msg == "" {
			msg = err.Error()
		}
		_ = s.jobRepo.UpdateStatus(jobID, "failed", msg)
		return
	}

	out, err = exec.Command("ffmpeg", "-i", rawPath, "-af", "loudnorm", "-ar", "44100", "-b:a", "192k", outPath).CombinedOutput()
	_ = os.Remove(rawPath)
	if err != nil {
		msg := truncate(string(out), 1000)
		if msg == "" {
			msg = err.Error()
		}
		_ = s.jobRepo.UpdateStatus(jobID, "failed", msg)
		return
	}

	f, err := os.Open(outPath)
	if err != nil {
		_ = s.jobRepo.UpdateStatus(jobID, "failed", err.Error())
		return
	}
	defer f.Close()

	fi, _ := f.Stat()
	var fileSize int64
	if fi != nil {
		fileSize = fi.Size()
	}

	tags := s.parseTagsFromReader(ownerID, f, url)

	track := &models.Track{
		OwnerID:     ownerID,
		Title:       tags.Title,
		ArtistID:    tags.ArtistID,
		AlbumID:     tags.AlbumID,
		Year:        tags.Year,
		TrackNumber: tags.TrackNumber,
		DiscNumber:  tags.DiscNumber,
		Source:      "youtube",
		SourceURL:   url,
		FilePath:    outPath,
		FileSize:    fileSize,
		MimeType:    "audio/mpeg",
	}

	if err := s.trackRepo.Create(track); err != nil {
		_ = s.jobRepo.UpdateStatus(jobID, "failed", err.Error())
		return
	}

	_ = s.jobRepo.SetTrackID(jobID, track.ID)
	_ = s.jobRepo.UpdateStatus(jobID, "complete", "")
}

type parsedTags struct {
	Title       string
	ArtistID    *uuid.UUID
	AlbumID     *uuid.UUID
	Year        int
	TrackNumber int
	DiscNumber  int
}

func (s *UploadService) buildTrackFromFile(ownerID uuid.UUID, file multipart.File, header *multipart.FileHeader, destPath string) (*models.Track, error) {
	tags := s.parseTagsFromReader(ownerID, file, header.Filename)

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return &models.Track{
		OwnerID:     ownerID,
		Title:       tags.Title,
		ArtistID:    tags.ArtistID,
		AlbumID:     tags.AlbumID,
		Year:        tags.Year,
		TrackNumber: tags.TrackNumber,
		DiscNumber:  tags.DiscNumber,
		Source:      "upload",
		FilePath:    destPath,
		FileSize:    header.Size,
		MimeType:    mimeType,
	}, nil
}

func (s *UploadService) parseTagsFromReader(ownerID uuid.UUID, r io.ReadSeeker, fallbackTitle string) parsedTags {
	m, err := tag.ReadFrom(r)
	if err != nil {
		return parsedTags{Title: fallbackTitle}
	}

	title := m.Title()
	if title == "" {
		title = fallbackTitle
	}

	artistID, _ := s.catalogSvc.findOrCreateArtist(ownerID, m.Artist())
	albumID, _ := s.catalogSvc.findOrCreateAlbum(ownerID, m.Album(), artistID, m.Year())

	trackNum, _ := m.Track()
	discNum, _ := m.Disc()

	return parsedTags{
		Title:       title,
		ArtistID:    artistID,
		AlbumID:     albumID,
		Year:        m.Year(),
		TrackNumber: trackNum,
		DiscNumber:  discNum,
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}
