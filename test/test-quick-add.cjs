// Smoke test de la carga rápida de tareas inmediatas (card "Mejorar carga de
// tareas inemediatas"): el campo superior crea la tarea con SOLO el nombre
// (Enter), usando fecha límite +7 días y prioridad media por defecto.
// Stubs de DOM + MyTickTick, sin browser (mismo patrón que test-record-modal.cjs).
const fs = require('fs');
const vm = require('vm');
const path = require('path');

const code = fs.readFileSync(path.join(__dirname, '..', 'web', 'js', 'immediate-tasks.js'), 'utf8');

let failures = 0;
function assert(cond, msg) {
  if (cond) {
    console.log('  ok - ' + msg);
  } else {
    failures++;
    console.log('  FAIL - ' + msg);
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
    appendChild(c) { this.children.push(c); return c; },
    querySelector() { return this.children[0] || null; },
    querySelectorAll() { return this._items || []; },
    getAttribute() { return null; },
    get innerHTML() { return this._innerHTML; },
    set innerHTML(v) { this._innerHTML = String(v); },
    get textContent() { return this._text; },
    set textContent(v) { this._text = String(v); },
    addEventListener(ev, fn) { (this.listeners[ev] = this.listeners[ev] || []).push(fn); },
    focus() { this._focused = true; },
    reset() { this.value = ''; },
    value: '',
    disabled: false
  };
}
function el(id) {
  if (!elements[id]) elements[id] = makeEl('div');
  return elements[id];
}

const quickInput = makeEl('input');
elements['quick-add-input'] = quickInput;
const quickBtn = makeEl('button');
elements['quick-add-btn'] = quickBtn;
const quickForm = makeEl('form');
elements['quick-add-form'] = quickForm;
elements['task-list'] = makeEl('div');

const documentStub = {
  addEventListener(ev, fn) { if (ev === 'DOMContentLoaded') documentStub._domReady = fn; },
  // La pantalla no muestra el botón "Nueva Tarea": la carga rápida es la
  // única entrada visible, y la inicialización debe tolerar que falte el botón.
  getElementById(id) { return id === 'create-task-btn' ? null : el(id); },
  querySelector() { return null; },
  createElement(tag) { return makeEl(tag) }
};

// ---- stub MyTickTick ----
const postCalls = [];
let postShouldFail = false;
const toasts = [];
const myTickTick = {
  fetchAPI(url, opts) {
    const method = (opts && opts.method) || 'GET';
    if (method === 'GET' && url === '/immediate-tasks') {
      return Promise.resolve([]);
    }
    if (method === 'POST' && url === '/immediate-tasks') {
      postCalls.push(opts.json);
      if (postShouldFail) return Promise.reject(new Error('boom'));
      return Promise.resolve({ id: 1, name: opts.json.name });
    }
    return Promise.resolve({});
  },
  attachToggle() {}, attachSwipe() {}, attachLongPress() {},
  confirm() { return Promise.resolve(true); },
  toast(msg, opts) { toasts.push({ msg: msg, kind: (opts && opts.kind) || 'default' }); }
};

const sandbox = {
  document: documentStub,
  MyTickTick: myTickTick,
  console,
  setTimeout, clearTimeout,
  Date
};
vm.createContext(sandbox);
vm.runInContext(code, sandbox);
const T = sandbox;

function isoPlusDays(days) {
  return new Date(Date.now() + days * 86400 * 1000).toISOString().split('T')[0];
}

function run() {
  // DOM ready → setup de listeners.
  documentStub._domReady();
  assert((quickForm.listeners['submit'] || []).length === 1, 'el form rápido tiene el listener de submit');

  const input = elements['quick-add-input'];
  const btn = elements['quick-add-btn'];

  // 1) Nombre vacío → no dispara POST.
  input.value = '   ';
  return T.submitQuickAdd({ preventDefault() {} }).then(function () {
    assert(postCalls.length === 0, 'con nombre vacío no se envía POST');

    // 2) Nombre + Enter → POST con solo el nombre relevante:
    //    dueDate = hoy+7 días (fin del día), priority = medium.
    input.value = 'Renovar pasaporte';
    const expectedDate = isoPlusDays(7) + 'T23:59:59Z';
    return T.submitQuickAdd({ preventDefault() {} }).then(function () {
      assert(postCalls.length === 1, 'con nombre se envía un solo POST');
      const payload = postCalls[0] || {};
      assert(payload.name === 'Renovar pasaporte', 'el payload lleva el nombre escrito');
      assert(payload.priority === 'medium', 'prioridad por defecto: media');
      assert(payload.dueDate === expectedDate, 'fecha límite por defecto: +7 días (' + payload.dueDate + ')');
      assert(input.value === '', 'el input se limpia tras crear');
      assert(btn.disabled === false, 'el botón queda habilitado de nuevo');
      const okToast = toasts.find((t) => t.kind === 'success');
      assert(!!okToast && okToast.msg.indexOf('Tarea agregada') === 0, 'toast de éxito al crear');

      // 3) Error de red → toast de error y el nombre SE QUEDA para reintentar.
      postShouldFail = true;
      input.value = 'Comprar leche';
      return T.submitQuickAdd({ preventDefault() {} }).then(function () {
        const errToast = toasts.find((t) => t.kind === 'error');
        assert(!!errToast && errToast.msg.indexOf('Error al agregar') === 0, 'toast de error al fallar');
        assert(input.value === 'Comprar leche', 'con error el input no se limpia (para reintentar)');
        postShouldFail = false;
      });
    });
  }).then(function () {
    if (failures > 0) {
      console.log(failures + ' fallo(s) en test-quick-add');
      process.exit(1);
    }
    console.log('test-quick-add: todo OK');
  }).catch(function (e) {
    console.log('test-quick-add: excepción inesperada: ' + (e && e.stack));
    process.exit(1);
  });
}

setTimeout(run, 10);
