FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
    curl \
    bash \
    docker.io \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY HLTV-Manager .
COPY frontend /app/frontend

RUN chmod +x ./HLTV-Manager

VOLUME /var/run/docker.sock:/var/run/docker.sock

USER root

CMD ["./HLTV-Manager"]
