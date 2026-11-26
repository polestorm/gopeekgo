# 0001 Zakres i sposób budowania

Status: przyjęta · 2026-11-25 (polecenie właściciela: osobny projekt zależny od `gomdgo` i `gomermaidgo`,
nazwa `gopeekgo`, licencja Apache 2.0)

## Kontekst

Właściciel chce oglądać dowolne pliki Markdown wynikiem `gomdgo` i `gomermaidgo` wprost w Finderze (spacja).
Do tej pory robiło to obce rozszerzenie (QLMarkdown) z innym parserem. macOS przyjmuje rozszerzenie Quick Look
tylko jako część aplikacji, napisane natywnie i działające w piaskownicy. Prototyp powstał w
`gomdgo/tools/quicklook` i został tu przeniesiony tego samego dnia, bez commita tam.

Rozważone sposoby połączenia rozszerzenia z kodem w Go:

1. **Uruchamianie programu `gomdgo-mermaid` z rozszerzenia.** Odrzucone: piaskownica utrudnia uruchamianie
   programów, a rodzina ich unika (`gomdgo`, ADR 0014).
2. **Archiwum C z Go łączone z rozszerzeniem** (`go build -buildmode=c-archive`). Wybrane.

## Decyzja

1. Osobne repozytorium i moduł `github.com/polestorm/gopeekgo`, zależne od wydanych wersji bibliotek, bez
   `replace`. Projekt jest konsumentem publicznego API, nie narzędziem wewnętrznym `gomdgo`.
2. `preview.Page(src, title)` buduje stronę w czystym Go i nie zwraca błędu. Konfiguracja `gomdgo`: GFM, bloki
   uwag, przypisy, identyfikatory nagłówków; surowy HTML zapisany jako tekst (polityka domyślna biblioteki).
3. Rozszerzenie oddaje systemowi gotowy HTML (podgląd oparty na danych, `QLIsDataBasedPreview`), bez własnego
   widoku i bez dostępu do sieci.
4. Budowanie skryptem `build.sh`, bez projektu Xcode; podpis lokalny ad hoc; instalacja do `~/Applications`.
5. Kolorowanie składni przyjdzie z biblioteki `gohighlightgo`. Pierwotnie: opcja dołączana znacznikiem
   `highlight`; zmienione w ADR 0003 na domyślnie włączone (znacznik `nohighlight` wyłącza).

## Skutki

- Obrazy o ścieżkach względnych nie są wczytywane; do zrobienia jako załączniki odpowiedzi podglądu.
- Budowane dla architektury komputera, na którym uruchomiono `build.sh`; rozpowszechnianie wymaga podpisu
  Developer ID, notaryzacji i binarki uniwersalnej (osobna decyzja).
- System wybiera jedno rozszerzenie na typ pliku; inne rozszerzenie dla Markdown trzeba wyłączyć.
- CI nie jest założone: testy `preview` potrzebują dostępu do prywatnych modułów rodziny.
