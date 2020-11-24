package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
)

func main() {
	dir := flag.String("dir", ".", "Directory to get JSON object files from")
	o := flag.String("o", "./result.json", "Path to file to write resulting JSON object files to")
	flag.Parse()

	var files []map[string]interface{}

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatal(err)
	}
	absOutput, err := filepath.Abs(*o)
	if err != nil {
		log.Fatal(err)
	}

	// open json object files
	err = filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal("filepath.Walk: ", err)
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".json" {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			log.Fatal("os.Open path: ", err)
		}

		var v map[string]interface{}
		if err = json.NewDecoder(f).Decode(&v); err != nil {
			log.Fatal("Decode: ", err)
		}
		files = append(files, v)
		return nil
	})
	if err != nil {
		log.Fatal("filepath.Walk returned: ", err)
	}

	// merge json objects
	result := make(map[string]interface{})
	for _, jsonObj := range files {
		for k, v := range jsonObj {
			result[k] = v
		}
	}

	to, err := os.OpenFile(absOutput, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Open output file: ", err)
	}

	output := json.NewEncoder(to)
	output.SetEscapeHTML(false)
	output.SetIndent("", "	")
	if err = output.Encode(result); err != nil {
		log.Fatal("Encode: ", err)
	}
}
