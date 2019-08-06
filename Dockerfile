FROM golang:1.12
LABEL maintainer="Nino <martino.khuangga@tokopedia.com>"

# Starting directory on docker machine
WORKDIR $GOPATH/src/github.com/martinock/devcamp-backend

# Copy the necessary files from this project onto the docker machine
COPY . .

# Echo the current directory (`pwd`) and lists its content (`ls`)
RUN \
    echo "\nThis is the current directory:" && \
    pwd && \
    echo "\nThese are files within this directory:" && \
    ls -a && \
    echo ""

# Build process
RUN go get -v

# Port to be used
EXPOSE 3000

# Entry point for this Dockerfile
CMD ["devcamp-backend"]
