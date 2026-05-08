// app.js — Vue 3 + Vuetify 3 application
const { createApp, ref, computed, onMounted, watch, provide, inject, nextTick } = Vue;
const { createVuetify } = Vuetify;

// Configure marked once. Render fenced ```mermaid blocks as <div class="mermaid">
// so the Mermaid library can pick them up after the markdown is injected.
//
// marked v12 passes a token object to renderer.code(token). Earlier versions
// passed (code, infostring, escaped). We accept both shapes defensively.
if (window.marked) {
  marked.use({
    gfm: true,
    breaks: false,
    renderer: {
      code(arg, infostring) {
        let lang, text;
        if (arg && typeof arg === 'object') {
          lang = (arg.lang || '').split(/\s+/)[0];
          text = arg.text || '';
        } else {
          lang = (infostring || '').split(/\s+/)[0];
          text = typeof arg === 'string' ? arg : '';
        }
        if (lang === 'mermaid') {
          // Mermaid will turn this into an SVG after we inject it into the DOM.
          return '<div class="mermaid">' + text + '</div>';
        }
        return false; // fall through to marked's default code rendering
      }
    }
  });
}

if (window.mermaid) {
  // Initial setup. Theme is reset before each render call to follow current app theme.
  mermaid.initialize({ startOnLoad: false, theme: 'default', securityLevel: 'loose' });
}

// Theme definitions — conservative auditor-credible palette
const lightTheme = {
  dark: false,
  colors: {
    background: '#f6f8fb',
    surface:    '#ffffff',
    primary:    '#0d4f78',
    secondary:  '#3a7ca5',
    info:       '#0e7490',
    success:    '#1e7e4a',
    warning:    '#b45309',
    error:      '#b91c1c',
    'on-background': '#1a202c',
    'on-surface':    '#1a202c',
    'on-primary':    '#ffffff'
  }
};

const darkTheme = {
  dark: true,
  colors: {
    background: '#0d1117',
    surface:    '#161b22',
    primary:    '#58a6ff',
    secondary:  '#7ee2b8',
    info:       '#79c0ff',
    success:    '#56d364',
    warning:    '#e3b341',
    error:      '#f85149',
    'on-background': '#c9d1d9',
    'on-surface':    '#c9d1d9',
    'on-primary':    '#0d1117'
  }
};

const vuetify = createVuetify({
  theme: {
    defaultTheme: 'light',
    themes: { light: lightTheme, dark: darkTheme }
  }
});

// ============================================================
// Component: DocCard — renders a single document with summary,
// analogy callout, key points, and link to the source.
// ============================================================
const DocCard = {
  props: {
    doc: { type: Object, required: true },
    emphasis: { type: String, default: 'tertiary' } // 'primary' | 'secondary' | 'tertiary'
  },
  setup() {
    const openDoc = inject('openDoc');
    return { openDoc };
  },
  computed: {
    category() {
      return window.CONTENT.CATEGORIES.find(c => c.id === this.doc.category) || {};
    }
  },
  template: `
    <v-card
      class="doc-card mb-3"
      :class="['emphasis-' + emphasis]"
      variant="outlined"
      :elevation="emphasis === 'primary' ? 1 : 0"
    >
      <v-card-item>
        <template #prepend>
          <v-tooltip :text="category.name" location="top">
            <template #activator="{ props }">
              <v-icon
                v-bind="props"
                :icon="category.icon || 'mdi-file-document-outline'"
                :color="emphasis === 'primary' ? 'primary' : 'on-surface'"
                size="large"
              ></v-icon>
            </template>
          </v-tooltip>
        </template>
        <v-card-title class="text-body-1 font-weight-medium">{{ doc.title }}</v-card-title>
        <v-card-subtitle class="pt-1">{{ doc.summary }}</v-card-subtitle>
      </v-card-item>

      <v-card-text class="pt-0">
        <div v-if="doc.analogy" class="analogy-callout mt-2">
          <span class="analogy-label">Analogy</span>
          {{ doc.analogy }}
        </div>

        <div v-if="doc.keyPoints && doc.keyPoints.length" class="mt-3">
          <p class="section-head mb-1">Key points</p>
          <ul class="keypoints">
            <li v-for="(kp, i) in doc.keyPoints" :key="i">{{ kp }}</li>
          </ul>
        </div>
      </v-card-text>

      <v-card-actions>
        <v-tooltip text="Read the full document in a dialog" location="bottom">
          <template #activator="{ props }">
            <v-btn
              v-bind="props"
              variant="flat"
              color="primary"
              size="small"
              prepend-icon="mdi-book-open-variant"
              @click="openDoc(doc.id)"
            >Read</v-btn>
          </template>
        </v-tooltip>
        <v-tooltip text="Open the raw markdown in a new tab" location="bottom">
          <template #activator="{ props }">
            <v-btn
              v-bind="props"
              :href="doc.path"
              target="_blank"
              variant="text"
              size="small"
              prepend-icon="mdi-open-in-new"
            >Source</v-btn>
          </template>
        </v-tooltip>
        <v-spacer></v-spacer>
        <v-chip
          v-if="emphasis === 'primary'"
          size="x-small"
          color="primary"
          variant="flat"
        >Emphasized</v-chip>
      </v-card-actions>
    </v-card>
  `
};

// ============================================================
// Component: DocViewer — modal dialog that fetches and renders
// the markdown source of a selected document.
// ============================================================
const DocViewer = {
  props: {
    docId: { type: String, default: null },
    open:  { type: Boolean, default: false },
    isDark: { type: Boolean, default: false }
  },
  emits: ['close'],
  data() {
    return {
      loading: false,
      error: null,
      htmlContent: '',
      isFileProtocol: typeof window !== 'undefined' && window.location && window.location.protocol === 'file:'
    };
  },
  computed: {
    doc() {
      if (!this.docId) return null;
      return window.CONTENT.DOCUMENTS.find(d => d.id === this.docId);
    },
    category() {
      if (!this.doc) return {};
      return window.CONTENT.CATEGORIES.find(c => c.id === this.doc.category) || {};
    },
    isOpen: {
      get() { return this.open; },
      set(v) { if (!v) this.$emit('close'); }
    }
  },
  watch: {
    docId(newId, oldId) {
      if (newId && newId !== oldId) this.loadDoc();
    },
    open(v) {
      if (v && this.docId && !this.htmlContent) this.loadDoc();
    },
    isDark() {
      // Re-render mermaid diagrams to pick up the new theme
      if (this.htmlContent && this.open) {
        this.$nextTick(() => this.runMermaid());
      }
    }
  },
  methods: {
    async loadDoc() {
      if (!this.doc) return;
      this.loading = true;
      this.error = null;
      this.htmlContent = '';
      try {
        const response = await fetch(this.doc.path);
        if (!response.ok) {
          throw new Error('HTTP ' + response.status + ' fetching ' + this.doc.path);
        }
        const md = await response.text();
        this.htmlContent = window.marked ? marked.parse(md) : ('<pre>' + this.escape(md) + '</pre>');
        this.$nextTick(() => this.runMermaid());
      } catch (e) {
        let msg = e.message || String(e);
        if (this.isFileProtocol) {
          msg += '\n\nThis page was opened via file:// — most browsers block fetch() for local files. ' +
                 'Serve the page over HTTP instead. From the repo root: python -m http.server 8000, ' +
                 'then visit http://localhost:8000/web/.';
        }
        this.error = msg;
      } finally {
        this.loading = false;
      }
    },
    runMermaid() {
      if (!window.mermaid) return;
      try {
        mermaid.initialize({
          startOnLoad: false,
          theme: this.isDark ? 'dark' : 'default',
          securityLevel: 'loose'
        });
        const nodes = document.querySelectorAll('.doc-viewer-body .mermaid');
        if (nodes.length) {
          // Reset processed state so re-render works after theme change
          nodes.forEach(n => { n.removeAttribute('data-processed'); });
          mermaid.run({ nodes: Array.from(nodes) });
        }
      } catch (e) {
        // Mermaid render failures are non-fatal
        console.warn('Mermaid render warning:', e);
      }
    },
    escape(s) {
      return s.replace(/[&<>]/g, c => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;' }[c]));
    }
  },
  template: `
    <v-dialog v-model="isOpen" max-width="1100" scrollable>
      <v-card class="doc-viewer-card" v-if="doc">
        <v-card-title class="d-flex align-center pa-3">
          <v-icon :icon="category.icon || 'mdi-file-document-outline'" :color="category.color || 'primary'" class="mr-3"></v-icon>
          <div class="flex-grow-1">
            <div class="text-body-1 font-weight-medium">{{ doc.title }}</div>
            <div class="text-caption" style="opacity: 0.66;">{{ category.name }}</div>
          </div>
          <v-tooltip text="Open raw markdown in a new tab" location="bottom">
            <template #activator="{ props }">
              <v-btn v-bind="props" :href="doc.path" target="_blank" variant="text" size="small" prepend-icon="mdi-open-in-new" class="mr-1">Source</v-btn>
            </template>
          </v-tooltip>
          <v-tooltip text="Close (Esc)" location="bottom">
            <template #activator="{ props }">
              <v-btn v-bind="props" icon="mdi-close" variant="text" @click="$emit('close')"></v-btn>
            </template>
          </v-tooltip>
        </v-card-title>
        <v-divider></v-divider>
        <v-alert
          v-if="isFileProtocol && !htmlContent && !loading"
          type="info"
          density="compact"
          variant="tonal"
          class="protocol-banner ma-2"
        >
          Opened via file:// — fetching local markdown may be blocked. Serve via <code>python -m http.server</code> for full functionality.
        </v-alert>
        <div class="doc-viewer-body markdown-content">
          <div v-if="loading" class="d-flex justify-center align-center" style="min-height: 200px;">
            <v-progress-circular indeterminate color="primary"></v-progress-circular>
          </div>
          <v-alert v-else-if="error" type="error" variant="tonal" class="ma-2">
            <div style="white-space: pre-wrap;">{{ error }}</div>
          </v-alert>
          <div v-else-if="htmlContent" v-html="htmlContent"></div>
          <div v-else class="text-center pa-6" style="opacity: 0.6;">
            <v-icon icon="mdi-file-document-outline" size="48"></v-icon>
            <div class="mt-2">No content loaded yet.</div>
          </div>
        </div>
      </v-card>
    </v-dialog>
  `
};

// ============================================================
// Component: Overview — the home page.
// ============================================================
const Overview = {
  setup() {
    const openDoc = inject('openDoc');
    return { openDoc };
  },
  computed: {
    stats() { return window.CONTENT.STATS; },
    threats() { return window.CONTENT.THREATS; },
    costComponents() { return window.CONTENT.COST_COMPONENTS; },
    costTiers() { return window.CONTENT.COST_TIERS; },
    costReductionOptions() { return window.CONTENT.COST_REDUCTION_OPTIONS; },
    auditStages() { return window.CONTENT.AUDIT_STAGES; }
  },
  methods: {
    docTitle(id) {
      const d = window.CONTENT.DOCUMENTS.find(x => x.id === id);
      return d ? d.title : id;
    }
  },
  template: `
    <div class="reading-width">
      <div class="hero">
        <div class="hero-title">Chain of Custody — v1.0</div>
        <div class="hero-tagline">Integrity-bearing AI-decision evidence for regulated institutions. Verified by examiners independently of any vendor.</div>
      </div>

      <v-row dense class="mb-4">
        <v-col cols="6" sm="3">
          <v-card variant="outlined" class="stat-card">
            <div class="stat-value">{{ stats.totalDocs }}</div>
            <div class="stat-label">Documents</div>
          </v-card>
        </v-col>
        <v-col cols="6" sm="3">
          <v-card variant="outlined" class="stat-card">
            <div class="stat-value">{{ stats.totalStakeholders }}</div>
            <div class="stat-label">Stakeholders</div>
          </v-card>
        </v-col>
        <v-col cols="6" sm="3">
          <v-card variant="outlined" class="stat-card">
            <div class="stat-value">{{ stats.primitives }}</div>
            <div class="stat-label">Primitives</div>
          </v-card>
        </v-col>
        <v-col cols="6" sm="3">
          <v-card variant="outlined" class="stat-card">
            <div class="stat-value">{{ stats.cuecs }}</div>
            <div class="stat-label">CUECs</div>
          </v-card>
        </v-col>
      </v-row>

      <v-card variant="outlined" class="mb-4">
        <v-card-item>
          <v-card-title class="text-body-1 font-weight-medium">What this is</v-card-title>
        </v-card-item>
        <v-card-text>
          <p>The corpus describes a <strong>chain-of-custody primitive set for AI-driven decisions in regulated systems</strong>. It captures every AI agent decision with cryptographic evidence that the record is authentic and unaltered, and that an examiner can verify both properties without trusting the institution.</p>
          <div class="analogy-callout mt-3">
            <span class="analogy-label">The big picture</span>
            Think of it as the FAA's flight-data recorder for AI. Every decision the agent makes is sealed in a tamper-evident box; an independent regulator can pry the box open at any time and confirm what's inside.
          </div>
          <p class="mt-3 mb-0">Pick a stakeholder from the left to see the documents most relevant to your role. The full corpus is always one click away.</p>
        </v-card-text>
      </v-card>

      <v-expansion-panels variant="accordion" class="mb-4">

        <!-- Cost panel -->
        <v-expansion-panel>
          <v-expansion-panel-title>
            <v-icon icon="mdi-currency-usd" size="small" class="mr-2"></v-icon>
            <span><strong>Annual cost</strong> — components, per-tier breakdown, reduction options</span>
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <p class="mb-3">The cost of running the chain is dominated by HSM operations. Software is open source. The remainder is institution-specific operational overhead. The table below shows the components; the per-tier breakdown shows how those components combine into the headline annual cost.</p>

            <p class="section-head mt-2 mb-2">Cost components — what makes up the annual figure</p>
            <v-table density="compact" class="mb-4">
              <thead>
                <tr><th style="width: 36px;"></th><th>Component</th><th style="width: 110px;">Share</th><th>Why</th></tr>
              </thead>
              <tbody>
                <tr v-for="c in costComponents" :key="c.name">
                  <td><v-icon :icon="c.icon" size="small" color="primary"></v-icon></td>
                  <td><strong>{{ c.name }}</strong></td>
                  <td>{{ c.share }}</td>
                  <td style="opacity: 0.85">{{ c.detail }}</td>
                </tr>
              </tbody>
            </v-table>

            <p class="section-head mt-4 mb-2">Per-tier breakdown — how the cost is achieved at each scale</p>
            <v-expansion-panels variant="accordion" multiple>
              <v-expansion-panel v-for="t in costTiers" :key="t.tier">
                <v-expansion-panel-title>
                  <v-icon :icon="t.icon" size="small" class="mr-2" color="primary"></v-icon>
                  <strong class="mr-2">{{ t.tier }}</strong>
                  <v-chip size="x-small" color="primary" variant="tonal">{{ t.annual }}</v-chip>
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <p class="text-caption mb-2" style="opacity: 0.7;">Configuration</p>
                  <p class="mb-3">{{ t.config }}</p>

                  <p class="text-caption mb-2" style="opacity: 0.7;">{{ t.workedExample.title }}</p>
                  <v-table density="compact" class="mb-3">
                    <tbody>
                      <tr v-for="(line, i) in t.workedExample.lines" :key="i" :class="{ 'font-weight-bold': line.total }">
                        <td>{{ line.label }}</td>
                        <td class="text-right">{{ line.cost }}</td>
                      </tr>
                    </tbody>
                  </v-table>

                  <p class="text-caption mb-2" style="opacity: 0.7;">Cost drivers at this tier</p>
                  <ul class="keypoints">
                    <li v-for="(d, i) in t.drivers" :key="i">{{ d }}</li>
                  </ul>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>

            <p class="section-head mt-4 mb-2">Cost reduction options</p>
            <v-table density="compact" class="mb-2">
              <thead>
                <tr><th>Option</th><th>Saving</th><th>Trade-off</th><th>Best for</th></tr>
              </thead>
              <tbody>
                <tr v-for="o in costReductionOptions" :key="o.option">
                  <td><strong>{{ o.option }}</strong></td>
                  <td>{{ o.saving }}</td>
                  <td style="opacity: 0.85">{{ o.tradeoff }}</td>
                  <td style="opacity: 0.85">{{ o.bestFor }}</td>
                </tr>
              </tbody>
            </v-table>

            <v-divider class="my-3"></v-divider>
            <p class="text-caption mb-0" style="opacity: 0.66;">Numbers are 2026 cloud reference pricing. Negotiated pricing, on-prem deployments, and existing volume commitments shift the figures. The full cost-model document has trajectory-over-time and additional detail.</p>
            <v-btn
              variant="tonal"
              size="small"
              prepend-icon="mdi-book-open-variant"
              class="mt-2"
              @click="openDoc('cost-model')"
            >Read the full cost model</v-btn>
          </v-expansion-panel-text>
        </v-expansion-panel>

        <!-- Audit process timeline -->
        <v-expansion-panel>
          <v-expansion-panel-title>
            <v-icon icon="mdi-timeline-clock-outline" size="small" class="mr-2"></v-icon>
            <span><strong>Audit process timeline</strong> — what auditor and auditee do at each stage</span>
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <p class="mb-3">From engagement notice through remediation. The chain composes with standard FFIEC examination and SOC engagement workflows; this timeline shows where the chain is touched at each stage. Each stage shows the <strong>auditor's actions and what they look for</strong> alongside the <strong>auditee's actions and what they provide</strong>.</p>

            <div class="audit-timeline">
              <div
                v-for="(stage, idx) in auditStages"
                :key="stage.id"
                class="audit-stage"
                :class="{ 'audit-stage-last': idx === auditStages.length - 1 }"
              >
                <div class="audit-stage-marker">
                  <v-avatar :color="stage.color" size="42">
                    <v-icon :icon="stage.icon" color="white"></v-icon>
                  </v-avatar>
                </div>
                <div class="audit-stage-body">
                  <div class="audit-stage-head">
                    <div class="audit-stage-name">{{ stage.name }}</div>
                    <v-chip size="x-small" variant="tonal" :color="stage.color">{{ stage.timing }}</v-chip>
                  </div>
                  <div class="audit-stage-summary">{{ stage.summary }}</div>

                  <v-row dense class="mt-2">
                    <v-col cols="12" md="6">
                      <div class="audit-side audit-side-auditor">
                        <div class="audit-side-label">
                          <v-icon icon="mdi-shield-search" size="small" class="mr-1"></v-icon>
                          Auditor
                        </div>
                        <div class="audit-subhead">Actions</div>
                        <ul class="keypoints">
                          <li v-for="(a, i) in stage.auditor.actions" :key="'a'+i">{{ a }}</li>
                        </ul>
                        <div v-if="stage.auditor.lookFor && stage.auditor.lookFor.length" class="audit-subhead mt-2">What they look for</div>
                        <ul v-if="stage.auditor.lookFor && stage.auditor.lookFor.length" class="keypoints">
                          <li v-for="(l, i) in stage.auditor.lookFor" :key="'lf'+i">{{ l }}</li>
                        </ul>
                      </div>
                    </v-col>
                    <v-col cols="12" md="6">
                      <div class="audit-side audit-side-auditee">
                        <div class="audit-side-label">
                          <v-icon icon="mdi-account-tie-outline" size="small" class="mr-1"></v-icon>
                          Auditee
                        </div>
                        <div class="audit-subhead">Actions</div>
                        <ul class="keypoints">
                          <li v-for="(a, i) in stage.auditee.actions" :key="'aa'+i">{{ a }}</li>
                        </ul>
                        <div v-if="stage.auditee.provide && stage.auditee.provide.length" class="audit-subhead mt-2">What they provide</div>
                        <ul v-if="stage.auditee.provide && stage.auditee.provide.length" class="keypoints">
                          <li v-for="(p, i) in stage.auditee.provide" :key="'pp'+i">{{ p }}</li>
                        </ul>
                      </div>
                    </v-col>
                  </v-row>

                  <div v-if="stage.relatedDocs && stage.relatedDocs.length" class="mt-2">
                    <span class="text-caption mr-2" style="opacity: 0.7;">Related:</span>
                    <v-btn
                      v-for="docId in stage.relatedDocs"
                      :key="docId"
                      variant="tonal"
                      size="x-small"
                      class="mr-1 mb-1"
                      prepend-icon="mdi-book-open-variant"
                      @click="openDoc(docId)"
                    >{{ docTitle(docId) }}</v-btn>
                  </div>
                </div>
              </div>
            </div>

            <div class="analogy-callout mt-4">
              <span class="analogy-label">Analogy</span>
              An audit is choreography. The auditor moves through pre-set positions on stage; the auditee must hit each cue. Skip a step or miss a beat, and the next step gets harder. The chain is the prop both performers reach for at each cue.
            </div>
          </v-expansion-panel-text>
        </v-expansion-panel>

        <!-- Threat model panel -->
        <v-expansion-panel>
          <v-expansion-panel-title>
            <v-icon icon="mdi-shield-alert-outline" size="small" class="mr-2"></v-icon>
            <span><strong>Threat model</strong> — eight adversaries with capability, defense, and residual risk</span>
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <p class="mb-3">The chain defends against eight adversary classes. Each class has a specific capability, goal, and defense. The residual risk is what remains after the defense applies — the threat model accepts these explicitly. Click any adversary to see the detail and links to the related documents.</p>

            <v-expansion-panels variant="accordion" multiple>
              <v-expansion-panel v-for="t in threats" :key="t.id">
                <v-expansion-panel-title>
                  <v-icon :icon="t.icon" size="small" class="mr-2"></v-icon>
                  <strong class="mr-2" style="font-family: monospace;">Adv {{ t.id }}</strong>
                  {{ t.name }}
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <p><strong>Capability.</strong> {{ t.capability }}</p>
                  <p><strong>Goal.</strong> {{ t.goal }}</p>
                  <p><strong>Defense.</strong> {{ t.defense }}</p>
                  <p><strong>Residual risk.</strong> {{ t.residual }}</p>
                  <v-divider class="my-3"></v-divider>
                  <p class="text-caption mb-2" style="opacity: 0.7;">Related documents (click to read)</p>
                  <v-btn
                    v-for="docId in t.relatedDocs"
                    :key="docId"
                    variant="tonal"
                    size="small"
                    class="mr-2 mb-2"
                    color="primary"
                    prepend-icon="mdi-book-open-variant"
                    @click="openDoc(docId)"
                  >{{ docTitle(docId) }}</v-btn>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>

            <v-divider class="my-3"></v-divider>
            <v-btn
              variant="tonal"
              size="small"
              prepend-icon="mdi-book-open-variant"
              @click="openDoc('threat-model')"
            >Read the full threat model</v-btn>
          </v-expansion-panel-text>
        </v-expansion-panel>

      </v-expansion-panels>
    </div>
  `
};

// ============================================================
// Component: StakeholderView — show emphasized + tertiary docs
// for a selected stakeholder.
// ============================================================
const StakeholderView = {
  components: { DocCard },
  props: {
    stakeholder: { type: Object, required: true }
  },
  computed: {
    primaryDocs() {
      const ids = this.stakeholder.primaryDocs || [];
      return ids.map(id => window.CONTENT.DOCUMENTS.find(d => d.id === id)).filter(Boolean);
    },
    otherDocs() {
      const primaryIds = new Set(this.stakeholder.primaryDocs || []);
      return window.CONTENT.DOCUMENTS.filter(d => !primaryIds.has(d.id));
    },
    otherDocsByCategory() {
      const groups = {};
      this.otherDocs.forEach(d => {
        if (!groups[d.category]) groups[d.category] = [];
        groups[d.category].push(d);
      });
      return groups;
    },
    categories() { return window.CONTENT.CATEGORIES; }
  },
  methods: {
    catName(id) {
      const c = window.CONTENT.CATEGORIES.find(x => x.id === id);
      return c ? c.name : id;
    },
    catIcon(id) {
      const c = window.CONTENT.CATEGORIES.find(x => x.id === id);
      return c ? c.icon : 'mdi-file-document-outline';
    }
  },
  template: `
    <div class="reading-width">
      <div class="hero">
        <v-icon :icon="stakeholder.icon" size="48" color="primary"></v-icon>
        <div class="hero-title mt-2">{{ stakeholder.name }}</div>
        <div class="hero-tagline">{{ stakeholder.tagline }}</div>
      </div>

      <v-card variant="outlined" class="mb-4">
        <v-card-text>{{ stakeholder.description }}</v-card-text>
      </v-card>

      <div v-if="primaryDocs.length">
        <p class="section-head">Most relevant for this role</p>
        <doc-card v-for="d in primaryDocs" :key="d.id" :doc="d" emphasis="primary"></doc-card>
      </div>

      <p class="section-head mt-6">Other documents available (click to expand)</p>
      <v-expansion-panels variant="accordion" multiple>
        <v-expansion-panel
          v-for="cat in categories"
          :key="cat.id"
          v-show="otherDocsByCategory[cat.id] && otherDocsByCategory[cat.id].length"
        >
          <v-expansion-panel-title>
            <v-icon :icon="cat.icon" size="small" class="mr-2"></v-icon>
            {{ cat.name }}
            <v-chip size="x-small" class="ml-2" variant="tonal">{{ (otherDocsByCategory[cat.id] || []).length }}</v-chip>
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <doc-card v-for="d in otherDocsByCategory[cat.id]" :key="d.id" :doc="d" emphasis="tertiary"></doc-card>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
    </div>
  `
};

// ============================================================
// Component: AllView — traditional all-categories browsing
// ============================================================
const AllView = {
  components: { DocCard },
  computed: {
    categories() { return window.CONTENT.CATEGORIES; },
    docsByCategory() {
      const groups = {};
      window.CONTENT.DOCUMENTS.forEach(d => {
        if (!groups[d.category]) groups[d.category] = [];
        groups[d.category].push(d);
      });
      return groups;
    }
  },
  template: `
    <div class="reading-width">
      <div class="hero">
        <div class="hero-title">All Documents</div>
        <div class="hero-tagline">The complete corpus, organized by category. Click a category to expand.</div>
      </div>

      <v-expansion-panels variant="accordion" multiple>
        <v-expansion-panel
          v-for="cat in categories"
          :key="cat.id"
          v-show="docsByCategory[cat.id] && docsByCategory[cat.id].length"
        >
          <v-expansion-panel-title>
            <v-icon :icon="cat.icon" size="small" class="mr-2"></v-icon>
            <strong>{{ cat.name }}</strong>
            <v-chip size="x-small" class="ml-2" variant="tonal">{{ (docsByCategory[cat.id] || []).length }}</v-chip>
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <doc-card v-for="d in docsByCategory[cat.id]" :key="d.id" :doc="d" emphasis="secondary"></doc-card>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
    </div>
  `
};

// ============================================================
// Main App
// ============================================================
const App = {
  components: { Overview, StakeholderView, AllView, DocViewer },
  setup() {
    const drawer = ref(true);
    const isDark = ref(window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches);
    const selectedStakeholderId = ref(null); // null = overview
    const stakeholders = window.CONTENT.STAKEHOLDERS;

    // Document viewer state — provided to descendants via inject
    const viewerOpen = ref(false);
    const viewerDocId = ref(null);
    const openDoc = (id) => {
      viewerDocId.value = id;
      viewerOpen.value = true;
    };
    const closeViewer = () => {
      viewerOpen.value = false;
    };
    provide('openDoc', openDoc);

    const selectedStakeholder = computed(() => {
      if (!selectedStakeholderId.value) return null;
      return stakeholders.find(s => s.id === selectedStakeholderId.value);
    });

    const viewMode = computed(() => {
      if (!selectedStakeholderId.value) return 'overview';
      if (selectedStakeholderId.value === 'all') return 'all';
      return 'stakeholder';
    });

    const toggleTheme = () => {
      isDark.value = !isDark.value;
    };

    // Persist theme & stakeholder in localStorage
    onMounted(() => {
      try {
        const savedTheme = localStorage.getItem('chain-theme');
        if (savedTheme === 'dark') isDark.value = true;
        if (savedTheme === 'light') isDark.value = false;
        const savedStakeholder = localStorage.getItem('chain-stakeholder');
        if (savedStakeholder) selectedStakeholderId.value = savedStakeholder;
      } catch (e) { /* ignore */ }
    });

    watch(isDark, (v) => {
      try { localStorage.setItem('chain-theme', v ? 'dark' : 'light'); } catch (e) {}
    });
    watch(selectedStakeholderId, (v) => {
      try {
        if (v) localStorage.setItem('chain-stakeholder', v);
        else localStorage.removeItem('chain-stakeholder');
      } catch (e) {}
    });

    const docCount = window.CONTENT.DOCUMENTS.length;

    return {
      drawer, isDark, selectedStakeholderId, selectedStakeholder, viewMode,
      stakeholders, toggleTheme, docCount,
      viewerOpen, viewerDocId, openDoc, closeViewer
    };
  },
  template: `
    <v-app :theme="isDark ? 'dark' : 'light'">
      <v-app-bar :elevation="1" density="comfortable">
        <v-app-bar-nav-icon @click="drawer = !drawer" aria-label="Toggle navigation"></v-app-bar-nav-icon>
        <v-app-bar-title class="font-weight-medium">
          <v-icon icon="mdi-link-variant" class="mr-2"></v-icon>
          Chain of Custody
        </v-app-bar-title>
        <v-spacer></v-spacer>

        <v-menu>
          <template #activator="{ props }">
            <v-tooltip text="Pick a stakeholder to focus the view" location="bottom">
              <template #activator="{ props: tooltipProps }">
                <v-btn v-bind="{ ...props, ...tooltipProps }" prepend-icon="mdi-account-search-outline" variant="tonal" class="mr-2">
                  {{ selectedStakeholder ? selectedStakeholder.name : 'Pick a role' }}
                </v-btn>
              </template>
            </v-tooltip>
          </template>
          <v-list density="compact" min-width="280">
            <v-list-item @click="selectedStakeholderId = null">
              <template #prepend>
                <v-icon icon="mdi-home-outline"></v-icon>
              </template>
              <v-list-item-title>Overview (home)</v-list-item-title>
            </v-list-item>
            <v-divider></v-divider>
            <v-list-item v-for="s in stakeholders" :key="s.id" @click="selectedStakeholderId = s.id">
              <template #prepend>
                <v-icon :icon="s.icon"></v-icon>
              </template>
              <v-list-item-title>{{ s.name }}</v-list-item-title>
              <v-list-item-subtitle>{{ s.tagline }}</v-list-item-subtitle>
            </v-list-item>
          </v-list>
        </v-menu>

        <v-tooltip :text="isDark ? 'Switch to light theme' : 'Switch to dark theme'" location="bottom">
          <template #activator="{ props }">
            <v-btn
              v-bind="props"
              :icon="isDark ? 'mdi-weather-sunny' : 'mdi-weather-night'"
              variant="text"
              @click="toggleTheme"
            ></v-btn>
          </template>
        </v-tooltip>
      </v-app-bar>

      <v-navigation-drawer v-model="drawer" :elevation="1">
        <v-list density="compact" nav>
          <v-list-subheader>VIEW</v-list-subheader>
          <v-list-item
            prepend-icon="mdi-home-outline"
            title="Overview"
            :active="!selectedStakeholderId"
            @click="selectedStakeholderId = null"
          ></v-list-item>
          <v-list-item
            prepend-icon="mdi-view-grid-outline"
            title="All documents"
            :active="selectedStakeholderId === 'all'"
            @click="selectedStakeholderId = 'all'"
          ></v-list-item>

          <v-divider class="my-2"></v-divider>
          <v-list-subheader>BY STAKEHOLDER</v-list-subheader>
          <v-list-item
            v-for="s in stakeholders.filter(x => x.id !== 'all')"
            :key="s.id"
            :prepend-icon="s.icon"
            :title="s.name"
            :active="selectedStakeholderId === s.id"
            @click="selectedStakeholderId = s.id"
            class="stakeholder-pick-row"
          >
            <v-tooltip activator="parent" location="end" :text="s.tagline" open-delay="500"></v-tooltip>
          </v-list-item>
        </v-list>

        <template #append>
          <div class="pa-3 text-caption" style="opacity: 0.6;">
            <div>Spec v1.0-final</div>
            <div>{{ stakeholders.length - 1 }} stakeholders · {{ docCount }} documents</div>
          </div>
        </template>
      </v-navigation-drawer>

      <v-main>
        <v-container fluid class="pa-md-6 pa-3">
          <overview v-if="viewMode === 'overview'"></overview>
          <all-view v-else-if="viewMode === 'all'"></all-view>
          <stakeholder-view v-else :stakeholder="selectedStakeholder"></stakeholder-view>
        </v-container>
      </v-main>

      <doc-viewer
        :doc-id="viewerDocId"
        :open="viewerOpen"
        :is-dark="isDark"
        @close="closeViewer"
      ></doc-viewer>
    </v-app>
  `
};

const app = createApp(App);
app.use(vuetify);
app.mount('#app');
