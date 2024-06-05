package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"mime/quotedprintable"
	"os"
)

func main() {
	if err := mainE(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func mainE() error {
	decodeP := flag.Bool("d", false, "decode")
	flag.Parse()

	src, err := getSrc(flag.Args())
	if err != nil {
		return err
	}

	if *decodeP {
		_, err := io.Copy(os.Stdout, quotedprintable.NewReader(src))
		return errors.Join(err, src.Close())
	}
	w := quotedprintable.NewWriter(os.Stdout)
	_, err = io.Copy(w, src)
	return errors.Join(err, src.Close(), w.Close())
}

func getSrc(args []string) (io.ReadCloser, error) {
	switch len(args) {
	case 0:
		return io.NopCloser(os.Stdin), nil
	case 1:
		f, err := os.Open(args[0])
		if err != nil {
			return nil, fmt.Errorf("failed to open file '%s': %w", args[0], err)
		}
		return f, nil
	default:
		return nil, errors.New("zero or one filename should be provided")
	}
}
