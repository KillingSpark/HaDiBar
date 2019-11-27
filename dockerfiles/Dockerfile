FROM busybox AS builder
RUN mkdir /empty-dir

FROM scratch

WORKDIR /app

# Now just add the binaries, settings and web resources
COPY server admin-client dockerfiles/settings.json /app/
COPY webapp/ /app/webapp/

# Copy empty directory (as mkdir is not available in scratch)
COPY --from=builder /empty-dir /sockets

EXPOSE 8080/tcp

ENV GIN_MODE=release \
	HADIBAR_PORT=8080 \
	HADIBAR_SOCKETPATH=/sockets \
	HADIBAR_LOGGINGLEVEL=DEBUG

ENTRYPOINT ["/app/server"]
VOLUME /data
VOLUME /logs
