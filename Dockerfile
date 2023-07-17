FROM alpine:latest

COPY ./release/alidns_linux_amd64 /bin/alidns

RUN mkdir -p /var/log/alidns

ENTRYPOINT ["alidns"]
CMD ["ddns", "-conf", "/etc/ddns.json"]
