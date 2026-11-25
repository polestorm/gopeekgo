# Testy i pomiary

## Testy

- `go test ./...`: pakiet `preview` (strona zawiera wynik Markdown i diagram; błąd diagramu nie przerywa
  podglądu; surowy HTML nie przechodzi do strony).
- Część macOS sprawdzamy ręcznie: `macos/build.sh install`, potem `qlmanage -p plik.md` albo spacja w Finderze.
  Stopka strony („gopeekgo · gomdgo + gomermaidgo · rozmiar · czas”) potwierdza, że podgląd pochodzi z tego
  rozszerzenia.

## Pomiar 2026-11-25: gopeekgo wobec QLMarkdown 1.5.7

Komputer: MacBook arm64, Go 1.26.6. Kod: prototyp sprzed przeniesienia (`gomdgo` w commicie `587ee59`,
`gomermaidgo` v0.1.0). Przewaga to `(QLMarkdown − gopeekgo) / QLMarkdown × 100%`.

### Wielkość

| | gopeekgo | QLMarkdown | przewaga |
|---|--:|--:|--:|
| cała aplikacja na dysku (`du`) | 6,6 MB | 64,0 MB | 90% |
| paczka rozszerzenia `.appex` | 6,5 MB | 11,2 MB | 42% |
| plik wykonywalny rozszerzenia, arm64 | 6,68 MB | 5,43 MB | −23% |

Plik wykonywalny QLMarkdown jest uniwersalny (x86_64 + arm64, razem 11,4 MB); w tabeli jest sam wycinek arm64.
Kolorowanie i diagramy leżą u niego poza tym plikiem: `libwrapper_highlight.dylib` 6,9 MB, `mermaid.min.js.gz`
1,6 MB. Nasz plik zawiera wszystko, łącznie ze środowiskiem uruchomieniowym Go.

### Szybkość i pamięć konwersji

Sposób: programy wiersza poleceń obu rozwiązań (`gomdgo-mermaid -standalone -alerts -footnotes -heading-ids`
i `qlmarkdown_cli` z paczki QLMarkdown), czas całego procesu, mediana z 15 uruchomień; pamięć szczytowa
z `/usr/bin/time -l`.

| plik | gopeekgo | QLMarkdown | przewaga |
|---|--:|--:|--:|
| czas, `gomdgo/docs/testing.md` (24 KB) | 4,2 ms | 21,5 ms | 80% |
| czas, próbka z diagramem `mermaid` (0,2 KB) | 2,9 ms | 16,5 ms | 82% |
| czas, 1,3 MB (40 × `testing.md` + `README.md`) | 14,9 ms | 199,8 ms | 93% |
| pamięć, próbka | 7,1 MB | 18,2 MB | 61% |
| pamięć, 1,3 MB | 31,8 MB | 60,3 MB | 47% |

Zastrzeżenia:
- nasz czas zawiera narysowanie diagramu; QLMarkdown wstawia tylko `<script>` z `mermaid.min.js`, a diagram
  rysuje JavaScript w oknie podglądu, więc jego czas przy diagramach jest w rzeczywistości dłuższy;
- QLMarkdown koloruje bloki kodu, my nie;
- to pomiar konwersji programem wiersza poleceń, nie czasu od spacji do obrazu; tego nie mierzono;
- po przeniesieniu plik wykonywalny rozszerzenia ma 6 676 112 B (budowany wtedy ze starej treści `gomdgo` v0.1.0; od tego samego dnia zależność to v0.1.1).

## Próba funkcji 2026-11-25: gopeekgo wobec QLMarkdown 1.5.7

Sposób: jeden plik z nagłówkiem YAML, tabelą, listą zadań, przekreśleniem, adresem `www.`, przypisem, uwagą
`> [!NOTE]`, diagramem `mermaid`, `<script>` i `<b>`, odnośnikiem `javascript:`, blokiem `go`, wzorami `$…$`
i `$$…$$`, skrótem `:smile:` i obrazem o ścieżce względnej. QLMarkdown: `qlmarkdown_cli` z ustawieniami
zainstalowanej aplikacji (`--show-settings`); gopeekgo: `preview.Page` z bieżącego kodu. Wynik to wiersze
o funkcjach w tabeli `README.md`.

- QLMarkdown, przy tych ustawieniach: wczytuje Mermaid 12.0.0 i MathJax 4.1.3 z `cdn.jsdelivr.net`
  (strona ma CSP dopuszczające tylko ten adres); uwag GitHub nie rozpoznaje, dopóki nie włączy się `--alert`;
  `<script>` usuwa, `<b>` przepuszcza, odnośnik `javascript:` zostawia; koloruje kod, zamienia `:smile:`
  na znak, nagłówek YAML pokazuje jako tabelę, istniejący obraz wstawia jako `data:`. Paczka jest podpisana
  Developer ID i notaryzowana, plik rozszerzenia jest uniwersalny; licencja GPL-3.0.
- gopeekgo: strona bez `<script>` i bez adresów zewnętrznych; diagram jako SVG; surowy HTML jako tekst;
  odnośnik `javascript:` usunięty (zostaje tekst).
- Znalezisko do zapisania dla `gomdgo`: nagłówek YAML (`---`, `title: front`, `---`) wychodzi jako linia
  pozioma i nagłówek `<h2>`. To zgodne z CommonMark, ale nagłówki YAML są częste w notatkach i dokumentacji;
  do decyzji, czy biblioteka ma je pomijać albo oddawać osobno (opcja), czy robi to strona podglądu.
