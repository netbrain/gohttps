
FROM scratch
ENTRYPOINT ["/gohttps"]

# Add the binary
ADD gohttps /
EXPOSE 8080
