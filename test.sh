#!/bin/bash

if [ -z "$1" ]; then
	echo "usage: $0 [count]"
	exit 1
fi

echo '53524653504b0002cafecafe0201d5b5d08d0020f98402504e504f0000000000000000000000000000000000' | xxd -r -p > tmp.pkt

for i in `seq 1 $1`; do
	echo "sending pkt #$i"
	cat tmp.pkt > /dev/udp/127.0.0.1/65200
done

rm tmp.pkt
