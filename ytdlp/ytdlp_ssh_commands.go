package ytdlp

import (
	"strings"
)

const DownloadedFileFormat = "%(title)s [%(id)s].%(ext)s"

// Converts provided preset to the ready to execute yt-dlp command
//
// All not relevant information such as comments and whitespaces is stripped down from the preset
func ytDlpCommand(args *CmdArgs) string {
	preset := loadPreset(args)
	var ytlpCmd []string
	for _, line := range preset {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		ytlpCmd = append(ytlpCmd, line)
	}
	return strings.Join(ytlpCmd, " ")
}

// Saves browser cookies to the file on remote ssh server.
//
// At the moment remote file name is hard coded in cookieFile
func saveCookiesToFileCommand(cookies string) string {
	return "echo -e '" + cookies + "' > " + cookieFile
}
