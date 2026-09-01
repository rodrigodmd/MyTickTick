// Smoke test de edición de tareas inmediatas: la card no muestra un botón
// "Editar" y la pulsación larga abre el formulario con la tarea seleccionada.
const fs = require('fs');
const vm = require('vm');
const path = require('path');

const code = fs.readFileSync(path.join(__dirname, '..', 'web', 'js', 'immediate-tasks.js'), 'utf8');
const page = fs.readFileSync(path.join(__dirname, '..', 'web', 'html', 'immediate-tasks.html'), 'utf8');

let failures = 0;
function assert(condition, message) {
  if (condition) {
    console.log('  ok - ' + message);
  } else {
    failures++;
    console.log('  FAIL - ' + message);
  }
}

const elements = {};
function makeEl(tag) {
  return {
    tagName: (tag || 'div').toUpperCase(),
    children: [],
    dataset: {},
    style: {},
    hidden: true,
    value: '',
    className: '',
    disabled: false,
    _innerHTML: '',
    _text: '',
    listeners: {},
    appendChild(child) { this.children.push(child); return child; },
    querySelector() { return null; },
    querySelectorAll() { return this._items || []; },
    getAttribute() { return null; },
    addEventListener(event, handler) {
      (this.listeners[event] = this.listeners[event] || []).push(handler);
    },
    focus() { this._focused = true; },
    reset() {},
    get innerHTML() { return this._innerHTML; },
    set innerHTML(value) { this._innerHTML = String(value); },
    get textContent() { return this._text; },
    set textContent(value) { this._text = String(value); }
  };
}
function el(id) {
  if (!elements[id]) elements[id] = makeEl('div');
  return elements[id];
}

const taskItem = makeEl('div');
taskItem.getAttribute = (name) => name === 'data-id' ? '42' : null;
taskItem.classList = { contains() { return false; } };
const taskList = makeEl('div');
taskList._items = [taskItem];
elements['task-list'] = taskList;

const documentStub = {
  addEventListener(event, handler) {
    if (event === 'DOMContentLoaded') this._domReady = handler;
  },
  getElementById(id) { return el(id); },
  querySelector() { return null; },
  createElement(tag) { return makeEl(tag); }
};

let longPressHandler = null;
const fetchCalls = [];
const task = {
  id: 42,
  name: 'Renovar pasaporte',
  description: 'Pedir turno',
  dueDate: '2026-09-10T23:59:59Z',
  priority: 'high',
  isCompleted: false
};
const myTickTick = {
  fetchAPI(url) {
    fetchCalls.push(url);
    if (url === '/immediate-tasks/42') return Promise.resolve(task);
    if (url === '/immediate-tasks') return Promise.resolve([task]);
    return Promise.resolve({});
  },
  attachToggle() {},
  attachSwipe() {},
  attachLongPress(item, options) {
    if (item === taskItem) longPressHandler = options.onLongPress;
  },
  confirm() { return Promise.resolve(true); },
  toast() {}
};

const sandbox = {
  document: documentStub,
  MyTickTick: myTickTick,
  console,
  setTimeout,
  clearTimeout,
  Date,
  alert() {}
};
vm.createContext(sandbox);
vm.runInContext(code, sandbox);

async function run() {
  sandbox.displayTasks([task]);

  assert(!/>\s*Editar\s*</i.test(taskList.innerHTML),
    'la card no renderiza un botón Editar');
  assert(!/>\s*Eliminar\s*</i.test(taskList.innerHTML),
    'la card pendiente no renderiza un botón Eliminar');
  assert(typeof longPressHandler === 'function',
    'la card registra la edición mediante pulsación larga');

  longPressHandler();
  await new Promise((resolve) => setImmediate(resolve));

  assert(fetchCalls.includes('/immediate-tasks/42'),
    'la pulsación larga carga la tarea seleccionada');
  assert(elements['edit-form'].hidden === false,
    'la pulsación larga abre el formulario de edición');
  assert(elements['edit-form'].dataset.taskId === '42',
    'el formulario conserva el id de la tarea seleccionada');
  assert(elements['edit-task-name'].value === task.name,
    'el formulario se completa con los datos de la tarea');
  assert(/id="delete-edit-btn"[^>]*>Eliminar<\/button>/.test(page),
    'el formulario de edición conserva su botón Eliminar');

  const completedTask = Object.assign({}, task, { isCompleted: true });
  vm.runInContext("viewState.status = 'completed'", sandbox);
  sandbox.displayTasks([completedTask]);
  assert(/>\s*Eliminar\s*</i.test(taskList.innerHTML),
    'la card completada renderiza el botón Eliminar');
  assert(/completed-delete-action/.test(taskList.innerHTML),
    'la eliminación rápida queda identificada como acción de completada');

  if (failures > 0) {
    console.log(failures + ' fallo(s) en test-immediate-task-edit');
    process.exit(1);
  }
  console.log('test-immediate-task-edit: todo OK');
}

run().catch(function (error) {
  console.log('test-immediate-task-edit: excepción inesperada: ' + (error && error.stack));
  process.exit(1);
});
