#!/bin/bash

for f in *.wav; do
	./wav2dstar.sh "$f"
done
