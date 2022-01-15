FROM alpine
ADD bin/amd64/httpserver /httpserver
RUN chmod 755 httpserver
EXPOSE 80
ENTRYPOINT /httpserver
