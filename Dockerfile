############################
# STEP 1 build executable binary
############################
FROM golang:1.22-alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/jlquiros/CeliMail/
COPY . .
# Build the binary with a fixed name, whatever the module is called.
RUN CGO_ENABLED=0 go build -o /mail
RUN pwd
RUN ls -lart
RUN mkdir ~/.aws && \
    cp ./cred/* ~/.aws && \
    ls -lart ~/.aws

############################
# STEP 2 build a small image
############################
#FROM scratch
# Copy our static executable.
#COPY --from=builder /go/src/jlquiros/CeliMail/ /go/bin/
#RUN ls -lart /go/bin/
EXPOSE 6565
# Run the hello binary.
ENTRYPOINT ["/mail"]
