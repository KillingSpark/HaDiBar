# iron/go is the alpine image with only ca-certificates added
FROM scratch

WORKDIR /app

# Now just add the binary
COPY hadibar hadibar
COPY settings.json settings.json
COPY webapp webapp

EXPOSE 8080/tcp

ENTRYPOINT ["./hadibar"]
VOLUME /app/data