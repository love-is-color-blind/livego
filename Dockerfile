
#FROM centos:latest:7
#RUN sudo yum install epel-release -y
#RUN sudo rpm --import http://li.nux.ro/download/nux/RPM-GPG-KEY-nux.ro
#RUN sudo rpm -Uvh http://li.nux.ro/download/nux/dextop/el7/x86_64/nux-dextop-release-0-5.el7.nux.noarch.rpm
#RUN sudo yum install ffmpeg ffmpeg-devel -y

FROM index.docker.io/jrottenberg/ffmpeg:4.1-centos
ADD rtsp-live-stream /app/rtsp-live-stream
RUN echo '' > /app/db.txt
RUN chmod  u+x /app/rtsp-live-stream
WORKDIR /app
EXPOSE 8888  1935
ENTRYPOINT ["/app/rtsp-live-stream"]