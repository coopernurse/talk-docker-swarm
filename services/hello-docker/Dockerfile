# Dockerfile for coopernurse/hello-docker
#
FROM debian:stretch

RUN apt update && apt install -y python && apt clean
CMD ["python", "-c", "print 'hello docker'"]
