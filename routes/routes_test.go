package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/arthur-rc18/Go-Redis/database"
	"github.com/arthur-rc18/Go-Redis/handlers"
	"github.com/arthur-rc18/Go-Redis/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var ID int

func RoutesSetup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	routes := gin.Default()
	return routes
}

func CreateBlockMock() {

	blockMock := "{\"id\":0,\"name\":\"block_test\",\"parentID\":0,\"centroid\":{\"type\":\"Point\",\"coordinates\":[-50.404319105,-17.8070625]},\"value\":1811.61}"

	redisClient := database.ConnectRedis()
	defer redisClient.Close()

	redisClient.Set(database.CTX, "block_test", blockMock, 0)

}

func DeleteBlockMock() {

	redisCLient := database.ConnectRedis()
	defer redisCLient.Close()

	redisCLient.Del(database.CTX, "block_test")

}

func TestGetAllBlocks(t *testing.T) {

	database.ConnectRedis()

	CreateBlockMock()
	defer DeleteBlockMock()
	t.Parallel()

	routes := RoutesSetup()
	routes.GET("/blocks", handlers.GetBlocksData)
	req, _ := http.NewRequest("GET", "/blocks", nil)
	res := httptest.NewRecorder()
	routes.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

}

func TestGetBlockByID(t *testing.T) {

	database.ConnectRedis()
	CreateBlockMock()
	defer DeleteBlockMock()
	t.Parallel()

	routes := RoutesSetup()
	routes.GET("/blocks/:id", handlers.GetBlockByID)
	path := "/blocks/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("GET", path, nil)
	res := httptest.NewRecorder()
	routes.ServeHTTP(res, req)

	var blocks models.Blocks
	json.Unmarshal(res.Body.Bytes(), &blocks)
	fmt.Println(&blocks)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestTreeID(t *testing.T) {

	database.ConnectRedis()
	CreateBlockMock()
	defer DeleteBlockMock()
	t.Parallel()

	routes := RoutesSetup()
	routes.GET("/blocks/:id", handlers.GetTreeID)
	path := "/blocks/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("GET", path, nil)
	res := httptest.NewRecorder()
	routes.ServeHTTP(res, req)

	var blocks models.Blocks
	json.Unmarshal(res.Body.Bytes(), &blocks)
	fmt.Println(&blocks)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestDeleteBlockByID(t *testing.T) {

	database.ConnectRedis()
	CreateBlockMock()
	t.Parallel()

	routes := RoutesSetup()
	routes.DELETE("/blocks/:id", handlers.DeleteBlockByID)
	path := "/blocks/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("DELETE", path, nil)
	res := httptest.NewRecorder()
	routes.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestUpdateBlockByID(t *testing.T) {

}
