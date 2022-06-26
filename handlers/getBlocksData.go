package handlers

import (
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/arthur-rc18/Go-Redis/models"
	"github.com/gin-gonic/gin"
)

func GetBlocksData(ctx *gin.Context) {

	blocks, err := models.GetBlocksData()

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"blocks": blocks,
		})
	}
}

func GetBlockByID(ctx *gin.Context) {

	blockId := ctx.Param("id")
	block, err := models.GetBlockByID(blockId)

	if err != nil {
		if reflect.DeepEqual(block, models.Blocks{}) {
			ctx.JSON(http.StatusNotFound, nil)
			return
		} else {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}

	ctx.JSON(http.StatusOK, block)
}

func GetTreeID(ctx *gin.Context) {

	id := ctx.Param("id")
	tree := models.GetTreeID(id)
	if reflect.DeepEqual(tree, models.Tree{}) {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	ctx.JSON(http.StatusOK, tree)

}

func DeleteBlockByID(ctx *gin.Context) {

	id := ctx.Param("id")

	err := models.DeleteBlockByID(id)

	if err != nil {
		if strings.Contains(err.Error(), "Non existant ID") {
			ctx.AbortWithStatus(http.StatusNotFound)
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"block": "deleted with success",
		})
	}

}

func UpdateBlockByID(ctx *gin.Context) {

	id := ctx.Param("id")
	var newBlock models.Blocks

	err := ctx.ShouldBindJSON(&newBlock)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	err = models.UpdateBlockByID(id, newBlock)
	if err == models.ErrBlockNotExists {
		ctx.JSON(http.StatusNotFound, gin.H{
			"data":  newBlock,
			"error": err.Error(),
		})
		return
	} else if err == models.ErrInvalidParentId {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data":  newBlock,
			"error": err.Error(),
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, newBlock)

}

func CreateBlock(ctx *gin.Context) {

	var block models.Blocks
	err := ctx.ShouldBindJSON(&block)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	err = models.CreateBlock(block)
	if err == models.ErrBlockAlreadyExists || err == models.ErrInvalidParentId {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"data":  block,
			"error": err.Error(),
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, block)
}
