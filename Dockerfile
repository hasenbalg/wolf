# syntax=docker/dockerfile:1

FROM golang:bullseye

# Set destination for COPY
WORKDIR /app
RUN apt update \
&& apt install wakeonlan curl -y \ 
&& apt clean \
&& rm -rf /var/lib/apt/lists/*
# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o ./wolf

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 3333
RUN ls -la
RUN echo http://localhost:3333
COPY config.yaml ./
COPY templates ./templates

# Run
CMD ["./wolf"]
