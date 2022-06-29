package models

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arthur-rc18/Go-Redis/database"
	geojson "github.com/paulmach/go.geojson"
)

var (
	c0 = Blocks{
		ID:       "C0:0",
		Name:     "Cliente A",
		ParentID: "0",
		Centroid: *geojson.NewPointGeometry([]float64{-14.745636148936388, -63.8203580290067}),
		Value:    25000,
	}
	f1 = Blocks{
		ID:       "F1:C0",
		Name:     "FAZENDA 1",
		ParentID: "C0",
		Centroid: *geojson.NewPointGeometry([]float64{-7.946470513501444, -74.68622424624228}),
		Value:    2000,
	}
)

var treeMock = Tree{
	Block: c0,
	Children: []Tree{
		{
			Block:    f1,
			Children: nil,
		},
	},
}

func MockTree(t *testing.T) {
	UnmockTree(t)
	db := database.DatabaseConnection()
	blocks := []Blocks{c0, f1}
	for _, block := range blocks {
		err := db.Set(database.Context, block.ID, block, 0).Err()
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func UnmockTree(t *testing.T) {
	db := database.DatabaseConnection()
	err := db.FlushAll(database.Context).Err()
	if err != nil {
		t.Error(err)
	}
}

func TestGetTreeById(t *testing.T) {
	t.Run("mocking tree test", func(t *testing.T) {
		MockTree(t)
		defer UnmockTree(t)

		got := GetTreeID("C0")

		assert.Equal(t, treeMock, got)
	})
	t.Run("key not presented in database test", func(t *testing.T) {
		got := GetTreeID("C0")
		assert.Equal(t, Tree{}, got)
		assert.NotEqual(t, treeMock, got)
	})
}
