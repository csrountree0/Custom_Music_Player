package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"musicapp/backend/internal/config"
	"musicapp/backend/internal/handlers"
	"musicapp/backend/internal/middleware"
	"musicapp/backend/internal/repository"
	"musicapp/backend/internal/services"
	"musicapp/backend/pkg/token"
)

func New(db *gorm.DB, cfg *config.Config) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	userRepo     := repository.NewUserRepository(db)
	artistRepo   := repository.NewArtistRepository(db)
	albumRepo    := repository.NewAlbumRepository(db)
	trackRepo    := repository.NewTrackRepository(db)
	jobRepo      := repository.NewUploadJobRepository(db)
	playlistRepo := repository.NewPlaylistRepository(db)

	tokenSvc    := token.NewService(cfg.JWTSecret, cfg.JWTExpiry)
	authSvc     := services.NewAuthService(userRepo, tokenSvc, cfg.JWTExpiry, cfg.RefreshTokenExpiry, db)
	catalogSvc  := services.NewCatalogService(artistRepo, albumRepo, trackRepo)
	uploadSvc   := services.NewUploadService(catalogSvc, trackRepo, jobRepo, cfg.AudioDir)
	playlistSvc := services.NewPlaylistService(playlistRepo, trackRepo)
	searchSvc   := services.NewSearchService(artistRepo, albumRepo, trackRepo)

	authH     := handlers.NewAuthHandler(authSvc)
	catalogH  := handlers.NewCatalogHandler(catalogSvc)
	uploadH   := handlers.NewUploadHandler(uploadSvc)
	playlistH := handlers.NewPlaylistHandler(playlistSvc)
	searchH   := handlers.NewSearchHandler(searchSvc)

	authMW := middleware.AuthRequired(tokenSvc)

	r := gin.Default()
	r.GET("/health", healthHandler(db))

	api := r.Group("/api")
	api.Use(middleware.RateLimit())
	api.GET("/health", healthHandler(db))

	auth := api.Group("/auth")
	auth.POST("/register", authH.Register)
	auth.POST("/login",    authH.Login)
	auth.POST("/refresh",  authH.Refresh)

	me := api.Group("/me", authMW)
	me.GET("",      authH.GetMe)
	me.PUT("",      authH.UpdateMe)
	me.POST("/logout", authH.Logout)

	me.GET("/artists",        catalogH.ListArtists)
	me.POST("/artists",       catalogH.CreateArtist)
	me.GET("/artists/:id",    catalogH.GetArtist)
	me.PUT("/artists/:id",    catalogH.UpdateArtist)
	me.DELETE("/artists/:id", catalogH.DeleteArtist)

	me.GET("/albums",        catalogH.ListAlbums)
	me.POST("/albums",       catalogH.CreateAlbum)
	me.GET("/albums/:id",    catalogH.GetAlbum)
	me.PUT("/albums/:id",    catalogH.UpdateAlbum)
	me.DELETE("/albums/:id", catalogH.DeleteAlbum)

	me.GET("/tracks",              catalogH.ListTracks)
	me.GET("/tracks/:id",          catalogH.GetTrack)
	me.GET("/tracks/:id/stream",   catalogH.StreamTrack)
	me.PUT("/tracks/:id",          catalogH.UpdateTrack)
	me.DELETE("/tracks/:id",       catalogH.DeleteTrack)
	me.POST("/tracks/upload",      uploadH.UploadFile)
	me.POST("/tracks/import",      uploadH.ImportYouTube)

	me.GET("/upload-jobs/:id", uploadH.GetUploadJob)

	me.GET("/playlists",             playlistH.ListPlaylists)
	me.POST("/playlists",            playlistH.CreatePlaylist)
	me.GET("/playlists/:id",         playlistH.GetPlaylist)
	me.PUT("/playlists/:id",         playlistH.UpdatePlaylist)
	me.DELETE("/playlists/:id",      playlistH.DeletePlaylist)
	me.GET("/playlists/:id/tracks",               playlistH.ListPlaylistTracks)
	me.POST("/playlists/:id/tracks",              playlistH.AddTrackToPlaylist)
	me.PUT("/playlists/:id/tracks/reorder",       playlistH.ReorderPlaylistTracks)
	me.DELETE("/playlists/:id/tracks/:trackID",   playlistH.RemoveTrackFromPlaylist)

	me.GET("/search", searchH.Search)

	return r
}

func healthHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "db": "unreachable", "error": err.Error()})
			return
		}
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "db": "unreachable", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "db": "connected"})
	}
}
