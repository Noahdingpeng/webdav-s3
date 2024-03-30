FROM goland:1.22.1 As builder

WORKDIR /app
COPY . /app
RUN gofmt -l .
RUN go get -d -v
RUN go build -o webdav -v .

FROM alpine:3.14.2
WORKDIR /app
COPY --from=builder /app/webdav /app/webdav
COPY --from=builder /app/conf/ /app/conf/
CMD [ "/app/webdav" ]