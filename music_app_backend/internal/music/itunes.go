package music

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"github.com/go-resty/resty/v2"
)

const (
	itunesSearchBaseURL = "https://itunes.apple.com/search"
	itunesLookupBaseURL = "https://itunes.apple.com/lookup"
	itunesUserAgent     = "music-app-backend/1.0"
)

// ITunesProvider implements the MusicProvider interface for iTunes Search API
type ITunesProvider struct {
	client *resty.Client
	config *ProviderConfig
}

// iTunesSearchResponse represents the response from iTunes Search API
type iTunesSearchResponse struct {
	ResultCount int               `json:"resultCount"`
	Results     []iTunesTrackData `json:"results"`
}

// iTunesTrackData represents track data from iTunes API
type iTunesTrackData struct {
	TrackID                 int     `json:"trackId"`
	ArtistID                int     `json:"artistId"`
	CollectionID            int     `json:"collectionId"`
	ArtistName              string  `json:"artistName"`
	CollectionName          string  `json:"collectionName"`
	TrackName               string  `json:"trackName"`
	CollectionCensoredName  string  `json:"collectionCensoredName"`
	TrackCensoredName       string  `json:"trackCensoredName"`
	ArtistViewURL           string  `json:"artistViewUrl"`
	CollectionViewURL       string  `json:"collectionViewUrl"`
	TrackViewURL            string  `json:"trackViewUrl"`
	PreviewURL              string  `json:"previewUrl"`
	ArtworkURL30            string  `json:"artworkUrl30"`
	ArtworkURL60            string  `json:"artworkUrl60"`
	ArtworkURL100           string  `json:"artworkUrl100"`
	CollectionPrice         float64 `json:"collectionPrice"`
	TrackPrice              float64 `json:"trackPrice"`
	ReleaseDate             string  `json:"releaseDate"`
	CollectionExplicitness  string  `json:"collectionExplicitness"`
	TrackExplicitness       string  `json:"trackExplicitness"`
	DiscCount               int     `json:"discCount"`
	DiscNumber              int     `json:"discNumber"`
	TrackCount              int     `json:"trackCount"`
	TrackNumber             int     `json:"trackNumber"`
	TrackTimeMillis         int64   `json:"trackTimeMillis"`
	Country                 string  `json:"country"`
	Currency                string  `json:"currency"`
	PrimaryGenreName        string  `json:"primaryGenreName"`
	ContentAdvisoryRating   string  `json:"contentAdvisoryRating"`
	WrapperType             string  `json:"wrapperType"`
	Kind                    string  `json:"kind"`
}

// NewITunesProvider creates a new iTunes provider
func NewITunesProvider(config *ProviderConfig) *ITunesProvider {
	client := resty.New()
	client.SetTimeout(config.Timeout)
	client.SetHeader("User-Agent", config.UserAgent)

	return &ITunesProvider{
		client: client,
		config: config,
	}
}

// GetName returns the provider name
func (i *ITunesProvider) GetName() string {
	return "itunes"
}

// SearchTracks searches for tracks on iTunes
func (i *ITunesProvider) SearchTracks(ctx context.Context, query string, page, size int, filters *SearchFilters) ([]Track, *PageInfo, error) {
	params := url.Values{}
	params.Set("term", query)
	params.Set("media", "music")
	params.Set("entity", "song")
	params.Set("limit", strconv.Itoa(size))
	params.Set("offset", strconv.Itoa((page-1)*size))

	// Apply filters
	if filters != nil {
		if filters.Genre != "" {
			params.Set("genreId", filters.Genre)
		}
		if filters.Explicit != nil {
			if *filters.Explicit {
				params.Set("explicit", "Yes")
			} else {
				params.Set("explicit", "No")
			}
		}
	}

	resp, err := i.client.R().
		SetContext(ctx).
		SetQueryParamsFromValues(params).
		Get(itunesSearchBaseURL)

	if err != nil {
		return nil, nil, NewProviderError(i.GetName(), "Failed to search tracks", "SEARCH_ERROR", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, nil, NewProviderError(i.GetName(), "API request failed", "API_ERROR", fmt.Errorf("status code: %d", resp.StatusCode()))
	}

	var searchResp iTunesSearchResponse
	if err := json.Unmarshal(resp.Body(), &searchResp); err != nil {
		return nil, nil, NewProviderError(i.GetName(), "Failed to parse response", "PARSE_ERROR", err)
	}

	tracks := make([]Track, 0, len(searchResp.Results))
	for _, result := range searchResp.Results {
		if result.Kind == "song" {
			track := i.convertToTrack(result)
			tracks = append(tracks, track)
		}
	}

	pageInfo := &PageInfo{
		Page:       page,
		Size:       size,
		Total:      int64(searchResp.ResultCount),
		HasNext:    len(tracks) == size,
		HasPrev:    page > 1,
		TotalPages: (searchResp.ResultCount + size - 1) / size,
	}

	return tracks, pageInfo, nil
}

// GetTrack gets a specific track by ID
func (i *ITunesProvider) GetTrack(ctx context.Context, trackID string) (*Track, error) {
	params := url.Values{}
	params.Set("id", trackID)
	params.Set("entity", "song")

	resp, err := i.client.R().
		SetContext(ctx).
		SetQueryParamsFromValues(params).
		Get(itunesLookupBaseURL)

	if err != nil {
		return nil, NewProviderError(i.GetName(), "Failed to get track", "GET_ERROR", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, NewProviderError(i.GetName(), "API request failed", "API_ERROR", fmt.Errorf("status code: %d", resp.StatusCode()))
	}

	var lookupResp iTunesSearchResponse
	if err := json.Unmarshal(resp.Body(), &lookupResp); err != nil {
		return nil, NewProviderError(i.GetName(), "Failed to parse response", "PARSE_ERROR", err)
	}

	if lookupResp.ResultCount == 0 {
		return nil, NewProviderError(i.GetName(), "Track not found", "NOT_FOUND", nil)
	}

	track := i.convertToTrack(lookupResp.Results[0])
	return &track, nil
}

// GetTopCharts gets top charts for a country (iTunes doesn't have a direct top charts API)
func (i *ITunesProvider) GetTopCharts(ctx context.Context, country string, page, size int) ([]Track, *PageInfo, error) {
	// iTunes doesn't have a direct top charts API, so we'll search for popular songs
	return i.SearchTracks(ctx, "pop", page, size, nil)
}

// GetCategories gets available music categories
func (i *ITunesProvider) GetCategories(ctx context.Context) ([]Category, error) {
	// iTunes doesn't provide a categories endpoint, return predefined categories
	categories := []Category{
		{ID: "1", Name: "Pop", Description: "Popular music"},
		{ID: "2", Name: "Rock", Description: "Rock music"},
		{ID: "3", Name: "Hip-Hop", Description: "Hip-Hop music"},
		{ID: "4", Name: "R&B", Description: "R&B music"},
		{ID: "5", Name: "Country", Description: "Country music"},
		{ID: "6", Name: "Electronic", Description: "Electronic music"},
		{ID: "7", Name: "Jazz", Description: "Jazz music"},
		{ID: "8", Name: "Classical", Description: "Classical music"},
	}

	return categories, nil
}

// GetPlaylistsByCategory gets playlists for a specific category (iTunes doesn't support this)
func (i *ITunesProvider) GetPlaylistsByCategory(ctx context.Context, categoryID string, page, size int) ([]PlaylistSummary, *PageInfo, error) {
	return nil, nil, NewProviderError(i.GetName(), "Playlists not supported", "NOT_SUPPORTED", nil)
}

// IsHealthy checks if the provider is healthy
func (i *ITunesProvider) IsHealthy(ctx context.Context) error {
	resp, err := i.client.R().
		SetContext(ctx).
		SetQueryParam("term", "test").
		SetQueryParam("limit", "1").
		Get(itunesSearchBaseURL)

	if err != nil {
		return NewProviderError(i.GetName(), "Health check failed", "HEALTH_ERROR", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return NewProviderError(i.GetName(), "Health check failed", "HEALTH_ERROR", fmt.Errorf("status code: %d", resp.StatusCode()))
	}

	return nil
}

// convertToTrack converts iTunes track data to our Track struct
func (i *ITunesProvider) convertToTrack(data iTunesTrackData) Track {
	// Get the highest resolution artwork URL
	artworkURL := data.ArtworkURL100
	if artworkURL == "" {
		artworkURL = data.ArtworkURL60
	}
	if artworkURL == "" {
		artworkURL = data.ArtworkURL30
	}

	// Replace artwork resolution for higher quality
	if artworkURL != "" {
		artworkURL = strings.Replace(artworkURL, "100x100bb", "600x600bb", 1)
	}

	return Track{
		ID:          strconv.Itoa(data.TrackID),
		Title:       data.TrackName,
		Artist:      data.ArtistName,
		Album:       data.CollectionName,
		Duration:    data.TrackTimeMillis,
		ArtworkURL:  artworkURL,
		PreviewURL:  data.PreviewURL,
		TrackNumber: data.TrackNumber,
		ReleaseDate: data.ReleaseDate,
		Genre:       data.PrimaryGenreName,
		Provider:    i.GetName(),
		ExternalURL: data.TrackViewURL,
		Explicit:    data.TrackExplicitness == "explicit",
	}
}
