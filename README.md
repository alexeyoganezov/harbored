# Harbored

Presentation software that allows to access presenter slides from mobile
devices connected to conference-hall Wi-Fi network.

Makes your slides available to people in back rows with sight issues.

![Harbored demo](https://i.ibb.co/9nGFg76/harbored.png)

*NB: there is also a "audience view" window for use on projector,
it's not displayed here*

Advantages:

- Cross-platform and open source
- Works locally and doesn't depend on internet access
- Allows speakers to use any slide editor they want

## How to use

0. Connect to a Wi-Fi
0. Run an application
0. Select a directory containing PDF slides
0. Perform additional setup steps and start a presentation
0. Ask an audience to connect to the same Wi-Fi network
0. Press a special button to display QR-code

Then by reading and opening a link from the QR-code audience will see currently
displayed slide that will be updated as you switch slides or presentations.

*Important notice*:
this application wan't tested in battle conditions so use it carefully.
If you encountered any kind of issues please feel free to provide feedback.

## Technical details

1. The application starts a web-server on port 8080, make sure it's not blocked
by firewall rules

2. Network routers may have a limit on a number of simultaneous clients,
make sure it meets your needs

3. Throughput of routers is also limited, make sure your PDF files are
optimized using
[Adobe Compress PDF](https://www.adobe.com/acrobat/online/compress-pdf.html)
or
[iLovePDF](https://www.ilovepdf.com/ru/compress_pdf)

## Installation

Just download and use a binary for your OS from
[release page](https://github.com/alexeyoganezov/harbored/releases).

Those binaries aren't signed and notarized so Windows and MacOS users may encounter
security warnings from OS and antivirus software.

## Building from source

Install [Go](https://golang.org/dl/) and [Node.js](https://nodejs.org/en/)

Clone the repository:

```
git clone https://github.com/alexeyoganezov/harbored.git
cd harbored
```

Build front-end:

```
cd client
npm install
npm run-script build
```

Build embeddable resources:

```
cd ../
go get github.com/leaanthony/mewn/cmd/mewn
mewn
```

Create a binary:

```
go mod download
go build
```

And run it when needed:

```
./harbored
```

If need a distributable use one of the following commands:

```
fyne package -os darwin
fyne package -os linux
fyne package -os windows
```
