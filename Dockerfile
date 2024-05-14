FROM alpine:3.18
LABEL maintainer="joel.malina@gmail.com"

ENV HOME /srv

COPY order $HOME/order
COPY swagger-ui $HOME/swagger-ui
COPY go.mod $HOME/go.mod
COPY order.sha256 $HOME/order.sha256

RUN find $HOME -type d -exec 'chmod' '555' '{}' ';' && \
    find $HOME -type f -exec 'chmod' '444' '{}' ';' && \
    find $HOME -type f -exec 'chown' 'root:root' '{}' ';' && \
    chmod 555 $HOME/order

USER nobody

WORKDIR $HOME
ENTRYPOINT ["./order"]