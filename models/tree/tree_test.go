package models

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arthur-rc18/Go-Redis/database"
	"github.com/arthur-rc18/Go-Redis/models"
	geojson "github.com/paulmach/go.geojson"
)

var (
	c0 = models.Blocks{
		ID:       "C0:0",
		Name:     "Cliente A",
		ParentID: "0",
		Centroid: *geojson.NewPointGeometry([]float64{-48.289546966552734, -18.931050694554795}),
		Value:    10000,
	}
	f1 = models.Blocks{
		ID:       "F1:C0",
		Name:     "FAZENDA 1",
		ParentID: "C0",
		Centroid: *geojson.NewPointGeometry([]float64{-52.9046630859375, -18.132801356084773}),
		Value:    1000,
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
	db := database.ConnectWithDB()
	blocks := []models.Blocks{c0, f1}
	for _, block := range blocks {
		err := db.Set(database.CTX, block.ID, block, 0).Err()
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func UnmockTree(t *testing.T) {
	db := database.ConnectWithDB()
	err := db.FlushAll(database.CTX).Err()
	if err != nil {
		t.Error(err)
	}
}

func TestGetTreeById(t *testing.T) {
	t.Run("mocked tree", func(t *testing.T) {
		MockTree(t)
		defer UnmockTree(t)

		got := GetTreeID("C0")

		assert.Equal(t, treeMock, got)
	})
	t.Run("nonexistent tree", func(t *testing.T) {
		got := GetTreeID("C0")
		assert.Equal(t, Tree{}, got)
		assert.NotEqual(t, treeMock, got)
	})
}
