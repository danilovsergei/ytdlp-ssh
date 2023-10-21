# Description
Its a go binary which executes yt-dlp on the remote host via ssh.
Handy features:
- Supports private key. So no need to type a password
- Supports user defined [presets](https://github.com/danilovsergei/ytdlp-ssh/tree/main#presets) with predefined yt-dlp flags
- Copies cookes from the client host to the remote host. So yt-dlp could login to the youtube-music account to download in maximum quality
- Prints remote output to console the same way like yt-dlp runs locally
- Automatically detects chrome profile to copy cookies.
- For users with multiple chrome profiles allows to find chrome profile by associated with it gmail account by provided --email
- [Improved m4a files spliting by chapter with preserving metadata](https://github.com/danilovsergei/yt-dlp-split-and-tag) 


# Usage
It uses custom m4a posprocesstor splitter because of the [yt-dlp bug in splitting m4a files](https://github.com/yt-dlp/yt-dlp/issues/8363).\
[SplitAndTag](https://github.com/danilovsergei/yt-dlp-split-and-tag) postprocesstor needs to be installed first on the remote ssh host per [README](https://github.com/danilovsergei/yt-dlp-split-and-tag).\
It gives correct metadata during splitting m4a files by chapters

Minimal required arguments:

```
go run main.go --dir="<REMOTE_HOST_DIR>" --url="<SUPPORTED_YT_DLP_URL>" --sshKey="<PRIVATE_KEY_PATH>" --sshHost="<USERNAME>@HOSTNAME:PORT"
```

--dir  - directrory on the remote host to save downloaded files\
--url - URL provided to yt-dlp\
--sshKey  - ssh private key to perform connection\
--sshHost - remote ssh host to connect to. Format is <USERNAME>@HOSTNAME:PORT> For example root@172.168.1.2 or root@172.168.1.2:22


In case there are multiple chrome profiles it's possible to specify email associated with chrome profile to find it:

```
go run main.go --dir="<REMOTE_HOST_DIR>" --url="<SUPPORTED_YT_DLP_URL>" --sshKey="<PRIVATE_KEY_PATH>" --sshHost="<USERNAME>@HOSTNAME:PORT" --email="myuser@gmail.com"
```

# Supported systems
## SSH Client
* Only linux is supported as client OS , mainly because cookie decryption implemnted only for Linux
* Only chrome supported as browser . mainly because of the cookies reading and decryption

## SSH Server
* Only *NIX systems are supported at the moment. Because of the shell specific commands for working with files

# Presets
Preset allows to pass predefined set of flags to the yt-dlp.
There are couple of predefined presets examples located at [ytdlp-ssh/presets](https://github.com/danilovsergei/ytdlp-ssh/tree/main/presets)

Pass the --preset flag to ytdlp-ssh

```
# By preset name
go run main.go --preset=m4a
# By absolute path
go run main.go --preset=/home/testuser/temp/m4a.preset
```

When short preset name passed like --preset=m4a presets could be located in :
- {ytdlp_ssh_binary_dir}/presets
- {current_dir}/presets
- local config dir, eg. /home/{username}}/.config/ytdlp-ssh/presets/

## Write new preset
See some predefined presets at [ytdlp-ssh/presets](https://github.com/danilovsergei/ytdlp-ssh/tree/main/presets).

Preset is just regular text file where each argument is on the new line for readability.\
yt-dlp path and arguments listed as is the same way they would be listed in terminal

```
// Downloads audio only in best available audio quality
// File format with highest available quality picked automatically

/usr/bin/yt-dlp
--retries infinite
// It sets output directory for split-chapters
-P "{{.OutDir}}"
-f bestaudio
--extract-audio
--cookies /tmp/ytdlp-cookies
--output="{{.OutDir}}/%(title)s [%(id)s].%(ext)s"
--add-metadata
--verbose
{{.Url}}
```

There is one exception though , eg. {{.Url}} and {{.OutDir}}

These are [Command line args](https://github.com/danilovsergei/ytdlp-ssh/blob/main/ytdlp/cmdargs.go) provided to ytdlp-ssh and exposed to the preset.\
Any [command line argument](https://github.com/danilovsergei/ytdlp-ssh/blob/main/ytdlp/cmdargs.go) could be accessed using {.<arg_name>} syntax. For details see [go templates](https://pkg.go.dev/text/template).\
All exposed command line args are defined in [cmdargs.go#CmdArgs](https://github.com/danilovsergei/ytdlp-ssh/blob/main/ytdlp/cmdargs.go)

Presets also support comments if text is prepended with //

