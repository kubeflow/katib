FROM debian:jessie

COPY pytorch-operator.v1beta1 /pytorch-operator.v1beta1
COPY pytorch-operator.v1beta2 /pytorch-operator.v1beta2

ENTRYPOINT ["/pytorch-operator", "-alsologtostderr"]
