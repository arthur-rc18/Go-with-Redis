package data

import (
	"fmt"

	"github.com/arthur-rc18/Go-Redis/database"
	"github.com/arthur-rc18/Go-Redis/models"
	geojson "github.com/paulmach/go.geojson"
)

func UpdateDatabase() {
	db := database.ConnectRedis()

	val := db.Do(database.Context, "FLUSHALL")

	fmt.Println(val)
}

func PopulateDatabase(blocksList []models.Blocks) {
	db := database.DatabaseConnection()
	if blocksList == nil {
		c0 := models.Blocks{
			ID:       "C0:0",
			Name:     "Cliente A",
			ParentID: "0",
			Centroid: *geojson.NewPointGeometry([]float64{37.894910837643636, -91.14305222361646}),
			Value:    10000,
		}
		f1 := models.Blocks{
			ID:       "F1:C0",
			Name:     "FAZENDA 1",
			ParentID: "C0",
			Centroid: *geojson.NewPointGeometry([]float64{0.22316755562051446, -72.51023909659871}),
			Value:    2000,
		}
		f2 := models.Blocks{
			ID:       "F2:C0",
			Name:     "FAZENDA 2",
			ParentID: "C0",
			Centroid: *geojson.NewPointGeometry([]float64{-13.535854466511525, -71.98289532885289}),
			Value:    3000,
		}
		f3 := models.Blocks{
			ID:       "F3:C0",
			Name:     "FAZENDA 3",
			ParentID: "0",
			Centroid: *geojson.NewPointGeometry([]float64{-15.068800064198802, 18.36867021159178}),
			Value:    4000,
		}
		b0 := models.Blocks{
			ID:       "B0:F1",
			Name:     "Bloco 0",
			ParentID: "F1",
			Centroid: *geojson.NewPointGeometry([]float64{23.60741499260964, 19.774920258913873}),
			Value:    100,
		}
		b1 := models.Blocks{
			ID:       "B1:F1",
			Name:     "BLOCK 1",
			ParentID: "F1",
			Centroid: *geojson.NewPointGeometry([]float64{26.15852660555876, 4.48195099428607}),
			Value:    200,
		}
		b2 := models.Blocks{
			ID:       "B2:F2",
			Name:     "BLOCK 2",
			ParentID: "F2",
			Centroid: *geojson.NewPointGeometry([]float64{59.6472786585547, 52.99757762689841}),
			Value:    300,
		}
		b3 := models.Blocks{
			ID:       "B3:F3",
			Name:     "BLOCK 3",
			ParentID: "F3",
			Centroid: *geojson.NewPointGeometry([]float64{61.87849168671919, 74.26710959264513}),
			Value:    400,
		}
		b4 := models.Blocks{
			ID:       "B4:F3",
			Name:     "BLOCK 4",
			ParentID: "F3",
			Centroid: *geojson.NewPointGeometry([]float64{-22.387531372602364, 127.52883013496952}),
			Value:    500,
		}
		b5 := models.Blocks{
			ID:       "B5:F3",
			Name:     "BLOCK 5",
			ParentID: "F3",
			Centroid: *geojson.NewPointGeometry([]float64{-26.54614340576769, 123.83742376074905}),
			Value:    600,
		}

		blocksList = []models.Blocks{c0, f1, f2, f3, b0, b1, b2, b3, b4, b5}
	}

	for _, block := range blocksList {
		err := db.Set(database.Context, block.ID, block, 0).Err()
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
