FROM alpine:3.10.1
ARG version="0.1.0"
LABEL maintainer="Haruaki TAMADA" \
    description="Purple cat, Purple cat, What do you see? I see the dependent libraries and their licenses!"

RUN    adduser -D purplecat \
    && cd /opt \
    && apk --no-cache add --update --virtual .builddeps curl tar \
    && curl -L https://github.com/tamadalab/purplecat/releases/download/v${version}/purplecat-${version}_linux_amd64.tar.gz -o /tmp/purplecat.tar.gz \
    && tar xvfz /tmp/purplecat.tar.gz \
    && ln -s /opt/purplecat-${version} /opt/purplecat \
    && rm /tmp/purplecat.tar.gz \
    && apk del --purge .builddeps

ENV HOME="/home/purplecat"
WORKDIR /home/purplecat
USER purplecat
ENTRYPOINT [ "/opt/purplecat/bin/purplecat" ]
