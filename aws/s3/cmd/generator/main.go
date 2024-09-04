package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	numberOfFiles := flag.Int("n", 100, "number of files to generate")
	flag.Parse()

	log.Printf("Generating %d files...", *numberOfFiles)
	i := 0
	for i < *numberOfFiles {
		createFile(i)
		i++
	}
	log.Println("Files generated.")
}

func createFile(i int) {
	f, err := os.Create(fmt.Sprintf("./tmp/file%d.txt", i))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("this is the file %d!", i))
}
