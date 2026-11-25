# gopeekgo — instrukcje dla Claude Code

Projekt: podgląd Quick Look plików Markdown na macOS, renderowany bibliotekami rodziny: `gomdgo`
(Markdown → HTML) i `gomermaidgo` (bloki `mermaid` → SVG). Kod Go jest łączony z rozszerzeniem jako archiwum C,
więc podgląd powstaje w procesie rozszerzenia: bez JavaScriptu, bez uruchamiania programów, bez sieci.
Cel właściciela: testować nasze biblioteki na każdym pliku, jaki ma pod ręką, i mieć jednego z pierwszych
prawdziwych konsumentów wydanych wersji.

Stan (2026-11-25): **działający prototyp**, przeniesiony z `gomdgo/tools/quicklook`. Zainstalowany lokalnie
w `~/Applications/Gopeekgo.app`; sprawdzony na pliku z tabelą, listą zadań, blokiem uwagi i diagramem.
Zacznij od `README.md`, potem `docs/adr/0001-scope-and-build.md`, `docs/roadmap.md`, `docs/testing.md`.

## Układ

- `preview/`: `preview.Page(src, title)` buduje kompletną stronę HTML (czyste Go, testy bez macOS);
  `style.css` jest wbudowany przez `embed`. `Page` nie zwraca błędu: błąd diagramu trafia na stronę w miejscu
  diagramu, panika biblioteki zamienia się w stronę ze śladem stosu (to sygnał błędu do zgłoszenia w bibliotece).
- `bridge/`: `GopeekgoPreview` i `GopeekgoFree` eksportowane do C; deklaracje dla Swifta są w `macos/bridge.h`
  i muszą się zgadzać z `bridge/main.go`.
- `macos/`: `PreviewProvider.swift` (rozszerzenie, podgląd oparty na danych: oddaje systemowi gotowy HTML),
  `App.swift` (aplikacja-nośnik), pliki `Info.plist` i uprawnienia (rozszerzenie działa w piaskownicy).
- `build.sh`: `go build -buildmode=c-archive`, `swiftc`, podpis ad hoc, instalacja i rejestracja
  (`lsregister`, `pluginkit`). Bez projektu Xcode. Rejestracja bywa spóźniona, dlatego skrypt ją ponawia.

Identyfikatory paczek: `com.polestorm.gopeekgo` i `com.polestorm.gopeekgo.extension`. Obsługiwany typ pliku:
`net.daringfireball.markdown`.

## Zależności

- `go.mod` wskazuje wydane wersje: `gomdgo` v0.1.0 i `gomermaidgo` v0.1.0, bez `replace`. To celowe: projekt
  sprawdza publiczne API tak, jak widzi je odbiorca. Dyrektywa `go 1.26.0` wynika z `go.mod` w `gomdgo` v0.1.0
  (obniżenie do 1.22 nastąpiło tam po wydaniu).
- Żeby obejrzeć niewydaną zmianę biblioteki, użyj lokalnie `go.work` (`go work init . ../gomdgo`) i nie
  commituj go; do zrobienia: przełącznik w `build.sh` (`docs/roadmap.md`).
- Kolorowanie składni przyjdzie z `../gohighlightgo` jako opcja kompilacji: plik `preview/highlight.go` ze
  znacznikiem `//go:build highlight` i `preview/nohighlight.go` z pustą mapą; `build.sh` dostanie przełącznik.
  Po dołączeniu trzeba zmierzyć rozmiar rozszerzenia i czas w obu wariantach (kryterium H5 w
  `../gohighlightgo/docs/adr/0001-scope-and-acceptance.md`). Biblioteka jest dziś szkieletem.

## Jak pracujemy

- brak funkcji albo błąd renderowania naprawiamy w bibliotece (`../gomdgo`, `../gomermaidgo`), nie obchodzimy
  tutaj; tu jest tylko strona podglądu, arkusz stylów i część macOS. Uwaga: w `gomdgo` każda zmiana plików `*.go`
  pakietu `markdown/` zeruje licznik fuzzingu, więc znaleziska zbieraj w paczki i zapisuj jako uwagi;
- kodu innych podglądów i parserów Markdown nie czytamy (zasada rodowodu rodziny, `../gomdgo/docs/provenance.md`);
  konkurencję uruchamiamy jako czarną skrzynkę i mierzymy;
- liczby są datowane i podpisane sposobem pomiaru; porównania podajemy jako przewagę `gopeekgo` w procentach,
  `(konkurent − gopeekgo) / konkurent × 100%`; wartość ujemną też zapisujemy; `README.md` uaktualniamy razem
  z `docs/testing.md`;
- po każdej zmianie części macOS: `./build.sh install` i obejrzenie podglądu (`qlmanage -p plik.md`); w liście
  `pluginkit -m -A -p com.apple.quicklook.preview` znak `+` oznacza rozszerzenie włączone, `-` wyłączone;
- nazwy plików, pakietów i identyfikatorów po angielsku (także w Swifcie); komentarze i dokumentacja po polsku;
  `README.md` po angielsku; teksty widoczne dla użytkownika aplikacji po polsku (do decyzji: angielski
  przy publikacji);
- każda decyzja = plik w `docs/adr/`; licencja: Apache 2.0;
- repozytorium `github.com/polestorm/gopeekgo` (prywatne); moduły `github.com/polestorm/*` są w `GOPRIVATE`;
  commit i push tylko na prośbę właściciela;
- opisy commitów po angielsku; przed commitem muszą przechodzić `go vet ./...` i `go test ./...`;
- opis commita kończy się jedną linią `Made with the help of Claude Code` (po pustym wierszu); nie dopisuj
  `Co-Authored-By` ani `Claude-Session` (decyzja właściciela z 2026-11-24, przyjęta w `gomdgo`).
