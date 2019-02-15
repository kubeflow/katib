FROM debian:jessie

COPY pytorch-operator.v2 /pytorch-operator.v2
COPY pytorch-operator.v1beta1 /pytorch-operator.v1beta1

ENTRYPOINT ["/pytorch-operator", "-alsologtostderr"]
