// Smoke test del modal de registro de valor (card "Arreglar ventana de
// registro diario"): al guardar, el modal debe cerrarse automáticamente
// (tras mostrar el feedback), y en caso de error debe seguir abierto.
// Stubs de DOM + MyTickTick + timers controlados, sin browser.
const fs = require('fs');
const vm = require('vm');
const path = require('path');

const code = fs.readFileSync(path.join(__dirname, '..', 'web', 'js', 'trackers.js'), 'utf8');

let failures = 0;
function assert(cond, msg) {
  if (cond) {
    console.log('  ok - ' + msg);
  } else {
    failures++;
    console.log('  FAIL - ' + msg);
  }
}

// ---- timers controlados ----
let now = 0;
let timerId = 0;
const timers = [];
function fakeSetTimeout(fn, ms) {
  const id = ++timerId;
  timers.push({ id, at: now + (ms || 0), fn, done: false });
  return id;
}
function fakeClearTimeout(id) {
  const t = timers.find((x) => x.id === id);
  if (t) t.done = true;
}
function flush(ms) {
  now += ms;
  for (const t of timers) {
    if (!t.done && t.at <= now) {
      t.done = true;
      t.fn();
    }
  }
}

// ---- stubs de DOM ----
const elements = {};
function makeEl(tag) {
  return {
    tagName: (tag || 'div').toUpperCase(),
    children: [],
    dataset: {},
    style: {},
    hidden: false,
    _text: '',
    _innerHTML: '',
    listeners: {},
    getContext() { return { _el: this }; },
    appendChild(c) { this.children.push(c); return c; },
    querySelector() { return this.children[0] || null; },
    querySelectorAll() { return this._items || []; },
    get innerHTML() { return this._innerHTML; },
    set innerHTML(v) { this._innerHTML = String(v); },
    get textContent() { return this._text; },
    set textContent(v) { this._text = String(v); },
    addEventListener(ev, fn) { (this.listeners[ev] = this.listeners[ev] || []).push(fn); },
    focus() {},
    reset() { this.value = ''; },
    value: '',
    disabled: false
  };
}
function el(id) {
  if (!elements[id]) elements[id] = makeEl('div');
  return elements[id];
}

const trackerItem = makeEl('div');
trackerItem.setAttribute = () => {};
const trackerList = makeEl('div');
trackerList._items = [Object.assign(trackerItem, { getAttribute: () => '1' })];
elements['tracker-list'] = trackerList;
elements['record-modal'] = makeEl('div');
elements['record-date'] = makeEl('input');
elements['record-date'].value = '2026-08-31';
elements['record-value'] = makeEl('input');

const documentStub = {
  addEventListener(ev, fn) { if (ev === 'DOMContentLoaded') documentStub._domReady = fn; },
  getElementById(id) { return el(id); },
  querySelector() { return null; },
  createElement(tag) { return makeEl(tag); }
};

// ---- stub MyTickTick ----
const TR = { id: 1, name: 'Peso', unit: 'kg', limitValue: 80, limitType: 'max', isActive: true };
let putShouldFail = false;
const fetchCalls = [];
const myTickTick = {
  fetchAPI(url, opts) {
    fetchCalls.push((opts && opts.method) + ' ' + url);
    const method = (opts && opts.method) || 'GET';
    if (method === 'GET' && url === '/trackers') return Promise.resolve([TR]);
    if (method === 'GET' && url === '/trackers/1') return Promise.resolve(TR);
    if (method === 'GET' && url === '/trackers/1/history') return Promise.resolve([]);
    if (method === 'PUT' && url === '/trackers/1/records') {
      if (putShouldFail) return Promise.reject(new Error('boom'));
      return Promise.resolve({ value: 76, date: '2026-08-31', isMet: false, deviation: 4 });
    }
    return Promise.resolve({});
  },
  attachToggle() {}, attachSwipe() {}, attachLongPress() {},
  confirm() { return Promise.resolve(true); },
  toast() {}
};

const sandbox = {
  document: documentStub,
  MyTickTick: myTickTick,
  console,
  setTimeout: fakeSetTimeout,
  clearTimeout: fakeClearTimeout
};
vm.createContext(sandbox);
vm.runInContext(code, sandbox);
const T = sandbox; // funciones top-level quedan en el contexto

function run() {
  // DOM ready → setup de listeners.
  documentStub._domReady();

  const modal = elements['record-modal'];
  const feedback = () => elements['record-feedback'];

  // 1) Abrir el modal de registro. El fetch de /trackers/1 resuelve en un
  //    microtask; esperar un macrotask para que corra y fije el placeholder.
  T.openRecordForm(1);
  return new Promise(function (res) { setTimeout(res, 0); }).then(function () {
  assert(modal.hidden === false, 'modal se abre al tocar la card');
  assert(modal.dataset.trackerId === '1', 'modal guarda el trackerId');
  // Card "Mejorar ejemplo tracker": el ejemplo del campo de valor debe ser el
  // threshold configurado del tracker (80), no un "Ej: 79.5" genérico.
  assert(elements['record-value'].placeholder === 'Ej: 80',
    'el placeholder del campo de valor usa el threshold del tracker (Ej: 80), no un ejemplo fijo');

  // 2) Guardar con éxito → el modal DEBE cerrarse solo (tras el feedback).
  elements['record-value'].value = '76';
  return T.submitRecord({ preventDefault() {} }).then(function () {
    flush(500);
    assert(feedback()._text.indexOf('Registro guardado') === 0, 'feedback de guardado se muestra antes de cerrar');
    assert(modal.hidden === false, 'el modal sigue abierto brevemente para ver el feedback');
    flush(800); // total 1300 ms > 1200 ms programados
    assert(modal.hidden === true, 'el modal se cierra automáticamente tras guardar');

    // 3) Reabrir con un cierre programado pendiente → no debe auto-cerrarse
    //    el modal recién abierto.
    elements['record-value'].value = '75';
    putShouldFail = false;
    T.openRecordForm(1);
    flush(100);
    assert(modal.hidden === false, 'modal reabierto no se cierra por timer pendiente');

    // 4) Guardar con error → el modal DEBE seguir abierto.
    putShouldFail = true;
    elements['record-value'].value = '74';
    return T.submitRecord({ preventDefault() {} }).then(function () {
      flush(2000);
      assert(modal.hidden === false, 'con error el modal sigue abierto para corregir');
      assert(feedback()._text.indexOf('Error al guardar') === 0, 'feedback de error se muestra');
      putShouldFail = false;
    });
  }).then(function () {
    if (failures > 0) {
      console.log(failures + ' fallo(s) en test-record-modal');
      process.exit(1);
    }
    console.log('test-record-modal: todo OK');
  }).catch(function (e) {
    console.log('test-record-modal: excepción inesperada: ' + (e && e.stack));
    process.exit(1);
  });
  });
}

// Dejar pasar los microtareas del openRecordForm (fetch del nombre) primero.
setTimeout(run, 10);
