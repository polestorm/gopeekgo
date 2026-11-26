import Cocoa

// The host application: macOS registers a Quick Look extension only as part of an application.
// It only shows a short note and quits when its window is closed.
final class AppDelegate: NSObject, NSApplicationDelegate {
    private var window: NSWindow?

    func applicationDidFinishLaunching(_ notification: Notification) {
        let text = NSTextField(wrappingLabelWithString:
            "The GoPeekGo Markdown preview (gomdgo + gomermaidgo + gohighlightgo) is installed.\n\n" +
            "Select a .md file in Finder and press Space. To turn the extension on or off, open " +
            "System Settings: General → Login Items & Extensions → Quick Look.")
        text.frame = NSRect(x: 24, y: 24, width: 432, height: 132)

        let window = NSWindow(contentRect: NSRect(x: 0, y: 0, width: 480, height: 180),
                              styleMask: [.titled, .closable], backing: .buffered, defer: false)
        window.title = "GoPeekGo"
        window.contentView?.addSubview(text)
        window.center()
        window.makeKeyAndOrderFront(nil)
        self.window = window
        NSApp.activate(ignoringOtherApps: true)
    }

    func applicationShouldTerminateAfterLastWindowClosed(_ sender: NSApplication) -> Bool { true }
}

@main
enum Main {
    static func main() {
        let app = NSApplication.shared
        let delegate = AppDelegate()
        app.delegate = delegate
        app.setActivationPolicy(.regular)
        app.run()
    }
}
