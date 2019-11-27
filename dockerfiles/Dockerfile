# iron/go is the alpine image with only ca-certificates added
FROM scratch

WORKDIR /app

# Now just add the binary
COPY server admin-client settings.json /app/
COPY webapp/ /app/webapp/

EXPOSE 8080/tcp

ENV GIN_MODE=release

ENTRYPOINT ["/app/server"]
VOLUME /app/data
VOLUME /app/sockets
