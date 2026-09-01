// Smoke test del carousel de trackers en dashboard.js
// Stubs de DOM + Chart + fetchAPI, luego corre loadDashboard y verifica slides.
const fs = require('fs');
const vm = require('vm');
const path = require('path');

const code = fs.readFileSync(path.join(__dirname, '..', 'web', 'js', 'dashboard.js'), 'utf8');

// ---- stubs ----
let createdCharts = [];
class Chart {
  constructor(ctx, cfg) { this.ctx = ctx; this.cfg = cfg; this.destroyed = false; createdCharts.push(this); }
  destroy() { this.destroyed = true; }
}

const elements = {};
function makeEl(tag) {
  return {
    tagName: (tag || 'div').toUpperCase(),
    children: [],
    style: {},
    className: '',
    _text: '',
    _innerHTML: '',
    clientWidth: 350,
    getContext() { return { _el: this }; },
    appendChild(c) { this.children.push(c); return c; },
    querySelector() { return this.children[0] || null; },
    scrollBy() { this._scrolled = (this._scrolled || 0) + arguments[0].left; },
    get innerHTML() { return this._innerHTML; },
    set innerHTML(v) { this._innerHTML = v; if (v === '') this.children = []; },
    get textContent() { return this._text; },
    set textContent(v) { this._text = String(v); },
    addEventListener() {},
    value: ''
  };
}

const trackerCarousel = makeEl('div');
elements['tracker-carousel'] = trackerCarousel;
elements['metric-month'] = makeEl('select');
elements['metric-month'].value = '8';
elements['metric-year'] = makeEl('select');
elements['metric-year'].value = '2026';

const documentStub = {
  addEventListener(ev, fn) { if (ev === 'DOMContentLoaded') documentStub._domReady = fn; },
  getElementById(id) {
    if (!elements[id]) elements[id] = makeEl('div');
    return elements[id];
  },
  createElement(tag) { return makeEl(tag); }
};

const trackers = [
  { id: 1, name: 'Peso', unit: 'kg', limitValue: 80, limitType: 'max' },
  { id: 2, name: 'Sueño', unit: 'h', limitValue: 8, limitType: 'min' },
  { id: 3, name: 'Sin registros', unit: '' },
  { id: 4, name: 'Sin límite', unit: 'L' }
];
const histories = {
  1: [
    { entryDate: '2026-08-01', value: 81.2, isMet: false },
    { entryDate: '2026-08-03', value: 80.9, isMet: false },
    { entryDate: '2026-08-05', value: 80.5, isMet: false }
  ],
  2: [
    { entryDate: '2026-08-01', value: 7, isMet: false },
    { entryDate: '2026-08-02', value: 8, isMet: true },
    { entryDate: '2026-08-04', value: 7.5, isMet: false }
  ],
  3: [],
  4: [
    { entryDate: '2026-08-01', value: 1.5 },
    { entryDate: '2026-08-02', value: 2 }
  ]
};

async function fetchAPI(endpoint) {
  if (/^\/metrics\?month=\d+&year=\d+$/.test(endpoint)) {
    return { monthly: { rate: 0.8, completed: 4, total: 5 }, monthlySeries: [], trackers: [] };
  }
  if (endpoint === '/immediate-tasks') return [{ isCompleted: true }];
  if (endpoint === '/trackers') return trackers;
  const m = endpoint.match(/^\/trackers\/(\d+)\/history$/);
  if (m) return histories[m[1]] || [];
  throw new Error('unexpected endpoint ' + endpoint);
}

// initFilters usa `new Date()` para el mes/año del select. Fijamos agosto
// 2026 para que el filtro coincida con `histories` (si no, el 1-sep deja
// todas las slides en empty-state).
const RealDate = Date;
function MockDate(...args) {
  if (args.length === 0) return new RealDate('2026-08-15T12:00:00Z');
  return new RealDate(...args);
}
MockDate.UTC = RealDate.UTC;
MockDate.now = () => RealDate.parse('2026-08-15T12:00:00Z');
MockDate.parse = RealDate.parse;

const sandbox = {
  document: documentStub,
  window: { addEventListener() {} },
  Chart,
  fetchAPI,
  console,
  Promise,
  Set,
  Date: MockDate,
  alert: (msg) => console.log('alert:', msg)
};
vm.createContext(sandbox);
vm.runInContext(code, sandbox);

(async () => {
  // 1) DOMContentLoaded → initFilters + loadDashboard
  await documentStub._domReady();
  await new Promise(r => setTimeout(r, 50));

  const slides = trackerCarousel.children;
  console.log('slides creadas:', slides.length);
  if (slides.length !== 4) throw new Error('esperaba 4 slides, hay ' + slides.length);

  const names = slides.map(s => s.children[0]._text);
  console.log('títulos slides:', JSON.stringify(names));
  if (names.join('|') !== 'Peso|Sueño|Sin registros|Sin límite') throw new Error('títulos inesperados: ' + names);

  // 2) Chart por tracker con datos: el dataset 0 es la línea de valores
  //    (solo line charts; el bar mensual no cuenta).
  const withChart = createdCharts.filter(c => !c.destroyed && c.cfg.type === 'line');
  console.log('charts line activos:', withChart.length, '(esperado 3: Peso, Sueño y Sin límite)');
  if (withChart.length !== 3) throw new Error('cantidad de charts line incorrecta');

  // 2a) Línea de límite (card: "marcar el límite actual en el tracker"):
  //     dataset punteado horizontal en el valor del límite, con etiqueta
  //     "Máximo: 80 kg" / "Mínimo: 8 h", y la escala Y lo cubre.
  const byName = {};
  for (const c of withChart) byName[c.cfg.data.datasets[0].label] = c;

  const peso = byName['Peso'];
  if (peso.cfg.data.datasets.length !== 2) throw new Error('Peso: se esperaba 2 datasets (valores + límite)');
  const pesoLimit = peso.cfg.data.datasets[1];
  if (pesoLimit.label !== 'Máximo: 80 kg') throw new Error('Peso: etiqueta de límite inesperada: ' + pesoLimit.label);
  if (JSON.stringify(pesoLimit.borderDash) !== '[6,4]') throw new Error('Peso: la línea de límite debe ser punteada');
  if (pesoLimit.pointRadius !== 0) throw new Error('Peso: la línea de límite no debe tener puntos');
  if (pesoLimit.borderColor !== '#22d3ee') throw new Error('Peso: el umbral debe ser cian, no el color por defecto: ' + pesoLimit.borderColor);
  if (!pesoLimit.data.every(v => v === 80)) throw new Error('Peso: la línea de límite no es horizontal en 80');
  if (peso.cfg.options.scales.y.beginAtZero) throw new Error('Peso: la escala Y no debe arrancar en 0 con límite');
  if (!(peso.cfg.options.scales.y.min <= 80 && peso.cfg.options.scales.y.max >= 80)) {
    throw new Error('Peso: la escala Y no cubre el límite 80');
  }
  const ttFilter = peso.cfg.options.plugins.tooltip.filter;
  if (!ttFilter) throw new Error('Peso: falta el filter del tooltip');
  if (ttFilter({ datasetIndex: 1 })) throw new Error('Peso: el tooltip no debe mostrar el dataset del límite');
  if (!ttFilter({ datasetIndex: 0 })) throw new Error('Peso: el tooltip sí debe mostrar el dataset de valores');

  const sueno = byName['Sueño'];
  if (sueno.cfg.data.datasets.length !== 2) throw new Error('Sueño: se esperaba 2 datasets');
  if (sueno.cfg.data.datasets[1].label !== 'Mínimo: 8 h') throw new Error('Sueño: etiqueta de límite inesperada: ' + sueno.cfg.data.datasets[1].label);
  if (!sueno.cfg.data.datasets[1].data.every(v => v === 8)) throw new Error('Sueño: la línea de límite no es horizontal en 8');
  if (!(sueno.cfg.options.scales.y.min <= 8 && sueno.cfg.options.scales.y.max >= 8)) {
    throw new Error('Sueño: la escala Y no cubre el límite 8');
  }

  const sinLimite = byName['Sin límite'];
  if (sinLimite.cfg.data.datasets.length !== 1) throw new Error('Sin límite: no debe tener dataset de límite');
  console.log('  línea de límite: Peso "Máximo: 80 kg" ✓, Sueño "Mínimo: 8 h" ✓, Sin límite sin línea ✓');

  // 2a-bis) Colores de cumplimiento: puntos y segmentos verde/rojo según isMet.
  const MET = '#10b981', MISS = '#ef4444', NEUTRAL = '#60a5fa';
  const pesoPts = peso.cfg.data.datasets[0].pointBackgroundColor;
  if (JSON.stringify(pesoPts) !== JSON.stringify([MISS, MISS, MISS])) {
    throw new Error('Peso: puntos deberían ser todos rojos (no cumple máx 80): ' + JSON.stringify(pesoPts));
  }
  const pesoSeg = peso.cfg.data.datasets[0].segment && peso.cfg.data.datasets[0].segment.borderColor;
  if (typeof pesoSeg !== 'function') throw new Error('Peso: falta segment.borderColor');
  if (pesoSeg({ p1DataIndex: 0 }) !== MISS) throw new Error('Peso: segmento 0 debería ser rojo');

  const suenoPts = sueno.cfg.data.datasets[0].pointBackgroundColor;
  if (JSON.stringify(suenoPts) !== JSON.stringify([MISS, MET, MISS])) {
    throw new Error('Sueño: puntos deberían ser rojo/verde/rojo: ' + JSON.stringify(suenoPts));
  }
  const suenoSeg = sueno.cfg.data.datasets[0].segment.borderColor;
  if (suenoSeg({ p1DataIndex: 0 }) !== MISS) throw new Error('Sueño: segmento hacia 7h debería ser rojo');
  if (suenoSeg({ p1DataIndex: 1 }) !== MET) throw new Error('Sueño: segmento hacia 8h debería ser verde');
  if (suenoSeg({ p1DataIndex: 2 }) !== MISS) throw new Error('Sueño: segmento hacia 7.5h debería ser rojo');
  if (sueno.cfg.data.datasets[1].borderColor !== '#22d3ee') {
    throw new Error('Sueño: el umbral debe ser cian: ' + sueno.cfg.data.datasets[1].borderColor);
  }

  const sinLimitePts = sinLimite.cfg.data.datasets[0].pointBackgroundColor;
  if (!sinLimitePts.every(c => c === NEUTRAL)) {
    throw new Error('Sin límite: la línea debe ser neutra, no verde/rojo: ' + JSON.stringify(sinLimitePts));
  }

  const pesoLabel = peso.cfg.options.plugins.tooltip.callbacks.label;
  const pesoTip = pesoLabel({ parsed: { y: 81.2 }, dataIndex: 0 });
  if (!pesoTip.includes('no cumplió')) throw new Error('Peso: tooltip debería decir no cumplió: ' + pesoTip);
  const suenoLabel = sueno.cfg.options.plugins.tooltip.callbacks.label;
  const suenoTip = suenoLabel({ parsed: { y: 8 }, dataIndex: 1 });
  if (!suenoTip.includes('cumplió') || suenoTip.includes('no cumplió')) {
    throw new Error('Sueño: tooltip del punto cumplido inesperado: ' + suenoTip);
  }

  if (sandbox.entryIsMet({ value: 79 }, { limitType: 'max', limitValue: 80 }, true) !== true) {
    throw new Error('fallback entryIsMet: 79 <= 80 max debería cumplir');
  }
  if (sandbox.entryIsMet({ value: 7 }, { limitType: 'min', limitValue: 8 }, true) !== false) {
    throw new Error('fallback entryIsMet: 7 < 8 min no debería cumplir');
  }
  if (sandbox.entryIsMet({ value: 2 }, { limitType: 'max', limitValue: 0 }, false) !== null) {
    throw new Error('fallback entryIsMet: sin límite debe devolver null');
  }
  console.log('  colores: Peso rojo (no cumple) ✓, Sueño rojo/verde/rojo ✓, umbral cian ✓, Sin límite neutro ✓');

  for (const c of withChart) {
    if (c.cfg.data.datasets[0].data.length !== c.cfg.data.labels.length) {
      throw new Error('dataset de valores con longitud distinta a los labels');
    }
    console.log(`  chart "${c.cfg.data.labels.length} puntos" labels=${c.cfg.data.labels[0]}..${c.cfg.data.labels[c.cfg.data.labels.length-1]}`);

    // 2b) Escala Y ajustada: sin beginAtZero, con min/max cerca de los datos
    const y = c.cfg.options.scales.y;
    const vals = c.cfg.data.datasets[0].data;
    const dmin = Math.min(...vals), dmax = Math.max(...vals);
    if (y.beginAtZero) throw new Error('la escala Y no debe arrancar en 0');
    if (typeof y.min !== 'number' || typeof y.max !== 'number') throw new Error('falta min/max en escala Y');
    if (!(y.min <= dmin && y.max >= dmax)) throw new Error(`min/max ${y.min}/${y.max} no cubren los datos ${dmin}/${dmax}`);
    if (y.min < dmin - (dmax - dmin)) throw new Error('min demasiado bajo (no está cerca del mínimo real)');
    console.log(`  escala Y: min=${y.min} max=${y.max} (datos ${dmin}..${dmax}) ✓`);
  }

  // 3) Slide sin registros muestra empty-state
  const empty = slides[2];
  if (!empty.children[1]._text.includes('Sin registros en este período')) throw new Error('falta empty-state');

  // 4) Filtro por rango: datos de julio no deben aparecer (mes=agosto)
  //    (histories ya vienen filtrados por el servidor-side filter; verificamos labels en rango)
  for (const c of withChart) {
    for (const l of c.cfg.data.labels) {
      if (l < '2026-08-01' || l > '2026-08-31') throw new Error('label fuera de rango: ' + l);
    }
  }

  // 5) Flechas del carousel
  sandbox.scrollTrackerCarousel(1);
  console.log('scrollBy aplicado:', trackerCarousel._scrolled, '(esperado 350)');
  if (trackerCarousel._scrolled !== 350) throw new Error('scrollBy no usó el ancho de la slide');

  // 6) Recarga (cambio de mes): charts viejos destruidos
  elements['metric-month'].value = '7';
  await sandbox.loadDashboard();
  const lineAlive = createdCharts.filter(c => !c.destroyed && c.cfg.type === 'line');
  const lineDestroyed = createdCharts.filter(c => c.destroyed && c.cfg.type === 'line').length;
  console.log('tras recargar mes: line charts activos =', lineAlive.length, ', destruidos =', lineDestroyed);
  if (lineAlive.length !== 0 || lineDestroyed !== 3) throw new Error('no se destruyeron los line charts previos');
  if (trackerCarousel.children.length !== 4) throw new Error('slides tras recarga incorrectas');
  // julio: el filter deja rows vacíos (historias son de agosto) → todos empty-state
  const allEmpty = trackerCarousel.children.every(s => !s.children[1] || s.children[1].tagName === 'P');
  if (!allEmpty) throw new Error('esperaba empty-states para julio');

  console.log('\n✅ SMOKE TEST OK — carousel renderiza 1 gráfico por tracker, línea de límite punteada, empty-state, flechas y recarga destruyen charts previos.');
})().catch(e => { console.error('❌ FALLO:', e.message); process.exit(1); });
