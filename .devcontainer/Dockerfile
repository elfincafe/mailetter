FROM docker.io/library/alpine:3.21.2


ENV GOPATH /usr/local/mailetter
ENV PATH ${GOPATH}/bin:$PATH
ENV PS1 "[\A]\u@mailtr:\W\$ "
ENV STAGE dev

RUN apk update && \
    apk add doas \
	        vim \
            go && \
    rm -rf /var/cache/apk/* && \
    adduser -u 2000 -D -h /home/dev -G users dev && \
    addgroup dev wheel && \
    echo -e "permit nopass :root\npermit nopass :wheel\n" > /etc/doas.d/doas.conf && \
	mkdir -p ${GOPATH}/bin ${GOPATH}/pkg && \
    chown -R zlkv:users ${GOPATH}/bin ${GOPATH}/pkg

USER dev
RUN go install golang.org/x/tools/cmd/goimports@latest
WORKDIR ${GOPATH}/src

