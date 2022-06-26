package models

import (
	"reflect"

	"github.com/arthur-rc18/Go-Redis/utils"
)

type Tree struct {
	Block    Blocks `json:"block,omitempty"`
	Children []Tree `json:"children,omitempty"`
}

func GetTreeID(id string) Tree {

	var tree Tree
	var blockId string
	var keysChildren []string
	if id == "0" {
		keysChildren = utils.GetKeys("*:" + id)
		tree.Block = Blocks{}
	} else {
		tree.Block, _ = GetBlockByID(id)
		if reflect.DeepEqual(tree.Block, Blocks{}) {
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
