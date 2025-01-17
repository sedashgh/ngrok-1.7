# ngrok - Introspected tunnels to localhost ([homepage](https://ngrok.com))
### "I want to securely expose a web server to the internet and capture all traffic for detailed inspection and replay"
![](https://ngrok.com/static/img/overview.png)

## What is ngrok?
ngrok is a reverse proxy that creates a secure tunnel between from a public endpoint to a locally running web service.
ngrok captures and analyzes all traffic over the tunnel for later inspection and replay.

## What can I do with ngrok?
- Expose any http service behind a NAT or firewall to the internet on a subdomain of ngrok.com
- Expose any tcp service behind a NAT or firewall to the internet on a random port of ngrok.com
- Inspect all http requests/responses that are transmitted over the tunnel
- Replay any request that was transmitted over the tunnel


## What is ngrok useful for?
- Temporarily sharing a website that is only running on your development machine
- Demoing an app at a hackathon without deploying
- Developing any services which consume webhooks (HTTP callbacks) by allowing you to replay those requests
- Debugging and understanding any web service by inspecting the HTTP traffic
- Running networked services on machines that are firewalled off from the internet


## Downloading and installing ngrok
ngrok has _no_ runtime dependencies. Just download a single binary for your platform and run it. Some premium features
are only available by creating an account on ngrok.com. If you need them, [create an account on ngrok.com](https://ngrok.com/signup).

- [Linux](https://dl.ngrok.com/linux_386/ngrok.zip)
- [Mac OSX](https://dl.ngrok.com/darwin_386/ngrok.zip)
- [Windows](https://dl.ngrok.com/windows_386/ngrok.zip)


## Developing on ngrok
[ngrok developer's guide](docs/DEVELOPMENT.md)

## Compile in docker for windows
Assuming you have cloned the repo to c:\projects\git\ngrok-1.7 in windows
docker run --rm  -v c:/projects/git/ngrok-1.7:/ngrok -w /ngrok -e GOOS=windows -e GOARCH=amd64 golang:1.4.3-cross make release-client
