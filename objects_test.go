package columnar

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// BenchmarkCollection/add-8         	 3165742	       344.1 ns/op	     430 B/op	       0 allocs/op
// BenchmarkCollection/fetch-to-8    	 2554174	       458.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkCollection/where-8       	 1402629	       858.4 ns/op	     336 B/op	      13 allocs/op
func BenchmarkCollection(b *testing.B) {
	players := loadPlayers()
	obj := Object{
		"name":   "Roman",
		"age":    35,
		"wallet": 50.99,
		"health": 100,
		"mana":   200,
	}

	b.Run("add", func(b *testing.B) {
		col := New()
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			col.Add(obj)
		}
	})

	b.Run("fetch-to", func(b *testing.B) {
		dst := make(Object, 8)
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			players.FetchTo(20, &dst)
		}
	})

	b.Run("where", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			players.Where(func(v interface{}) bool {
				return v.(string) == "human"
			}, "race")
		}
	})
}

func TestCollection(t *testing.T) {
	obj := Object{
		"name":   "Roman",
		"age":    35,
		"wallet": 50.99,
		"health": 100,
		"mana":   200,
	}

	col := New()
	idx := col.Add(obj)

	{ // Find the object by its index
		obj, ok := col.Fetch(idx)
		assert.True(t, ok)
		assert.Equal(t, "Roman", obj["name"])
	}

	{ // Remove the object
		col.Remove(idx)
		obj, ok := col.Fetch(idx)
		assert.False(t, ok)
		assert.Nil(t, obj)
	}

	{ // Add a new one, should replace
		idx := col.Add(obj)
		obj, ok := col.Fetch(idx)
		assert.Equal(t, uint32(0), idx)
		assert.True(t, ok)
		assert.Equal(t, "Roman", obj["name"])
	}
}

// loadPlayers loads a list of players from the fixture
func loadPlayers() *Collection {
	b, err := os.ReadFile("fixtures/players.json")
	if err != nil {
		panic(err)
	}

	var players []Object
	if err := json.Unmarshal(b, &players); err != nil {
		panic(err)
	}

	out := New()
	for _, p := range players {
		out.Add(p)
	}
	return out
}