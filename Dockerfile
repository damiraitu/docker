FROM golang
COPY . /gofolder
WORKDIR /gofolder
RUN go mod download
ENV GO111MODULE=on\
    CGO_ENABLED=0\
    GOOS=linux\
    GOARCH=amd64
RUN go build /gofolder/cmd/web
RUN chmod +x /gofolder/cmd/web
CMD ["./web"]
EXPOSE 7070