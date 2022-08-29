package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/xerrors"
)

const header = `// Code generated by gen/enum. DO NOT EDIT.
package database
`

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	dump, err := os.Open("dump.sql")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %s must be run in the database directory with dump.sql present\n", os.Args[0])
		return err
	}
	defer dump.Close()

	var uniqueConstraints []string

	s := bufio.NewScanner(dump)
	query := ""
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		switch {
		case strings.HasPrefix(line, "--"):
		case line == "":
		case strings.HasSuffix(line, ";"):
			query += line
			if isUniqueConstraint(query) {
				uniqueConstraints = append(uniqueConstraints, query)
			}
			query = ""
		default:
			query += line + " "
		}
	}
	if err = s.Err(); err != nil {
		return err
	}

	return writeContents("unique_constraint.go", uniqueConstraints, generateUniqueConstraints)
}

func isUniqueConstraint(query string) bool {
	return strings.Contains(query, "UNIQUE")
}

func generateUniqueConstraints(queries []string) ([]byte, error) {
	s := &bytes.Buffer{}

	_, _ = fmt.Fprint(s, header)
	_, _ = fmt.Fprint(s, `
// UniqueConstraint represents a named unique constraint on a table.
type UniqueConstraint string

// UniqueConstraint enums.
const (
`)
	for _, query := range queries {
		name := ""
		switch {
		case strings.Contains(query, "ALTER TABLE") && strings.Contains(query, "ADD CONSTRAINT"):
			name = strings.Split(query, " ")[6]
		case strings.Contains(query, "CREATE UNIQUE INDEX"):
			name = strings.Split(query, " ")[3]
		default:
			return nil, xerrors.Errorf("unknown unique constraint format: %s", query)
		}
		_, _ = fmt.Fprintf(s, "\tUnique%s UniqueConstraint = %q // %s\n", nameFromSnakeCase(name), name, query)
	}
	_, _ = fmt.Fprint(s, ")\n")

	return s.Bytes(), nil
}

func writeContents[T any](dest string, arg T, fn func(T) ([]byte, error)) error {
	b, err := fn(arg)
	if err != nil {
		return err
	}
	err = os.WriteFile(dest, b, 0o600)
	if err != nil {
		return err
	}
	cmd := exec.Command("goimports", "-w", dest)
	return cmd.Run()
}

func nameFromSnakeCase(s string) string {
	var ret string
	for _, ss := range strings.Split(s, "_") {
		switch ss {
		case "id":
			ret += "ID"
		case "ids":
			ret += "IDs"
		case "jwt":
			ret += "JWT"
		case "idx":
			ret += "Index"
		default:
			ret += strings.Title(ss)
		}
	}
	return ret
}
