
FROM scratch
ENTRYPOINT ["/goredir"]

# Add the binary
ADD goredir /
EXPOSE 8080
