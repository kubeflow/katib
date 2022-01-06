# FROM python:3.8-slim-buster

FROM python:3.9

RUN apt-get update
RUN apt-get install ffmpeg libsm6 libxext6 -y

WORKDIR /app
COPY . /app
RUN pip install --no-cache-dir -r requirements.txt