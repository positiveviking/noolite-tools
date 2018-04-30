package engine

import (
	"log"
	"sync"
	"time"

	"github.com/venushome/noolite/rx"
	"github.com/venushome/noolite/tx"
)

type Engine struct {
	Output chan Event
	Input  chan Command

	rxEngine *rx.RxEngine
	txEngine *tx.TxEngine

	stopped bool
	done    sync.WaitGroup
}

type Command struct {
	Channel uint
	Action  string
	Value   []int
}

type Event struct {
	Channel uint
	Action  string
	Value   []int
}

func InitEngine() *Engine {
	e := Engine{
		Output: make(chan Event),
		Input:  make(chan Command),
	}
	e.rxEngine, _ = rx.NewRxEngine()
	e.txEngine, _ = tx.NewTxEngine()

	return &e
}

func (this *Engine) Run() {
	this.done.Add(2)
	go this.recoverable(this.listenRxCycle)
	go this.recoverable(this.writeTxCicle)
}

func (this *Engine) Stop() {
	this.stopped = true
	this.done.Wait()

	this.rxEngine.Close()
	this.txEngine.Close()

	this.rxEngine.Exit()
	this.txEngine.Exit()
}

func (this *Engine) listenRxCycle() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Rx cycle recover: %v", r)
			return
		}
	}()

	if err := this.rxEngine.Open(); err != nil {
		return
	}

	for {
		response, err := this.rxEngine.Read(time.Second)
		if err != nil {
			this.rxEngine.Close()
			return
		}

		if this.stopped {
			return
		}

		if len(response) == 0 {
			continue
		}

		ev := Event{
			Channel: uint(response.Channel()),
			Action:  rxActions[response.Command()],
		}
		if response.DataLen() > 0 {
			for _, b := range response.Data() {
				ev.Value = append(ev.Value, int(b))
			}
		}

		this.Output <- ev
	}
}

func (this *Engine) writeTxCicle() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Rx cycle recover: %v", r)
			return
		}
	}()

	if err := this.txEngine.Open(); err != nil {
		return
	}

	for {
		if this.stopped {
			return
		}

		select {
		case msg, ok := <-this.Input:
			if !ok {
				return
			}
			cmdType, ok := txActions[msg.Action]
			if !ok {
				continue
			}
			cmd := tx.Command{
				Channel: byte(msg.Channel),
				Type:    cmdType,
			}

			switch len(msg.Value) {
			case 1:
				cmd.SetRGB(byte(msg.Value[0]), byte(msg.Value[0]), byte(msg.Value[0]))
			case 3:
				cmd.SetRGB(byte(msg.Value[0]), byte(msg.Value[2]), byte(msg.Value[3]))
			}

			if err := this.txEngine.Write(cmd); err != nil {
				this.txEngine.Close()
				return
			}

		case <-time.After(time.Second):
			continue
		}
	}

}

func (this *Engine) recoverable(f func()) {
	defer this.done.Done()

	firstStart := true
	for {
		if this.stopped {
			return
		}
		if !firstStart {
			time.Sleep(time.Second * 3)
		}
		f()
		firstStart = false
	}
}
