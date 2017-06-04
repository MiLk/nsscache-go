FROM ubuntu:16.04

RUN apt-get update \
 && apt-get install -y libnss-cache  \
 && rm -rf /var/lib/apt/lists/*

COPY nsswitch.conf /etc/nsswitch.conf

CMD ["getent"]
