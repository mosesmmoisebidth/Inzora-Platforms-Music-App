package music

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// SpotifyProvider implements the MusicProvider interface for Spotify Web API
type SpotifyProvider struct {
	config      *ProviderConfig
	httpClient  *http.Client
	accessToken string
	tokenExpiry time.Time
	clientID    string
	clientSecret string
}

// SpotifyTrack represents a track from Spotify API
type SpotifyTrack struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Artists      []SpotifyArtist     `json:"artists"`
	Album        SpotifyAlbum        `json:"album"`
	DurationMs   int64               `json:"duration_ms"`
	Explicit     bool                `json:"explicit"`
	Popularity   int                 `json:"popularity"`
	PreviewURL   *string             `json:"preview_url"`
	TrackNumber  int                 `json:"track_number"`
	ExternalUrls SpotifyExternalUrls `json:"external_urls"`
}

type SpotifyArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SpotifyAlbum struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Images      []SpotifyImage        `json:"images"`
	ReleaseDate string                `json:"release_date"`
	TotalTracks int                   `json:"total_tracks"`
}

type SpotifyImage struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type SpotifyExternalUrls struct {
	Spotify string `json:"spotify"`
}

type SpotifySearchResponse struct {
	Tracks SpotifyTracksResponse `json:"tracks"`
}

type SpotifyTracksResponse struct {
	Items    []SpotifyTrack `json:"items"`
	Total    int            `json:"total"`
	Limit    int            `json:"limit"`
	Offset   int            `json:"offset"`
	Previous *string        `json:"previous"`
	Next     *string        `json:"next"`
}

type SpotifyCategory struct {
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	Icons []SpotifyImage `json:"icons"`
}

type SpotifyCategoriesResponse struct {
	Categories struct {
		Items []SpotifyCategory `json:"items"`
	} `json:"categories"`
}

type SpotifyPlaylist struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Images      []SpotifyImage        `json:"images"`
	Tracks      SpotifyTracksTotal    `json:"tracks"`
	ExternalUrls SpotifyExternalUrls  `json:"external_urls"`
	Owner       SpotifyUser           `json:"owner"`
}

type SpotifyTracksTotal struct {
	Total int `json:"total"`
}

type SpotifyUser struct {
	DisplayName string `json:"display_name"`
}

type SpotifyPlaylistsResponse struct {
	Items []SpotifyPlaylist `json:"items"`
	Total int               `json:"total"`
	Limit int               `json:"limit"`
	Offset int              `json:"offset"`
}

type SpotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// NewSpotifyProvider creates a new Spotify provider
func NewSpotifyProvider(config *ProviderConfig, clientID, clientSecret string) *SpotifyProvider {
	return &SpotifyProvider{
		config:       config,
		httpClient:   &http.Client{Timeout: config.Timeout},
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

// GetName returns the provider name
func (s *SpotifyProvider) GetName() string {
	return "spotify"
}

// ensureToken ensures we have a valid access token
func (s *SpotifyProvider) ensureToken(ctx context.Context) error {
	if s.accessToken != "" && time.Now().Before(s.tokenExpiry) {
		return nil
	}

	// Get client credentials token
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, "POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	// Set authorization header
	auth := base64.StdEncoding.EncodeToString([]byte(s.clientID + ":" + s.clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get access token: %d %s", resp.StatusCode, string(body))
	}

	var tokenResp SpotifyTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	s.accessToken = tokenResp.AccessToken
	s.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second) // Expire 60s early

	return nil
}

// makeRequest makes an authenticated request to Spotify API
func (s *SpotifyProvider) makeRequest(ctx context.Context, endpoint string) (*http.Response, error) {
	if err := s.ensureToken(ctx); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.spotify.com/v1"+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("User-Agent", s.config.UserAgent)

	return s.httpClient.Do(req)
}

// SearchTracks searches for tracks
func (s *SpotifyProvider) SearchTracks(ctx context.Context, query string, page, size int, filters *SearchFilters) ([]Track, *PageInfo, error) {
	if size > 50 {
		size = 50 // Spotify API limit
	}

	offset := (page - 1) * size
	
	params := url.Values{}
	params.Set("q", query)
	params.Set("type", "track")
	params.Set("limit", strconv.Itoa(size))
	params.Set("offset", strconv.Itoa(offset))
	params.Set("market", "US")

	endpoint := "/search?" + params.Encode()
	resp, err := s.makeRequest(ctx, endpoint)
	if err != nil {
		return nil, nil, NewProviderError("spotify", "Search failed", "SEARCH_ERROR", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, NewProviderError("spotify", fmt.Sprintf("API error: %d %s", resp.StatusCode, string(body)), "API_ERROR", nil)
	}

	var searchResp SpotifySearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, nil, NewProviderError("spotify", "Failed to decode response", "DECODE_ERROR", err)
	}

	tracks := make([]Track, len(searchResp.Tracks.Items))
	for i, spotifyTrack := range searchResp.Tracks.Items {
		tracks[i] = s.convertTrack(spotifyTrack)
	}

	pageInfo := &PageInfo{
		Page:       page,
		Size:       size,
		Total:      int64(searchResp.Tracks.Total),
		HasNext:    searchResp.Tracks.Next != nil,
		HasPrev:    searchResp.Tracks.Previous != nil,
		TotalPages: (searchResp.Tracks.Total + size - 1) / size,
	}

	return tracks, pageInfo, nil
}

// GetTrack gets a specific track by ID
func (s *SpotifyProvider) GetTrack(ctx context.Context, trackID string) (*Track, error) {
	endpoint := "/tracks/" + trackID + "?market=US"
	resp, err := s.makeRequest(ctx, endpoint)
	if err != nil {
		return nil, NewProviderError("spotify", "Get track failed", "GET_TRACK_ERROR", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, NewProviderError("spotify", "Track not found", "NOT_FOUND", nil)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, NewProviderError("spotify", fmt.Sprintf("API error: %d %s", resp.StatusCode, string(body)), "API_ERROR", nil)
	}

	var spotifyTrack SpotifyTrack
	if err := json.NewDecoder(resp.Body).Decode(&spotifyTrack); err != nil {
		return nil, NewProviderError("spotify", "Failed to decode response", "DECODE_ERROR", err)
	}

	track := s.convertTrack(spotifyTrack)
	return &track, nil
}

// GetTopCharts gets top charts for a country
func (s *SpotifyProvider) GetTopCharts(ctx context.Context, country string, page, size int) ([]Track, *PageInfo, error) {
	if size > 50 {
		size = 50
	}

	offset := (page - 1) * size
	
	// Use Spotify's "Top 50" playlist for the country
	playlistID := "37i9dQZEVXbLRQDuF5jeBp" // Global Top 50 as fallback
	if country == "US" {
		playlistID = "37i9dQZEVXbLRQDuF5jeBp"
	}

	params := url.Values{}
	params.Set("limit", strconv.Itoa(size))
	params.Set("offset", strconv.Itoa(offset))
	params.Set("market", country)

	endpoint := "/playlists/" + playlistID + "/tracks?" + params.Encode()
	resp, err := s.makeRequest(ctx, endpoint)
	if err != nil {
		return nil, nil, NewProviderError("spotify", "Top charts failed", "TOP_CHARTS_ERROR", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, NewProviderError("spotify", fmt.Sprintf("API error: %d %s", resp.StatusCode, string(body)), "API_ERROR", nil)
	}

	var tracksResp SpotifyTracksResponse
	if err := json.NewDecoder(resp.Body).Decode(&tracksResp); err != nil {
		return nil, nil, NewProviderError("spotify", "Failed to decode response", "DECODE_ERROR", err)
	}

	tracks := make([]Track, len(tracksResp.Items))
	for i, spotifyTrack := range tracksResp.Items {
		tracks[i] = s.convertTrack(spotifyTrack)
	}

	pageInfo := &PageInfo{
		Page:       page,
		Size:       size,
		Total:      int64(tracksResp.Total),
		HasNext:    tracksResp.Next != nil,
		HasPrev:    tracksResp.Previous != nil,
		TotalPages: (tracksResp.Total + size - 1) / size,
	}

	return tracks, pageInfo, nil
}

// GetCategories gets available music categories
func (s *SpotifyProvider) GetCategories(ctx context.Context) ([]Category, error) {
	params := url.Values{}
	params.Set("limit", "50")
	params.Set("country", "US")

	endpoint := "/browse/categories?" + params.Encode()
	resp, err := s.makeRequest(ctx, endpoint)
	if err != nil {
		return nil, NewProviderError("spotify", "Get categories failed", "CATEGORIES_ERROR", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, NewProviderError("spotify", fmt.Sprintf("API error: %d %s", resp.StatusCode, string(body)), "API_ERROR", nil)
	}

	var categoriesResp SpotifyCategoriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&categoriesResp); err != nil {
		return nil, NewProviderError("spotify", "Failed to decode response", "DECODE_ERROR", err)
	}

	categories := make([]Category, len(categoriesResp.Categories.Items))
	for i, spotifyCategory := range categoriesResp.Categories.Items {
		iconURL := ""
		if len(spotifyCategory.Icons) > 0 {
			iconURL = spotifyCategory.Icons[0].URL
		}
		
		categories[i] = Category{
			ID:      spotifyCategory.ID,
			Name:    spotifyCategory.Name,
			IconURL: iconURL,
		}
	}

	return categories, nil
}

// GetPlaylistsByCategory gets playlists for a specific category
func (s *SpotifyProvider) GetPlaylistsByCategory(ctx context.Context, categoryID string, page, size int) ([]PlaylistSummary, *PageInfo, error) {
	if size > 50 {
		size = 50
	}

	offset := (page - 1) * size
	
	params := url.Values{}
	params.Set("limit", strconv.Itoa(size))
	params.Set("offset", strconv.Itoa(offset))
	params.Set("country", "US")

	endpoint := "/browse/categories/" + categoryID + "/playlists?" + params.Encode()
	resp, err := s.makeRequest(ctx, endpoint)
	if err != nil {
		return nil, nil, NewProviderError("spotify", "Get category playlists failed", "CATEGORY_PLAYLISTS_ERROR", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, NewProviderError("spotify", fmt.Sprintf("API error: %d %s", resp.StatusCode, string(body)), "API_ERROR", nil)
	}

	var playlistsResp struct {
		Playlists SpotifyPlaylistsResponse `json:"playlists"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&playlistsResp); err != nil {
		return nil, nil, NewProviderError("spotify", "Failed to decode response", "DECODE_ERROR", err)
	}

	playlists := make([]PlaylistSummary, len(playlistsResp.Playlists.Items))
	for i, spotifyPlaylist := range playlistsResp.Playlists.Items {
		coverURL := ""
		if len(spotifyPlaylist.Images) > 0 {
			coverURL = spotifyPlaylist.Images[0].URL
		}

		playlists[i] = PlaylistSummary{
			ID:          spotifyPlaylist.ID,
			Title:       spotifyPlaylist.Name,
			Description: spotifyPlaylist.Description,
			CoverURL:    coverURL,
			TrackCount:  spotifyPlaylist.Tracks.Total,
			Provider:    "spotify",
			ExternalURL: spotifyPlaylist.ExternalUrls.Spotify,
			Creator:     spotifyPlaylist.Owner.DisplayName,
		}
	}

	pageInfo := &PageInfo{
		Page:       page,
		Size:       size,
		Total:      int64(playlistsResp.Playlists.Total),
		HasNext:    playlistsResp.Playlists.Offset+playlistsResp.Playlists.Limit < playlistsResp.Playlists.Total,
		HasPrev:    playlistsResp.Playlists.Offset > 0,
		TotalPages: (playlistsResp.Playlists.Total + size - 1) / size,
	}

	return playlists, pageInfo, nil
}

// IsHealthy checks if the provider is healthy
func (s *SpotifyProvider) IsHealthy(ctx context.Context) error {
	if err := s.ensureToken(ctx); err != nil {
		return NewProviderError("spotify", "Health check failed", "HEALTH_CHECK_ERROR", err)
	}

	// Try a simple API call
	resp, err := s.makeRequest(ctx, "/browse/categories?limit=1")
	if err != nil {
		return NewProviderError("spotify", "Health check failed", "HEALTH_CHECK_ERROR", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NewProviderError("spotify", "Health check failed", "HEALTH_CHECK_ERROR", fmt.Errorf("status: %d", resp.StatusCode))
	}

	return nil
}

// convertTrack converts a Spotify track to our Track format
func (s *SpotifyProvider) convertTrack(spotifyTrack SpotifyTrack) Track {
	artistNames := make([]string, len(spotifyTrack.Artists))
	for i, artist := range spotifyTrack.Artists {
		artistNames[i] = artist.Name
	}

	artworkURL := ""
	if len(spotifyTrack.Album.Images) > 0 {
		artworkURL = spotifyTrack.Album.Images[0].URL
	}

	previewURL := ""
	if spotifyTrack.PreviewURL != nil {
		previewURL = *spotifyTrack.PreviewURL
	}

	return Track{
		ID:          spotifyTrack.ID,
		Title:       spotifyTrack.Name,
		Artist:      strings.Join(artistNames, ", "),
		Album:       spotifyTrack.Album.Name,
		Duration:    spotifyTrack.DurationMs,
		ArtworkURL:  artworkURL,
		PreviewURL:  previewURL,
		TrackNumber: spotifyTrack.TrackNumber,
		ReleaseDate: spotifyTrack.Album.ReleaseDate,
		Provider:    "spotify",
		ExternalURL: spotifyTrack.ExternalUrls.Spotify,
		Explicit:    spotifyTrack.Explicit,
		Popularity:  spotifyTrack.Popularity,
	}
}
