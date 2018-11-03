# iron/go is the alpine image with only ca-certificates added
FROM iron/go

WORKDIR /app

# Now just add the binary
ADD HaDiBar /app/
ADD settings.json /app/

EXPOSE 8080/tcp

ENTRYPOINT ["./HaDiBar"]