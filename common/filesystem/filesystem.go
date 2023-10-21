package filesystem

import (
	"log"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
)

const ytdlpAppName = "ytdlp-ssh"

// Checks file or dir exists on file system
func IsFileExists(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			log.Fatalf("Failed to find chrome profile: %s", err)
		}
	}
	return true
}

// Returns ytdlp-ssh config dir. Eg. ~/.config/ytdlp-ssh on Linux.
func YtdlpSshConfigDir() string {
	return filepath.Join(configdir.LocalConfig(), ytdlpAppName)
}
