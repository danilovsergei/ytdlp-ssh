#!/bin/bash
project_dir=$(dirname `realpath "$0"`)
presets_dir=$(echo $project_dir"/presets")
bin=$(echo $project_dir"/bin/ytdlp-ssh")
timestamp=$(date -d "@$(date +%s)" +"%y-%m-%d")
cd "$project_dir/cmd/ytdlp-ssh"

go build -o $bin/ytdlp-ssh
cp -r $presets_dir $bin"/"

cd  $project_dir"/bin"
zip -r $timestamp-release.zip ytdlp-ssh 