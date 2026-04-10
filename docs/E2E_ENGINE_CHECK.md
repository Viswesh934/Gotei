# E2E Engine Check (HTML/CSS -> Layout -> PDF)

Date: 2026-04-10

## Pipeline
1. HTML parse: `internal/dom/parser.go` (`golang.org/x/net/html`)
2. CSS parse + selector match: `internal/css/parser.go`, `internal/css/cascade.go` (WebRender parser/selector)
3. Computed styles to layout tree: `internal/css/tree.go`, `internal/layout/tree.go`
4. Layout: `internal/layout/layout.go`, `internal/layout/flexbox.go`
5. PDF paint: `internal/render/pdf.go`, `internal/render/borders.go`

## Status by property family
- Font family: Partial. Mapped to PDF core fonts (`Arial`, `Times`, `Courier`). Custom families not embedded yet.
- Font size/weight/style: Supported.
- Word spacing: Supported in render pass (letter + word spacing path).
- Text align: Supported (`left`, `center`, `right`, `justify`, `start`, `end`).
- Text decoration: Basic lines supported (`underline`, `line-through`, `overline`).
- Text transform: Supported (`uppercase`, `lowercase`, `capitalize`, `full-width`).
- Basic named colors: Supported.
- Extended named colors: Supported via `css_theme.json` load.
- Hex color formats: Supported for `#RGB`, `#RGBA`, `#RRGGBB`, `#RRGGBBAA` (alpha ignored in current PDF color path).
- Text shadow: Supported (no blur kernel; offset + colorized draw).
- Border radius: Supported in border drawing via `RoundedRect`.
- Padding models: Supported in style parse + layout.
- Flexbox: Partial implementation; common row/column and justify/align cases supported, not full spec parity.

## Known consistency limits
- CSS units are parsed numerically with simplified semantics; `%`/`em`/`rem` are not fully resolved against CSS context.
- No custom font embedding yet, so many font families collapse to core PDF fonts.
- Flexbox algorithm is simplified and may differ from browser results on complex/nested cases.

## Next hardening steps
1. Introduce unit-aware value model (store value + unit, resolve at layout time).
2. Add optional custom font registration/embedding for true font-family fidelity.
3. Add golden tests for key example outputs per property bucket.
4. Extend flexbox test matrix with nested and mixed-size containers.
