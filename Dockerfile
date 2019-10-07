FROM index.docker.io/jrottenberg/ffmpeg
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN  apk add openjdk8 && rm -rf /var/cache/apk/*

ADD ../docker /app

WORKDIR /app
EXPOSE 7777 7001 7002 1935
ENTRYPOINT ["java","-Djava.security.egd=file:/dev/./urandom","-Duser.timezone=GMT+08","-jar","app-jar-with-dependencies.jar"]
# docker run xxxx sh -c "ls &&  ls"