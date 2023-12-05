## GGoS
HTTP/2 DDoS method that uses legitimate headers and TLS configurations, and focuses on compliance with Cloudflare HTTP DDoS rules.

### Usage
```sh
go run main.go url 120 100 proxies.txt 100000
```
Thread number this high since goroutines work a little bit different than usual threads

### Acknowledgements
Thanks to [@t101804](https://github.com/t101804) for [CloudFlare Bypass](https://github.com/t101804/Priv8Bypass)
