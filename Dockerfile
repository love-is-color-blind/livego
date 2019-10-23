
#FROM centos:latest:7
#RUN sudo yum install epel-release -y
#RUN sudo rpm --import http://li.nux.ro/download/nux/RPM-GPG-KEY-nux.ro
#RUN sudo rpm -Uvh http://li.nux.ro/download/nux/dextop/el7/x86_64/nux-dextop-release-0-5.el7.nux.noarch.rpm
#RUN sudo yum install ffmpeg ffmpeg-devel -y

FROM index.docker.io/jrottenberg/ffmpeg:4.1-centos
ADD rtsp-live-stream /rtsp-live-stream
ADD index.html /index.html
RUN echo '' > /db.txt
RUN chmod  u+x rtsp-live-stream
EXPOSE 8888  1935
ENV LIVE_STREAM_REDIRECT_SERVER localhost
ENTRYPOINT ["./rtsp-live-stream"]