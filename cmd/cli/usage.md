# Usage for cli interface for gnotes

Basic documentation for gnotes, if you are not a develper, you probaly should
not be using this. 

This program is deisgned by develper(s) (so the ui is not ideal), for a develpers
needs.

## What is gnotes

GNOTES is a s3-syncing, encrypted note app that can be access accross multiable
platforms.

## Configuring

To configure, create a new configuration file at `~/.config/gnotes/config.ini`
with this template, then change all these templated `{}` values into your own:

```
[settings]
notes_dir = ${HOME}/.config/gnotes

# editor, you can change with others you like.
editor = vim

[s3]
active = true
bucket = {YOUR_S3_BUCKET}
endpoint = {YOUR_S3_ENDPOINT}
region = {YOUR_S3_REGION}
accesskey = {ACCESS_KEY}
secretkey = {SECRET_KEY}

crypt_key = {16_BIT_KEY}
user_id = {ONE_TIME_GENERATED_UUID}
```

* `user_id` should be a uuid, and it cannot change after initalization.
* `crypt_key` should be 16 bits (16 chars len), and encryption is enforced.
* `editor` should be a terminal editor of your choice.
* Make sure `s3/active` is true, gnotes is not designed to work without syncing to s3

After your configuration file is complete, run:

```
$ gnotes -R -s
```

**Only run once to initalize, otherwise you may lose data.**

Once ran, create a new note just so theres something to save into s3.
Then save and press `q` to quit gnotes, you should see it upload some
files.

### Creating new notes

To create new notes, first select a note folder (todo, not impmented), then
press enter for the first item to create a new note, add content with your
editor, save and exit. Make sure to quit gnotes succesfuly otherwise those
changes may not be saved.


### Deleting notes

To delete a note, select it, and delete all content from your editor. Make sure
**all** lines are removed. Save and exit. You should see your deleted note not
there. Quite gnotes to save your changes.


