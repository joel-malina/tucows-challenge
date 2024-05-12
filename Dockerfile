FROM accelbyte/alpine:3.11
LABEL maintainer="joel.malina@gmail.com"

ENV HOME /srv

COPY orderservice $HOME/orderservice
COPY swagger-ui $HOME/swagger-ui
COPY go.mod $HOME/go.mod
COPY orderservice.sha256 $HOME/orderservice.sha256

RUN find $HOME -type d -exec 'chmod' '555' '{}' ';' && \
    find $HOME -type f -exec 'chmod' '444' '{}' ';' && \
    find $HOME -type f -exec 'chown' 'root:root' '{}' ';' && \
    chmod 555 $HOME/orderservice

USER nobody

WORKDIR $HOME
ENTRYPOINT ["./orderservice"]