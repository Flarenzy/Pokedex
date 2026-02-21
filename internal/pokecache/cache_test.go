package pokecache

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

const first_twenty_resp = `{"count":1089,"next":"https://pokeapi.co/api/v2/location-area/?offset=20&limit=20","previous":null,"results":[{"name":"canalave-city-area","url":"https://pokeapi.co/api/v2/location-area/1/"},{"name":"eterna-city-area","url":"https://pokeapi.co/api/v2/location-area/2/"},{"name":"pastoria-city-area","url":"https://pokeapi.co/api/v2/location-area/3/"},{"name":"sunyshore-city-area","url":"https://pokeapi.co/api/v2/location-area/4/"},{"name":"sinnoh-pokemon-league-area","url":"https://pokeapi.co/api/v2/location-area/5/"},{"name":"oreburgh-mine-1f","url":"https://pokeapi.co/api/v2/location-area/6/"},{"name":"oreburgh-mine-b1f","url":"https://pokeapi.co/api/v2/location-area/7/"},{"name":"valley-windworks-area","url":"https://pokeapi.co/api/v2/location-area/8/"},{"name":"eterna-forest-area","url":"https://pokeapi.co/api/v2/location-area/9/"},{"name":"fuego-ironworks-area","url":"https://pokeapi.co/api/v2/location-area/10/"},{"name":"mt-coronet-1f-route-207","url":"https://pokeapi.co/api/v2/location-area/11/"},{"name":"mt-coronet-2f","url":"https://pokeapi.co/api/v2/location-area/12/"},{"name":"mt-coronet-3f","url":"https://pokeapi.co/api/v2/location-area/13/"},{"name":"mt-coronet-exterior-snowfall","url":"https://pokeapi.co/api/v2/location-area/14/"},{"name":"mt-coronet-exterior-blizzard","url":"https://pokeapi.co/api/v2/location-area/15/"},{"name":"mt-coronet-4f","url":"https://pokeapi.co/api/v2/location-area/16/"},{"name":"mt-coronet-4f-small-room","url":"https://pokeapi.co/api/v2/location-area/17/"},{"name":"mt-coronet-5f","url":"https://pokeapi.co/api/v2/location-area/18/"},{"name":"mt-coronet-6f","url":"https://pokeapi.co/api/v2/location-area/19/"},{"name":"mt-coronet-1f-from-exterior","url":"https://pokeapi.co/api/v2/location-area/20/"}]}"`
const second_twenty_resp = `{"count":1089,"next":"https://pokeapi.co/api/v2/location-area/?offset=40&limit=20","previous":"https://pokeapi.co/api/v2/location-area/?offset=0&limit=20","results":[{"name":"mt-coronet-1f-route-216","url":"https://pokeapi.co/api/v2/location-area/21/"},{"name":"mt-coronet-1f-route-211","url":"https://pokeapi.co/api/v2/location-area/22/"},{"name":"mt-coronet-b1f","url":"https://pokeapi.co/api/v2/location-area/23/"},{"name":"great-marsh-area-1","url":"https://pokeapi.co/api/v2/location-area/24/"},{"name":"great-marsh-area-2","url":"https://pokeapi.co/api/v2/location-area/25/"},{"name":"great-marsh-area-3","url":"https://pokeapi.co/api/v2/location-area/26/"},{"name":"great-marsh-area-4","url":"https://pokeapi.co/api/v2/location-area/27/"},{"name":"great-marsh-area-5","url":"https://pokeapi.co/api/v2/location-area/28/"},{"name":"great-marsh-area-6","url":"https://pokeapi.co/api/v2/location-area/29/"},{"name":"solaceon-ruins-2f","url":"https://pokeapi.co/api/v2/location-area/30/"},{"name":"solaceon-ruins-1f","url":"https://pokeapi.co/api/v2/location-area/31/"},{"name":"solaceon-ruins-b1f-a","url":"https://pokeapi.co/api/v2/location-area/32/"},{"name":"solaceon-ruins-b1f-b","url":"https://pokeapi.co/api/v2/location-area/33/"},{"name":"solaceon-ruins-b1f-c","url":"https://pokeapi.co/api/v2/location-area/34/"},{"name":"solaceon-ruins-b2f-a","url":"https://pokeapi.co/api/v2/location-area/35/"},{"name":"solaceon-ruins-b2f-b","url":"https://pokeapi.co/api/v2/location-area/36/"},{"name":"solaceon-ruins-b2f-c","url":"https://pokeapi.co/api/v2/location-area/37/"},{"name":"solaceon-ruins-b3f-a","url":"https://pokeapi.co/api/v2/location-area/38/"},{"name":"solaceon-ruins-b3f-b","url":"https://pokeapi.co/api/v2/location-area/39/"},{"name":"solaceon-ruins-b3f-c","url":"https://pokeapi.co/api/v2/location-area/40/"}]}`

func TestCache(t *testing.T) {
	t.Parallel()
	inputs := []struct {
		key string
		val []byte
	}{
		{key: "https://pokeapi.co/api/v2/location-area/", val: []byte(first_twenty_resp)},
		{key: "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20", val: []byte(second_twenty_resp)},
	}
	expected := []struct {
		key string
		val []byte
	}{
		{key: "https://pokeapi.co/api/v2/location-area/", val: []byte(first_twenty_resp)},
		{key: "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20", val: []byte(second_twenty_resp)},
	}
	cache := NewCache(20 * time.Second)
	defer cache.Done()
	for _, input := range inputs {
		err := cache.Add(input.key, input.val)
		if err != nil {
			t.Errorf("error adding key %s: %v", input.key, err)
			t.Fatal()
		}
	}
	for _, exp := range expected {
		val, err := cache.Get(exp.key)
		if err != nil {
			t.Errorf("error getting key %s: %v", exp.key, err)
			t.Fatal()
		}
		if !bytes.Equal(val, exp.val) {
			t.Errorf("key %s: expected %v, got %v", exp.key, string(exp.val), string(val))
			t.Fatal()
		}
	}
}

func TestEmptyCache(t *testing.T) {
	t.Parallel()
	cache := NewCache(20 * time.Second)
	v, err := cache.Get("Random Key")
	if err == nil {
		t.Error("expected error, got none")
		t.Fatal()
	}
	if v != nil {
		t.Errorf("expected nil, got %v", v)
		t.Fatal()
	}
	if !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("expected ErrKeyNotFound, got %v", err)
		t.Fatal()
	}
	if !strings.Contains(err.Error(), "Random Key") {
		t.Errorf("expected key in error, got %s", err.Error())
		t.Fatal()
	}
}

func TestDoubleEntry(t *testing.T) {
	t.Parallel()
	cache := NewCache(20 * time.Second)
	defer cache.Done()
	inputs := []struct {
		key string
		val []byte
	}{
		{key: "https://pokeapi.co/api/v2/location-area/", val: []byte(first_twenty_resp)},
		{key: "https://pokeapi.co/api/v2/location-area/", val: []byte(first_twenty_resp)},
	}
	err := cache.Add(inputs[0].key, inputs[0].val)
	if err != nil {
		t.Errorf("error adding key %s: %v", inputs[0].key, err)
		t.Fatal()
	}
	err1 := cache.Add(inputs[1].key, inputs[1].val)
	if err1 != nil {
		if !errors.Is(err1, ErrKeyExists) {
			t.Errorf("expected ErrKeyExists, got %v", err1)
			t.Fatal()
		}
	}
}

func TestReapLoop(t *testing.T) {
	t.Parallel()
	cache := NewCache(3 * time.Second)
	defer cache.Done()
	inputs := []struct {
		key string
		val []byte
	}{
		{key: "https://pokeapi.co/api/v2/location-area/", val: []byte(first_twenty_resp)},
		{key: "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20", val: []byte(second_twenty_resp)},
	}
	for _, input := range inputs {
		err := cache.Add(input.key, input.val)
		if err != nil {
			t.Errorf("error adding key %s: %v", input.key, err)
			t.Fatal()
		}
	}
	_, err := cache.Get(inputs[0].key)
	if err != nil {
		t.Errorf("error getting key %s: %v", inputs[0].key, err)
		t.Fatal()
	}
	time.Sleep(4 * time.Second)
	val, err1 := cache.Get(inputs[0].key)
	if err1 != nil {
		if !errors.Is(err1, ErrKeyNotFound) {
			t.Errorf("expected ErrKeyNotFound, got %v", err1)
			t.Fatal()
		}
	}
	if val != nil {
		t.Errorf("key %s: expected nil, got %v", inputs[0].key, string(val))
		t.Fatal()
	}
}
