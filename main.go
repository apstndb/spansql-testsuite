package main

import (
	"cloud.google.com/go/spanner/spansql"
	"context"
	"encoding/csv"
	"github.com/samber/lo"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalln(err)
	}
}

func run(ctx context.Context) error {
	csvWriter := csv.NewWriter(os.Stdout)
	csvWriter.Write([]string{"filename", "success", "text"})
	for _, tcase := range []struct {
		dir      string
		stmtType string
	}{
		{"./memefish/testdata/input/ddl", "DDL"},
		{"./memefish/testdata/input/query", "Query"},
		{"./memefish/testdata/input/dml", "DML"},
	} {
		dirPath := tcase.dir
		dir, err := os.ReadDir(dirPath)
		if err != nil {
			return err
		}
		for _, f := range dir {
			if f.IsDir() {
				continue
			}
			fpath := filepath.Join(dirPath, f.Name())
			b, err := os.ReadFile(fpath)
			if err != nil {
				return err
			}
			text := string(b)
			switch tcase.stmtType {
			case "DDL":
				_, err = spansql.ParseDDL(fpath, text)
			case "DML":
				_, err = spansql.ParseDML(fpath, text)
			case "Query":
				_, err = spansql.ParseQuery(text)
			}
			csvWriter.Write([]string{
				lo.Must(filepath.Rel("./memefish/testdata/input", fpath)),
				strconv.FormatBool(err == nil),
				text})
		}
	}
	csvWriter.Flush()
	return nil
}
