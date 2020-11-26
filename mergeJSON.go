package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
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

		data, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal("Error reading: ", err)
		}

		// Remove Byte Order Mark (BOM) from data.
		// The BOM identifies that the text is UTF-8 encoded, but it should be
		// removed before decoding
		data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

		var v map[string]interface{}
		if err = json.Unmarshal(data, &v); err != nil {
			log.Fatalf("Error unmarshaling %s: %v\n", path, err)
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
