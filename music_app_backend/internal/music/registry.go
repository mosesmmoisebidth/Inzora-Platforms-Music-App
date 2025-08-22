package music

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ProviderRegistry manages multiple music providers
type ProviderRegistry struct {
	providers   map[string]MusicProvider
	enabledOnly []string
	mu          sync.RWMutex
}

// NewProviderRegistry creates a new provider registry
func NewProviderRegistry(enabledProviders []string) *ProviderRegistry {
	return &ProviderRegistry{
		providers:   make(map[string]MusicProvider),
		enabledOnly: enabledProviders,
	}
}

// Register adds a provider to the registry
func (r *ProviderRegistry) Register(provider MusicProvider) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := provider.GetName()
	
	// Check if this provider should be enabled
	enabled := false
	for _, enabledProvider := range r.enabledOnly {
		if enabledProvider == name {
			enabled = true
			break
		}
	}

	if !enabled {
		return fmt.Errorf("provider %s is not enabled", name)
	}

	r.providers[name] = provider
	return nil
}

// GetProvider returns a specific provider by name
func (r *ProviderRegistry) GetProvider(name string) (MusicProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, exists := r.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}

	return provider, nil
}

// GetEnabledProviders returns all enabled providers
func (r *ProviderRegistry) GetEnabledProviders() []MusicProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	providers := make([]MusicProvider, 0, len(r.providers))
	for _, provider := range r.providers {
		providers = append(providers, provider)
	}

	return providers
}

// GetProviderNames returns names of all enabled providers
func (r *ProviderRegistry) GetProviderNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}

	return names
}

// SearchAllProviders searches across all enabled providers
func (r *ProviderRegistry) SearchAllProviders(ctx context.Context, query string, page, size int, filters *SearchFilters) (map[string][]Track, map[string]*PageInfo, []error) {
	providers := r.GetEnabledProviders()
	results := make(map[string][]Track)
	pageInfos := make(map[string]*PageInfo)
	errors := make([]error, 0)

	// Use channels for concurrent searches
	type result struct {
		provider string
		tracks   []Track
		pageInfo *PageInfo
		err      error
	}

	resultChan := make(chan result, len(providers))

	// Start searches concurrently
	for _, provider := range providers {
		go func(p MusicProvider) {
			tracks, pageInfo, err := p.SearchTracks(ctx, query, page, size, filters)
			resultChan <- result{
				provider: p.GetName(),
				tracks:   tracks,
				pageInfo: pageInfo,
				err:      err,
			}
		}(provider)
	}

	// Collect results
	for i := 0; i < len(providers); i++ {
		res := <-resultChan
		if res.err != nil {
			errors = append(errors, res.err)
		} else {
			results[res.provider] = res.tracks
			pageInfos[res.provider] = res.pageInfo
		}
	}

	return results, pageInfos, errors
}

// HealthCheckAll checks health of all providers
func (r *ProviderRegistry) HealthCheckAll(ctx context.Context) map[string]error {
	providers := r.GetEnabledProviders()
	results := make(map[string]error)

	// Use channels for concurrent health checks
	type healthResult struct {
		provider string
		err      error
	}

	resultChan := make(chan healthResult, len(providers))

	// Start health checks concurrently
	for _, provider := range providers {
		go func(p MusicProvider) {
			err := p.IsHealthy(ctx)
			resultChan <- healthResult{
				provider: p.GetName(),
				err:      err,
			}
		}(provider)
	}

	// Collect results
	for i := 0; i < len(providers); i++ {
		res := <-resultChan
		results[res.provider] = res.err
	}

	return results
}

// MusicService wraps the provider registry with additional functionality
type MusicService struct {
	registry *ProviderRegistry
	config   *ProviderConfig
}

// NewMusicService creates a new music service
func NewMusicService(enabledProviders []string, timeout time.Duration, cacheTTL time.Duration, spotifyClientID, spotifyClientSecret string) *MusicService {
	config := &ProviderConfig{
		Timeout:   timeout,
		RateLimit: 100,
		CacheTTL:  cacheTTL,
		UserAgent: "music-app-backend/1.0",
	}

	registry := NewProviderRegistry(enabledProviders)

	service := &MusicService{
		registry: registry,
		config:   config,
	}

	// Initialize enabled providers
	service.initializeProviders(enabledProviders, spotifyClientID, spotifyClientSecret)

	return service
}

// initializeProviders initializes all enabled providers
func (m *MusicService) initializeProviders(enabledProviders []string, spotifyClientID, spotifyClientSecret string) {
	for _, providerName := range enabledProviders {
		switch providerName {
		case "itunes":
			itunesProvider := NewITunesProvider(m.config)
			if err := m.registry.Register(itunesProvider); err != nil {
				// Log error but don't fail
				continue
			}
		case "spotify":
			if spotifyClientID != "" && spotifyClientSecret != "" {
				spotifyProvider := NewSpotifyProvider(m.config, spotifyClientID, spotifyClientSecret)
				if err := m.registry.Register(spotifyProvider); err != nil {
					// Log error but don't fail
					continue
				}
			}
		}
	}
}

// SearchTracks searches for tracks using a specific provider or all providers
func (m *MusicService) SearchTracks(ctx context.Context, provider, query string, page, size int, filters *SearchFilters) ([]Track, *PageInfo, error) {
	if provider != "" {
		// Search using specific provider
		p, err := m.registry.GetProvider(provider)
		if err != nil {
			return nil, nil, err
		}
		return p.SearchTracks(ctx, query, page, size, filters)
	}

	// Search using all providers and combine results
	allResults, allPageInfos, errors := m.registry.SearchAllProviders(ctx, query, page, size, filters)
	
	if len(allResults) == 0 {
		if len(errors) > 0 {
			return nil, nil, errors[0]
		}
		return []Track{}, &PageInfo{Page: page, Size: size}, nil
	}

	// Combine results from all providers
	var combinedTracks []Track
	var totalResults int64

	for providerName, tracks := range allResults {
		combinedTracks = append(combinedTracks, tracks...)
		if pageInfo, exists := allPageInfos[providerName]; exists {
			totalResults += pageInfo.Total
		}
	}

	// Create combined page info
	pageInfo := &PageInfo{
		Page:       page,
		Size:       size,
		Total:      totalResults,
		HasNext:    len(combinedTracks) == size,
		HasPrev:    page > 1,
		TotalPages: int(totalResults+int64(size)-1) / size,
	}

	return combinedTracks, pageInfo, nil
}

// GetTrack gets a specific track from a provider
func (m *MusicService) GetTrack(ctx context.Context, provider, trackID string) (*Track, error) {
	p, err := m.registry.GetProvider(provider)
	if err != nil {
		return nil, err
	}
	return p.GetTrack(ctx, trackID)
}

// GetTopCharts gets top charts from a provider
func (m *MusicService) GetTopCharts(ctx context.Context, provider, country string, page, size int) ([]Track, *PageInfo, error) {
	p, err := m.registry.GetProvider(provider)
	if err != nil {
		return nil, nil, err
	}
	return p.GetTopCharts(ctx, country, page, size)
}

// GetCategories gets categories from all providers
func (m *MusicService) GetCategories(ctx context.Context) (map[string][]Category, error) {
	providers := m.registry.GetEnabledProviders()
	results := make(map[string][]Category)

	for _, provider := range providers {
		categories, err := provider.GetCategories(ctx)
		if err == nil {
			results[provider.GetName()] = categories
		}
	}

	return results, nil
}

// GetPlaylistsByCategory gets playlists by category from a provider
func (m *MusicService) GetPlaylistsByCategory(ctx context.Context, provider, categoryID string, page, size int) ([]PlaylistSummary, *PageInfo, error) {
	p, err := m.registry.GetProvider(provider)
	if err != nil {
		return nil, nil, err
	}
	return p.GetPlaylistsByCategory(ctx, categoryID, page, size)
}

// HealthCheck checks the health of all providers
func (m *MusicService) HealthCheck(ctx context.Context) map[string]error {
	return m.registry.HealthCheckAll(ctx)
}

// GetProviderNames returns names of all enabled providers
func (m *MusicService) GetProviderNames() []string {
	return m.registry.GetProviderNames()
}
