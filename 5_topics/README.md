## Topics

Although using the `direct` exchange improved our system, it still has limitations - it can't do routing based on multiple criteria.

In our logging system we might want to subscribe to not only logs based on severity, but also based on the source which emitted the log. You might know this concept from the [`syslog`](http://en.wikipedia.org/wiki/Syslog) unix tool, which routes logs based on both severity (info/warn/crit...) and facility (auth/cron/kern...).

## How to run

To receive all the logs:

```
go run receive_logs_topic.go "#"
```

To receive all logs from the facility kern:

```
go run receive_logs_topic.go "kern.*"
```

Or if you want to hear only about critical logs:

```
go run receive_logs_topic.go "*.critical"
```

You can create multiple bindings:

```
go run receive_logs_topic.go "kern.*" "*.critical"
```

And to emit a log with a routing key kern.critical type:

```
go run emit_log_topic.go "kern.critical" "A critical kernel error"
```
