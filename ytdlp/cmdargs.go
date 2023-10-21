package ytdlp

import "ytlpd-ssh/sshclient"

// Encapsulates command line args provided to ytdlp-ssh
type CmdArgs struct {
	sshclient.SshCreds
	// Dowload url
	Url string
	// Output directory to save download
	OutDir string
	// email address used to find chrome profile
	Email string
	// ytdlp-ssh preset name
	Preset string
}
