FROM golang:1.22

WORKDIR /app

# Dependencies
COPY go.mod go.sum Makefile ./
RUN make mod-download

RUN make clean

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Source code
COPY . .

RUN make build

# Hot-Reload
CMD ["make", "run"]