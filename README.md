# gohttpd
A lightweight http file server writen by golang,DO NOT need any configuration. 

[![Build Status](https://travis-ci.org/jijinggang/gohttpd.svg)](https://travis-ci.org/jijinggang/gohttpd)

## Simple Usage

	gohttpd -dir root_dir -p port

for example:

	gohttpd -dir /var/www -p 80
	gohttpd -dir /var/www -p 443 -key key.pem -cert cert.pem
	gohttpd -cache false

## Command Line
	Usage of gohttpd.exe:
		-cache
			if cache is false, tell broswer dont't cache file (default true)
		-cert string
			cert.pem file， using by https
		-dir string
			http root folder (default ".")	
		-key string
			key.pem file， using by https
		-p string
			Address of the http server (default "80")

## Download
<https://github.com/jijinggang/gohttpd/blob/master/gohttpd.zip?raw=true>

Online Build

<http://gobuild.io/download/github.com/jijinggang/gohttpd>

