
#include <stdio.h>
#include <stdlib.h>
#include "../bridge-c/gnotes-bridge.h"

int main() {
	char* foo = Download("");
	printf("Download response: %s\n", foo);
	free(foo);

	foo = List("");
	printf("List response: %s\n", foo);
	free(foo);

	foo = Save("notes_changed=yes");
	printf("Save response: %s\n", foo);
	free(foo);

	return 0;
}

