package arr

import (
	"context"
	"net/http"
	"testing"
)

func TestSonarrLookupSeriesUsesSeriesLookupEndpoint(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/api/v3/series/lookup" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("term"); got != "Silo" {
			t.Fatalf("term = %q", got)
		}
		return response(http.StatusOK, `[{"title":"Silo","year":2023,"tvdbId":403245}]`), nil
	})
	sonarr := NewSonarr(NewClient("http://sonarr.local", "key", &http.Client{Transport: transport}))

	results, err := sonarr.LookupSeries(context.Background(), "Silo")
	if err != nil {
		t.Fatalf("LookupSeries() error = %v", err)
	}
	if len(results) != 1 || results[0].Title != "Silo" || results[0].TVDBID != 403245 {
		t.Fatalf("results = %#v", results)
	}
}

func TestRadarrLookupMovieUsesMovieLookupEndpoint(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/api/v3/movie/lookup" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("term"); got != "Arrival" {
			t.Fatalf("term = %q", got)
		}
		return response(http.StatusOK, `[{"title":"Arrival","year":2016,"tmdbId":329865}]`), nil
	})
	radarr := NewRadarr(NewClient("http://radarr.local", "key", &http.Client{Transport: transport}))

	results, err := radarr.LookupMovie(context.Background(), "Arrival")
	if err != nil {
		t.Fatalf("LookupMovie() error = %v", err)
	}
	if len(results) != 1 || results[0].Title != "Arrival" || results[0].TMDBID != 329865 {
		t.Fatalf("results = %#v", results)
	}
}

func TestCommonReadOnlyEndpoints(t *testing.T) {
	seen := map[string]bool{}
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		seen[r.URL.Path] = true
		switch r.URL.Path {
		case "/api/v3/system/status":
			return response(http.StatusOK, `{"version":"4.0.0"}`), nil
		case "/api/v3/health":
			return response(http.StatusOK, `[]`), nil
		case "/api/v3/rootfolder":
			return response(http.StatusOK, `[{"id":1,"path":"/media/tv"}]`), nil
		case "/api/v3/qualityprofile":
			return response(http.StatusOK, `[{"id":2,"name":"HD-1080p"}]`), nil
		default:
			t.Fatalf("unexpected path = %q", r.URL.Path)
			return nil, nil
		}
	})
	service := NewService(NewClient("http://arr.local", "key", &http.Client{Transport: transport}))

	if status, err := service.Status(context.Background()); err != nil || status.Version != "4.0.0" {
		t.Fatalf("Status() = %#v, %v", status, err)
	}
	if health, err := service.Health(context.Background()); err != nil || len(health) != 0 {
		t.Fatalf("Health() = %#v, %v", health, err)
	}
	if roots, err := service.RootFolders(context.Background()); err != nil || roots[0].Path != "/media/tv" {
		t.Fatalf("RootFolders() = %#v, %v", roots, err)
	}
	if profiles, err := service.QualityProfiles(context.Background()); err != nil || profiles[0].Name != "HD-1080p" {
		t.Fatalf("QualityProfiles() = %#v, %v", profiles, err)
	}

	for _, path := range []string{"/api/v3/system/status", "/api/v3/health", "/api/v3/rootfolder", "/api/v3/qualityprofile"} {
		if !seen[path] {
			t.Fatalf("did not call %s", path)
		}
	}
}
