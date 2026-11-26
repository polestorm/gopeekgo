// Declarations of the functions exported by ../bridge/main.go (the Go archive), for the Swift code.
#ifndef GOPEEKGO_BRIDGE_H
#define GOPEEKGO_BRIDGE_H

char *GopeekgoPreview(const char *src, int n, const char *title, int *outLen);
void GopeekgoFree(char *p);

#endif
