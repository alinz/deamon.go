package deamon

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type SummonType int

const (
	_      SummonType = iota
	Call              // Call happens once at the beginning of the program
	Recall            // Recall happens when the program is requested to be reloaded
	Kill              // Kill happens when the program is killed
)

type Summoner interface {
	Summon(ctx context.Context, summonType SummonType) error
}

type SummonerFunc func(ctx context.Context, summonType SummonType) error

func (f SummonerFunc) Summon(ctx context.Context, summonType SummonType) error {
	return f(ctx, summonType)
}

func Summoning(ctx context.Context, summoner Summoner) error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)

	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		close(signalChan)
		cancel()
	}()

	err := summoner.Summon(ctx, Call)
	if err != nil {
		return err
	}

	for {
		select {
		case s := <-signalChan:
			switch s {
			case syscall.SIGHUP:
				err := summoner.Summon(ctx, Recall)
				if err != nil {
					return err
				}
			case os.Interrupt:
				return summoner.Summon(ctx, Kill)
			}
		case <-ctx.Done():
			return context.Canceled
		}
	}
}
