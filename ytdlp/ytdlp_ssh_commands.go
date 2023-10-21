package ytdlp

import (
	"strings"
)

const DownloadedFileFormat = "%(title)s [%(id)s].%(ext)s"

// Loads yt-dlp command from provided preset.
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

func saveCookiesToFileCommand(cookies string) string {
	return "echo -e '" + cookies + "' > " + cookieFile
}
