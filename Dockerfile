FROM golang:alpine

RUN apk update && apk add postgresql-dev
RUN apk add git
RUN apk add --no-cache gcc musl-dev
COPY ./ /go/src/github.com/nlittlepoole/differential
WORKDIR /go/src/github.com/nlittlepoole/differential

ENV PATH="$PATH:$GOROOT/bin:$GOPATH/bin"

RUN go get -u github.com/microo8/plgo/plgo
RUN ../../../../bin/plgo .

FROM postgres:9.6-alpine

RUN apk add --update make

COPY ./ /go/src/github.com/nlittlepoole/differential
WORKDIR /go/src/github.com/nlittlepoole/differential
COPY initdifferential.sql /docker-entrypoint-initdb.d
COPY --from=0 /go/src/github.com/nlittlepoole/differential/build build


RUN cd build && make install


