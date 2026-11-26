# Code blocks

The 21 files of `tools/compare/inputs/corpus` in gohighlightgo v0.1.2, without their first lines, one fenced
code block per file, in the order of their paths.

## c/basics.txt

````c
#include <stdio.h>
#include "local.h"
  #  define MAX(a, b) ((a) > (b) ? (a) : (b))
#if defined(DEBUG) && DEBUG > 1
#endif

/* Counts characters.
   Not nested: /* stays inside */
static const char *names[] = {"a\tb", u8"zażółć", L"wide", "hex \x41\x4142 oct \101\0 uni \u00e9 \U0001F600"};

typedef struct point { int32_t x, y; } point;

int main(int argc, char **argv) {
    size_t n = sizeof(names) / sizeof *names;
    unsigned long long big = 18'446'744'073'709'551'615ULL;
    double d = 1.5e-3 + .5f + 0x1.8p1 + 0b1010 + 017;
    char c = '\'', q = 'x', w = L'\n';
    bool ok = true && !false;
    void *p = nullptr;
    if (argc > 1 && p == NULL) {
        printf("%zu %c\n", n, c);
    }
    const char *s = "continued \
on the next line";
    for (int i = 0; i < 10; ++i) { big >>= 1; big ^= ~0u; }
    return ok ? 0 : -1;
}
````

## cpp/basics.txt

````cpp
#include <vector>
#pragma once

namespace demo {

template <typename T>
concept Number = requires(T a) { a + a; };

class Widget final : public Base {
public:
    explicit Widget(int n) noexcept : n_(n) {}
    ~Widget() override = default;
    auto size() const -> std::size_t { return static_cast<std::size_t>(n_); }
    bool operator==(const Widget&) const = default;
private:
    int n_ = 1'000;
};

}  // namespace demo

int main() {
    auto raw = R"(a "quoted" \n stays)";
    auto delim = u8R"xy(
multi-line )" still inside
)xy";
    std::vector<char32_t> v{U'ż', u'a', L'b'};
    if (v.empty() and not false) co_return;
    throw std::runtime_error("bad\n");
}
````

## css/styles.txt

````css
@import url("theme.css");
@charset "utf-8";

:root {
  --accent: #0a84ff;
  --gap: 1.5rem;
}

body, .page > main + aside ~ footer {
  margin: 0 auto !important;
  font: italic bold 12px/30px Georgia, serif;
  color: rgb(10 20 30 / 50%);
  background: url(img/bg.png) no-repeat, var(--accent);
  width: calc(100% - 2 * var(--gap));
  content: "\201C quote \"x\"";
  transition: opacity .3s ease-in-out, -webkit-transform 0.5s;
  z-index: -1;
}

a:hover, a::after, #header .nav li:nth-child(2n+1) {
  display: inline-block;
}

@media screen and (max-width: 600px) {
  .grid { grid-template-columns: repeat(2, 1fr); }
}
````

## diff/patch.txt

````diff
diff --git a/highlight/table.go b/highlight/table.go
index 1234567..89abcde 100644
--- a/highlight/table.go
+++ b/highlight/table.go
@@ -10,7 +10,8 @@ func scan(w *writer, src []byte) {
 	for i := 0; i < len(src); {
-		b := src[i]
+		b := src[i] // current byte
+		next := byte(0)
--- this deleted line starts with two dashes
+++ and this added line with two pluses
 		if isSpace(b) {
\ No newline at end of file
new file mode 100644
--- /dev/null
+++ b/docs/new.md
@@ -0,0 +1 @@
+# New
2c2
< old line
---
> new line
Binary files a/logo.png and b/logo.png differ
````

## go/basics.txt

````go
// Package demo shows the basics.
package demo

import (
	"fmt"
	"strings"
)

/* A block comment
   over two lines. */
type Point struct {
	X, Y int
	Name string
}

func (p *Point) Scale(k float64) error {
	if p == nil || k < 0 {
		return fmt.Errorf("bad scale %v", k)
	}
	p.X = int(float64(p.X) * k)
	p.Y *= 2
	return nil
}

func Equal[T comparable](a, b T) bool { return a == b }

func main() {
	pts := make([]Point, 0, 4)
	pts = append(pts, Point{X: 1, Y: 2})
	for i, p := range pts {
		go func() { defer recover(); _ = p.Scale(1.5) }()
		fmt.Println(i, len(pts), strings.ToUpper (p.Name))
	}
	var ch chan<- struct{}
	select {
	case ch <- struct{}{}:
	default:
	}
	switch x := any(true); x.(type) {
	case bool:
		fallthrough
	default:
		break
	}
	const c = iota
	goto end
end:
	var m map[string]interface{}
	_, _ = m, c
	clear(m)
	println(min(1, 2), max(3, 4), cap(pts), new(uintptr), real(complex(1, 2)), imag(2i))
	copy(pts, pts)
	delete(m, "k")
	close(make(chan error))
	if false || true {
		panic(nil)
	}
	type T = rune
	var _ = T(0) &^ 3
	var _ byte = ^uint8(1) % 7
	print(^0 | 1<<2 >> 1)
	continue
}
````

## go/edge.txt

````go
package edge

func zażółć() {}

func Größe(ä int) bool { return ä <= 1 && ä >= 0 || ä != 2 }

func html() string { return "<script>alert('x') & </script>" + `<b>"quoted"</b>` }

var _ = a<b && c>d // <!-- not a comment --> &amp;

var s = "unterminated
var r = 'x
x := Größe (1); y := zażółć	()
z := f() + g () - h
	()
/* unterminated comment
until the end of the file <&>
````

## go/literals.txt

````go
package literals

var ints = []int{0, 42, 4_2, 0600, 0_600, 0o600, 0O600, 0xBadFace, 0xBad_Face, 0x_67_7a_2f_cc_40_c6, 0b1011, 0B_1_0}

var floats = []float64{0., 72.40, 072.40, 2.71828, 1.e+0, 6.67428e-11, 1E6, .25, .12345E+5, 1_5., 0.15e+0_2, 0x1p-2, 0x2.p10, 0x1.Fp+0, 0X.8p-0, 0X_1FFFP-16}

var imaginary = []complex128{0i, 0123i, 0o123i, 0xabci, 0.i, 2.71828i, 1.e+0i, 6.67428e-11i, 1E6i, .25i, .12345E+5i, 0x1p-2i}

var runes = []rune{'a', 'ä', '本', '\t', '\000', '\007', '\377', '\x07', '\xff', '\u12e4', '\U00101234', '\'', '\\'}

var strs = []string{
	"",
	"Hello, world!\n",
	"日本語",
	"\u65e5本\U00008a9e",
	"\xff\u00FF",
	"quote \" and backslash \\ and \q",
	`raw \n stays`,
	`multi
line raw string`,
}

var x = a[1:2] + b.c - 3.0*d

// Comments do not nest: /* inside a line comment is plain text.
/* outer /* inner */ x := 1 /* after */
/* multi-line
block comment */
````

## html/page.txt

````html
<!DOCTYPE html>
<html lang=pl>
<head>
  <meta charset="utf-8">
  <title>Zażółć &mdash; test</title>
  <style>
    body { color: #333; margin: 0 auto; }
    /* </b> is not the end */
  </style>
  <script type="module">
    const html = "<b>" + '</style>';
    if (a < b && c > d) { console.log(`${a}`); }
  </SCRIPT>
</head>
<body class="main" data-x='1' hidden>
  <p>Text with <a href="https://example.com/?a=1&amp;b=2">a link</a> &amp; more.</p>
  <img src=logo.png alt="">
  <br/>
</body>
</html>
````

## java/basics.txt

````java
package demo;

import java.util.List;

/** Javadoc comment. */
@FunctionalInterface
public sealed interface Shape permits Circle {
    double area();
}

public record Circle(double r) implements Shape {
    @Override
    public double area() { return Math.PI * r * r; }
}

class Main {
    private static final long BIG = 1_000_000L;
    public static void main(String[] args) throws Exception {
        var $name = "tab\t octal \101 unicode \u0041 \uuuu0042";
        char c = '\'';
        float f = .5f + 1e3f + 0x1.0p3;
        String block = """
            Hello, "block"
            end\
            """;
        boolean ok = args.length > 0 ? true : false;
        Object o = null;
        if (o instanceof String s && !ok) {
            System.out.println(s + $name + c + f + block);
        }
        int x = switch (args.length) { case 0 -> 1; default -> { yield 2; } };
    }
}
````

## javascript/basics.txt

````javascript
import { readFile } from "node:fs/promises";

/* Block comment */
export default class Store extends Base {
    static #count = 0;
    @logged
    async load(path) {
        const data = await readFile(path, 'utf8');
        let $el = document.querySelector(`#item-${this.id}`);
        const msg = `Hello ${user.name}, you have ${count > 1 ? `${count} items` : "one item"} {not code}
second line \` and \${not}`;
        const obj = { a: { b: 1 } };
        const t = `${ { a: 1 }.a }`;
        const big = 9007199254740993n + 0x1Fn + 0b1010 + 0o17 + 1_000 + .5e-3;
        if (data === undefined || data == null || Number.isNaN(NaN) || x !== Infinity) {
            throw new Error("bad \x41 \u0041 \u{1F600}");
        }
        const re = /'not a string/; // a regular expression after =
        const half = (a + b) / 2 / c, ratio = arr[0] / x.length;
        total /= 2;
        if (/^[a-z/]+$/i.test(path) && [/x/, /\//g].length) {}
        const after = typeof /y/ === 'object' ? /z/u : null;
        return typeof data === 'string' ? data : void 0;
    }
}
````

## json/basics.txt

````json
{
  "name": "gohighlightgo",
  "version" : 1.5e3,
  "negative": -42,
  "flags": [true, false, null],
  "escapes": "tab\t quote\" uni \u00e9 slash\/",
  /* JSONC comment */
  "nested": {"a": [1, 2, {"b": "c"}]},
  "key with \"quotes\"": "value"
}
````

## markdown/readme.txt

`````markdown
# gohighlightgo

Setext heading
==============

Text with `code`, ``code with ` inside``, an unmatched ` backtick, a [link](https://example.com/a)
and <https://example.com> or <user@example.com>, escaped \*stars\* and &amp; entity.

> Quote with **bold** text
> > nested quote

- item one
* item two
  1. ordered
  2) other
- [ ] task
- [x] done task

[ref]: https://example.com "Title"

---

```go
func main() { fmt.Println("hi") }
```

~~~
plain code block <&>
~~~

````markdown
```go
not highlighted as Go inside Markdown
```
````

<div class="note">HTML block</div>

```unknown
kept as a string
`````

## pascal/basics.txt

````pascal
unit Demo;

{$IFDEF DEBUG}
{$R *.res}
{$ENDIF}

interface

uses
  System.SysUtils, System.Classes;

type
  TShape = class abstract(TObject)
  strict private
    FName: string;
  public
    constructor Create(const AName: string); virtual;
    function Area: Double; virtual; abstract;
    property Name: string read FName write FName;
  end;

implementation

{ Block comment }
(* Another block comment *)
(*$WARNINGS OFF*)

constructor TShape.Create(const AName: string);
var
  I: Integer;
  S: UnicodeString;
begin
  inherited Create;
  FName := AName + 'It''s' + #13#10 + #$0D + 'end';
  for I := 1 to 10 do
    Inc(I, $FF + %1010 + &777 + 1_000 + 1.5E3);
  if Assigned(Self) and not (Length(S) > 0) then
    WriteLn(Format('%d', [I]))
  else
    Exit;
  case I of
    1..10: Result := nil;
  end;
  BEGIN End;
end;

end.
````

## python/basics.txt

````python
"""Module docstring
over two lines."""
from dataclasses import dataclass
import re


@dataclass(frozen=True)
class Point:
    x: int = 0
    y: float = 1.5e-3

    def norm(self) -> float:
        return (self.x ** 2 + self.y ** 2) ** .5


async def main(*args, **kwargs):
    p = Point(1, 2)
    names = [str(n) for n in range(10) if n % 2 == 0]
    raw = r"C:\new\dir\"" + R'\d+' + rb"\x00"
    text = f"{p.x!r:>10} done\n" + b'bytes\x41' + u"uni \u00e9 \N{BULLET} \U0001F600"
    big = 1_000_000 + 0x_FF + 0o17 + 0b1010 + 3.14j
    if p is not None and isinstance(p, Point) or not True:
        print(len(names), sorted(names), sep=", ")
    lam = lambda x: x @ x
    try:
        await other()
    except (ValueError, TypeError) as err:
        raise RuntimeError("bad") from err
    s = 'line continued \
here'
    doc = '''single
quotes'''
    return NotImplemented
````

## rust/basics.txt

````rust
#![allow(dead_code)]
use std::collections::HashMap;

/* outer /* nested */ still a comment */

#[derive(Debug, Clone)]
pub struct Point<'a> {
    name: &'a str,
    x: f64,
}

impl<'a> Point<'a> {
    pub fn new(name: &'a str) -> Self {
        Self { name, x: 1.5e3 }
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut map: HashMap<String, Vec<u8>> = HashMap::new();
    let c = 'a';
    let nl = '\n';
    let quote = '\'';
    let b = b'x';
    let ok = Some(1u32).is_some() && true;
    let raw = r#"raw "string" with \n"#;
    let raw2 = br##"bytes "#" inside"##;
    let s = "multi
line \u{1F600} \x41";
    let bytes = b"abc\0";
    for i in 0..10 {
        let t = (i, 2);
        println!("{} {}", t.0, 1.max(2));
    }
    'outer: loop {
        if map.len() != 0 { break 'outer; }
        break;
    }
    let n = 0xff_u8 + 0o7 + 0b1 + 1_000i64 as u8;
    Ok(())
}
````

## sh/script.txt

````sh
#!/usr/bin/env bash
set -euo pipefail

NAME="world" COUNT=3 ./run.sh --verbose
readonly DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

greet() {
    local who=${1:-$NAME}
    echo "Hello, $who! Exit code: $? in \$HOME=$HOME" 'single $not' `date`
}

if [[ -f "$DIR/config" ]]; then
    source "$DIR/config"
elif test -d /tmp; then
    grep -r 'pattern' . | sort | uniq -c > out.txt 2>&1
else
    exit 1
fi

for f in *.txt; do
    wc -l "$f" && echo done
done

cat <<EOF
Here document with $NAME
EOF
cat <<-'END'
	literal $NAME
	END
count=$((COUNT + 1)) # arithmetic
echo escaped\ space
````

## sql/basics.txt

````sql
/* Block comment /* nested */ still comment */
CREATE TABLE IF NOT EXISTS "order items" (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL DEFAULT 'it''s',
    price DECIMAL(10, 2) CHECK (price >= 0),
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

select o.id, count(*) as n, coalesce(sum(o.price), 0.0) AS total, my_func(o.name)
from "order items" o
left join `select` s on s.id = o.id
where o.name like N'%zażółć%' and o.price between 1 and 1e3 or o.flag is null
group by o.id having count(*) > 1
order by total desc
limit 10 offset 20;

insert into t (a, b) values (X'FF', 'multi
line');
````

## toml/cargo.txt

````toml
title = "TOML example"
"quoted key" = 'literal \n stays'

[owner]
name = "Tom Preston-Werner"
dob = 1979-05-27T07:32:00-08:00
time = 07:32:00

[database.servers]
enabled = true
ports = [ 8000, 8001, 8002 ]
temp_targets = { cpu = 79.5, case = 72.0 }
site."google.com" = true

[[products]]
hex = 0xDEAD_BEEF
oct = 0o755
bin = 0b1101
float = -3.1415e-2
special = [ inf, -inf, nan ]
desc = """
Multi-line with "quotes" and \t escape \u00e9.
"""
raw = '''
C:\Users\no\escapes
'''
````

## typescript/basics.txt

````typescript
interface Point { readonly x: number; y?: number }

type Shape = { kind: "circle"; r: number } | { kind: "square"; a: bigint };

declare module "lib" {
    export function f<T extends keyof Point>(k: T): Point[T];
}

namespace Geometry {
    export abstract class Base {
        protected abstract area(): number;
    }
}

enum Color { Red = 1, Green }

function isString(x: unknown): x is string {
    return typeof x === "string";
}

const conf = { debug: true } satisfies Record<string, boolean>;
let n = (value as any) as never;
const label = `Point ${p.x}`;
````

## xml/jpk.txt

````xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE note SYSTEM "note.dtd">
<tns:JPK xmlns:tns="http://jpk.mf.gov.pl/wzor/2022/02/17/02171/"
         xmlns:etd='http://crd.gov.pl/xml/schematy/dziedzinowe/mf/2022/01/05/eD/DefinicjeTypy/'>
  <tns:Naglowek>
    <tns:KodFormularza kodSystemowy="JPK_V7M (2)" wersjaSchemy="1-0E">JPK_VAT</tns:KodFormularza>
    <tns:DataWytworzeniaJPK>2024-01-31T10:00:00</tns:DataWytworzeniaJPK>
  </tns:Naglowek>
  <tns:Opis>Zażółć &amp; gęślą &#261; &#x105; jaźń &lt;tag&gt; & a < b</tns:Opis>
  <tns:P_15>1230.00</tns:P_15>
  <empty attr="a &quot;quoted&quot; value"/>
  <![CDATA[ <not> a & tag ]]>
  <!-- multi-line
       comment -->
</tns:JPK>
````

## yaml/config.txt

````yaml
%YAML 1.2
---
name: gohighlightgo
version: 1.5
count: 42
hex: 0x1F
octal: 0o17
infinite: -.inf
empty: ~
enabled: true
disabled: False
nothing: null
date: 2024-01-31
url: http://example.com:8080/path # comment after a value
"quoted key": "line\n with \"escape\" and \u00e9"
'single': 'it''s fine'
list:
  - first
  - second: nested
  - [a, b, {c: 1, d: "x"}]
? complex key
: complex value
anchor: &base
  a: 1
merged: *base
tagged: !!str 123
custom: !mytag value
literal: |
  line one
    indented line

  after blank line
folded: >-
  folded text
  continues
next: done
multi: "first
  second"
...
````
