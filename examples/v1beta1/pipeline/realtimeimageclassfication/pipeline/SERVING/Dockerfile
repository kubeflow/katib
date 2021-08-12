FROM ubuntu:latest

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
  && apt-get install -y python3-pip python3-dev \
  && cd /usr/local/bin \
  && ln -s /usr/bin/python3 python \
  && pip3 install --upgrade pip

RUN apt-get -yq install --no-install-recommends \
    libsm6 \
    libxext6 \
    libxrender-dev \
    tesseract-ocr \
    tesseract-ocr-ben \
    libtesseract-dev \
   libleptonica-dev \
   tesseract-ocr-all \
   libxrender-dev \
   protobuf-compiler \
   libtesseract-dev \
   libleptonica-dev \
   autoconf automake libtool \
   autoconf-archive \
   zlib1g-dev \
   libicu-dev \
   libgtk2.0-dev \
   libopencv-* \
   libpango1.0-dev \
   libcairo2-dev \
   pkg-config \
   gcc \
   g++  \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*
RUN mkdir /app
WORKDIR /app

ADD requirements.txt /app/requirements.txt
RUN pip3 install -r requirements.txt

RUN pip install opencv-contrib-python

ADD serving.py /app/serving.py
RUN chmod +x /app/serving.py



ENTRYPOINT ["python"]
CMD ["/app/serving.py"]




