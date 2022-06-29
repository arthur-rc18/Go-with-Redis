package models

import (
	"reflect"
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
		keysChildren = KeysRedis("*:" + id)
		tree.Block = Blocks{}
	} else {
		tree.Block, _ = GetBlockByID(id)
		if reflect.DeepEqual(tree.Block, Blocks{}) {
			return Tree{}
		}
		blockId = GetBlockId(tree.Block.ID)
		keysChildren = KeysRedis("*:" + blockId)
	}

	for _, keyChild := range keysChildren {
		tree.Children = append(tree.Children, GetTreeID(GetBlockId(keyChild)))
	}
	return tree

}
