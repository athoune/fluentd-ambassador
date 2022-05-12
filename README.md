# Fluentd ambassador

Send all fluentd messages to a [Redis STREAMS](https://redis.io/docs/manual/data-types/streams/).

The Ambassador is tiny, brainless, and stateless.

The real job is done by some Redis STREAMS consumers.

```

a VM
+-----------------------+
|                       |
| Docker -> Ambassador -+-> Redis STREAMS
|                       |
+-----------------------+

```
