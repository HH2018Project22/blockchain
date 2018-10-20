package blockchain

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/boltdb/bolt"
)

var (
	blocksBucket = []byte("blocks")
)

type Blockchain struct {
	blocks []*Block
}

func (bc *Blockchain) Blocks() []*Block {
	return bc.blocks
}

func (bc *Blockchain) AddEvent(event Event) bool {
	if !event.Validate(bc) {
		return false
	}
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(event, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
	return true
}

func (bc *Blockchain) AddBlock(block *Block) bool {
	if !block.Validate(bc) {
		return false
	}
	bc.blocks = append(bc.blocks, block)
	return true
}

func (bc *Blockchain) ListPrescriptions() []*Prescription {
	prescriptions := make([]*Prescription, 0)
	for _, b := range bc.blocks {
		if b.Event.Type() == PrescriptionEventType {
			pe := b.Event.(*PrescriptionEvent)
			prescriptions = append(prescriptions, pe.Prescription)
		}
	}
	return prescriptions
}

func (bc *Blockchain) Save(path string) error {
	db, err := bolt.Open(path, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(blocksBucket)
		if err != nil {
			return err
		}
		for _, b := range bc.blocks {
			data, err := json.Marshal(b)
			if err != nil {
				return err
			}
			if err := bucket.Put(b.Hash, data); err != nil {
				return err
			}
		}
		return nil
	})
}

func LoadBlockchain(path string) (*Blockchain, error) {

	db, err := bolt.Open(path, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	defer db.Close()

	blocks := []*Block{}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(blocksBucket)
		if bucket == nil {
			return errors.New("unknown bucket")
		}
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			b := &Block{}
			if err := json.Unmarshal(v, b); err != nil {
				return err
			}
			blocks = append(blocks, b)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Blockchain{blocks}, nil
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		blocks: []*Block{
			NewGenesisBlock(),
		},
	}
}
