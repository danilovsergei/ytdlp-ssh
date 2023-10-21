package ytdlp

import (
	"log"
	"ytlpd-ssh/cookie"
	"ytlpd-ssh/sshclient"
)

const cookieFile = "/tmp/ytdlp-cookies"

// Executes yt-dlp command over ssh
//
//	yt-dlp receives cookies copied from chrome browser on ssh client
func ExecuteYtDlp(cmdArgs *CmdArgs) {
	cookies := cookie.ParseCookies(cmdArgs.Email)
	conn := sshclient.ConnectWithKey(cmdArgs.SshCreds)

	if !conn.IsFileOrDirExists(cmdArgs.OutDir) {
		log.Fatalf("Failed to download. Downloads directory %s does not exist", cmdArgs.OutDir)
	}
	conn.ExecCommand(saveCookiesToFileCommand(cookies))
	log.Println("Successfully copied cookies to the remote host")
	ytDlpCmd := ytDlpCommand(cmdArgs)
	log.Printf("ytdlp command: %s\n", ytDlpCmd)
	conn.ExecCommand(ytDlpCmd)
}
