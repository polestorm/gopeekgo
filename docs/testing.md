# Testy i pomiary

## Testy

- `go test ./...`: pakiet `preview` (strona zawiera wynik Markdown i diagram; błąd diagramu nie przerywa
  podglądu; surowy HTML nie przechodzi do strony; bloki kodu są kolorowane, a blok w nieznanym języku zostaje
  bezpiecznym tekstem; wersja w `preview.Version` zgadza się z `CFBundleShortVersionString` obu `Info.plist`).
- `go test -tags nohighlight ./...`: to samo bez kolorowania (test sprawdza, że blok kodu jest zwykłym tekstem).
- Część macOS sprawdzamy ręcznie: `macos/build.sh install`, potem `qlmanage -p plik.md` albo spacja w Finderze.
  Stopka strony („gopeekgo 0.1.0 · gomdgo + gomermaidgo + gohighlightgo · rozmiar · czas”) potwierdza, że podgląd
  pochodzi z tego rozszerzenia i z którego wariantu.

## Pomiary porównawcze: gopeekgo wobec QLMarkdown

Pomiar jest powtarzalny: `go run ./tools/bench` (czas, pamięć, rozmiar) i `go run ./tools/features` (wiersze
o funkcjach) w katalogu głównym, na plikach z `testdata/bench/` i `testdata/features/`. Wsad
(z sumami kontrolnymi), metoda i ograniczenia są opisane w [`benchmark.md`](benchmark.md) (po angielsku, bo ma
go móc przeczytać i powtórzyć konkurencja). Każda liczba porównania w `README.md` musi pochodzić z wiersza tego
rozdziału; pomiarów wykonanych inaczej nie publikujemy. Przewaga to `(QLMarkdown − gopeekgo) / QLMarkdown × 100%`.

### 2026-11-26, wersja 0.1.0 z kolorowaniem (`gohighlightgo` v0.1.2), kod roboczy na `a6ff34b`

Środowisko: Apple M4, macOS 27.0.1, zasilanie sieciowe, Go 1.26.6, `gomdgo` v0.1.1, `gomermaidgo` v0.1.0,
`gohighlightgo` v0.1.2, QLMarkdown 1.5.7 z ustawieniami jak niżej, 15 uruchomień na komórkę. Nowy wsad
`code.md` (21 bloków kodu w 19 językach) i jego 87 kopii; wsad z kodem pokazuje koszt kolorowania, którego
poprzednie pliki prawie nie miały.

| wsad | rozmiar | gopeekgo, mediana (min) | QLMarkdown, mediana (min) | przewaga | pamięć gopeekgo | pamięć QLMarkdown | przewaga |
|---|--:|--:|--:|--:|--:|--:|--:|
| `document.md` | 24 499 B | 3,3 ms (3,2 ms) | 21,4 ms (20,8 ms) | 84% | 8,1 MB | 21,7 MB | 63% |
| `diagram.md` | 202 B | 3,3 ms (3,2 ms) | 16,7 ms (16,3 ms) | 80% | 8,8 MB | 18,2 MB | 52% |
| 55 × `document.md` | 1 347 445 B | 15,7 ms (15,2 ms) | 213,6 ms (210,5 ms) | 93% | 32,8 MB | 59,5 MB | 45% |
| `code.md` | 15 429 B | 3,5 ms (3,4 ms) | 47,9 ms (47,3 ms) | 93% | 8,7 MB | 21,1 MB | 59% |
| 87 × `code.md` | 1 342 323 B | 31,2 ms (30,6 ms) | 2 222,1 ms (2 201,0 ms) | 99% | 42,8 MB | 57,2 MB | 25% |

| rozmiar na dysku | gopeekgo | QLMarkdown | przewaga |
|---|--:|--:|--:|
| aplikacja, suma rozmiarów plików | 7 258 656 B | 63 896 352 B | 89% |
| paczka rozszerzenia `.appex` | 7 177 500 B | 11 453 601 B | 37% |
| plik wykonywalny rozszerzenia, wycinek arm64 | 7 174 048 B | 5 426 192 B | −32% |

Ustawienia QLMarkdown w czasie pomiaru (`qlmarkdown_cli --show-settings`; zapisuje je też samo narzędzie):
kolorowanie włączone, Mermaid 12.0.0 i MathJax 4.1.3 wczytywane z `cdn.jsdelivr.net`, uwagi GitHub (`--alert`)
wyłączone, emoji jako znaki, nagłówek YAML dla wszystkich plików, obrazy lokalne wstawiane do strony,
`--tag-filter` włączony.

#### Kryterium H5 (`gohighlightgo`): ten sam pomiar z `-tags nohighlight`

| | z kolorowaniem | `nohighlight` | różnica |
|---|--:|--:|--:|
| plik wykonywalny rozszerzenia, arm64 | 7 174 048 B | 6 676 352 B | +497 696 B (+7,5%) |
| aplikacja, suma rozmiarów plików | 7 258 656 B | 6 761 578 B | +497 078 B |
| `document.md` | 3,3 ms | 2,6 ms | +0,7 ms |
| `code.md` | 3,5 ms | 2,4 ms | +1,1 ms |
| 87 × `code.md` (1,3 MB kodu) | 31,2 ms | 10,5 ms | +20,7 ms |
| pamięć, 87 × `code.md` | 42,8 MB | 18,2 MB | +24,6 MB |

Na małych plikach koszt to głównie rozruch (definicje języków wczytywane leniwie, rejestracja nazw w mapie
`CodeBlocks`). Wynik w wariancie `nohighlight` wobec QLMarkdown: 88%, 85%, 93%, 95%, 100% krócej i 71%, 61%,
49%, 72%, 68% mniej pamięci w tej samej kolejności wsadów; plik wykonywalny −23%.

### 2026-11-26, commit `a6ff34b` + narzędzie pomiarowe (bez kolorowania)

Środowisko: Apple M4, macOS 27.0.1, zasilanie sieciowe, Go 1.26.6, `gomdgo` v0.1.1, `gomermaidgo` v0.1.0,
QLMarkdown 1.5.7, 15 uruchomień na komórkę.

| wsad | rozmiar | gopeekgo, mediana (min) | QLMarkdown, mediana (min) | przewaga | pamięć gopeekgo | pamięć QLMarkdown | przewaga |
|---|--:|--:|--:|--:|--:|--:|--:|
| `document.md` | 24 499 B | 3,6 ms (2,9 ms) | 22,4 ms (21,9 ms) | 84% | 6,3 MB | 21,7 MB | 71% |
| `diagram.md` | 202 B | 2,7 ms (2,6 ms) | 17,5 ms (17,1 ms) | 85% | 7,0 MB | 18,2 MB | 62% |
| 55 × `document.md` | 1 347 445 B | 14,4 ms (14,1 ms) | 214,3 ms (211,3 ms) | 93% | 30,6 MB | 59,4 MB | 49% |

| rozmiar na dysku | gopeekgo | QLMarkdown | przewaga |
|---|--:|--:|--:|
| aplikacja, suma rozmiarów plików | 6 739 672 B | 63 896 352 B | 89% |
| paczka rozszerzenia `.appex` | 6 679 628 B | 11 453 601 B | 42% |
| plik wykonywalny rozszerzenia, wycinek arm64 | 6 676 176 B | 5 426 192 B | −23% |

Kolorowanie i diagramy leżą w QLMarkdown poza plikiem wykonywalnym (`libwrapper_highlight.dylib` 6,9 MB,
`mermaid.min.js.gz` 1,6 MB); nasz plik zawiera wszystko, łącznie ze środowiskiem uruchomieniowym Go.

### Pomiar wstępny 2026-11-25 (zastąpiony)

Pierwsze liczby (4,2 ms / 21,5 ms na 24 KB; 14,9 ms / 199,8 ms na 1,3 MB; aplikacja 6,6 MB / 64,0 MB) pochodziły
z doraźnego skryptu, z programu `gomdgo-mermaid` zamiast kodu tego repozytorium i z plików spoza repozytorium,
więc nie dało się ich powtórzyć. Zostały zastąpione pomiarem z 2026-11-26 i nie są już nigdzie cytowane.

## Próba funkcji: gopeekgo wobec QLMarkdown

`go run ./tools/features`: każda funkcja to osobny zamrożony plik w `testdata/features/`, a regułę, po której
narzędzie poznaje wynik na stronie, zapisano w `tools/features/main.go`. Reguły czytają treść strony bez
elementów `<style>` i `<script>`, żeby nazwa klasy w arkuszu nie udawała treści (tak było przy pierwszej
wersji narzędzia: alert QLMarkdown wyszedł „yes” z jego CSS).

### 2026-11-26, wersja 0.1.0, QLMarkdown 1.5.7 z ustawieniami jak wyżej

| funkcja | plik | gopeekgo | QLMarkdown |
|---|---|---|---|
| tabele, listy zadań, przekreślenie, adresy `www.`, przypisy, identyfikatory nagłówków | `table.md` … `heading-id.md` | tak | tak |
| uwagi GitHub (`> [!NOTE]`) | `alert.md` | tak | nie (przy tych ustawieniach; `--alert on` je włącza) |
| diagramy Mermaid | `mermaid.md` | SVG na stronie | Mermaid.js w podglądzie |
| kolorowanie składni | `code.md` | tak | tak |
| wzory | `math.md` | nie | tak, skrypt w podglądzie |
| nagłówek YAML | `front-matter.md` | czytany jako Markdown | pokazany jako metadane |
| skróty emoji | `emoji.md` | nie | tak |
| obraz o ścieżce względnej | `image.md` | zostaje adres | wstawiony do strony |
| `<script>` w źródle | `raw-script.md` | pokazany jako tekst | pokazany jako tekst |
| inny surowy HTML (`<b>`) | `raw-tag.md` | pokazany jako tekst | zinterpretowany jako HTML |
| odnośnik `javascript:` | `javascript-link.md` | usunięty | zostaje |
| skrypty na stronach, skąd wczytywane | wszystkie | brak | 4 elementy `<script>`; `cdn.jsdelivr.net` |

Poza narzędziem, sprawdzone 2026-11-25 ręcznie (`codesign`, `spctl`, `lipo`, GitHub API): paczka QLMarkdown jest
podpisana Developer ID i notaryzowana, plik rozszerzenia jest uniwersalny (x86_64 + arm64), licencja GPL-3.0,
3 795 gwiazdek i 124 rozwidlenia. Poprawka wobec próby z 2026-11-25: QLMarkdown nie usuwa `<script>`, tylko
pokazuje go jako tekst (`&lt;script>`), tak jak gopeekgo.

Znalezisko dla `gomdgo` (do paczki uwag): nagłówek YAML (`---`, `title: front`, `---`) wychodzi jako linia
pozioma i nagłówek `<h2>`. To zgodne z CommonMark, ale nagłówki YAML są częste w notatkach i dokumentacji;
do decyzji, czy biblioteka ma je pomijać albo oddawać osobno (opcja), czy robi to strona podglądu.
