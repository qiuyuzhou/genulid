package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "genulid",
		Usage: "Generate ULID",
		Action: func(cCtx *cli.Context) error {
			var t time.Time
			if tt := cCtx.Timestamp("time"); tt != nil {
				t = *tt
			}

			id, err := generate(t, cCtx.Bool("zero"))
			if err != nil {
				return err
			}

			fmt.Println(id.String())

			return nil
		},
		Flags: []cli.Flag{
			&cli.TimestampFlag{
				Name:    "time",
				Aliases: []string{"t"},
				Usage:   "when generating, use the specified time instead of now",
				Layout:  time.RFC3339,
			},
			&cli.BoolFlag{
				Name:    "zero",
				Aliases: []string{"z"},
				Usage:   "when generating, fix entropy to all-zeroes",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	return
}

func generate(t time.Time, zero bool) (id ulid.ULID, err error) {
	var entropy io.Reader

	if zero {
		entropy = zeroReader{}
	} else {
		entropy = rand.Reader
	}

	if t.IsZero() {
		t = time.Now()
	}

	id, err = ulid.New(ulid.Timestamp(t), entropy)
	if err != nil {
		return
	}

	return
}

type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}
