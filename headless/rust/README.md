# README

## 运行

```sh
cargo run config.json
```

## 问题


存在自动断开连接的问题，迁移go(thread 'main' panicked at 'called `Result::unwrap()` on an `Err` value: The event waited for never came)

存在连接问题(thread 'main' panicked at 'called `Result::unwrap()` on an `Err` value: Unable to make method calls because underlying connection is closed)
