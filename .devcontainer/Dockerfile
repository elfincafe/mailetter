FROM docker.io/library/golang:1.22-alpine3.19

LABEL version="0.1.0"
LABEL description="MaiLetter Laibrary for Go"
LABEL maintainer="Takamichi Seki <takamichi.seki@tsukinoha.jp>"

ENV GOROOT /usr/local/go
ENV GOPATH /usr/local/mailetter
ENV PATH $GOPATH/bin:$PATH

RUN mv /go ${GOROOT} && \
    apk update && \
    apk add doas \
	        vim && \
    adduser -u 2000 -G users -D -H zlkv && \
    echo -e "permit nopass :root\npermit nopass :wheel\n" > /etc/doas.conf && \
	mkdir -p ${GOPATH}/bin ${GOPATH}/pkg && \
    chown -R zlkv:users ${GOPATH}

USER zlkv
ENV PS1="[\A]\u@\h:\W\$ "
WORKDIR ${GOPATH}/src
CMD ["/bin/sh"]

