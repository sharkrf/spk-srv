#!/bin/bash

gen() {
	echo "generating $1 $2.wav..."
#	pico2wave -w "$1 $2.wav" "$2"

	echo "$2" | text2wave -o "$1 $2.wav"
	sox "$1 $2.wav" outfile.wav trim =0.1 -0.1 vol 7 dB compand 0.3,1 30:-70,-30,-10 -7 -60 0.2
	mv outfile.wav "$1 $2.wav"

	# We take more params as pre-silence and post-silence
	if [ "$3" != "" ] || [ "$4" != "" ]; then
		sox -n -r 16k -b 16 -c 1 silence-pre.wav trim 0 $3
		sox -n -r 16k -b 16 -c 1 silence-post.wav trim 0 $4
		sox silence-pre.wav "$1 $2.wav" silence-post.wav outfile.wav
		mv outfile.wav "$1 $2.wav"
		rm silence-pre.wav silence-post.wav
	fi
}

gen PA alpha
gen PB bravo
gen PC charlie
gen PD delta
gen PE echo
gen PF foxtraat
gen PG golf
gen PH hotel
gen PI india
gen PJ juliet
gen PK kilo
gen PL lima
gen PM mike
gen PN november
gen PO oscar
gen PP papa
gen PQ quebec
gen PR romeo
gen PS sierra
gen PT tango
gen PU uniform
gen PV victor
gen PW whiskey
gen PX x-ray
gen PY yankee
gen PZ zulu

gen OS openspot 0.5 0
gen CT "connected to" 0 0.2
gen SV server
gen HB homebrew
gen MM "m m d v m"
gen FC "f c s"
gen YS "y s f"
gen SR "shark r f"
gen IP "i p"
gen BM brandmeister
gen DP "d m r plus"
gen RM room 0.3 0
gen CL client
gen CD connected
gen CO "trying to connect to"
gen CN "trying to connect"
gen WC "waiting for connection"
gen RF reflector 0 0.2
gen ST static
gen DN dynamic
gen DC disconnected
gen VE active
gen RO profile
gen RY ready
gen TG talkgroup
gen GS talkgroups
gen ND and
gen LK linked 0.3 0
gen GR "group call"
gen RI "private call"
gen RE "call routing is active" 0.3 0.3
gen DT dot
gen DR "i p address"
gen CP "access point"
gen WI "wai-fi"
gen NE network
gen IN internet
gen UN "un reachable"
gen SP special
gen CE connector
gen NX "n x d n"
gen NF "not found."
gen RQ "requested"
gen P2 "p 25"
gen N0 hundred
gen BT battery
gen RC percent
gen CG charging
gen TA ey
gen TP p
gen TM m
gen TI "time is"
gen TO oh
gen BC broadcast

for i in {B..Z}; do
	gen "A$i" $i
done

gen AA ay

for i in {0..99}; do
	n=`printf "%.2d" $i`
	gen $n $i
done
