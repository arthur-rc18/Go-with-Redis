package handlers

import (
	"log"
	"net/http"
	"reflect"

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

	idToDelete := ctx.Param("id")
	err := models.DeleteBlockByID(idToDelete)
	if err == models.NonExistantBlock {
		ctx.JSON(http.StatusNotFound, nil)
		return
	} else if err == models.InvalidParent {
		ctx.JSON(http.StatusBadRequest, err)
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.JSON(http.StatusOK, nil)

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
	if err == models.NonExistantBlock {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	} else if err == models.InvalidParent {
		ctx.JSON(http.StatusBadRequest, gin.H{
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
	if err == models.ErrBlockExisted || err == models.InvalidParent {
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
