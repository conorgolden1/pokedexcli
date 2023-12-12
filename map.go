package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"internal/pokecache"
)

type Map struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (m *Map) nextMap(c *pokecache.Cache) (Map, error) {
	return getMap(m.Next, c)
}

func (m *Map) prevMap(c *pokecache.Cache) (Map, error) {
	str, ok := m.Previous.(string)
	if ok {
		return getMap(str, c)
	}
	return Map{}, errors.New(fmt.Sprintf("Previous is invalid: %v", str))
}

func (m *Map) Print() {
	for _, result := range m.Results {
		fmt.Println(result.Name)
	}
}

func getMap(url string, c *pokecache.Cache) (Map, error) {
	var j []byte
	var err error
	if data, ok := c.Get(url); ok {
		return mapUnmarshal(data)
	}
	if url == "" {
		j, err = get("https://pokeapi.co/api/v2/location-area/")
	} else {
		j, err = get(url)
	}
	if err != nil {
		return Map{}, err
	}
	c.Add(url, j)
	return mapUnmarshal(j)
}

func mapUnmarshal(data []byte) (Map, error) {
	mapObject := Map{}
	err := json.Unmarshal(data, &mapObject)
	if err != nil {
		fmt.Println(mapObject)
		return Map{}, err
	}
	return mapObject, nil
}
