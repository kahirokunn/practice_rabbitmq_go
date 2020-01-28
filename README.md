
rabbitmqの起動

```
$ docker pull rabbitmq
$ docker run -d --hostname my-rabbit --name some-rabbit -p 5672:5672 rabbitmq
```

goのパッケージのインストール

```
$ go get github.com/streadway/amqp
```

メッセージ配信

```
go run publisher.go
```

メッセージ受信 (複数起動してみるといい)

```
go run consumer.go
```
