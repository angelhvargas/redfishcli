FROM go:latest

COPY cmd/ /redfishcli/
CMD ["go", "build" "."]