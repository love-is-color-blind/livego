FROM index.docker.io/jrottenberg/ffmpeg

ADD livego /app/livego
ADD livego.cfg /app/livego.cfg

WORKDIR /app
VOLUME /app/db.txt
EXPOSE 7777 7001 7002 1935
ENTRYPOINT ["/app/livego"]