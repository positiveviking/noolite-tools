package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/akamensky/argparse"
	"github.com/venushome/hands.noolite/connection"
	"github.com/venushome/hands.noolite/engine"
)

type argT struct {
	NatsURL   string
	System    string
	Subsystem string
}

func run(args argT) error {
	eng := engine.InitEngine()
	con := connection.NewNatsConnector(args.NatsURL)
	con.EngineInput = eng.Input
	con.EngineOutput = eng.Output
	con.System = args.System
	con.Subsystem = args.Subsystem

	con.Run()
	eng.Run()

	fmt.Println("Hands.noolite service started")
	fmt.Println("Press ctrl+c for exit")

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	s := <-sigc
	fmt.Printf("\nGot signal %+v\n", s)
	con.Stop()
	eng.Stop()
	fmt.Println("Hands.noolite service stopped")

	return nil
}

func main() {
	parser := argparse.NewParser("head", "Venushome head server")
	natsURL := parser.String("n", "nats", &argparse.Options{
		Required: true,
		Help:     "url of nats server",
	})
	system := parser.String("s", "system", &argparse.Options{
		Required: false,
		Help:     "system identifier",
		Default:  "hands",
	})
	subsystem := parser.String("u", "subsystem", &argparse.Options{
		Required: false,
		Help:     "subsystem identifier",
		Default:  "noolite",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatalf("Argument error: %+v", err)
	}
	run(argT{
		NatsURL:   *natsURL,
		System:    *system,
		Subsystem: *subsystem,
	})
}
