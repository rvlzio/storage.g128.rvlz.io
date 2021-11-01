FROM python:3.10.0-slim-buster
RUN mkdir /app
WORKDIR /app
ADD requirements.txt ./
RUN pip install -r requirements.txt
ADD . .
