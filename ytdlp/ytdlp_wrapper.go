package ytdlp

import (
	"log"
	"ytlpd-ssh/cookie"
	"ytlpd-ssh/sshclient"
)

const DefaultFileFormat = "%(title)s [%(id)s].%(ext)s"
const cookieFile = "/tmp/ytdlp-cookies"

type CmdArgs struct {
	sshclient.SshCreds
	Url        string
	Dir        string
	FileFormat string
}

// Executes yt-dlp command over ssh
//
//	yt-dlp receives cookies copied from chrome browser on ssh client
func ExecuteYtDlp(cmdArgs CmdArgs) {
	cookies := cookie.ParseCookies()
	conn := sshclient.ConnectWithKey(cmdArgs.SshCreds)

	if !conn.IsFileOrDirExists(cmdArgs.Dir) {
		log.Fatalf("Failed to download. Downloads directory %s does not exist", cmdArgs.Dir)
	}
	conn.ExecCommand(saveCookiesToFileCommand(cookies))
	conn.ExecCommand(ytDlpCommand(cmdArgs))
}
