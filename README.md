# rmq-br

cli can dump rabbitMQ messages to files and upload from files.


Examle: send all files from data_folder

```bash
rmq-br.exe send -u amqp://user:passwork@server:5672/vhost -x amq.topic -k routing-key -d data_folder  --quite true
```


Examle: read all messages from queue_name to data_folder
`Important` - message 'll be acknoleged!

```bash

rmq-br.exe dump -u amqp://user:passwork@server:5672/vhost -q queue_name -d data_folder
```


