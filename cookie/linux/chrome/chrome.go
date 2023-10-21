package chrome

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"ytlpd-ssh/common/filesystem"

	"github.com/kirsle/configdir"
)

// JSONN file in chrome profile with main profile settings
const preferencesFile = "Preferences"

type chromePreferences struct {
	AccountInfo []accountInfo `json:"account_info"`
}

type accountInfo struct {
	Email string `json:"email"`
}

// Finds chrome profile folder by provided email associated with profile.
//
// Returns last modified profile in case no matches for email or email is not provided
func FindChromeProfile(email string) string {
	chromeParentDir := filepath.Join(configdir.LocalConfig(), "google-chrome")
	files, err := os.ReadDir(chromeParentDir)
	if err != nil {
		log.Fatalf("Failed to find chrome profile : %s", err)
	}
	var mostRecentProfile string
	var mostRecentModifyDate int64
	for _, file := range files {
		profileDir := filepath.Join(chromeParentDir, file.Name())
		if isChromeProfile, stat := isChromeProfileDir(profileDir); isChromeProfile {
			if hasEmail(profileDir, email) {
				return profileDir
			}
			if stat.ModTime().UnixMicro() > mostRecentModifyDate {
				mostRecentModifyDate = stat.ModTime().UnixMicro()
				mostRecentProfile = profileDir
			}
		}
	}
	return mostRecentProfile
}

func isChromeProfileDir(path string) (bool, fs.FileInfo) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, nil
	}
	if !stat.IsDir() {
		return false, nil
	}
	return filesystem.IsFileExists(filepath.Join(path, preferencesFile)), stat
}

func hasEmail(profileDir, email string) bool {
	// Never match by email if no email provided
	if email == "" {
		return false
	}
	content, err := os.ReadFile(filepath.Join(profileDir, preferencesFile))
	// file does not exist is expected behavior and just use empty configuration
	preferences := chromePreferences{}
	if err != nil {
		log.Fatalf("Failed to find chrome profile : %s", err)
	}

	if err = json.Unmarshal(content, &preferences); err != nil {
		log.Fatalf("Failed to find chrome profile : %s", err)
	}
	if len(preferences.AccountInfo) == 0 {
		return false
	}
	// email could present in multiple profiles. It looks like in happens when
	// when chrome synced to email1 account opens sites logged as email2.
	// In that case cookies synced to both profiles : with email1 and email2
	for _, info := range preferences.AccountInfo {
		if strings.EqualFold(info.Email, email) {
			return true
		}
	}
	return false
}
