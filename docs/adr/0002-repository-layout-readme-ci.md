# 0002 Układ repozytorium, README i CI na wzór bibliotek rodziny

Status: przyjęta · 2026-11-25 (polecenie właściciela: README i układ jak w publicznych `gomdgo`,
`gomermaidgo`, `gopdfgo`)

## Kontekst

Publiczne repozytoria rodziny mają wspólny schemat. W korzeniu są tylko `LICENSE`, `NOTICE`, `README.md`,
`go.mod` (i `go.sum`, gdy są zależności) oraz `.gitignore`; reszta leży w katalogach. Każde ma
`.github/workflows/test.yml`. README po angielsku ma stały porządek: kilka zdań o tym, co projekt robi, krótki
przykład, „Why …” z wypunktowanymi przewagami i akapitem o słabościach, spis treści, tabela porównania
z jednym–trzema konkurentami (🥇 przy najlepszej wartości, przypisy ze sposobem pomiaru i datą), „Status”,
„Documentation”, „How the comparison was measured”, „License”.

ADR 0001 odłożył CI, bo moduły rodziny były prywatne. `gomdgo` i `gomermaidgo` są dziś publiczne;
2026-11-25 sprawdzono, że `go mod download` przez `proxy.golang.org` i `sum.golang.org` daje sumy zgodne
z `go.sum`, a `go vet ./...` i `go test ./...` przechodzą bez dostępu prywatnego.

## Decyzja

1. `build.sh` przenosimy do `macos/build.sh` (buduje część macOS; skrypt przechodzi do korzenia repozytorium,
   więc katalog `build/` zostaje tam, gdzie był).
2. README według schematu rodziny, z jednym konkurentem (QLMarkdown). Ponieważ zasada projektu każe podawać
   przewagę w procentach, tabela ma dodatkową kolumnę „gopeekgo advantage” przy wierszach liczbowych. Wiersze
   o funkcjach pochodzą z próby na jednym pliku obiema programami (`docs/testing.md`), nie z dokumentacji.
3. CI jak w bibliotekach: Go 1.22 i najnowsze, Linux i macOS; `gofmt`, `go vet`, `go test`, a na macOS także
   zbudowanie aplikacji (`macos/build.sh` bez instalacji).

## Skutki

- Punkt ADR 0001 o braku CI jest nieaktualny.
- Polecenia w dokumentacji brzmią `macos/build.sh install` zamiast `./build.sh install`.
- Porównanie z innymi podglądami (`docs/roadmap.md`, punkt 5) dopisuje kolejne kolumny tej samej tabeli.
