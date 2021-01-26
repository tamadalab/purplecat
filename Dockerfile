FROM alpine:3.10.1
ARG version="0.3.3"
LABEL maintainer="Haruaki TAMADA" \
    description="Purple cat, Purple cat, What do you see? I see the dependent libraries and their licenses!"

RUN    adduser -D -h /home/purplecat purplecat \
    && cd /opt    \
    && apk update \
    && apk --no-cache add --update --virtual .builddeps curl tar \
    && curl -s -L https://github.com/tamadalab/purplecat/releases/download/v${version}/purplecat-${version}_linux_amd64.tar.gz -o /tmp/purplecat.tar.gz \
#    && curl -s -L https://www.dropbox.com/s/b87at7bjn87n191/purplecat-${version}_linux_amd64.tar.gz?dl=0 -o /tmp/purplecat.tar.gz \
    && tar xvfz /tmp/purplecat.tar.gz \
    && ln -s /opt/purplecat-${version} /opt/purplecat \
    && rm -rf /tmp/purplecat.tar.gz /opt/purplecat/{README.md,LICENSE,completions} \
    && apk del --purge .builddeps

ENV HOME="/home/purplecat"
WORKDIR /home/purplecat
USER purplecat
ENTRYPOINT [ "/opt/purplecat/bin/purplecat" ]
# CMD /opt/purplecat/bin/purplecat --server --port=$PORT
