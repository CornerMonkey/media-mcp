package arr

import "context"

type Service struct {
	client *Client
}

type Sonarr struct {
	*Service
}

type Radarr struct {
	*Service
}

func NewService(client *Client) *Service {
	return &Service{client: client}
}

func NewSonarr(client *Client) *Sonarr {
	return &Sonarr{Service: NewService(client)}
}

func NewRadarr(client *Client) *Radarr {
	return &Radarr{Service: NewService(client)}
}

type Status struct {
	Version string `json:"version"`
	OSName  string `json:"osName,omitempty"`
	AppName string `json:"appName,omitempty"`
}

type HealthCheck struct {
	Source  string `json:"source"`
	Type    string `json:"type"`
	Message string `json:"message"`
	WikiURL string `json:"wikiUrl,omitempty"`
}

type RootFolder struct {
	ID         int    `json:"id"`
	Path       string `json:"path"`
	Accessible bool   `json:"accessible"`
	FreeSpace  int64  `json:"freeSpace,omitempty"`
}

type QualityProfile struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	UpgradeAllowed bool   `json:"upgradeAllowed"`
}

type Series struct {
	ID      int    `json:"id,omitempty"`
	TVDBID  int    `json:"tvdbId,omitempty"`
	IMDBID  string `json:"imdbId,omitempty"`
	Title   string `json:"title"`
	Year    int    `json:"year,omitempty"`
	Status  string `json:"status,omitempty"`
	Network string `json:"network,omitempty"`
}

type Movie struct {
	ID     int    `json:"id,omitempty"`
	TMDBID int    `json:"tmdbId,omitempty"`
	IMDBID string `json:"imdbId,omitempty"`
	Title  string `json:"title"`
	Year   int    `json:"year,omitempty"`
	Status string `json:"status,omitempty"`
}

func (s *Service) Status(ctx context.Context) (Status, error) {
	var out Status
	err := s.client.Get(ctx, "system/status", nil, &out)
	return out, err
}

func (s *Service) Health(ctx context.Context) ([]HealthCheck, error) {
	var out []HealthCheck
	err := s.client.Get(ctx, "health", nil, &out)
	return out, err
}

func (s *Service) RootFolders(ctx context.Context) ([]RootFolder, error) {
	var out []RootFolder
	err := s.client.Get(ctx, "rootfolder", nil, &out)
	return out, err
}

func (s *Service) QualityProfiles(ctx context.Context) ([]QualityProfile, error) {
	var out []QualityProfile
	err := s.client.Get(ctx, "qualityprofile", nil, &out)
	return out, err
}

func (s *Sonarr) LookupSeries(ctx context.Context, term string) ([]Series, error) {
	var out []Series
	err := s.client.Get(ctx, "series/lookup", map[string]string{"term": term}, &out)
	return out, err
}

func (s *Sonarr) ListSeries(ctx context.Context) ([]Series, error) {
	var out []Series
	err := s.client.Get(ctx, "series", nil, &out)
	return out, err
}

func (r *Radarr) LookupMovie(ctx context.Context, term string) ([]Movie, error) {
	var out []Movie
	err := r.client.Get(ctx, "movie/lookup", map[string]string{"term": term}, &out)
	return out, err
}

func (r *Radarr) ListMovies(ctx context.Context) ([]Movie, error) {
	var out []Movie
	err := r.client.Get(ctx, "movie", nil, &out)
	return out, err
}
