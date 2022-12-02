package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

type Relations struct {
	Relations []*Relation `json:"relations"`
}

type Relation struct {
	Subject struct {
		ID   string `json:"id,omitempty"`
		Type string `json:"type,omitempty"`
		Key  string `json:"key,omitempty"`
	} `json:"subject"`
	Relation string `json:"relation"`
	Object   struct {
		ID   string `json:"id,omitempty"`
		Type string `json:"type,omitempty"`
		Key  string `json:"key,omitempty"`
	} `json:"object"`
}

func main() {
	var (
		input  string
		output string
		flip   bool
	)

	flag.StringVarP(&input, "input", "i", "", "relation tuples")
	flag.StringVarP(&output, "output", "o", "", "dot output file")
	flag.BoolVarP(&flip, "flip", "f", false, "invert directionality (sub->obj to obj->sub)")
	flag.Parse()

	fmt.Fprintf(os.Stderr, "convert %s into %s\n", input, output)

	if exists, _ := fileExists(input); !exists {
		fmt.Fprintf(os.Stderr, "input %s not found\n", input)
		os.Exit(1)
	}

	r, err := os.Open(input)
	if err != nil {
		log.Fatalln(err)
	}

	relations, err := readInput(r)
	if err != nil {
		log.Fatalln(err)
	}

	w := os.Stdout
	if output != "" {
		w, err = os.Create(output)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if err := convert(w, relations, flip); err != nil {
		log.Fatalln(err)
	}
}

func readInput(r io.Reader) ([]*Relation, error) {
	dec := json.NewDecoder(r)

	var relations Relations
	if err := dec.Decode(&relations); err != nil {
		return nil, err
	}

	return relations.Relations, nil
}

func convert(w io.Writer, relations []*Relation, flip bool) error {

	if _, err := w.Write([]byte("digraph G {\n")); err != nil {
		return err
	}

	for _, r := range relations {
		// a -> b [label="  a to b" labeltooltip="this is a tooltip"];
		if _, err := w.Write([]byte(
			iff(flip,
				fmt.Sprintf("\"%s:%s\" -> \"%s:%s\" [label=%q];\n",
					r.Object.Type, r.Object.Key,
					r.Subject.Type, r.Subject.Key,
					r.Relation,
				),
				fmt.Sprintf("\"%s:%s\" -> \"%s:%s\" [label=%q];\n",
					r.Subject.Type, r.Subject.Key,
					r.Object.Type, r.Object.Key,
					r.Relation,
				),
			),
		),
		); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte("}\n")); err != nil {
		return err
	}

	return nil
}

func fileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, errors.Wrapf(err, "failed to stat file '%s'", path)
	}
}

func iff[T any](cond bool, valTrue, valFalse T) T {
	if cond {
		return valTrue
	}
	return valFalse
}
