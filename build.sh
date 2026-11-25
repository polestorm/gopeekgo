#!/bin/sh
# Buduje aplikację "Gopeekgo" z rozszerzeniem Quick Look dla plików Markdown.
#
#   ./build.sh             buduje build/Gopeekgo.app
#   ./build.sh install     buduje, przenosi do ~/Applications i rejestruje rozszerzenie
#   ./build.sh uninstall   wyrejestrowuje rozszerzenie i usuwa aplikację
#
# Wymaga Go z cgo i narzędzi Xcode (swiftc, codesign). Podpis jest lokalny (ad hoc); inną tożsamość
# podaje zmienna SIGN_IDENTITY.
set -eu

cd "$(dirname "$0")"

APP_NAME=Gopeekgo
EXT_NAME=GopeekgoExtension
EXT_ID=com.polestorm.gopeekgo.extension
MIN_MACOS=13.0
ARCH=$(uname -m)
TARGET="$ARCH-apple-macos$MIN_MACOS"
BUILD=build
APP="$BUILD/$APP_NAME.app"
APPEX="$APP/Contents/PlugIns/$EXT_NAME.appex"
INSTALLED="$HOME/Applications/$APP_NAME.app"
LSREGISTER=/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister
SIGN_IDENTITY=${SIGN_IDENTITY:--}

uninstall() {
	if [ -d "$INSTALLED" ]; then
		pluginkit -r "$INSTALLED/Contents/PlugIns/$EXT_NAME.appex" 2>/dev/null || true
		"$LSREGISTER" -u "$INSTALLED" 2>/dev/null || true
		rm -rf "$INSTALLED"
	fi
	qlmanage -r >/dev/null 2>&1 || true
}

if [ "${1:-}" = uninstall ]; then
	uninstall
	echo "usunięto $INSTALLED"
	exit 0
fi

rm -rf "$BUILD"
mkdir -p "$APP/Contents/MacOS" "$APPEX/Contents/MacOS"

# 1. Biblioteka Go jako archiwum C.
CGO_ENABLED=1 MACOSX_DEPLOYMENT_TARGET=$MIN_MACOS CGO_CFLAGS="-mmacosx-version-min=$MIN_MACOS" \
	go build -trimpath -buildmode=c-archive -o "$BUILD/libbridge.a" ./bridge

# 2. Rozszerzenie: punkt wejścia daje system (_NSExtensionMain), klasę wskazuje Info.plist.
xcrun swiftc -O -swift-version 5 -target "$TARGET" -parse-as-library -application-extension \
	-module-name "$EXT_NAME" -import-objc-header macos/bridge.h \
	macos/PreviewProvider.swift "$BUILD/libbridge.a" \
	-framework Quartz -framework CoreFoundation -framework Security -lresolv \
	-Xlinker -e -Xlinker _NSExtensionMain \
	-o "$APPEX/Contents/MacOS/$EXT_NAME"
cp macos/Extension-Info.plist "$APPEX/Contents/Info.plist"

# 3. Aplikacja-nośnik.
xcrun swiftc -O -swift-version 5 -target "$TARGET" -parse-as-library -module-name "$APP_NAME" \
	macos/App.swift -o "$APP/Contents/MacOS/$APP_NAME"
cp macos/App-Info.plist "$APP/Contents/Info.plist"

# 4. Podpis: najpierw rozszerzenie, potem aplikacja. Rozszerzenie musi działać w piaskownicy.
codesign --force --sign "$SIGN_IDENTITY" --options runtime --entitlements macos/Extension.entitlements "$APPEX"
codesign --force --sign "$SIGN_IDENTITY" --options runtime --entitlements macos/App.entitlements "$APP"
echo "zbudowano $APP"

if [ "${1:-}" = install ]; then
	uninstall
	mkdir -p "$HOME/Applications"
	cp -R "$APP" "$INSTALLED"
	# Druga kopia paczki o tym samym identyfikatorze myli rejestr systemu, więc zostaje tylko zainstalowana.
	rm -rf "$APP"
	"$LSREGISTER" -f "$INSTALLED"
	# Rejestracja bywa spóźniona względem lsregister, więc ponawiamy, aż system pokaże rozszerzenie.
	registered=no
	for attempt in 1 2 3 4 5 6 7 8 9 10; do
		pluginkit -a "$INSTALLED/Contents/PlugIns/$EXT_NAME.appex"
		sleep 1
		if pluginkit -m -A -i "$EXT_ID" 2>/dev/null | grep -q "$EXT_ID"; then
			registered=yes
			break
		fi
	done
	pluginkit -e use -i "$EXT_ID"
	qlmanage -r >/dev/null 2>&1 || true
	echo "zainstalowano $INSTALLED"
	if [ "$registered" = no ]; then
		echo "uwaga: system nie pokazał rozszerzenia $EXT_ID; uruchom aplikację raz albo powtórz instalację" >&2
		exit 1
	fi
fi
