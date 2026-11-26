# Plan

| # | Zakres | Stan |
|---|---|---|
| 0 | prototyp: strona podglądu, most C, rozszerzenie, `macos/build.sh`, instalacja lokalna | gotowy (2026-11-25) |
| 1 | obejrzenie na dokumentach rzeczywistych: tryb ciemny, długie tabele, szerokie diagramy, duże pliki; lista znalezisk dla bibliotek | nie rozpoczęty |
| 2 | obrazy o ścieżkach względnych (załączniki odpowiedzi podglądu) | nie rozpoczęty |
| 3 | przełącznik w `macos/build.sh` do budowania z lokalnych kopii bibliotek (`go.work`) | nie rozpoczęty |
| 4 | kolorowanie składni z `gohighlightgo`, domyślnie włączone, znacznik `nohighlight` wyłącza (ADR 0003); pomiar rozmiaru i czasu w obu wariantach; warunek wydania | gotowy z `gohighlightgo` v0.1.2, H5 zmierzone (2026-11-26); wydanie czeka na decyzję właściciela |
| 5 | porównanie z konkurencją na rynku (QLMarkdown, podgląd systemowy, inne): rozmiar, czas, pamięć, funkcje; tabela w `README.md` | QLMarkdown powtarzalnie (`tools/bench`, `tools/features`, 2026-11-26); inne podglądy do dodania |
| 6 | czas od spacji do obrazu (nie tylko konwersja) | nie rozpoczęty |
| 7 | rozpowszechnianie: binarka uniwersalna, Developer ID, notaryzacja | do decyzji |
| 8 | CI: `go vet`, `go test`, budowanie aplikacji na macOS (ADR 0002) | plik gotowy; przepływ wyłączony w GitHubie (`gh workflow disable`) do nowego okresu rozliczeniowego, 2026-11-25 |
| 9 | upublicznienie: wersja 0.1.0, komentarze po angielsku, `NOTICE` z czcionkami OFL, porównanie z wsadem (ADR 0004) | gotowe do przeglądu (2026-11-26); zmiana widoczności i znacznik `v0.1.0` na decyzję właściciela |
| 10 | teksty aplikacji-nośnika po angielsku przy publikacji (`CLAUDE.md`) | do decyzji |
