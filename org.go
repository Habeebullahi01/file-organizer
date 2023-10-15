package main

import (
	"flag"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	srcDir := flag.String("s", "", "Source directory")
	targetDir := flag.String("t", "", "Target directory")
	sourceFileType := flag.String("sft", ".jpg", "The type of file to rename")
	prefix := flag.String("p", "", "The prefix to be added to every renamed file.")
	// change file extension
	targetFileType := flag.String("tft", "", "New file extension of renamed files.")
	// sortCriteria := flag.String("sc", "", "The criteria to sort by")
	flag.Parse()

	log.Default()
	log.Println("This is ORG")
	// dirContent, _ := os.ReadDir(*srcDir)

	// Ensure existence of source directory
	if _, err := os.ReadDir(*srcDir); os.IsNotExist(err) || *srcDir == "" {
		log.Fatalln("The specified source directory does not exist.")
	}

	// Use a folder named "renamed" inside the source directory if a target directory was not specified
	if *targetDir == "" {
		log.Println("No target directory")
		*targetDir = filepath.Join(*srcDir, "renamed")
	}

	// Ensure existence of target directory
	if _, err := os.ReadDir(*targetDir); os.IsNotExist(err) {
		os.Mkdir(*targetDir, os.ModeDevice)
		log.Println("The specified target directory did not exist, but has now been created.")
	}

	// Format prefix
	if *prefix != "" {
		*prefix = *prefix + " "
	}

	// Ensure target file type
	if *targetFileType == "" {
		log.Println("no target file type")
		*targetFileType = *sourceFileType
	}

	// Ensure that target file type is a valid extension
	if strings.Split(*targetFileType, "")[0] != "." {
		*targetFileType = "." + *targetFileType
	}

	log.Printf("Source Directory: %s", *srcDir)
	log.Printf("Target Directory: %s", *targetDir)

	// No of file worked on
	index := 0

	// No to be used as file name
	rindex := 0

	filepath.WalkDir(*srcDir, func(path string, d fs.DirEntry, err error) error {

		// File has to meet 3 conditions:
		// 		Not be a directory
		// 		have the extension specified by filetype
		// 		reside directly under srcDir
		if !d.IsDir() && filepath.Ext(path) == *sourceFileType && filepath.Dir(path) == *srcDir {
			con, _ := os.ReadFile(path)
			rindex += 1

			// Create filename
			fname := filepath.Join(*targetDir, *prefix+strconv.Itoa(rindex)+*targetFileType)

			// Closure function to change file name
			uniquer := func() {
				rindex += 1

				fname = filepath.Join(*targetDir, *prefix+strconv.Itoa(rindex)+*targetFileType)
			}

			// Loop to ensure uniquness of new file name
			for !testUnique(fname) {
				uniquer()
			}

			if err := os.WriteFile(fname, con, 0600); err != nil {
				log.Println("An error occured")
				log.Println(err.Error())
			} else {
				index += 1
				log.Printf("File %d done.", index)
			}

		}
		return nil
	})
	log.Printf("Successfully renamed %d files. Files are in %s", index, *targetDir)

	os.Exit(0)

}

// Tests uniquness of file name
func testUnique(name string) bool {
	if _, err := os.ReadFile(name); err == nil {
		return false
	} else {
		return true
	}
}
