#!/bin/bash
go-bindata -nocompress \
	voices/v0/dmr voices/v0/dstar voices/v0/p25 \
	voices/v1/srf-male-en/dmr voices/v1/srf-male-en/dstar voices/v1/srf-male-en/p25 \
	voices/v1/srf-female-en/dmr voices/v1/srf-female-en/dstar voices/v1/srf-female-en/p25
