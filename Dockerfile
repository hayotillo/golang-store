FROM golang:1.20

ARG port

WORKDIR /code

COPY go.mod .

RUN go mod download && go mod verify

RUN apt-get update -y && apt-get upgrade -y \
&& apt-get install -y util-linux\
&& rm -rf /var/lib/apt/lists/*

COPY . .

RUN chmod +x ./start.sh

EXPOSE $port

ENTRYPOINT ["/bin/bash", "./start.sh"]