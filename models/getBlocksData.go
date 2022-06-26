package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/arthur-rc18/Go-Redis/database"
	"github.com/arthur-rc18/Go-Redis/utils"

	geojson "github.com/paulmach/go.geojson"
)

const defaultPattern string = "*:*"

var (
	ErrBlockAlreadyExists = errors.New("this key already exists on database")
	ErrInvalidParentId    = errors.New("invalid parent id or parent doesn't exists")
	ErrBlockNotExists     = errors.New("key not exists on database")
)

type Blocks struct {
	ID       string           `json:"id,omitempty" `
	Name     string           `json:"name,omitempty" `
	ParentID string           `json:"parentID,omitempty" `
	Centroid geojson.Geometry `json:"centroid,omitempty"`
	Value    float64          `json:"value,omitempty" `
}

func (b Blocks) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

func (b Blocks) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &b)
}

func GetBlocksData() ([]Blocks, error) {

	redisClient := database.ConnectRedis()
	defer redisClient.Close()

	blocksList := []Blocks{}

	keys, err := redisClient.Keys(database.CTX, "*").Result()
	if err != nil {
		log.Println(err)
		return []Blocks{}, err
	}

	var block Blocks
	for _, result := range keys {

		value, err := redisClient.Get(database.CTX, result).Result()
		if err != nil {
			log.Println(err.Error())
			return blocksList, err
		}
		fmt.Println(value)
		json.Unmarshal([]byte(value), &block)
		fmt.Println(block)

		blocksList = append(blocksList, block)
	}

	return blocksList, nil

}

func GetBlockByID(key string) (Blocks, error) {

	db := database.ConnectWithDB()

	blockKey := utils.GetKeys(key + ":*")
	if len(blockKey) != 1 {
		return Blocks{}, nil
	}
	result, err := db.Get(database.CTX, blockKey[0]).Result()
	if err != nil {
		return Blocks{}, err
	}
	var block Blocks
	if err := block.UnmarshalBinary([]byte(result)); err != nil {
		fmt.Println(err.Error(), err)
		return Blocks{}, err
	}
	return block, nil

}

func DeleteBlockByID(key string) error {

	db := database.ConnectWithDB()
	checkBlockKey := utils.GetKeys(key + ":*")
	if len(checkBlockKey) != 1 {
		return ErrBlockNotExists
	}
	blockKey := checkBlockKey[0]
	childrenKeys := utils.GetKeys("*:" + utils.GetIndividualBlockId(blockKey))
	if len(childrenKeys) == 0 {
		err := db.Del(database.CTX, blockKey).Err()
		return err
	}

	childrenBlocks, err := getChildren(childrenKeys)
	if err != nil {
		return err
	}

	block, _ := GetBlockByID(utils.GetIndividualBlockId(blockKey))
	for _, childBlock := range childrenBlocks {
		childBlock.ID = utils.UpdatedBlockId(childBlock.ID, block.ParentID)
		childBlock.ParentID = block.ParentID
		err := CreateBlock(childBlock)
		if err != nil {
			return err
		}
	}
	keysToDelete := append(childrenKeys, block.ID)
	err = db.Del(database.CTX, keysToDelete...).Err()

	return err

}

func CreateBlock(block Blocks) error {
	if existentKeys := utils.GetKeys(block.ID); len(existentKeys) != 0 {
		return ErrBlockAlreadyExists
	}
	return setBlock(block)
}

func UpdateBlockByID(key string, block Blocks) error {

	if checkBlockKey := utils.GetKeys(key + ":*"); len(checkBlockKey) != 1 {
		return ErrBlockNotExists
	}
	return setBlock(block)

}

func setBlock(block Blocks) error {
	if block.ParentID != "0" {
		parentBlock, _ := GetBlockByID(block.ParentID)
		if reflect.DeepEqual(parentBlock, Blocks{}) {
			return ErrInvalidParentId
		}
	}

	db := database.ConnectWithDB()
	err := db.Set(database.CTX, block.ID, block, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func getChildren(childrenKeys []string) ([]Blocks, error) {
	db := database.ConnectWithDB()
	result, err := db.MGet(database.CTX, childrenKeys...).Result()
	if err != nil {
		return nil, err
	}

	var childrenBlocks []Blocks
	for _, item := range result {
		var childBlock Blocks
		err := childBlock.UnmarshalBinary([]byte(fmt.Sprint(item)))
		if err != nil {
			return nil, err
		}
		childrenBlocks = append(childrenBlocks, childBlock)
	}
	return childrenBlocks, nil
}
