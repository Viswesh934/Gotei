# Gotei

**Gotei** is a lightweight HTML → PDF rendering engine written in Go.
It focuses on **deterministic document layout** instead of relying on heavy headless browsers.

---

## 🚀 Why Gotei?

Most HTML-to-PDF solutions depend on browser engines (like Chromium), which are:

* Heavy
* Slow to start
* Hard to scale

Gotei takes a different approach:

> Parse HTML → Build layout → Render directly to PDF

---

## ⚙️ Features (Current)

* ✅ HTML → DOM parsing
* ✅ Basic layout engine (vertical flow)
* ✅ Text wrapping
* ✅ PDF rendering
* ✅ Simple HTTP API
* ✅ Docker-ready

---

## 📦 Project Structure

```
gotei/
├── cmd/server          # HTTP server entrypoint
├── internal/
│   ├── dom/            # HTML parsing
│   ├── layout/         # layout engine
│   ├── render/         # PDF rendering
│   └── engine/         # orchestration
├── pkg/api             # request/response types
├── plugins/            # future extensions
├── assets/fonts        # fonts (future)
├── testdata            # sample inputs
├── Dockerfile
└── README.md
```

---

## 🧪 Run Locally

```bash
go run ./cmd/server
```

Test:

```bash
curl -X POST localhost:8080/render \
  -H "Content-Type: application/json" \
  -d '{"html":"<div><p>Hello from Gotei</p></div>"}' \
  --output out.pdf
```

Download `out.pdf` and open it locally.

---

## 🐳 Run with Docker

```bash
docker build -t gotei .
docker run -p 8080:8080 gotei
```

---

## 🔌 API

### `POST /render`

```json
{
  "html": "<div><p>Hello World</p></div>"
}
```

Response:

```
application/pdf
```

---

## 🧠 How it works

```
HTML
 → DOM
 → Layout Tree
 → Positioned Boxes
 → PDF Rendering
```

---

## ⚠️ Current Limitations

* No CSS support (yet)
* No flexbox/grid
* Basic text layout only
* No images (yet)
* Single-page rendering (pagination WIP)

---

## 🛣️ Roadmap

* [ ] Margins & padding
* [ ] Pagination (multi-page support)
* [ ] Image rendering
* [ ] Basic styling (inline styles)
* [ ] Plugin system

---

## 🎯 Vision

Gotei aims to be:

> A fast, predictable, and developer-friendly alternative to browser-based PDF generation.

---

## 🤝 Contributing

PRs and ideas are welcome.
Keep it simple, predictable, and focused on document rendering.

---

## 📄 License

MIT (or choose later)
