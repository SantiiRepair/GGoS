# GGoS
HTTP/2 DDoS Method using legit headers aswell as tls configs and checks against mainly cloudflares http ddos rules

usage example: ```go run main.go url 120 100 proxies.txt 100000``` thread number this high since goroutines work a little bit different than usual threads

thanks to https://github.com/t101804/Priv8Bypass for cf bypass
