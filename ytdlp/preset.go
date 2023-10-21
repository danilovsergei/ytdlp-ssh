package ytdlp

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"ytlpd-ssh/common/filesystem"
)

const presetsDir = "presets"

// Loads preset from args.Preset.
//
// Uses either absolutel path or tries to guess preset path if only preset name provided
// Searches preset by name in {binary dir}/presets, {current dir}/presets or ~/.config/ytdlp-ssh/presets

func loadPreset(args *CmdArgs) []string {
	presetPath := args.Preset
	if !filepath.IsAbs(presetPath) {
		if filepath.Ext(presetPath) != ".preset" {
			presetPath = presetPath + ".preset"
		}
		presetPath = absolutePresetPath(presetPath)
	}
	preset, err := template.ParseFiles(presetPath)
	if err != nil {
		log.Fatalf("Failed to load preset %s: %s\n", args.Preset, err)
	}
	var buf bytes.Buffer
	err = preset.Execute(&buf, args)
	if err != nil {
		log.Fatalf("Failed to load preset %s: %s\n", args.Preset, err)
	}
	return strings.Split(buf.String(), "\n")
}

// Gets absolute preset path from short preset name .

// Searches preset by name in {binary dir}/presets, {current dir}/presets or ~/.config/ytdlp-ssh/presets
// m4a transformed to like of ~/.config/ytdlp-ssh/presets/m4a.preset

// Returns empty string if preset not found
func absolutePresetPath(presetName string) string {
	if absPath := getPresetInBinaryDir(presetName); absPath != "" {
		return absPath
	}

	if absPath := getPresetInCurrentDir(presetName); absPath != "" {
		return absPath
	}

	if absPath := getPresetInConfigDir(presetName); absPath != "" {
		return absPath
	}

	log.Fatalf("Failed to locate preset %s\n", presetName)
	return ""
}

// Returns absolute preset path if presetName found in provided dir/presetsDir/ dir or either empty string
func presetPath(dir, presetName string) string {
	fullPath := filepath.Join(dir, presetsDir, presetName)
	if filesystem.IsFileExists(fullPath) {
		return fullPath
	}
	return ""
}

// Returns absolute preset path if presetName found in {current dir}/presetsDir/ or either empty string
func getPresetInCurrentDir(presetName string) string {
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to load preset %s: %s\n", presetName, err)
	}
	return presetPath(currDir, presetName)
}

// Returns absolute preset path if presetName found in {current binary}/presetsDir/ or either empty string
func getPresetInBinaryDir(presetName string) string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to load preset %s: %s\n", presetName, err)
	}
	return presetPath(filepath.Dir(exePath), presetName)
}

// Returns absolute preset path if presetName found in {user config dir}/presetsDir/ or either empty string
func getPresetInConfigDir(presetName string) string {
	configDir := filesystem.YtdlpSshConfigDir()
	return presetPath(configDir, presetName)
}
