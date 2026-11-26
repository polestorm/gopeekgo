#!/bin/sh
# Builds the "Gopeekgo" application with the Quick Look extension for Markdown files.
#
#   macos/build.sh             builds build/Gopeekgo.app
#   macos/build.sh install     builds, moves it to ~/Applications and registers the extension
#   macos/build.sh uninstall   unregisters the extension and removes the application
#
# Syntax highlighting (gohighlightgo) is on by default; TAGS=nohighlight macos/build.sh … builds without it
# (ADR 0003). The bundle version comes from macos/*-Info.plist; a test keeps it equal to preview.Version.
#
# Requires Go with cgo and the Xcode tools (swiftc, codesign). The signature is local (ad hoc); the variable
# SIGN_IDENTITY gives another identity.
set -eu

cd "$(dirname "$0")/.."

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
TAGS=${TAGS:-}

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
	echo "removed $INSTALLED"
	exit 0
fi

rm -rf "$BUILD"
mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources" "$APPEX/Contents/MacOS"

# 1. The Go library as a C archive.
CGO_ENABLED=1 MACOSX_DEPLOYMENT_TARGET=$MIN_MACOS CGO_CFLAGS="-mmacosx-version-min=$MIN_MACOS" \
	go build -trimpath -tags "$TAGS" -buildmode=c-archive -o "$BUILD/libbridge.a" ./bridge

# 2. The extension: the system provides the entry point (_NSExtensionMain), Info.plist names the class.
xcrun swiftc -O -swift-version 5 -target "$TARGET" -parse-as-library -application-extension \
	-module-name "$EXT_NAME" -import-objc-header macos/bridge.h \
	macos/PreviewProvider.swift "$BUILD/libbridge.a" \
	-framework Quartz -framework CoreFoundation -framework Security -lresolv \
	-Xlinker -e -Xlinker _NSExtensionMain \
	-o "$APPEX/Contents/MacOS/$EXT_NAME"
cp macos/Extension-Info.plist "$APPEX/Contents/Info.plist"

# 3. The host application.
xcrun swiftc -O -swift-version 5 -target "$TARGET" -parse-as-library -module-name "$APP_NAME" \
	macos/App.swift -o "$APP/Contents/MacOS/$APP_NAME"
cp macos/App-Info.plist "$APP/Contents/Info.plist"

# The licences go with the application: ours (Apache 2.0) and those of the fonts gomermaidgo embeds (OFL).
cp LICENSE NOTICE "$APP/Contents/Resources/"
FONTS=$(go list -m -f '{{.Dir}}' github.com/polestorm/gomermaidgo)/fonts
cp "$FONTS"/OFL-*.txt "$APP/Contents/Resources/"

# 4. Signing: the extension first, then the application. The extension must run in the sandbox.
codesign --force --sign "$SIGN_IDENTITY" --options runtime --entitlements macos/Extension.entitlements "$APPEX"
codesign --force --sign "$SIGN_IDENTITY" --options runtime --entitlements macos/App.entitlements "$APP"
echo "built $APP"

if [ "${1:-}" = install ]; then
	uninstall
	mkdir -p "$HOME/Applications"
	cp -R "$APP" "$INSTALLED"
	# A second copy of a bundle with the same identifier confuses the system registry, so only the installed one stays.
	rm -rf "$APP"
	"$LSREGISTER" -f "$INSTALLED"
	# Registration can lag behind lsregister, so it is repeated until the system lists the extension.
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
	echo "installed $INSTALLED"
	if [ "$registered" = no ]; then
		echo "warning: the system does not list the extension $EXT_ID; open the application once or install again" >&2
		exit 1
	fi
fi
