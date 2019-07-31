# iron/go is the alpine image with only ca-certificates added
FROM scratch

WORKDIR /app

# Now just add the binary
COPY server admin-client settings.json webapp /app/

EXPOSE 8080/tcp

ENTRYPOINT ["/app/server"]
VOLUME /app/data
VOLUME /app/sockets
