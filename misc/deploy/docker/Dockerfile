FROM  alpine:3.4
RUN   set -x && \
      apk update && \
      apk add --no-cache ansible openssl ca-certificates && \
      rm -rf /tmp/* && \
      rm -rf /var/cache/apk/*
ARG   CACHEBUST=1
RUN   mkdir /wegigo && cd /wegigo && \
      wget -q https://github.com/rfancn/ansible-libs/archive/master.zip && \
      unzip master.zip "ansible-libs-master/ansible.cfg" && \
      unzip master.zip "ansible-libs-master/roles/ensure_kubernetes/*" && \
      unzip master.zip "ansible-libs-master/roles/configure_kubernetes/*" && \
      mv ansible-libs-master/* /wegigo && \
      rm -fr ansible-libs-master master.zip
COPY  wegigo /wegigo/
ENV   PATH /wegigo:$PATH
WORKDIR /wegigo
ENTRYPOINT ["wegigo", "deploy", "-b", "0.0.0.0"]