package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	timeout := 5 * time.Second
	client := NewClient(timeout, 5 * time.Minute)

	if client.httpClient.Timeout != timeout {
		t.Errorf("got timeout %v, want %v", client.httpClient.Timeout, timeout)
	}
	if client.baseURL != "https://pokeapi.co/api/v2" {
		t.Errorf("got baseURL %q, want %q", client.baseURL, "https://pokeapi.co/api/v2")
	}
}

func TestGetAreas_Mock(t *testing.T) {
	t.Run("default_page", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/location-area/" {
				t.Errorf("got path %s, want /location-area/", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"count": 2,
				"next": "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
				"previous": null,
				"results": [
					{"name": "canalave-city-area", "url": "https://pokeapi.co/api/v2/location-area/1/"},
					{"name": "eterna-city-area", "url": "https://pokeapi.co/api/v2/location-area/2/"}
				]
			}`))
		}))
		defer server.Close()

		client := NewClient(5 * time.Second, 5 * time.Minute)
		client.baseURL = server.URL

		areas, err := client.GetAreas(nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(areas.Results) != 2 || areas.Results[0].Name != "canalave-city-area" {
			t.Errorf("unexpected areas: %+v", areas)
		}
	})

	t.Run("next_page", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"count": 2,
				"next": null,
				"previous": "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
				"results": [{"name": "pastoria-city-area", "url": "https://pokeapi.co/api/v2/location-area/3/"}]
			}`))
		}))
		defer server.Close()

		client := NewClient(5 * time.Second, 5 * time.Minute)
		pageURL := server.URL + "/location-area/?offset=20&limit=20"

		areas, err := client.GetAreas(&pageURL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if areas.Previous == nil || areas.Results[0].Name != "pastoria-city-area" {
			t.Errorf("unexpected pagination result: %+v", areas)
		}
	})

	t.Run("server_error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewClient(5 * time.Second, 5 * time.Minute)
		client.baseURL = server.URL

		if _, err := client.GetAreas(nil); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{broken`))
		}))
		defer server.Close()

		client := NewClient(5 * time.Second, 5 * time.Minute)
		client.baseURL = server.URL

		if _, err := client.GetAreas(nil); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestGetAreas_Live(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live test")
	}

	client := NewClient(10 * time.Second, 5 * time.Minute)
	areas, err := client.GetAreas(nil)
	if err != nil {
		t.Fatalf("live request failed: %v", err)
	}
	if areas.Count == 0 || len(areas.Results) == 0 {
		t.Fatalf("empty live response: %+v", areas)
	}
}
