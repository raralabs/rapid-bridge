FROM golang:1.23.4-bullseye As build

WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends git


COPY go.mod ./
COPY go.sum ./


# Add the private key arg
ARG SSH_PRIVATE_KEY
RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/
# Make the root foler for our ssh 
RUN mkdir -p /root/.ssh && \
chmod 0700 /root/.ssh && \
ssh-keyscan github.com > /root/.ssh/known_hosts

RUN echo "$SSH_PRIVATE_KEY" >> /root/.ssh/id_rsa && \
chmod 600 /root/.ssh/id_rsa && \
# ssh-add /root/.ssh/id_rsa  && \
go mod download && \
rm -rf /root/.ssh/


RUN go mod download
COPY . ./

RUN --mount=type=cache,target="/root/.cache/go-build" go build  -o rapid-bridge cmd/main.go



# new stage
FROM debian:buster

WORKDIR /app

RUN apt-get update && \
  apt-get install -y --no-install-recommends libc6 ca-certificates && \
  update-ca-certificates




COPY  --from=build app/rapid-bridge ./server

CMD ["./rapid-bridge", "init", "server"]

