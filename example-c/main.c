#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include "../bridge-c/gnotes-bridge.h"

char* catstr(char* str1, char* str2);

int main() {
	char* configStr = "config=config.ini new_note=no";

	InitApp(configStr);

	char* foo = Download("");
	printf("Download response: %s\n", foo);
	free(foo);

	foo = NewNote("");
	printf("NewNote response: %s\n", foo);
	free(foo);

	// Can only list once the notes are saved
	foo = List("");
	printf("List response: %s\n", foo);
	free(foo);

	foo = Save("notes_changed=yes");
	printf("Save response: %s\n", foo);
	free(foo);

	return 0;
}

char* catstr(char* str1, char* str2) {
	char* ret = malloc(strlen(str1) + strlen(str2) + 2);

	strcpy(ret, str1);
	strcat(ret, str2);

	return ret;
}

