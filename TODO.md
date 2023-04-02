# TODO:

## v1
 - [ ] Add "Save on ctrl-c exit" to config.ini
 - [x] Sort notes by date
 - [x] Dont upload notes if it did not change
 - [x] Note title should be the first line of the note
 - [x] Encrypt the notes and store the encryption key in config.ini (maybe)
 - [ ] If config.ini is empty, then should put a default config there
 - [x] Add version flag
 - [ ] Add download/upload progress, or output to show that its uploading/downloading
 - [ ] Make the app self contained, and not part of the ui
 - [x] Should not store all attachments in the one archive (slow to download)
 - [x] Use https://github.com/h2non/filetype for file type detection
 - [ ] Eventally have "note folders" to store notes in a diffrent tarball (to decrease the download size)
 - [ ] Add a --disable-encryption flag to disable the note encryption (will download all note objects and decrypt them)
 - [ ] Add debugging logs
 - [x] Better recover on crash or fail to upload
 - [ ] Still use the appLock incase there is more then one app active (at least for local app)

## v2
 - [x] Loop through all notes and compair checksum to see it it needs to be uploaded, instead of having a array to track that
 - [x] Remove downloaded attachment encryption file after
 - [x] Open a note, edit, close, open it again and the changes seem to me missing or redownloaded (should fix by re-updating a note with a function right away)
 - [ ] Maybe should always update the index.json file after anything was changed? (not when the app closes)
 - [ ] Creating an empty note, and exiting will call delete from s3 when it was not created (not really an issue)
 - [x] Cannot download attachments
 - [x] Fix sorting issues
 - [ ] Fix c bindings
 - [ ] Create an ios app to use this
 - [x] Use env parsing to get the variable in the config file, like HOME
 - [ ] Use datasize.ByteSize type for file bytes (maybe)
 - [ ] Test changing a note, closing the app without uploading the index.json file and see how it recovers (should use local file maybe?) or just update the index.json right away
 - [x] Add flags to generate user id and 16 crypt key
 - [ ] Fix cmd/cli directory names to be go install-able
 - [ ] Should autoclean not tracked notes... maybe, since that are not uploaded anyway
 - [ ] Add flag to autoclean not tracked notes, ^^^ replaces above item
 - [ ] Should be a way to delete whole folders
 - [ ] When uploading attachments, should be a way to specify which folder it should be uploaded to, not just the current/last selected
 - [ ] Add window view for errors
 - [ ] Deleting a note folder does not delete the directory, maybe thats okay for backup

## v3 (not even started)
 - [ ] Be diff based for even faster performance (maybe, probaly not needed at all, v2 is fine)

