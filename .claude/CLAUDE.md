# gopeekgo — instrukcje dla Claude Code

Projekt: podgląd Quick Look plików Markdown na macOS, renderowany bibliotekami rodziny: `gomdgo`
(Markdown → HTML), `gomermaidgo` (bloki `mermaid` → SVG) i `gohighlightgo` (kolorowanie bloków kodu). Kod Go jest łączony z rozszerzeniem jako archiwum C,
więc podgląd powstaje w procesie rozszerzenia: bez JavaScriptu, bez uruchamiania programów, bez sieci.
Cel właściciela: testować nasze biblioteki na każdym pliku, jaki ma pod ręką, i mieć jednego z pierwszych
prawdziwych konsumentów wydanych wersji.

Stan (2026-11-26): **wersja 0.1.0 przygotowana do upublicznienia** (ADR 0004), bez znacznika i jeszcze
prywatna. Kolorowanie z `gohighlightgo` v0.1.2 włączone; polecenie `cmd/gopeekgo-render` daje tę samą stronę
co rozszerzenie. Zainstalowany lokalnie w `~/Applications/Gopeekgo.app`.
Zacznij od `README.md`, potem `docs/adr/0001-scope-and-build.md`, `docs/roadmap.md`, `docs/testing.md`.
Układ repozytorium i README trzymają schemat publicznych bibliotek rodziny (ADR 0002): w korzeniu tylko
`LICENSE`, `NOTICE`, `README.md`, `go.mod`, `go.sum`, `.gitignore`; nowe pliki idą do katalogów.

## Układ

- `preview/`: `preview.Page(src, title)` buduje kompletną stronę HTML (czyste Go, testy bez macOS);
  `style.css` jest wbudowany przez `embed`. `Page` nie zwraca błędu: błąd diagramu trafia na stronę w miejscu
  diagramu, panika biblioteki zamienia się w stronę ze śladem stosu (to sygnał błędu do zgłoszenia w bibliotece).
- `bridge/`: `GopeekgoPreview` i `GopeekgoFree` eksportowane do C; deklaracje dla Swifta są w `macos/bridge.h`
  i muszą się zgadzać z `bridge/main.go`.
- `macos/`: `PreviewProvider.swift` (rozszerzenie, podgląd oparty na danych: oddaje systemowi gotowy HTML),
  `App.swift` (aplikacja-nośnik), pliki `Info.plist` i uprawnienia (rozszerzenie działa w piaskownicy).
- `macos/build.sh`: `go build -buildmode=c-archive`, `swiftc`, podpis ad hoc, instalacja i rejestracja
  (`lsregister`, `pluginkit`). Bez projektu Xcode. Rejestracja bywa spóźniona, dlatego skrypt ją ponawia.
  `TAGS=nohighlight` buduje bez kolorowania; `LICENSE`, `NOTICE` i teksty OFL czcionek trafiają do `Resources`.
- `cmd/gopeekgo-render`: ta sama strona co rozszerzenie, z wiersza poleceń (pomiary, oglądanie w przeglądarce).
- `tools/bench`, `tools/features`: powtarzalne porównanie z QLMarkdown (`docs/benchmark.md`).
- Wersja: `preview.Version`; `CFBundleShortVersionString` w `macos/*-Info.plist` musi być taki sam (test).
- Komentarze dokumentacyjne pakietów i poleceń po angielsku (pkg.go.dev, ADR 0004); reszta komentarzy po polsku.

Identyfikatory paczek: `com.polestorm.gopeekgo` i `com.polestorm.gopeekgo.extension`. Obsługiwany typ pliku:
`net.daringfireball.markdown`.

## Zależności

- `go.mod` wskazuje wydane wersje: `gomdgo` v0.1.1, `gomermaidgo` v0.1.0 i `gohighlightgo` v0.1.2, bez `replace`. To celowe: projekt
  sprawdza publiczne API tak, jak widzi je odbiorca. Wersji `gomdgo` v0.1.0 nie używamy: jej znacznik został
  przesunięty po wydaniu i ma dwie różne treści (`../gomdgo/docs/adr/0017-version-0.1.1-and-immutable-tags.md`).
- Żeby obejrzeć niewydaną zmianę biblioteki, użyj lokalnie `go.work` (`go work init . ../gomdgo`) i nie
  commituj go; do zrobienia: przełącznik w `macos/build.sh` (`docs/roadmap.md`).
- Kolorowanie składni z `../gohighlightgo` jest **domyślnie włączone** (opt-out, ADR 0003): `preview/highlight.go`
  z `//go:build !nohighlight`, `preview/nohighlight.go` z `//go:build nohighlight`. Testy uruchamiaj w obu
  wariantach (`go test ./...` i `go test -tags nohighlight ./...`). H5 zmierzone 2026-11-26 (`docs/testing.md`);
  po każdym nowym wydaniu `gohighlightgo` (np. kolejne języki): `go get …@vX`, testy, `tools/bench` w obu wariantach.
- **Znacznik wersji i zmiana widoczności repozytorium tylko na decyzję właściciela** (ADR 0003, punkt 3; ADR 0004).

## Jak pracujemy

- brak funkcji albo błąd renderowania naprawiamy w bibliotece (`../gomdgo`, `../gomermaidgo`), nie obchodzimy
  tutaj; tu jest tylko strona podglądu, arkusz stylów i część macOS. Uwaga: w `gomdgo` każda zmiana plików `*.go`
  pakietu `markdown/` zeruje licznik fuzzingu, więc znaleziska zbieraj w paczki i zapisuj jako uwagi;
- kodu innych podglądów i parserów Markdown nie czytamy (zasada rodowodu rodziny, `../gomdgo/docs/provenance.md`);
  konkurencję uruchamiamy jako czarną skrzynkę i mierzymy;
- **porównania muszą być powtarzalne** (reguła właściciela z 2026-11-26): każda liczba i każdy wiersz porównania
  z konkurencją w `README.md` ma wskazany wsad i metodę, a pomiar da się powtórzyć jednym poleceniem na plikach
  z repozytorium. Czas, pamięć i rozmiar: `go run ./tools/bench` na `testdata/bench/`, opis w `docs/benchmark.md`,
  wyniki w `docs/testing.md`. Pomiarów doraźnych nie publikujemy. Wiersze o funkcjach: `go run ./tools/features` na `testdata/features/` (jeden plik na funkcję);
- liczby są datowane i podpisane sposobem pomiaru; porównania podajemy jako przewagę `gopeekgo` w procentach,
  `(konkurent − gopeekgo) / konkurent × 100%`; wartość ujemną też zapisujemy; `README.md` uaktualniamy razem
  z `docs/testing.md`;
- po każdej zmianie części macOS: `macos/build.sh install` i obejrzenie podglądu (`qlmanage -p plik.md`); w liście
  `pluginkit -m -A -p com.apple.quicklook.preview` znak `+` oznacza rozszerzenie włączone, `-` wyłączone;
- nazwy plików, pakietów i identyfikatorów po angielsku (także w Swifcie); komentarze i dokumentacja po polsku;
  `README.md` po angielsku; teksty widoczne dla użytkownika aplikacji po polsku (do decyzji: angielski
  przy publikacji);
- każda decyzja = plik w `docs/adr/`; licencja: Apache 2.0;
- repozytorium `github.com/polestorm/gopeekgo` (do upublicznienia, ADR 0004); `gomdgo`, `gomermaidgo`
  i `gohighlightgo` są publiczne;
  commit i push tylko na prośbę właściciela;
- opisy commitów po angielsku; przed commitem muszą przechodzić `go vet ./...` i `go test ./...`;
- opis commita kończy się jedną linią `Made with the help of Claude Code` (po pustym wierszu); nie dopisuj
  `Co-Authored-By` ani `Claude-Session` (decyzja właściciela z 2026-11-24, przyjęta w `gomdgo`).
