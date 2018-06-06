# TM Search for memoQ Server
[![Go Report Card](https://goreportcard.com/badge/github.com/pczajkowski/tmsearch)](https://goreportcard.com/report/github.com/pczajkowski/tmsearch)

This is a proof-of-concept tool (hobby project) which utilizes [memoQ server Resources API](https://www.memoq.com/en/the-memoq-apis/memoq-server-resources-api).

It provides simple HTML interface which of course can be improved. There's also logging mechanism which collects requestor's IP, phrase he was searching for, target language and number of served results. Logs are saved in *log* subfolder in separate *.log* files (one per day) in CSV format.

You just need to build it and make sure that subfolders **html** and **log** are present in the same location as your binary. You'll also need **secrets.json**, just make sure you fill it with proper credentials. Account used needs to be able to list TMs on your server and read their content, of course. It's using only standard GO packages, so there are no external dependencies.

Usage is simple. To get started just ***launch compiled binary with *-b* switch followed by the URL of your Resources API***. Now just navigate to *localhost/* in your browser and start searching your TMs. You may also want to adjust *html/languages.json* to be more relevant to your environment.

Optional parameters are as follows:

- *h* - if you want to serve it under hostname different than *localhost*
- *p* - if you want to serve it on port different than *80*

You can also navigate to *localhost/tms* to list all your TMs or to *localhost/tms?lang=fre-FR* to list TMs for given language.

**This app was designed to be used on local network or via VPN, so it lacks any security which would be necessary when exposed to Internet. It was also never tested under heavy load. You're free to use it however you wish, but I take no responsibility for any possible damage caused by it.**
