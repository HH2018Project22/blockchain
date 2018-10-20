package blockchain

import (
	"bytes"
	"encoding/json"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/boltdb/bolt"
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

func (bc *Blockchain) AddEvent(event Event) error {
	if err := event.Validate(bc); err != nil {
		return err
	}
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(event, prevBlock.Hash)
	return bc.add(newBlock)
}

func (bc *Blockchain) AddBlock(block *Block) error {
	if err := block.Validate(bc); err != nil {
		return err
	}
	return bc.add(block)
}

func (bc *Blockchain) add(block *Block) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	for _, b := range bc.blocks {
		if bytes.Compare(b.PrevBlockHash, block.PrevBlockHash) == 0 {
			return errors.New("parent already used")
		}
	}
	if bc.beforeAddBlockHook != nil {
		if err := bc.beforeAddBlockHook(block); err != nil {
			return err
		}
	}
	bc.blocks = append(bc.blocks, block)
	return nil
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

func (bc *Blockchain) FindPrescription(prescriptionHash []byte) *Prescription {
	for _, b := range bc.blocks {
		if b.Event.Type() == PrescriptionEventType {
			pe := b.Event.(*PrescriptionEvent)
			if bytes.Compare(pe.Prescription.Hash(), prescriptionHash) == 0 {
				return pe.Prescription
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
