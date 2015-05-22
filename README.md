### Usage

```
./bench --url https://github.com/waniji/benchmark-tool --count 10 --worker 3
  --url             : アクセスするURL
  --count           : URLにアクセスする回数
  --worker          : 同時アクセス数
  --basic-auth-user : BASIC認証に使用するユーザー
  --basic-auth-pass : BASIC認証に使用するパスワード
```

### Result

```
URL: https://github.com/waniji/benchmark-tool
Total Access Count: 10
Concurrency: 3
--------------------------------------------------
Response Time: 930 msec, Status: 200 OK
Response Time: 931 msec, Status: 200 OK
Response Time: 932 msec, Status: 200 OK
Response Time: 796 msec, Status: 200 OK
Response Time: 807 msec, Status: 200 OK
Response Time: 837 msec, Status: 200 OK
Response Time: 760 msec, Status: 200 OK
Response Time: 796 msec, Status: 200 OK
Response Time: 831 msec, Status: 200 OK
Response Time: 721 msec, Status: 200 OK
--------------------------------------------------
Total Time           : 3251 msec
Average Response Time: 834 msec
Minimum Response Time: 721 msec
Maximum Response Time: 932 msec
```

