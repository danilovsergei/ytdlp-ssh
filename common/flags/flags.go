package flags

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
	"ytlpd-ssh/common/filesystem"
)

// Returns true if flag is provided from command line and false otherwise
func isFlagSet(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// Sets flags from config file
//
// Flag is set only if config exists and flag has not been provided from command line
func SetFlagsFromConfig(configPath string) {
	if !filesystem.IsFileExists(configPath) {
		return
	}
	log.Printf("Reading flags config from %s\n", configPath)
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// Create a scanner
	scanner := bufio.NewScanner(file)

	// Read the file line by line
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "=")
		if len(parts) != 2 || isFlagSet(parts[0]) {
			continue
		}
		flag.Set(parts[0], parts[1])

	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading file:", err)
		return
	}
}
