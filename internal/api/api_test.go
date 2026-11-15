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

func TestGetAreaInfo_Mock(t *testing.T) {
	t.Run("success_and_cache", func(t *testing.T) {
		calls := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			if r.URL.Path != "/location-area/pastoria-city-area" {
				t.Errorf("got path %s, want /location-area/pastoria-city-area", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": 1,
				"name": "pastoria-city-area",
				"pokemon_encounters": [
					{
						"pokemon": {"name": "tentacool", "url": "https://pokeapi.co/api/v2/pokemon/72/"}
					},
					{
						"pokemon": {"name": "pelipper", "url": "https://pokeapi.co/api/v2/pokemon/279/"}
					}
				]
			}`))
		}))
		defer server.Close()

		client := NewClient(5*time.Second, 5*time.Minute)
		client.baseURL = server.URL

		// Перший запит — іде на сервер
		info, err := client.GetAreaInfo("pastoria-city-area")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if info.Name != "pastoria-city-area" || len(info.PokemonEncounters) != 2 {
			t.Errorf("unexpected area info: %+v", info)
		}
		if info.PokemonEncounters[0].Pokemon.Name != "tentacool" {
			t.Errorf("got pokemon %q, want tentacool", info.PokemonEncounters[0].Pokemon.Name)
		}
		if calls != 1 {
			t.Errorf("expected 1 server call, got %d", calls)
		}

		// Другий запит — має взятися з кешу без звернення до сервера
		cachedInfo, err := client.GetAreaInfo("pastoria-city-area")
		if err != nil {
			t.Fatalf("unexpected error on cached call: %v", err)
		}
		if cachedInfo.Name != "pastoria-city-area" {
			t.Errorf("unexpected cached info: %+v", cachedInfo)
		}
		if calls != 1 {
			t.Errorf("expected cached call not to hit server, got %d calls", calls)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}))
		defer server.Close()

		client := NewClient(5*time.Second, 5*time.Minute)
		client.baseURL = server.URL

		if _, err := client.GetAreaInfo("unknown-area"); err == nil {
			t.Fatal("expected error for 404, got nil")
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{broken`))
		}))
		defer server.Close()

		client := NewClient(5*time.Second, 5*time.Minute)
		client.baseURL = server.URL

		if _, err := client.GetAreaInfo("broken-area"); err == nil {
			t.Fatal("expected error for invalid json, got nil")
		}
	})
}

func TestGetAreaInfo_Live(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live test")
	}

	client := NewClient(10*time.Second, 5*time.Minute)
	info, err := client.GetAreaInfo("canalave-city-area")
	if err != nil {
		t.Fatalf("live request failed: %v", err)
	}
	if info.Name != "canalave-city-area" || len(info.PokemonEncounters) == 0 {
		t.Fatalf("unexpected live response: %+v", info)
	}
}

func TestGetPokemonInfo_Mock(t *testing.T) {
	t.Run("success_and_cache", func(t *testing.T) {
		calls := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			if r.URL.Path != "/pokemon/pikachu" {
				t.Errorf("got path %s, want /pokemon/pikachu", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": 25,
				"name": "pikachu",
				"base_experience": 112,
				"height": 4,
				"weight": 60,
				"stats": [
					{"base_stat": 35, "stat": {"name": "hp"}},
					{"base_stat": 55, "stat": {"name": "attack"}}
				],
				"types": [
					{"type": {"name": "electric"}}
				]
			}`))
		}))
		defer server.Close()

		client := NewClient(5*time.Second, 5*time.Minute)
		client.baseURL = server.URL

		// Перший запит — іде на мок-сервер
		pokemon, err := client.GetPokemonInfo("pikachu")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if pokemon.Name != "pikachu" || pokemon.BaseExperience != 112 {
			t.Errorf("unexpected pokemon data: %+v", pokemon)
		}
		if len(pokemon.Stats) != 2 || pokemon.Stats[0].BaseStat != 35 || pokemon.Stats[0].Stat.Name != "hp" {
			t.Errorf("unexpected stats: %+v", pokemon.Stats)
		}
		if len(pokemon.Types) != 1 || pokemon.Types[0].Type.Name != "electric" {
			t.Errorf("unexpected types: %+v", pokemon.Types)
		}
		if calls != 1 {
			t.Errorf("expected 1 server call, got %d", calls)
		}

		// Другий запит — має взятися з кешу без звернення до сервера
		cachedPokemon, err := client.GetPokemonInfo("pikachu")
		if err != nil {
			t.Fatalf("unexpected error on cached call: %v", err)
		}
		if cachedPokemon.ID != 25 {
			t.Errorf("unexpected cached pokemon: %+v", cachedPokemon)
		}
		if calls != 1 {
			t.Errorf("expected cached call not to hit server, got %d calls", calls)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}))
		defer server.Close()

		client := NewClient(5*time.Second, 5*time.Minute)
		client.baseURL = server.URL

		if _, err := client.GetPokemonInfo("missingno"); err == nil {
			t.Fatal("expected error for 404, got nil")
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{broken`))
		}))
		defer server.Close()

		client := NewClient(5*time.Second, 5*time.Minute)
		client.baseURL = server.URL

		if _, err := client.GetPokemonInfo("broken"); err == nil {
			t.Fatal("expected error for invalid json, got nil")
		}
	})
}

func TestGetPokemonInfo_Live(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live test")
	}

	client := NewClient(10*time.Second, 5*time.Minute)
	pokemon, err := client.GetPokemonInfo("magikarp")
	if err != nil {
		t.Fatalf("live request failed: %v", err)
	}
	if pokemon.Name != "magikarp" || pokemon.BaseExperience == 0 || len(pokemon.Stats) == 0 {
		t.Fatalf("unexpected live response: %+v", pokemon)
	}
}
