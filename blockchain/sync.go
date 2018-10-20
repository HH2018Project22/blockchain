package blockchain

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

func CreateBlockSyncHook(endpoint string) BlockHookFunc {
	return func(block *Block) error {

		data, err := json.Marshal(block)
		if err != nil {
			return err
		}

		res, err := http.Post(endpoint, "json/application", bytes.NewBuffer(data))
		if err != nil {
			return err
		}

		if res.StatusCode != http.StatusCreated {
			return errors.New("could not propagate block")
		}

		return nil
	}
}
