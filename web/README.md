# Chain-of-Custody Stakeholder Navigator

> Static Vue 3 + Vuetify 3 web app that organizes the chain-of-custody documentation by stakeholder. Pick a role; the view emphasizes documents most relevant to that role. The full corpus is always one click away.

## Run it

The app is a single static page with no build step required. **Serve it via HTTP** so the in-app markdown reader can fetch documents.

### Recommended — serve with a local HTTP server

From the repo root:

```
# Python (built in)
python -m http.server 8000
# Then visit http://localhost:8000/web/

# Node (http-server)
npx http-server -p 8000
# Then visit http://localhost:8000/web/
```

A static server lets the in-app dialog fetch and render the `.md` files cleanly.

### Open the file directly (degraded mode)

You can also open `index.html` in any modern browser via `file://`. The static stakeholder navigation works fine, but **the in-app document viewer cannot fetch local markdown files over `file://`** (browser security restriction). The app shows a banner explaining this and the "Source" link still opens the raw markdown in a new tab.

## Features

- **Stakeholder picker** — choose from 16 roles via the app bar dropdown or the navigation drawer
- **In-app markdown reader** — click "Read" on any document card to open a scrollable dialog with the rendered markdown
- **Mermaid diagram support** — diagrams in the docs render as SVGs in the dialog
- **Light / dark theme** — toggle in the app bar; persists in localStorage
- **Tooltips** — on icons, theme toggle, navigation items, and dialog actions
- **Categorized browsing** — collapsed expansion panels keep visual load small
- **Tables and charts** — convergence trajectory, cost-by-tier, threat-model adversaries on the overview page

## Files

- `index.html` — application shell; loads Vue, Vuetify, content, and app
- `content.js` — **all the data lives here**: stakeholders, documents, categories, analogies, trajectory data
- `app.js` — Vue 3 components (Overview, StakeholderView, AllView, DocCard, TrajectoryChart)
- `style.css` — minimal polish on top of Vuetify

## How to add a stakeholder

In `content.js`, append to the `STAKEHOLDERS` array:

```js
{
  id: 'new-role',
  name: 'New Role Name',
  icon: 'mdi-some-icon',
  tagline: 'Short one-line description.',
  description: 'Longer description shown on the role page.',
  primaryDocs: ['doc-id-1', 'doc-id-2', 'doc-id-3']
}
```

The `primaryDocs` array uses document `id`s from the `DOCUMENTS` array. Documents listed here are emphasized for the role; all other documents remain accessible under "Other documents available".

## How to add a document

In `content.js`, append to the `DOCUMENTS` array:

```js
{
  id: 'unique-doc-id',
  title: 'Document Title',
  category: 'operations',           // matches a CATEGORIES id
  path: '../docs/your-document.md', // relative path to source
  summary: 'One-sentence summary.',
  analogy: 'Optional analogy to clarify dry procedures.',
  keyPoints: [
    'First key point',
    'Second key point',
    'Third key point'
  ]
}
```

Then add the document's `id` to the `primaryDocs` of any stakeholders for whom the document is most relevant.

## How to add a category

In `content.js`, append to the `CATEGORIES` array:

```js
{ id: 'new-cat', name: 'New Category', icon: 'mdi-icon-name', color: 'primary' }
```

Documents reference the category by `id`.

## Theme

The app supports light and dark themes. The theme is persisted in `localStorage` between visits. The default follows the OS preference.

## Stakeholder selection

The selected stakeholder is persisted in `localStorage`. Click "Pick a role" in the app bar (or the navigation drawer) to switch.

## Conventions

The text in `content.js` uses plain prose with bullets and analogies. Vuetify renders the structured components (cards, tables, expansion panels, tooltips). When adding documents, keep summaries to one or two sentences and key points to 3–5 items. Analogies are optional but valued — they make dry procedures stick.

## Browser support

Modern browsers (last 2 years). Tested on Chrome, Firefox, Edge, Safari.
