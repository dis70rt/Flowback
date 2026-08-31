package seed

import (
	"log"

	"github.com/razorpay/razorpay-go"
)

type Task interface {
	Name() string
	Execute(client *razorpay.Client, keyID string)
}

type Engine struct {
	Client *razorpay.Client
	KeyID  string
}

func New(keyID, secret string) *Engine {
	return &Engine{
		Client: razorpay.NewClient(keyID, secret),
		KeyID:  keyID,
	}
}

func (e *Engine) Execute(task Task) {
	log.Printf("\n========================================\n")
	log.Printf("[SEEDER] Starting Task: %s\n", task.Name())
	log.Printf("========================================\n")
	
	task.Execute(e.Client, e.KeyID)
	
	log.Printf("[SEEDER] Task '%s' Complete!\n", task.Name())
	log.Printf("========================================\n\n")
}
