FROM alpine:3.12.4

ARG KUBECTL_VERSION="v1.19.3"

RUN apk add --update openssl
RUN wget https://storage.googleapis.com/kubernetes-release/release/$KUBECTL_VERSION/bin/linux/amd64/kubectl \
  && chmod +x ./kubectl && mv ./kubectl /usr/local/bin/kubectl

COPY ./hack/cert-generator.sh /app/cert-generator.sh
RUN chmod +x /app/cert-generator.sh

WORKDIR /app
ENTRYPOINT ["sh", "./cert-generator.sh"]
