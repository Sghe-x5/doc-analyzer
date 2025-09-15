#!/bin/bash

rm -rf proto/third_party
mkdir -p proto/third_party

cd proto/third_party && git clone https://github.com/googleapis/googleapis.git && cd ../..

mv proto/third_party/googleapis/google proto
mv proto/third_party/googleapis/grafeas proto

rm -rf proto/third_party