FROM busybox AS builder
RUN mkdir /empty-dir

FROM scratch

# Set a random and static user id for the hadibar user (and group).
# Operators may have to configure write access of additional volume mounts for this user id.
ARG HADIBAR_USER_ID=30159

WORKDIR /app

# Now just add the binaries, settings and web resources
COPY server admin-client dockerfiles/settings.json /app/
COPY webapp/ /app/webapp/

# Copy empty directory (as mkdir is not available in scratch)
COPY --from=builder --chown=${HADIBAR_USER_ID}:${HADIBAR_USER_ID} /empty-dir /sockets
COPY --from=builder --chown=${HADIBAR_USER_ID}:${HADIBAR_USER_ID} /empty-dir /data
COPY --from=builder --chown=${HADIBAR_USER_ID}:${HADIBAR_USER_ID} /empty-dir /logs

USER ${HADIBAR_USER_ID}:${HADIBAR_USER_ID}

EXPOSE 8080/tcp

ENV GIN_MODE=release \
	HADIBAR_PORT=8080 \
	HADIBAR_SOCKETPATH=/sockets \
	HADIBAR_LOGGINGLEVEL=DEBUG

ENTRYPOINT ["/app/server"]
VOLUME /data
VOLUME /logs
