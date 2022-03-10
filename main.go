package main

import (
	"fmt"
	"io"
	"os"
	"sort"

	yaml "gopkg.in/yaml.v2"
)

type ReleaseEntries map[string][]string

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("missing args")
	}

	entries, err := decodeReleaseFile(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Add header to table
	table := "| Name | Version | Backend Dev | Notes |\n| ---- | ------- | ----------- | ----- |\n"

	// Sort chart entries by name
	sortedEntries := make([]string, 0, len(entries))
	for chart := range entries {
		sortedEntries = append(sortedEntries, chart)
	}
	sort.Strings(sortedEntries)

	// Add rows to table
	for _, chart := range sortedEntries {
		versionsStr := ""
		for _, version := range entries[chart] {
			versionsStr += string(version) + "<br/>"
		}
		table += fmt.Sprintf("| %s | %s |||\n", chart, versionsStr)
	}
	fmt.Print(table)
}

func decodeReleaseFile(path string) (ReleaseEntries, error) {
	var entries ReleaseEntries
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := decodeYAMLFile(file, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func decodeYAMLFile(r io.Reader, target interface{}) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, target)
}
