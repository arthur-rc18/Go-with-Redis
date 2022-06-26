package models

import (
	"testing"

	"github.com/arthur-rc18/Go-Redis/database"
	"github.com/stretchr/testify/assert"

	geojson "github.com/paulmach/go.geojson"
)

var blockMock = Blocks{
	ID:       "C3:0",
	Name:     "Block test",
	ParentID: "0",
	Centroid: *geojson.NewPointGeometry([]float64{-12.50830530855398 - 52.64695717817407}),
}

func mockBlock() {
	db := database.ConnectWithDB()
	db.Set(database.CTX, blockMock.ID, blockMock, 0)
}

func UnmockBlock() {

	db := database.ConnectWithDB()
	db.FlushAll(database.CTX)
}

func TestGetAllBlocks(t *testing.T) {

	t.Parallel()
	t.Run("Test getting all blocks", func(t *testing.T) {
		mockBlock()
		defer UnmockBlock()
		got, _ := GetBlocksData()
		assert.Equal(t, []Blocks{blockMock}, got)
	})
}

func TestGetBlockByID(t *testing.T) {
	t.Parallel()
	t.Run("Test getting an existed block", func(t *testing.T) {
		mockBlock()
		defer UnmockBlock()
		got, _ := GetBlockByID("C3")
		assert.Equal(t, blockMock, got)
	})

	t.Run("Test getting inexistent block", func(t *testing.T) {
		got, _ := GetBlockByID("C3")
		assert.Equal(t, Blocks{}, got)

	})
}

func TestCreateBlock(t *testing.T) {
	t.Parallel()
	t.Run("insert existent key", func(t *testing.T) {
		mockBlock()
		defer UnmockBlock()

		err := CreateBlock(blockMock)
		assert.Error(t, err)
	})

	t.Run("insert new block", func(t *testing.T) {
		UnmockBlock()
		err := CreateBlock(blockMock)
		if err != nil {
			t.Error(err)
		}
		gotBlock, _ := GetBlockByID("C3")
		assert.Equal(t, blockMock, gotBlock)
	})
}

func TestUpdateBlock(t *testing.T) {

	updatedBlock := blockMock
	t.Parallel()
	t.Run("valid update block", func(t *testing.T) {
		mockBlock()
		defer UnmockBlock()
		err := UpdateBlockByID("C3", updatedBlock)
		if err != nil {
			t.Error(err)
		}
		gotBlock, _ := GetBlockByID("C3")
		assert.Equal(t, gotBlock, updatedBlock)
	})

	t.Run("Invalid block", func(t *testing.T) {
		UnmockBlock()
		err := UpdateBlockByID("C3", updatedBlock)
		assert.Error(t, err)
	})

}

func TestDeleteBlock(t *testing.T) {
	t.Run("existent block", func(t *testing.T) {
		mockBlock()
		defer UnmockBlock()
		err := DeleteBlockByID("C3")
		if err != nil {
			t.Error(err)
		}
		gotBlock, _ := GetBlockByID("C3")
		assert.Equal(t, Blocks{}, gotBlock)
	})

	t.Run("nonexistent block", func(t *testing.T) {
		UnmockBlock()
		err := DeleteBlockByID("C3")
		assert.Error(t, err)
	})
}
