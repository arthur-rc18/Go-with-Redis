package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/arthur-rc18/Go-Redis/database"

	geojson "github.com/paulmach/go.geojson"
)

type Blocks struct {
	ID       string           `json:"id,omitempty" `
	Name     string           `json:"name,omitempty" `
	ParentID string           `json:"parentID,omitempty" `
	Centroid geojson.Geometry `json:"centroid,omitempty"`
	Value    float64          `json:"value,omitempty" `
}

func KeysRedis(pattern string) []string {
	db := database.DatabaseConnection()
	result, err := db.Keys(database.Context, pattern).Result()
	if err != nil {
		return nil
	}
	return result
}

func GetBlockId(compositeKey string) string {
	return strings.Split(compositeKey, ":")[0]
}

func UpdatedBlockId(key, parentKey string) string {
	blockKey := GetBlockId(key)
	return blockKey + ":" + parentKey
}

func (b Blocks) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Blocks) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &b)
}

func GetBlocksData() ([]Blocks, error) {

	redisClient := database.ConnectRedis()
	defer redisClient.Close()

	blocksList := []Blocks{}

	keys, err := redisClient.Keys(database.Context, "*").Result()
	if err != nil {
		log.Println(err)
		return []Blocks{}, err
	}

	var block Blocks
	for _, result := range keys {

		value, err := redisClient.Get(database.Context, result).Result()
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

	db := database.DatabaseConnection()

	blockKey := KeysRedis(key + ":*")
	if len(blockKey) != 1 {
		return Blocks{}, NonExistantBlock
	}
	result, err := db.Get(database.Context, blockKey[0]).Result()
	if err != nil {
		return Blocks{}, err
	}

	var block Blocks

	if err := block.UnmarshalBinary([]byte(result)); err != nil {
		fmt.Println(err.Error(), err)
		return Blocks{}, err
	}
	fmt.Println(block)
	return block, nil

}

func DeleteBlockByID(key string) error {

	db := database.DatabaseConnection()
	checkBlockKey := KeysRedis(key + ":*")
	if len(checkBlockKey) != 1 {
		return NonExistantBlock
	}
	blockKey := checkBlockKey[0]
	childrenKeys := KeysRedis("*:" + GetBlockId(blockKey))
	if len(childrenKeys) == 0 {
		err := db.Del(database.Context, blockKey).Err()
		return err
	}

	childrenBlocks, err := children(childrenKeys)
	if err != nil {
		return err
	}

	block, err := GetBlockByID(GetBlockId(blockKey))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	for _, childBlock := range childrenBlocks {
		childBlock.ID = UpdatedBlockId(childBlock.ID, block.ParentID)
		childBlock.ParentID = block.ParentID
		err := CreateBlock(childBlock)
		if err != nil {
			return err
		}
	}
	keysToDelete := append(childrenKeys, block.ID)
	err = db.Del(database.Context, keysToDelete...).Err()

	return err

}

func CreateBlock(block Blocks) error {
	if existentKeys := KeysRedis(block.ID); len(existentKeys) != 0 {
		return ErrBlockExisted
	}
	return setBlock(block)
}

func UpdateBlockByID(key string, block Blocks) error {

	if existentKeys := KeysRedis(key + ":*"); len(existentKeys) != 1 {
		return NonExistantBlock
	}
	return setBlock(block)

}

func setBlock(block Blocks) error {
	if block.ParentID != "0" {
		parentBlock, _ := GetBlockByID(block.ParentID)
		if reflect.DeepEqual(parentBlock, Blocks{}) {
			return InvalidParent
		}
	}

	db := database.DatabaseConnection()
	err := db.Set(database.Context, block.ID, block, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func children(childrenKeys []string) ([]Blocks, error) {
	db := database.DatabaseConnection()
	result, err := db.MGet(database.Context, childrenKeys...).Result()
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

var (
	ErrBlockExisted  = errors.New("already existed key")
	InvalidParent    = errors.New("error with parent id")
	NonExistantBlock = errors.New("nonExistant key on database")
)
