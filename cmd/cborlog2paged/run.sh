#!/bin/sh

export ENV_PAGE_SIZE4LOG=Xs
export ENV_PAGE_SIZE4LOG=Sm
export ENV_PAGE_SIZE4LOG=Md
#export ENV_PAGE_SIZE4LOG=Lg
#export ENV_PAGE_SIZE4LOG=Xl

input=./sample.d/input.jsonl

cat \
	"${input}" \
	"${input}" \
	"${input}" \
	"${input}" |
	json2map2cbor |
	./cborlog2paged |
	python3 \
		-m uv \
		tool \
		run \
		cbor2 \
		--sequence
