ARG cuda_version=10.0
ARG cudnn_version=7
FROM pytorch/pytorch:1.0-cuda${cuda_version}-cudnn${cudnn_version}-runtime

ENV TARGET_DIR /opt/darts-cnn-cifar10

ADD examples/v1beta1/trial-images/darts-cnn-cifar10 ${TARGET_DIR}
WORKDIR  ${TARGET_DIR}

RUN chgrp -R 0 ${TARGET_DIR} \
  && chmod -R g+rwX ${TARGET_DIR}

ENTRYPOINT ["python3", "-u", "run_trial.py"]
