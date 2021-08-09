package main

import (
	"context"
	"flag"
	"github.com/gari8/wildyam"
	"github.com/gari8/wildyam/domain"
)

func main() {
	var recepter domain.Recepter
	flag.Parse()
	if len(flag.Args()) > 0 {
		arg := flag.Arg(0)
		switch {
		case arg == "run" || arg == "r":
			recepter.SubCmd = domain.Run
		case arg == "help" || arg == "h":
			recepter.SubCmd = domain.Help
		default:
			recepter.SubCmd = domain.Guide
		}
	} else {
		recepter.SubCmd = domain.Guide
	}

	ctx := context.Background()
	dm := wildyam.NewDrawManager(recepter)
	dm.Draw(ctx)
}
