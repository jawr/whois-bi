FROM rabbitmq:management-alpine

COPY dev/rabbitmq.conf /etc/rabbitmq/rabbitmq.conf

ENTRYPOINT ["docker-entrypoint.sh"]
EXPOSE 25672/tcp 4369/tcp 5671/tcp
CMD ["rabbitmq-server"]
