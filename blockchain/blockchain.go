package blockchain

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/btcsuite/btcutil/base58"
)

var (
	blocksBucket = []byte("blocks")
)

type BlockHookFunc func(block *Block) error

type Blockchain struct {
	blocks             []*Block
	mutex              sync.RWMutex
	beforeAddBlockHook BlockHookFunc
}

func (bc *Blockchain) Blocks() []*Block {
	return bc.blocks
}

func (bc *Blockchain) AddEvent(event Event) (*Block, error) {
	if err := event.Validate(bc); err != nil {
		return nil, err
	}
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(event, prevBlock.Hash)
	return bc.add(newBlock)
}

func (bc *Blockchain) AddBlock(block *Block) error {
	if err := block.Validate(bc); err != nil {
		return err
	}
	_, err := bc.add(block)
	return err
}

func (bc *Blockchain) add(block *Block) (*Block, error) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	for _, b := range bc.blocks {
		if bytes.Compare(b.PrevBlockHash, block.PrevBlockHash) == 0 {
			return nil, errors.New("parent already used")
		}
	}
	if bc.beforeAddBlockHook != nil {
		if err := bc.beforeAddBlockHook(block); err != nil {
			return nil, err
		}
	}
	bc.blocks = append(bc.blocks, block)
	return block, nil
}

func (bc *Blockchain) ListPrescriptions() []*HashedPrescription {
	hashedPrescriptions := make([]*HashedPrescription, 0)
	for _, b := range bc.blocks {
		if b.Event.Type() == PrescriptionEventType {
			pe := b.Event.(*PrescriptionEvent)
			hashedPrescription := &HashedPrescription{
				Hash:         base58.Encode(b.Hash),
				Prescription: pe.Prescription,
			}
			hashedPrescriptions = append(hashedPrescriptions, hashedPrescription)
		}
	}
	return hashedPrescriptions
}

func (bc *Blockchain) FindPrescriptionNotificationEvents(prescriptionHash []byte) []*NotificationEvent {
	events := make([]*NotificationEvent, 0)
	for _, b := range bc.blocks {
		if b.Event.Type() != NotificationEventType {
			continue
		}
		notificationEvent := b.Event.(*NotificationEvent)
		if bytes.Compare(prescriptionHash, notificationEvent.PrescriptionHash) != 0 {
			continue
		}
		events = append(events, notificationEvent)
	}
	return events
}

func (bc *Blockchain) FindPrescriptionBlock(prescriptionHash []byte) *Block {
	for _, b := range bc.blocks {
		if b.Event.Type() == PrescriptionEventType {
			if bytes.Compare(b.Hash, prescriptionHash) == 0 {
				return b
			}
		}
	}
	return nil
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

func LoadBlockchain(path string, beforeAddBlockHook BlockHookFunc) (*Blockchain, error) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	db, err := bolt.Open(path, 0666, &bolt.Options{
		Timeout: 1 * time.Second,
	})
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

	bc := &Blockchain{
		blocks:             blocks,
		beforeAddBlockHook: beforeAddBlockHook,
	}

	bc.sortBlocks()

	return bc, nil

}

func (bc *Blockchain) sortBlocks() {
	sort.Slice(bc.blocks, func(i, j int) bool {
		a := bc.blocks[i]
		b := bc.blocks[j]
		return a.Timestamp < b.Timestamp
	})
}

func NewBlockchain(beforeAddBlockHook BlockHookFunc) *Blockchain {
	return &Blockchain{
		blocks: []*Block{
			NewGenesisBlock(),
		},
		beforeAddBlockHook: beforeAddBlockHook,
	}
}
