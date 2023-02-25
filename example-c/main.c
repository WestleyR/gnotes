#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>

#include "../bridge-c/gnotes-bridge.h"

char config_file[] = "/home/westley/.config/gnotes/config.ini";

int main() {

  // Download the note index
  Download(config_file);

  char test_note[] = "Notes/3f454501-2460-43df-8bee-2f446a1b6b1a/content";

  DownloadNote(config_file, test_note);

  printf("C: waiting 10 seconds in case you want to change %s to test upload\n", test_note);
  sleep(10);

  // Save will save all local notes that were changed, and save the index file.
  Save(config_file);

  return 0;
}
