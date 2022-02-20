# gnotes - Terminal based S3 syncing note app

**A WIP app. This is my local code and may not work for you. But you may
find this useful.**

Before using, you need access to a S3 server. Otherwise it will only save notes
locally.

_[screenshot here...]_

## Installation

```
$ git clone https://github.com/WestleyR/gnotes
$ cd gnotes/
$ make
$ cp gnotes ~/.local/bin  # or your preferred path
```

## Setting up gnotes for your S3 server

The user config file for gnotes is located in: _(subject to change)_

```
${HOME}/.config/wst.gnotes/config.ini
```

And as an example, for dreamhost, it should look like this:

```ini
[settings]
notes_dir = ${HOME}/.config/wst.gnotes
editor = vim

[encrypt]
enable = false
key = 16_bit_key______

[s3]
active = true
file = notes.tar.gz
bucket = gnotes
endpoint = https://objects-us-east-1.dream.io
region = us-east-1
accesskey = KEY
secretkey = KEY
```

