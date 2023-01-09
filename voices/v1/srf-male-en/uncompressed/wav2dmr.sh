#!/bin/bash

# Converts given file to AMBE.

if [ -z "$1" ]; then
	echo "usage: $0 [.wav file]"
	exit
fi

mkdir -p ../dmr
ffmpeg -i "$1" -f s16be -acodec pcm_s16be -ar 8000 /tmp/play.raw
~/work/sharkrf/a3k/build/Debug/a3k -d /dev/ttyUSB0 -q -20 -w 20 -m pcm2dmr -i /tmp/play.raw -o "../dmr/${1/.wav/}.ambe"
rm -f /tmp/play.raw
