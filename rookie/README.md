# Description

`ytdlp-rookie` is the binary built from `main.rs` code
It extracts cookies from the chrome browser and outputs it to console in [netscape cookie file format](https://curl.haxx.se/rfc/cookie_spec.html)\
It's bundled in the release and `ytdlp-ssh` runs it automatically

But could be use as standalone binary as well
```
ytdlp-rookie > /tmp/cookies
yt-dlp --cookies /tmp/cookies
```

Build for Linux only at the moment.
