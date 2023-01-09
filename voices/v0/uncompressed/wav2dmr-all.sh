#!/bin/bash

for f in *.wav; do
	./wav2dmr.sh "$f"
done
