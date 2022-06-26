package models

import (
	"reflect"

	"github.com/arthur-rc18/Go-Redis/models"
	"github.com/arthur-rc18/Go-Redis/utils"
)

type Tree struct {
	Block    models.Blocks `json:"block,omitempty"`
	Children []Tree        `json:"children,omitempty"`
}

func GetTreeID(id string) Tree {

	var tree Tree
	var blockId string
	var keysChildren []string
	if id == "0" {
		keysChildren = utils.GetKeys("*:" + id)
		tree.Block = models.Blocks{}
	} else {
		tree.Block, _ = models.GetBlockByID(id)
		if reflect.DeepEqual(tree.Block, models.Blocks{}) {
			return Tree{}
		}
		blockId = utils.GetIndividualBlockId(tree.Block.ID)
		keysChildren = utils.GetKeys("*:" + blockId)
	}

	for _, keyChild := range keysChildren {
		tree.Children = append(tree.Children, GetTreeID(utils.GetIndividualBlockId(keyChild)))
	}
	return tree

}
