FROM golang:alpine

RUN mkdir -p src/csv_builder
WORKDIR /go/src/csv_builder

COPY wpe_merge wpe_merge.go /go/src/csv_builder/

# If you need to copy other files
# COPY input.csv /go/src/csv_builder/