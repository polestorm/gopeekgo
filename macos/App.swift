import Cocoa

// Aplikacja-nośnik: macOS rejestruje rozszerzenie Quick Look tylko jako część aplikacji.
// Sama pokazuje jedynie krótką informację i kończy pracę po zamknięciu okna.
final class AppDelegate: NSObject, NSApplicationDelegate {
    private var window: NSWindow?

    func applicationDidFinishLaunching(_ notification: Notification) {
        let text = NSTextField(wrappingLabelWithString:
            "Podgląd Markdown (gomdgo + gomermaidgo + gohighlightgo) jest zainstalowany.\n\n" +
            "Zaznacz plik .md w Finderze i naciśnij spację. Rozszerzenie włączysz i wyłączysz " +
            "w Ustawieniach systemowych: Ogólne → Rzeczy otwierane przy logowaniu i rozszerzenia → Quick Look.")
        text.frame = NSRect(x: 24, y: 24, width: 432, height: 132)

        let window = NSWindow(contentRect: NSRect(x: 0, y: 0, width: 480, height: 180),
                              styleMask: [.titled, .closable], backing: .buffered, defer: false)
        window.title = "Gopeekgo"
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
