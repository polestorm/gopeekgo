// Deklaracje funkcji eksportowanych z ../bridge/main.go (archiwum Go) dla kodu w Swifcie.
#ifndef GOPEEKGO_BRIDGE_H
#define GOPEEKGO_BRIDGE_H

char *GopeekgoPreview(const char *src, int n, const char *title, int *outLen);
void GopeekgoFree(char *p);

#endif
