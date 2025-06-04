FROM python:3.13.4

RUN mkdir /api
WORKDIR /api

RUN apt update && \
    apt install -y postgresql-client

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY . .
