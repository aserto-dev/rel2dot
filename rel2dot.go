package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	dsc3 "github.com/aserto-dev/go-directory/aserto/directory/common/v3"

	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

type Relations struct {
	// Relations []*Relation `json:"relations"`
	Relations []*dsc3.Relation `json:"relations"`
}

type Entity struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
}

func (e *Entity) String() string {
	return fmt.Sprintf("%s:%s", e.Type, e.ID)
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

	fmt.Fprintf(os.Stderr, "convert %s into %s\n", input, iff(output == "", "stdout", output))

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

func readInput(r io.Reader) ([]*dsc3.Relation, error) {
	dec := json.NewDecoder(r)

	var relations Relations
	if err := dec.Decode(&relations); err != nil {
		return nil, err
	}

	return relations.Relations, nil
}

func convert(w io.Writer, relations []*dsc3.Relation, flip bool) error {

	if _, err := w.Write([]byte("digraph G {\n")); err != nil {
		return err
	}

	for _, r := range relations {
		if _, err := w.Write([]byte(
			iff(flip,
				fmt.Sprintf("\"%s:%s\" -> \"%s:%s\" [label=%q];\n",
					r.ObjectType, r.ObjectId,
					r.SubjectType, r.SubjectId,
					r.Relation,
				),
				fmt.Sprintf("\"%s:%s\" -> \"%s:%s\" [label=%q];\n",
					r.SubjectType, r.SubjectId,
					r.ObjectType, r.ObjectId,
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
