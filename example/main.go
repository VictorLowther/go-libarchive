package main

import (
	"fmt"
	"os"

	ar "github.com/VictorLowther/go-libarchive"
)

func printContents(filename string) {
	fmt.Println("file ", filename)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error while opening file:\n %s\n", err)
		return
	}
	reader, err := ar.NewReader(file)
	if err != nil {
		fmt.Printf("Error on NewReader\n %s\n", err)
	}
	defer reader.Close()
	for {
		entry, err := reader.Next()
		if err != nil {
			fmt.Printf("Error on reader.Next():\n%s\n", err)
			return
		}
		fmt.Printf("Name %s, Size %d\n", entry.PathName(),entry.Size())
	}
}

func main() {
	for _, filename := range os.Args[1:] {
		printContents(filename)
	}
}
