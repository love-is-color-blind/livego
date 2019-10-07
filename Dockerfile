FROM index.docker.io/jrottenberg/ffmpeg:4.1-centos

ADD livego /app/livego
ADD livego.cfg /app/livego.cfg
RUN ls /app
RUN chmod 777 /app/livego
WORKDIR /app
VOLUME /app/db.txt
EXPOSE 7777 7001 7002 1935
ENTRYPOINT ["/app/livego"]