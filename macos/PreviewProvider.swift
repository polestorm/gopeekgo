import Cocoa
import Quartz
import UniformTypeIdentifiers

// Rozszerzenie Quick Look: czyta plik Markdown i oddaje stronę HTML zbudowaną przez gomdgo
// i gomermaidgo (archiwum Go, funkcje z bridge.h). System sam wyświetla HTML, więc nie ma tu widoku.
final class PreviewProvider: QLPreviewProvider, QLPreviewingController {
    func providePreview(for request: QLFilePreviewRequest) async throws -> QLPreviewReply {
        let url = request.fileURL
        let size = CGSize(width: 920, height: 1000)
        return QLPreviewReply(dataOfContentType: .html, contentSize: size) { reply in
            let source = try Data(contentsOf: url)
            reply.stringEncoding = .utf8
            reply.title = url.lastPathComponent
            return renderPage(source: source, title: url.lastPathComponent)
        }
    }
}

private func renderPage(source: Data, title: String) -> Data {
    var length: Int32 = 0
    let page: UnsafeMutablePointer<CChar>? = source.withUnsafeBytes { raw in
        title.withCString { titlePointer in
            GopeekgoPreview(raw.bindMemory(to: CChar.self).baseAddress, Int32(raw.count), titlePointer, &length)
        }
    }
    guard let page else { return Data() }
    defer { GopeekgoFree(page) }
    return Data(bytes: page, count: Int(length))
}
