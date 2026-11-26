# 0004 Przygotowanie do upublicznienia: wersja, kolorowanie, porównanie, licencje

Status: przyjęta · 2026-11-26 (polecenie właściciela: ujednolicić repozytorium z publicznymi `gomdgo`,
`gomermaidgo`, `gopdfgo`, `gohighlightgo` przed upublicznieniem)

## Kontekst

Właściciel chce wkrótce upublicznić `gopeekgo`. Publiczne repozytoria rodziny mają: numer wersji w tabeli
porównania i w poleceniach, komentarze dokumentacyjne po angielsku (pokazuje je pkg.go.dev; `gohighlightgo`
v0.1.2), porównanie powtarzalne na zamrożonym wsadzie z opisaną metodą, mało plików w katalogu głównym,
licencję Apache 2.0 i `NOTICE` z listą cudzych składników. `gohighlightgo` wydało v0.1.2 i jest publiczne,
więc warunek z ADR 0003, punkt 4, dla dopisania zależności jest spełniony.

## Decyzja

1. **Wersja 0.1.0** w jednym miejscu: `preview.Version`. Pokazują ją stopka podglądu, `gopeekgo-render -version`
   i nagłówek tabeli w README; `CFBundleShortVersionString` w obu `macos/*-Info.plist` musi być taki sam
   (`TestVersionMatchesBundles`). Znacznika `v0.1.0` nie nadajemy bez decyzji właściciela (ADR 0003, punkt 3).
2. **Kolorowanie** z `gohighlightgo` v0.1.2 w `go.mod`, domyślnie włączone; adapter z
   `gohighlightgo/docs/adapters/gopeekgo.patch`. `macos/build.sh` przyjmuje `TAGS=nohighlight`.
3. **Komentarze dokumentacyjne po angielsku** w pakietach i poleceniach, które pokaże pkg.go.dev (`preview`,
   `bridge`, `cmd/gopeekgo-render`, `tools/bench`, `tools/features`); pozostałe komentarze, dokumentacja
   wewnętrzna i ADR zostają po polsku, jak w bibliotekach.
4. **Porównanie z wsadem i metodą**: czas, pamięć i rozmiar z `tools/bench` na `testdata/bench/` (doszedł wsad
   z kodem `code.md`, bo teraz kolorujemy), wiersze o funkcjach z `tools/features` na `testdata/features/`
   (jeden plik na funkcję, reguła rozpoznania w kodzie). Oba narzędzia zapisują ustawienia QLMarkdown, bo od
   nich zależy wynik. Opis po angielsku w `docs/benchmark.md`, wyniki w `docs/testing.md`.
5. **Medale** tylko tam, gdzie jedna wartość jest lepsza bez względu na upodobania czytelnika (czas, pamięć,
   rozmiar, obecność funkcji, brak skryptów, popularność, sposób dystrybucji); remisy i wiersze zależne od gustu
   (sposób rysowania diagramów, surowy HTML) bez medalu. Kolumna przewagi w procentach tylko przy liczbach.
6. **Licencje**: `NOTICE` wymienia biblioteki rodziny i czcionki Arimo i Noto Emoji (SIL OFL 1.1), które
   `gomermaidgo` wbudowuje w nasz plik wykonywalny; `macos/build.sh` kopiuje `LICENSE`, `NOTICE` i teksty OFL do
   `Gopeekgo.app/Contents/Resources`.
7. **README** krótkie, jak w `gohighlightgo`: bez sekcji „Status”; instalacja i użycie w `docs/usage.md`.

## Skutki

- W katalogu głównym zostają `LICENSE`, `NOTICE`, `README.md`, `go.mod`, `go.sum` i `.gitignore`. `go.sum` jest
  konieczny, bo mamy zależności (jak `gopdfgo`); mniej plików się nie da bez rezygnacji z modułów.
- Kolorowanie kosztuje 497 696 B pliku wykonywalnego i 0,7–1,1 ms na małym pliku (H5, `docs/testing.md`);
  przewaga rozmiaru pliku wykonywalnego wobec QLMarkdown spadła z −23% do −32%.
- Teksty aplikacji-nośnika są dalej po polsku; przy publikacji do decyzji (`CLAUDE.md`).
- `go install …@latest` z README zadziała dopiero po upublicznieniu repozytorium i nadaniu znacznika.
