## NAME

benchmark-tool - Web server benchmarking tool

## USAGE

```
benchmark-tool [options]
```

## OPTIONS

```
--config-file, -f    設定ファイル
--url, -u            アクセスするURL
--count, -c "0"      URLにアクセスする回数
--worker, -w "0"     同時アクセス数
--basic-auth-user    BASIC認証に使用するユーザー
--basic-auth-pass    BASIC認証に使用するパスワード
--format "simple"    実行結果の出力フォーマット
--version, -v        print the version
```

## CONFIG FILE EXAMPLE

```json
{
    "url": [
        "http://github.com",
        "http://www.google.com"
    ],
    "count": 30,
    "worker": 10,
    "basic-auth-user": "user",
    "basic-auth-pass": "pass"
}
```

## RESULT

```
30 / 30 [==============================================================================] 100.00 % 2s

Total Access Count: 30
Concurrency: 10
Total Time: 2548 msec

[all]
Success: 30
Failure: 0
Average Response Time: 695 msec
Minimum Response Time: 85 msec
Maximum Response Time: 1857 msec

[http://github.com]
Success: 15
Failure: 0
Average Response Time: 1287 msec
Minimum Response Time: 1009 msec
Maximum Response Time: 1857 msec

[http://www.google.com]
Success: 15
Failure: 0
Average Response Time: 104 msec
Minimum Response Time: 85 msec
Maximum Response Time: 136 msec
```

## AUTHOR

Makoto Sasaki

