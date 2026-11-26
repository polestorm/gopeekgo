# 0003 Kolorowanie składni domyślnie włączone; wydanie czeka na gohighlightgo

Status: przyjęta · 2026-11-25 (decyzja właściciela; zmienia punkt 5 ADR 0001)

## Kontekst

ADR 0001 zakładał kolorowanie jako opcję dołączaną przy kompilacji (znacznik `highlight`). Podgląd ogląda
człowiek, jego HTML nie trafia do dalszego przetwarzania, a konkurencja (QLMarkdown) koloruje domyślnie.
Biblioteka `gohighlightgo` jest w budowie w osobnym repozytorium.

## Decyzja

1. Kolorowanie składni jest w `gopeekgo` **domyślnie włączone** (opt-out). Zwykłe budowanie dołącza
   `gohighlightgo`; znacznik `nohighlight` je wyłącza: `preview/highlight.go` ma `//go:build !nohighlight`,
   a `preview/nohighlight.go` (`//go:build nohighlight`) daje pustą mapę.
2. Skrypt budujący dostaje przełącznik, który przekazuje `-tags nohighlight`. Służy też do pomiaru: rozmiar
   rozszerzenia i czas konwersji mierzymy w obu wariantach (kryterium H5 w `gohighlightgo`).
3. **Wydanie `gopeekgo` czeka, aż `gohighlightgo` będzie dojrzałe** (spełnione jego kryteria przyjęcia i wydana
   wersja). Do tego czasu `gopeekgo` zostaje prototypem budowanym lokalnie, bez znacznika wersji.
4. `go.mod` wskazuje wyłącznie wydane wersje bibliotek (ADR 0001, punkt 1), więc zależność od `gohighlightgo`
   dopisujemy dopiero po jego pierwszym wydaniu. Wcześniejsze próby robimy lokalnie przez `go.work`, bez commita.

## Skutki

- Pozostali odbiorcy (`gomdgo/cmd/gomdgo-mermaid`, `gopressgo`) zostają przy dołączaniu na życzenie; ta decyzja
  dotyczy tylko `gopeekgo`. Sama biblioteka `gomdgo` nie zna `gohighlightgo`.
- ADR 0001 w `gohighlightgo` opisuje u odbiorców znacznik `highlight`; dla `gopeekgo` obowiązuje ten dokument.
