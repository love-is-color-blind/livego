
#FROM centos:latest:7
#RUN sudo yum install epel-release -y
#RUN sudo rpm --import http://li.nux.ro/download/nux/RPM-GPG-KEY-nux.ro
#RUN sudo rpm -Uvh http://li.nux.ro/download/nux/dextop/el7/x86_64/nux-dextop-release-0-5.el7.nux.noarch.rpm
#RUN sudo yum install ffmpeg ffmpeg-devel -y

FROM index.docker.io/jrottenberg/ffmpeg:4.1-centos
ADD livego.cfg /app/livego.cfg
ADD rtsp-live-stream /app/rtsp-live-stream
WORKDIR /app
VOLUME /app/db.txt
EXPOSE 7777 7001 7002 1935
ENTRYPOINT ["/app/rtsp-live-stream"]