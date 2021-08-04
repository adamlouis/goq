# goq

`goq` is a queue server written in golang

![a view of the goq web interface](/static/web.png?raw=true "a view of the goq web interface")

I am using `goq` to coordinate asynchronous jobs for some personal projects.

It currently uses sqlite as a data store, though could be modified to use any arbitrary storage layer.

It uses [adamlouis/gsg](https://github.com/adamlouis/gsg) to generate server boilerplate code.

It supports adding scheduled jobs to the queue via [robfig/cron](https://github.com/robfig/cron).

The `/api/*` routes provide a json api and other routes provide an web interface using [html/template](https://pkg.go.dev/html/template).
