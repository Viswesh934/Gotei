# EFFORT 3: CHAPTER DEFEAT

Date: 2026-04-10
Status: Completed (working CSS engine in current scope)

## Why this chapter exists
This chapter records the pivot from maintaining a custom CSS parser/selector stack to integrating WebRender parser/selector components while keeping our own layout and PDF rendering pipeline.

The goal was simple:
- stop rebuilding fragile CSS parsing behavior,
- adopt a stronger selector/cascade foundation,
- keep product velocity focused on output quality.

## The pivot
We moved from custom selector parsing/matching to:
- `github.com/benoitkugler/webrender/css/parser`
- `github.com/benoitkugler/webrender/css/selector`

And retained our internal pipeline:
- HTML parse -> CSS cascade -> layout tree -> PDF rendering.

## What changed
### 1) CSS parsing and selector matching
- Replaced custom CSS rule parsing with WebRender parser tokens and declaration parsing.
- Replaced custom selector matcher with WebRender matcher.
- Added selector specificity and source-order based cascading.

### 2) DOM compatibility for selector engine
- Added DOM mirroring to `golang.org/x/net/html` nodes so complex selector matching works correctly with tree relationships.

### 3) Cascade and inheritance stabilization
- Removed duplicate legacy inheritance/defaulting passes that were reintroducing defaults unexpectedly.
- Fixed inheritance semantics for text-facing properties when child values still matched defaults.

### 4) Rendering fidelity fixes
- Text alignment uses content box geometry (not outer box).
- Added `word-spacing` draw behavior.
- Added support for `text-align: start/end`.
- Added support for `text-transform: full-width`.
- Border radius now draws with rounded primitives.

### 5) Flexbox behavior hardening
- Fixed row centering issues by improving intrinsic width handling.
- Improved column cross-axis sizing/alignment behavior.
- Updated example fixtures to make flex-column behavior visually obvious.

### 6) Diagnostics and observability
- Added toggleable debug mode via `GOTEI_DEBUG=1`.
- Added logs for parse/cascade/layout/render stages and text draw coordinates.
- Added paint-stage color resolution logs to prove final RGB values.

## Critical bug defeated in this chapter
A deep initialization bug caused named colors (including `red`) to silently fail at render conversion:
- Theme JSON unmarshal failed due to type mismatch in spacing config (`-0.5` in JSON vs `[]int` in Go struct).
- Because loader init failed, `NamedColors` stayed empty.
- Cascade logs looked correct, but paint fell back to black.

Fix:
- Changed spacing schema to accept decimals (`[]float64`).
- Verified runtime color resolution now resolves named colors (basic and extended).

## Output path consistency
A source of confusion was duplicate output locations when running examples from different directories.
Current direction is to keep output generation deterministic and repo-scoped under `examples/output` (and corresponding example output folders).

## Result
We now have a working CSS engine for the current product scope:
- robust parser/selector foundation,
- stable cascade behavior,
- improved layout and render consistency,
- verified named color resolution,
- repeatable debug instrumentation.

## Remaining limits (known and accepted for now)
- Unit semantics for `%`, `em`, `rem` are still simplified in several paths.
- Custom font embedding is not implemented; font family mapping remains limited to PDF core fonts.
- Flexbox is practical but not full browser-spec parity for edge-heavy layouts.

## Chapter close
Effort 3, Chapter Defeat, represents the point where we stopped fighting CSS fundamentals and started shipping with a reliable integration base.

## Next step
The next chapter begins with real PDF output, not just style fixtures. The goal is to render documents that look like actual PDFs should look:
- tables with borders, spacing, and nested content
- images with sizing and placement
- multi-section documents with headings, paragraphs, lists, and page flow
- mixed content that exercises layout, paint, and pagination together

This is the ideal validation stage for the engine because it moves beyond isolated CSS cases and into complete document rendering.
