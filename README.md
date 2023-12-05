## GGoS
HTTP/2 DDoS method that uses legitimate headers and TLS configurations, and focuses on compliance with Cloudflare HTTP DDoS rules.

### Usage

For massive attacks on URLs protected by Cloudflare:

```sh
go run cloudflare.go url 120 100 proxies.txt 100000
```
Thread number this high since goroutines work a little bit different than usual threads

### Acknowledgements
Thanks to [@t101804](https://github.com/t101804) for [Cloudflare Bypass](https://github.com/t101804/Priv8Bypass)

### Notice
This was cloned from [Dark Anubis Repo](https://github.com/darkanubis0100/http2go) to be updated and removed the existing bugs and that one was forked from the [original](https://github.com/NetworkDir/http2go) made by [@udbnt](https://github.com/udbnt)
