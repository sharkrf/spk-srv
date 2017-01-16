# spk-srv

Voice announcement server. Plays voice announcement .ambe file sequence when
requested with a request UDP packet.

# Installing

```
go get github.com/sharkrf/spk-srv
go generate github.com/sharkrf/spk-srv
go install github.com/sharkrf/spk-srv
```

# About the voice announcement files

spk-srv has a pregenerated set of announcement files generated with
[Festival](http://www.festvox.org/festival/) using the
[voice_cmu_us_bdl_cg](http://festvox.org/packed/festival/2.4/voices/festvox_cmu_us_bdl_cg.tar.gz)
voice package.
