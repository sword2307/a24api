#!/bin/bash

docker build --tag a24api .
docker run --rm a24api $@
