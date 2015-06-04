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
Total Access Count: 6
Concurrency: 2
--------------------------------------------------
Time: 100 msec, Status: 200 OK, URL: http://www.google.com
Time: 1503 msec, Status: 200 OK, URL: http://github.com
Time: 1408 msec, Status: 200 OK, URL: http://github.com
Time: 77 msec, Status: 200 OK, URL: http://www.google.com
Time: 90 msec, Status: 200 OK, URL: http://www.google.com
Time: 999 msec, Status: 200 OK, URL: http://github.com
--------------------------------------------------
Total Time           : 2508 msec
Average Response Time: 696 msec
Minimum Response Time: 77 msec
Maximum Response Time: 1503 msec
```

## AUTHOR

Makoto Sasaki

