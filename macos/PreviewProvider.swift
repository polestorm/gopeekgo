import Cocoa
import Quartz
import UniformTypeIdentifiers

// The Quick Look extension: reads a Markdown file and returns the HTML page built by gomdgo, gomermaidgo
// and gohighlightgo (the Go archive, functions from bridge.h). The system shows the HTML, so there is no view.
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
