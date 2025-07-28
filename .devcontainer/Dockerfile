FROM golang:1.24.0-alpine3.21


ENV GOPATH /usr/local/mailetter
ENV PATH ${GOPATH}/bin:$PATH
ENV PS1 "[\A]\u@mailttr:\W\$ "
ENV TZ "Aisa/Tokyo"
ENV STAGE dev

RUN apk update && \
    apk add doas \
            git \
            tzdata \
	        vim && \
    rm -rf /var/cache/apk/* && \
    adduser -u 2000 -D -h /home/dev -G users dev && \
    addgroup dev wheel && \
    echo -e "permit nopass :root\npermit nopass :wheel\n" > /etc/doas.d/doas.conf && \
	mkdir -p ${GOPATH}/bin ${GOPATH}/pkg && \
    chown -R dev:users ${GOPATH}/bin ${GOPATH}/pkg

USER dev
RUN go install golang.org/x/tools/cmd/goimports@latest && \
    go install github.com/golang/mock/mockgen@latest
WORKDIR ${GOPATH}/src

