package main

import (
	"encoding/json"
	"fmt"
	"internal/pokecache"
)

type Location struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func (l *Location) Print() {
	fmt.Printf("Found pokemon:\n")
	for _, p := range l.PokemonEncounters {
		fmt.Printf("- %v\n", p.Pokemon.Name)
	}
}

func GetLocation(locationname string, c *pokecache.Cache) (Location, error) {
	if data, ok := c.Get(locationname); ok {
		return locationUnmarshal(data)
	}
	baseurl := "https://pokeapi.co/api/v2/location-area/"

	j, err := get(baseurl + locationname)
	if err != nil {
		return Location{}, err
	}
	c.Add(locationname, j)
	return locationUnmarshal(j)
}

func locationUnmarshal(data []byte) (Location, error) {
	location := Location{}
	err := json.Unmarshal(data, &location)
	if err != nil {
		fmt.Println(location)
		return Location{}, err
	}
	return location, nil
}
