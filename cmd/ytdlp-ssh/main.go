package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"ytlpd-ssh/common/filesystem"
	"ytlpd-ssh/common/flags"
	"ytlpd-ssh/sshclient"
	"ytlpd-ssh/ytdlp"
)

const flagsConfig = "flags.ini"

func parseFlags() ytdlp.CmdArgs {
	var dir string
	var url string
	var sshKey string
	var sshHost string
	var email string
	var preset string

	flag.StringVar(&dir, "dir", "", "output directory to save download")
	flag.StringVar(&url, "url", "", "download url")
	flag.StringVar(&sshKey, "sshKey", "", "ssh key file path")
	flag.StringVar(&sshHost, "sshHost", "", "ssh user@hostname to login")
	flag.StringVar(&email, "email", "", "email to find chrome profile if there are multiple chrome profiles exist")
	flag.StringVar(&preset, "preset", "best_audio", "preset file with predefined yt-dlp flags.\n"+
		"Could be either absolute path or just preset name, eg. m4a.\n"+
		"Preset by name will be searched in {binary dir}/presets, {current dir}/presets or ~/.config/ytdlp-ssh/presets")

	flag.Parse()
	if flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	flags.SetFlagsFromConfig(filepath.Join(filesystem.YtdlpSshConfigDir(), flagsConfig))

	if dir == "" {
		log.Fatalln("--dir is empty. Downloads dir must be provided")
	}
	if url == "" {
		log.Fatalln("--url is empty. Download url must be provided")
	}
	if sshKey == "" {
		log.Fatalln("--sshKey is empty. Ssh key must be provided to execute commands via ssh")
	}
	if len(strings.Split(sshHost, "@")) != 2 {
		log.Fatalf("Hostname format must be user@hostname:port. But %s provided \n", sshHost)
	}

	user, host := parseUserHost(sshHost)
	return ytdlp.CmdArgs{
		Url:      url,
		OutDir:   dir,
		Email:    email,
		SshCreds: sshclient.SshCreds{User: user, Host: host, KeyFile: sshKey},
		Preset:   preset,
	}
}

func parseUserHost(hostFlag string) (user, host string) {
	userHost := strings.Split(hostFlag, "@")
	return userHost[0], userHost[1]
}

func main() {
	cmdArgs := parseFlags()
	ytdlp.ExecuteYtDlp(&cmdArgs)
}
