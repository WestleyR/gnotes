# gnotes - Terminal based S3 syncing note app

**A WIP app. This is my local code and may not work for you. But you may
find this useful.**

Before using, you access to a S3 server. Otherwise it will only save notes
locally.

As of now, these notes are not encrypted. Encryption would be easy to implemented,
but I dont need encryption right now.

_[screenshot here...]_

## Installation

_todo..._

```
$ go install ...
```

Or

```
$ git clone ... \
  cd gnotes/    \
  go build      \
  cp gnotes ~/.local/bin  # or your preferred path
```

## Setting up gnotes for your S3 server

The user config file for gnotes is located in: _(subject to change)_

```
${HOME}/.config/wst.gnotes/config.ini
```

_**NOTE:** You will need to restart the app after making changes to the config file._


And as an example, for dream host, it should look like this:

```ini
[settings]
editor = vim

[s3]
active = true
bucket = gnotes
endpoint = https://objects-us-east-1.dream.io
region = us-east-1
savefile = gnotes.json
accesskey = <ACCESS_KEY>
secretkey = <SECRET_KEY>
```

Make sure you set `active = true`. Otherwise gnotes will only use
local storage.

