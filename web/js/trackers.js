// MyTickTick - Trackers Diarios (sección 10 del plan)
//
// Sección 10:
//   10.1 Lista de trackers (target, umbral, estado, últimos registros)
//   10.2 Formulario de creación
//   10.3 Configuración de límite unilateral (mín o máx)
//   10.4 Registro de valor diario (modal)
//   10.5 Visualización de cumplimiento / no cumplimiento
//   10.6 Visualización de desviación (valor - target)
//
// Usa el cliente fetchAPI de app.js y los estilos .task-form / .modal
// compartidos con las secciones 8 y 9.
document.addEventListener('DOMContentLoaded', function () {
    loadTrackers();
    setupCreateTrackerButton();
    setupCreateForm();
    setupRecordModal();
    setupHistoryModal();
    setupEditForm();
});

// Nombres de meses en español para mostrar fechas.
const MONTH_NAMES = ['enero', 'febrero', 'marzo', 'abril', 'mayo', 'junio',
    'julio', 'agosto', 'septiembre', 'octubre', 'noviembre', 'diciembre'];

// ---------- 10.1 - Lista de trackers ----------

// Cargar lista de trackers + resumen de últimos 3 registros de cada uno
// (spec daily-trackers: "Ver historial de tracker → muestra todos los
// registros anteriores con fechas").
async function loadTrackers() {
    const container = document.getElementById('tracker-list');
    if (!container) return;
    container.innerHTML = '<p class="loading-state">Cargando trackers…</p>';

    let trackers = [];
    try {
        trackers = (await MyTickTick.fetchAPI('/trackers')) || [];
    } catch (error) {
        console.error('Error loading trackers:', error);
        container.innerHTML = '<div class="error-state">Error al cargar los trackers: ' +
            escapeHtml(error.message) + '</div>';
        return;
    }

    // Últimos 3 registros de cada tracker para el resumen (10.5 / 10.6).
    // Si falla un historial no se bloquea la lista: se asume sin registros.
    const summaries = await Promise.all(trackers.map(function (t) {
        return MyTickTick.fetchAPI('/trackers/' + t.id + '/history')
            .then(function (h) { return (h || []).slice(0, 3); })
            .catch(function (e) {
                console.warn('No se pudo cargar el historial de ' + t.id + ':', e.message);
                return [];
            });
    }));

    displayTrackers(trackers, summaries);
}

// Mostrar trackers en la lista con su configuración y últimos registros.
function displayTrackers(trackers, summaries) {
    const container = document.getElementById('tracker-list');
    if (!container) return;

    if (!trackers || trackers.length === 0) {
        container.innerHTML = '<div class="empty-state">' +
            'No hay trackers todavía. Creá el primero: peso, sueño, agua…' +
            '</div>';
        return;
    }

    const html = trackers.map(function (t, i) {
        const unit = t.unit ? ' ' + escapeHtml(t.unit) : '';
        // Límite unilateral (mín o máx).
        const limitWord = (t.limitType === 'min') ? 'mínimo' : 'máximo';
        const last3 = (summaries && summaries[i]) || [];

        // "Registrado hoy": el primer registro (el más reciente, el backend
        // los devuelve por fecha descendente) cae en la fecha de hoy.
        const latest = last3[0];
        const registeredToday = latest && latest.entryDate === toLocalISODate(new Date());

        // 10.5 / 10.6 - últimos registros con cumplimiento y desviación
        const recent = last3.length
            ? '<div class="tracker-recent">' + last3.map(function (e) {
                return recentBadge(e);
              }).join('') + '</div>'
            : '<p class="task-repeats">Todavía sin registros.</p>';

        return `
        <div class="tracker-item ${registeredToday ? 'tracker-item--done' : ''}" data-id="${t.id}">
            <h3>${escapeHtml(t.name)}</h3>
            <div class="tracker-config">
                <span>Límite ${limitWord}: ${formatNumber(t.limitValue)}${unit}</span>
                ${registeredToday ? '<span class="status-badge status-badge--done" title="Ya cargaste el valor de hoy">✓ Registrado hoy</span>' : ''}
                ${t.isActive ? '' : '<span class="status-badge status-badge--pending">Pausado</span>'}
            </div>
            ${recent}
            <div class="tracker-actions">
                <button class="btn" onclick="openHistory(${t.id})">Historial</button>
                <button class="btn btn-danger" onclick="deleteTracker(${t.id})">Eliminar</button>
            </div>
        </div>
    `;
    }).join('');

    container.innerHTML = html;
    bindTrackerInteractions(container);
}

// 10.7 - toggle/swipe + Deshacer para registrar y revertir valores diarios.
// El registro se hace tocando la card en sí (sin botón "Registrar valor"):
// cualquier clic/tap sobre el ítem abre el modal de registro del día.
// Aquí además agregamos swipe para eliminar el tracker (con toast Deshacer)
// de forma rápida sin confirm().
// El botón "Editar" se quitó de la card: el LONG PRESS (sostener ~500ms en
// cualquier parte de la card) abre la ventana de edición del tracker
// (nombre/límite/unidad), igual que en las secciones de tareas.
function bindTrackerInteractions(container) {
    if (!container) return;
    const items = container.querySelectorAll('.tracker-item[data-id]');
    items.forEach(function (item) {
        const id = Number(item.getAttribute('data-id'));
        MyTickTick.attachToggle(item, { onClick: function () { openRecordForm(id); } });
        MyTickTick.attachSwipe(item, { onSwipe: function () { deleteTracker(id); }, color: 'red' });
        MyTickTick.attachLongPress(item, { onLongPress: function () { openEditForm(id); } });
    });
}

// ---------- 10.2 / 10.3 - Formulario de creación ----------

function openCreateForm() {
    const section = document.getElementById('create-form');
    if (!section) return;
    section.hidden = false;
    clearCreateFormFeedback();
    const nameInput = document.getElementById('tracker-name');
    if (nameInput && nameInput.focus) nameInput.focus();
}

function closeCreateForm() {
    const section = document.getElementById('create-form');
    if (section) section.hidden = true;
}

function setCreateFormFeedback(message, kind) {
    const el = document.getElementById('tracker-form-feedback');
    if (!el) return;
    el.textContent = message;
    el.hidden = false;
    el.className = 'form-feedback ' + (kind === 'success' ? 'success' : (kind === 'error' ? 'error' : ''));
}

function clearCreateFormFeedback() {
    const el = document.getElementById('tracker-form-feedback');
    if (el) {
        el.textContent = '';
        el.hidden = true;
        el.className = 'form-feedback';
    }
}

function showFieldError(id, show) {
    const el = document.getElementById(id);
    if (el) el.hidden = !show;
}

// Validación de nombre (requerido, spec data-model).
function validateTrackerName(prefix) {
    const nameInput = document.getElementById(prefix + '-name');
    const name = nameInput ? nameInput.value.trim() : '';
    showFieldError(prefix + '-name-error', name === '');
    if (name === '') {
        if (nameInput) nameInput.focus();
        return null;
    }
    return name;
}

// Validación de límite (10.3): número finito (mín o máx).
function validateTrackerLimit(prefix) {
    const input = document.getElementById(prefix + '-limit');
    const raw = input ? input.value.trim() : '';
    const value = Number(raw);
    const ok = raw !== '' && isFinite(value);
    showFieldError(prefix + '-limit-error', !ok);
    if (!ok) {
        if (input) input.focus();
        return null;
    }
    return value;
}

// Submit del formulario de creación: valida y crea el tracker (POST).
// Límite unilateral: limitValue + limitType (min | max).
async function submitCreateTracker(event) {
    if (event && event.preventDefault) event.preventDefault();

    const name = validateTrackerName('tracker');
    if (name === null) return;
    const limit = validateTrackerLimit('tracker');
    if (limit === null) return;

    const typeRadio = document.querySelector('input[name="limitType"]:checked');
    const limitType = typeRadio ? typeRadio.value : 'max';

    const unitInput = document.getElementById('tracker-unit');
    const unit = unitInput ? unitInput.value.trim() : '';

    const submitBtn = document.getElementById('submit-create-tracker-btn');
    if (submitBtn) submitBtn.disabled = true;
    setCreateFormFeedback('Creando tracker…', '');

    try {
        const tracker = await MyTickTick.fetchAPI('/trackers', {
            method: 'POST',
            json: {
                name: name,
                limitValue: limit,
                limitType: limitType,
                unit: unit
            }
        });
        console.log('Tracker creado:', tracker);
        const form = document.getElementById('create-tracker-form');
        if (form && form.reset) form.reset();
        setCreateFormFeedback('Tracker creado. Ya podés registrar el primer valor.', 'success');
        closeCreateForm();
        loadTrackers();
    } catch (error) {
        console.error('Error creating tracker:', error);
        setCreateFormFeedback('Error al crear el tracker: ' + error.message, 'error');
    } finally {
        if (submitBtn) submitBtn.disabled = false;
    }
}

function setupCreateTrackerButton() {
    const btn = document.getElementById('create-tracker-btn');
    if (btn) {
        btn.addEventListener('click', function () {
            const section = document.getElementById('create-form');
            if (section && !section.hidden) closeCreateForm();
            else openCreateForm();
        });
    }
    const cancelBtn = document.getElementById('cancel-create-tracker-btn');
    if (cancelBtn) cancelBtn.addEventListener('click', closeCreateForm);
}

function setupCreateForm() {
    const form = document.getElementById('create-tracker-form');
    if (form) form.addEventListener('submit', submitCreateTracker);
}

// ---------- 10.4 - Registro de valor diario ----------

// Abrir el modal de registro precargado con la fecha de hoy.
function openRecordForm(trackerId) {
    const modal = document.getElementById('record-modal');
    if (!modal) return;
    modal.dataset.trackerId = String(trackerId);

    // Si quedó un cierre programado por un guardado anterior, lo cancelamos
    // para que no cierre el modal recién abierto.
    if (modal._closeTimer) {
        clearTimeout(modal._closeTimer);
        modal._closeTimer = null;
    }

    // Título + contexto: target y rango cumplido del tracker.
    MyTickTick.fetchAPI('/trackers/' + trackerId).then(function (t) {
        const nameEl = document.getElementById('record-tracker-name');
        if (nameEl) nameEl.textContent = 'Registrar valor — ' + (t ? t.name : '');
        const hint = document.getElementById('record-target-hint');
        if (hint && t) {
            const unit = t.unit ? ' ' + escapeHtml(t.unit) : '';
            const word = t.limitType === 'min' ? 'mínimo' : 'máximo';
            hint.textContent = 'Límite ' + word + ': ' + formatNumber(t.limitValue) + unit;
        }
        // Placeholder del campo de valor: en vez de un "Ej: 79.5" genérico
        // (que no tenía relación con el tracker), se muestra como ejemplo el
        // límite/threshold configurado para ESTE tracker (76 para peso, 60
        // para sueño, 16 para ayuno). Si no tiene límite numérico, se deja el
        // placeholder original.
        const valueInput = document.getElementById('record-value');
        if (valueInput && t && isFinite(Number(t.limitValue)) && t.limitValue !== '') {
            valueInput.placeholder = 'Ej: ' + formatNumber(t.limitValue);
        }
    }).catch(function (e) {
        console.warn('No se pudo cargar el tracker para el registro:', e.message);
    });

    const dateInput = document.getElementById('record-date');
    if (dateInput) dateInput.value = toLocalISODate(new Date());

    const valueInput = document.getElementById('record-value');
    if (valueInput) {
        valueInput.value = '';
        valueInput.focus();
    }
    const notesInput = document.getElementById('record-notes');
    if (notesInput) notesInput.value = '';

    setRecordFeedback('', '');
    closeCreateForm();
    modal.hidden = false;
}

function closeRecordModal() {
    const modal = document.getElementById('record-modal');
    if (modal) modal.hidden = true;
}

function setRecordFeedback(message, kind) {
    const el = document.getElementById('record-feedback');
    if (!el) return;
    el.textContent = message;
    el.hidden = message === '';
    el.className = 'form-feedback ' + (kind === 'success' ? 'success' : (kind === 'error' ? 'error' : ''));
}

// Validación del valor del registro (10.4): número finito.
function validateRecordValue() {
    const input = document.getElementById('record-value');
    const raw = input ? input.value.trim() : '';
    const value = Number(raw);
    const ok = raw !== '' && isFinite(value);
    showFieldError('record-value-error', !ok);
    if (!ok) {
        if (input) input.focus();
        return null;
    }
    return value;
}

// Fecha del registro: YYYY-MM-DD local (por defecto hoy).
function recordDateValue() {
    const input = document.getElementById('record-date');
    const value = input && input.value ? input.value : toLocalISODate(new Date());
    return value;
}

// Guardar el registro: PUT /trackers/{id}/records (upsert por fecha).
// Si ya hay un registro para esa fecha, lo ACTUALIZA (permite corregir un valor
// marcado por error); si no, lo crea. El backend calcula isMet (10.5) y
// deviation (10.6); se muestra el resultado de inmediato.
async function submitRecord(event) {
    if (event && event.preventDefault) event.preventDefault();

    const value = validateRecordValue();
    if (value === null) return;
    const date = recordDateValue();

    const modal = document.getElementById('record-modal');
    const trackerId = modal ? modal.dataset.trackerId : null;
    if (!trackerId) return;

    const notesInput = document.getElementById('record-notes');
    const notes = notesInput ? notesInput.value.trim() : '';

    const submitBtn = document.getElementById('submit-record-btn');
    if (submitBtn) submitBtn.disabled = true;
    setRecordFeedback('Guardando registro…', '');

    try {
        const entry = await MyTickTick.fetchAPI('/trackers/' + trackerId + '/records', {
            method: 'PUT',
            json: { value: value, date: date, notes: notes }
        });
        console.log('Registro guardado (upsert):', entry);

        // 10.5 / 10.6 - feedback inmediato con cumplimiento y desviación.
        const dev = entry ? entry.deviation : 0;
        const met = entry ? entry.isMet : false;
        setRecordFeedback(
            'Registro guardado. ' + (met ? '✓ Cumplió el target' : '✗ No cumplió el target') +
            ' · Desviación: ' + formatDeviation(dev),
            met ? 'success' : 'error'
        );
        const form = document.getElementById('record-form');
        if (form && form.reset) form.reset();
        if (document.getElementById('record-date')) {
            document.getElementById('record-date').value = date;
        }
        // Refrescar el resumen de la lista (últimos registros).
        loadTrackers();
        // Cerrar el modal a los 1.2 s para que se vea el feedback (card
        // "Registrado hoy" + badge de cumplimiento). El error NO cierra:
        // deja el modal abierto para que el usuario lo corrija.
        if (modal) {
            if (modal._closeTimer) clearTimeout(modal._closeTimer);
            modal._closeTimer = setTimeout(function () {
                modal._closeTimer = null;
                closeRecordModal();
            }, 1200);
        }
    } catch (error) {
        console.error('Error creating record:', error);
        setRecordFeedback('Error al guardar el registro: ' + error.message, 'error');
    } finally {
        if (submitBtn) submitBtn.disabled = false;
    }
}

function setupRecordModal() {
    const modal = document.getElementById('record-modal');
    if (!modal) return;

    const closeBtn = document.getElementById('close-record');
    if (closeBtn) closeBtn.addEventListener('click', closeRecordModal);
    const cancelBtn = document.getElementById('cancel-record-btn');
    if (cancelBtn) cancelBtn.addEventListener('click', closeRecordModal);
    const form = document.getElementById('record-form');
    if (form) form.addEventListener('submit', submitRecord);

    modal.addEventListener('click', (e) => { if (e.target === modal) closeRecordModal(); });
}

// ---------- 10.5 / 10.6 - Historial: cumplimiento y desviación ----------

// Insignia de un registro: ✓/✗ + desviación con signo (10.5 / 10.6).
function recentBadge(entry) {
    const met = entry.isMet;
    return `<span class="status-badge ${met ? 'status-badge--done' : 'met-no'}" title="${escapeHtml(entry.entryDate)}">` +
        (met ? '✓ ' : '✗ ') +
        formatNumber(entry.value) +
        ' (' + formatDeviation(entry.deviation) + ')</span>';
}

// Desviación con signo explícito: +2.5, -1, 0.
function formatDeviation(deviation) {
    const d = Number(deviation);
    if (!isFinite(d)) return '—';
    const rounded = Math.round(d * 100) / 100;
    if (rounded > 0) return '+' + formatNumber(rounded);
    return formatNumber(rounded);
}

// Número sin ceros de basura: 80 → "80", 79.5 → "79.5".
function formatNumber(value) {
    const n = Number(value);
    if (!isFinite(n)) return String(value);
    return String(Math.round(n * 1000) / 1000);
}

// Fecha YYYY-MM-DD → "15 agosto 2026".
function formatDate(isoDate) {
    const parts = String(isoDate || '').split('-');
    if (parts.length !== 3) return String(isoDate || '');
    const day = parseInt(parts[2], 10);
    const month = parseInt(parts[1], 10);
    const year = parts[0];
    const name = MONTH_NAMES[month - 1];
    return day + ' ' + (name || '') + ' ' + year;
}

// Abrir el historial completo del tracker (10.5 / 10.6):
// cada registro con fecha, valor, cumplimiento y desviación,
// más recientes primero (orden del backend).
async function openHistory(trackerId) {
    const modal = document.getElementById('history-modal');
    if (!modal) return;

    const list = document.getElementById('history-list');
    const nameEl = document.getElementById('history-tracker-name');
    const hintEl = document.getElementById('history-range-hint');
    if (nameEl) nameEl.textContent = 'Historial';
    if (hintEl) hintEl.textContent = '';
    if (list) list.innerHTML = '<p class="loading-state">Cargando historial…</p>';

    let tracker = null;
    let history = [];
    try {
        [tracker, history] = await Promise.all([
            MyTickTick.fetchAPI('/trackers/' + trackerId),
            MyTickTick.fetchAPI('/trackers/' + trackerId + '/history')
        ]);
    } catch (e) {
        console.error('Error loading history:', e);
        if (list) {
            list.innerHTML = '<div class="error-state">Error al cargar el historial: ' +
                escapeHtml(e.message) + '</div>';
        }
        modal.hidden = false;
        return;
    }

    if (nameEl && tracker) nameEl.textContent = 'Historial de ' + tracker.name;

    // Contexto: límite unilateral y % de cumplimiento total.
    if (hintEl && tracker && history.length) {
        const unit = tracker.unit ? ' ' + tracker.unit : '';
        const word = tracker.limitType === 'min' ? 'mínimo' : 'máximo';
        const metCount = history.filter((h) => h.isMet).length;
        const pct = Math.round((metCount / history.length) * 100);
        hintEl.textContent = 'Límite ' + word + ': ' + formatNumber(tracker.limitValue) + unit +
            ' · ' + metCount + ' de ' + history.length + ' registros cumplidos (' + pct + '%)';
    }

    history = history || [];
    if (list) {
        if (history.length === 0) {
            list.innerHTML = '<div class="empty-state">Todavía no hay registros. ' +
                'Tocá la card del tracker para cargar el primero.</div>';
        } else {
            list.innerHTML = history.map(function (entry) {
                const met = entry.isMet;
                const badge = met
                    ? '<span class="status-badge status-badge--done">✓ Cumplió</span>'
                    : '<span class="status-badge met-no">✗ No cumplió</span>';
                return `
                <div class="history-entry ${met ? 'met' : 'not-met'}">
                    <span class="history-month-year">${formatDate(entry.entryDate)}</span>
                    <span>Valor: ${formatNumber(entry.value)} · Desviación: ${formatDeviation(entry.deviation)}</span>
                    <span class="history-entry-actions">${badge}</span>
                </div>`;
            }).join('');
        }
    }

    closeCreateForm();
    modal.hidden = false;
}

function closeHistoryModal() {
    const modal = document.getElementById('history-modal');
    if (modal) modal.hidden = true;
}

function setupHistoryModal() {
    const modal = document.getElementById('history-modal');
    if (!modal) return;

    const closeBtn = document.getElementById('close-history');
    if (closeBtn) closeBtn.addEventListener('click', closeHistoryModal);

    modal.addEventListener('click', (e) => { if (e.target === modal) closeHistoryModal(); });

    // Tecla Escape cierra el modal visible.
    document.addEventListener('keydown', (e) => {
        if (e.key !== 'Escape') return;
        const rec = document.getElementById('record-modal');
        if (rec && !rec.hidden) closeRecordModal();
        const edit = document.getElementById('edit-modal');
        if (edit && !edit.hidden) closeEditForm();
        if (!modal.hidden) closeHistoryModal();
    });
}

// ---------- Eliminar (con confirmación + undo) ----------
// Pedir confirmación antes de borrar porque se pierde todo el historial
// (registraciones diarias). Si el usuario confirma, el toast "Deshacer"
// sigue disponible como red de seguridad.

async function deleteTracker(id) {
    // Cargar el tracker (para el mensaje y para el undo).
    let tracker = null;
    try {
        tracker = await MyTickTick.fetchAPI('/trackers/' + id);
    } catch (e) {
        console.error('No se pudo cargar el tracker para eliminar:', e);
    }

    // Confirmar (12.x / pedido del usuario): evita borrar por accidente.
    var name = tracker ? tracker.name : '';
    const ok = await MyTickTick.confirm({
        title: '¿Eliminar tracker' + (name ? ': ' + name : '') + '?',
        message: 'Se eliminará el tracker y TODO su historial de registros diarios. Esta acción no se puede deshacer por el sistema (el botón "Deshacer" solo recrea el tracker, no los registros).',
        confirmLabel: 'Eliminar',
        cancelLabel: 'Cancelar',
        kind: 'danger'
    });
    if (!ok) return;

    try {
        await MyTickTick.fetchAPI('/trackers/' + id, { method: 'DELETE' });
        closeCreateForm();

        MyTickTick.toast('Tracker eliminado', {
            kind: 'error',
            actionLabel: 'Deshacer',
            onUndo: function () {
                if (!tracker) { loadTrackers(); return; }
                MyTickTick.fetchAPI('/trackers', {
                    method: 'POST',
                    json: {
                        name: tracker.name,
                        limitValue: tracker.limitValue,
                        limitType: tracker.limitType || 'max',
                        unit: tracker.unit || ''
                    }
                }).then(loadTrackers).catch(function (e) {
                    console.error('No se pudo restaurar el tracker:', e);
                    MyTickTick.toast('No se pudo restaurar el tracker', { kind: 'error' });
                });
            },
            timeout: 8000
        });

        loadTrackers();
    } catch (error) {
        console.error('Error deleting tracker:', error);
        MyTickTick.toast('Error al eliminar el tracker: ' + error.message, { kind: 'error' });
    }
}

// ---------- Editar tracker (sin tocar el historial) ----------
// El backend PUT /api/trackers/{id} solo actualiza la fila del tracker
// (nombre, límite, tipo, unidad). Los registros diarios viven en
// tracker_entry_dbs y no se tocan: el historial queda intacto y los
// cumplimientos de los días anteriores conservan el isMet/deviation
// calculados en su momento contra el límite de entonces.

function openEditForm(trackerId) {
    const modal = document.getElementById('edit-modal');
    if (!modal) return;
    modal.dataset.trackerId = String(trackerId);

    // Precarga con los datos actuales del tracker.
    MyTickTick.fetchAPI('/trackers/' + trackerId).then(function (t) {
        if (!t) return;
        const nameEl = document.getElementById('edit-tracker-name');
        if (nameEl) nameEl.textContent = 'Editar — ' + t.name;

        const nameInput = document.getElementById('edit-tracker-name-input');
        if (nameInput) nameInput.value = t.name || '';

        const limitInput = document.getElementById('edit-tracker-limit');
        if (limitInput) limitInput.value = t.limitValue;

        const typeRadio = document.querySelector('input[name="edit-limitType"][value="' + (t.limitType || 'max') + '"]');
        if (typeRadio) typeRadio.checked = true;

        const unitInput = document.getElementById('edit-tracker-unit');
        if (unitInput) unitInput.value = t.unit || '';

        setEditFeedback('', '');
        if (nameInput && nameInput.focus) nameInput.focus();
    }).catch(function (e) {
        console.error('No se pudo cargar el tracker para editar:', e);
        MyTickTick.toast('Error al cargar el tracker: ' + e.message, { kind: 'error' });
    });

    clearEditErrors();
    closeCreateForm();
    modal.hidden = false;
}

function closeEditForm() {
    const modal = document.getElementById('edit-modal');
    if (modal) modal.hidden = true;
}

function setEditFeedback(message, kind) {
    const el = document.getElementById('edit-feedback');
    if (!el) return;
    el.textContent = message;
    el.hidden = message === '';
    el.className = 'form-feedback ' + (kind === 'success' ? 'success' : (kind === 'error' ? 'error' : ''));
}

function clearEditErrors() {
    ['edit-tracker-name-error', 'edit-tracker-limit-error'].forEach(function (id) {
        const el = document.getElementById(id);
        if (el) el.hidden = true;
    });
}

function validateEditName() {
    const el = document.getElementById('edit-tracker-name-input');
    const v = el ? el.value.trim() : '';
    const ok = v !== '';
    document.getElementById('edit-tracker-name-error').hidden = ok;
    if (!ok && el) el.focus();
    return ok ? v : null;
}

function validateEditLimit() {
    const el = document.getElementById('edit-tracker-limit');
    const raw = el ? el.value.trim() : '';
    const v = Number(raw);
    const ok = raw !== '' && isFinite(v);
    document.getElementById('edit-tracker-limit-error').hidden = ok;
    if (!ok && el) el.focus();
    return ok ? v : null;
}

async function submitEditTracker(event) {
    if (event && event.preventDefault) event.preventDefault();

    const name = validateEditName();
    if (name === null) return;
    const limit = validateEditLimit();
    if (limit === null) return;
    const typeRadio = document.querySelector('input[name="edit-limitType"]:checked');
    const limitType = typeRadio ? typeRadio.value : 'max';
    const unitEl = document.getElementById('edit-tracker-unit');
    const unit = unitEl ? unitEl.value.trim() : '';

    const modal = document.getElementById('edit-modal');
    const trackerId = modal ? modal.dataset.trackerId : null;
    if (!trackerId) return;

    const btn = document.getElementById('submit-edit-tracker-btn');
    if (btn) btn.disabled = true;
    setEditFeedback('Guardando…', '');

    try {
        const updated = await MyTickTick.fetchAPI('/trackers/' + trackerId, {
            method: 'PUT',
            json: {
                name: name,
                limitValue: limit,
                limitType: limitType,
                unit: unit,
                isActive: true
            }
        });
        console.log('Tracker actualizado:', updated);
        setEditFeedback('Cambios guardados. El historial quedó intacto.', 'success');
        closeEditForm();
        loadTrackers();
    } catch (error) {
        console.error('Error actualizando tracker:', error);
        setEditFeedback('Error al guardar los cambios: ' + error.message, 'error');
    } finally {
        if (btn) btn.disabled = false;
    }
}

function setupEditForm() {
    const modal = document.getElementById('edit-modal');
    if (!modal) return;
    modal.addEventListener('click', function (e) { if (e.target === modal) closeEditForm(); });
    const closeBtn = document.getElementById('close-edit');
    if (closeBtn) closeBtn.addEventListener('click', closeEditForm);
    const cancelBtn = document.getElementById('cancel-edit-tracker-btn');
    if (cancelBtn) cancelBtn.addEventListener('click', closeEditForm);
    const form = document.getElementById('edit-tracker-form');
    if (form) form.addEventListener('submit', submitEditTracker);
}

// ---------- Utilidades ----------

// Fecha local → YYYY-MM-DD (sin desfase UTC, a diferencia de toISOString).
function toLocalISODate(d) {
    const y = d.getFullYear();
    const m = String(d.getMonth() + 1).padStart(2, '0');
    const day = String(d.getDate()).padStart(2, '0');
    return y + '-' + m + '-' + day;
}

// Escape de HTML para evitar inyección al renderizar nombres.
function escapeHtml(value) {
    return String(value)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}
