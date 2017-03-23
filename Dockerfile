FROM    alpine:3.4
RUN mkdir /blog
COPY journey /blog
COPY content /blog/content
COPY config.json /blog
COPY built-in /blog/built-in
WORKDIR /blog
ENV MYSQL_HOST=127.0.0.1
ENV MYSQL_USER=root
ENV MYSQL_PASS=admin
ENV MYSQL_PORT=3306
ENV MYSQL_DATABASE=blog
COPY docker-entrypoint.sh /bin/
RUN chmod 755 /bin/docker-entrypoint.sh
CMD ["sh","/bin/docker-entrypoint.sh"]



