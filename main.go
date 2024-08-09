package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AleksBekker/project-api/database"
	_ "github.com/joho/godotenv/autoload"
)

func run(ctx context.Context, args []string, lookupEnv func(string) (string, bool), errStream io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := log.New(errStream, "", log.LstdFlags | log.Lshortfile)

	clargs, err := parseArgv(args)
	if err != nil {
		log.Println("invalid command-line arguments")
		return err
	}

	db, err := db.FromEnv(lookupEnv)
	if err != nil {
		return err
	}
	defer db.Close()

	errorChannel := make(chan error)
	server := NewServer(clargs.addr, logger, db)
	go func(l *log.Logger) {
		var err error
		if err = server.Run(); err != nil && errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		errorChannel <- err
		close(errorChannel)
	}(logger)

	<-ctx.Done() // blocks until server stops

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutDownCtx); err != nil {
		return errors.New(fmt.Sprintf("server shutdown failed: %+s\n", err))
	}

	return nil
}

func main() {
	if err := run(context.Background(), os.Args, os.LookupEnv, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "unhandled error encountered: %+s\n", err.Error())
		os.Exit(1)
	}
}

type clArgs struct {
	addr string
}

func parseArgv(args []string) (*clArgs, error) {
	clargs, flagSet := new(clArgs), flag.NewFlagSet("Program Flags", flag.ContinueOnError)
	flagSet.StringVar(&clargs.addr, "addr", ":3000", "this API's address")
	err := flagSet.Parse(args[1:])
	return clargs, err
}
