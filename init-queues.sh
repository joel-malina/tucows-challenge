#!/bin/bash
# Wait until RabbitMQ is fully up
until rabbitmqctl wait /var/lib/rabbitmq/mnesia/rabbit\@$(hostname).pid; do
  echo "Waiting for RabbitMQ to start..."
  sleep 10
done

rabbitmqadmin -u user -p password declare queue name=payment_request durable=true
rabbitmqadmin -u user -p password declare queue name=payment_response durable=true
