package ytdlp

import (
	"log"
	"path/filepath"
	"strings"
)

func ytDlpCommand(args CmdArgs) string {
	var fileName = DefaultFileFormat
	if args.FileFormat != "" {
		fileName = args.FileFormat
	}
	cmd := strings.Join([]string{
		"/usr/bin/youtube-dl",
		"--retries infinite",
		// It sets output directory for split-chapters
		"-P " + "\"" + args.Dir + "\"",
		// That would download OPUS from youtube
		// Opus is higher quality and better but not supported by direct play at Plex
		//"-f bestaudio",
		"-f bestaudio[ext=m4a]",
		"--extract-audio",
		"--cookies /tmp/ytdlp-cookies",
		"--output=\"" + filepath.Join(args.Dir, fileName) + "\"",
		// Extracting thumbnail Is causing trouble for opus youtube videos
		// But works fine for m4a
		// "--embed-thumbnail",
		"--add-metadata",
		"--verbose",
		// custom post processor to split by chapter and tag
		"--use-postprocessor SplitAndTag:when=after_move",
		args.Url,
	}, " ")
	log.Printf("\ncommand: %s\n", cmd)
	return cmd
}

func saveCookiesToFileCommand(cookies string) string {
	return "echo -e '" + cookies + "' > " + cookieFile
}
