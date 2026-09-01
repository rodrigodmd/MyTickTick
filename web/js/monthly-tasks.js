// MyTickTick - Tareas Mensuales (usa el cliente fetchAPI de app.js)
document.addEventListener('DOMContentLoaded', function () {
    loadTasks();
    setupCreateTaskButton();
    setupCreateForm();
    setupEditForm();
    setupHistoryModal();
});

// Nombres de meses en español para mostrar fechas/historial.
const MONTH_NAMES = ['enero', 'febrero', 'marzo', 'abril', 'mayo', 'junio',
    'julio', 'agosto', 'septiembre', 'octubre', 'noviembre', 'diciembre'];

// Periodo actual (mes/año) según la fecha local del navegador.
function currentPeriod() {
    const d = new Date();
    return { month: d.getMonth() + 1, year: d.getFullYear() };
}

// Cargar lista de tareas + estado de cada una para el mes en curso
// (spec monthly-tasks: "muestra todas las tareas mensuales configuradas
// con su estado actual").
async function loadTasks() {
    const container = document.getElementById('task-list');
    container.innerHTML = '<p class="loading-state">Cargando tareas…</p>';
    try {
        const tasks = (await MyTickTick.fetchAPI('/monthly-tasks')) || [];
        const statuses = await Promise.all(tasks.map((t) => currentStatus(t.id)));
        displayTasks(tasks, statuses);
    } catch (error) {
        console.error('Error loading tasks:', error);
        container.innerHTML = '<div class="error-state">Error al cargar las tareas: ' + escapeHtml(error.message) + '</div>';
    }
}

// Estado del mes en curso para una tarea: 'pending' o 'completed' (+ fecha).
// Si no se puede leer el historial no bloquea la lista: se asume pendiente.
async function currentStatus(taskId) {
    const fallback = { state: 'pending', completedAt: null };
    try {
        const history = (await MyTickTick.fetchAPI('/monthly-tasks/' + taskId + '/history')) || [];
        const p = currentPeriod();
        const rec = history.find((h) => h.month === p.month && h.year === p.year);
        if (!rec) return fallback;
        return rec.completed
            ? { state: 'completed', completedAt: rec.completedAt }
            : fallback;
    } catch (e) {
        console.warn('No se pudo cargar el estado de la tarea ' + taskId + ':', e.message);
        return fallback;
    }
}

// Mostrar tareas en la lista con su estado del mes en curso.
// Interacción (8.5): toggle (clic en el ítem) o swipe para completar,
// sin confirmación. "Deshacer" viaja por el toast.
function displayTasks(tasks, statuses) {
    const container = document.getElementById('task-list');

    if (!tasks || tasks.length === 0) {
        container.innerHTML = '<div class="empty-state">No hay tareas mensuales todavía. Creá la primera.</div>';
        return;
    }

    const html = tasks.map((task, i) => {
        const st = statuses[i] || { state: 'pending', completedAt: null };
        const done = st.state === 'completed';
        const badge = done
            ? '<span class="status-badge status-badge--done">Cumplida' +
              (st.completedAt ? ' · ' + formatDate(st.completedAt) : '') + '</span>'
            : '<span class="status-badge status-badge--pending">Pendiente este mes</span>';
        const checkbox = done
            ? '<span class="task-toggle task-toggle--on" aria-hidden="true">✓</span>'
            : '<span class="task-toggle" aria-hidden="true"></span>';
        return `
        <div class="task-item ${done ? 'completed' : ''}" data-id="${task.id}"
             data-completable="${done ? 'false' : 'true'}">
            <div class="task-item__head">
                ${checkbox}
                <h3>${escapeHtml(task.name)}</h3>
            </div>
            <p>${escapeHtml(task.description || '')}</p>
            <div class="task-meta">
                ${badge}
                <span class="task-repeats">Se repite cada mes</span>
            </div>
        </div>
    `;
    }).join('');

    container.innerHTML = html;
    bindListInteractions(container);
}

// Vincular toggle (clic) + swipe + long press sobre cada item de la lista.
// Mensuales: toggle o swipe → completar el mes en curso (solo si está pendiente).
// Los botones Historial/Editar/Eliminar se quitaron de la card: el LONG PRESS
// (sostener ~500ms en cualquier parte de la card) abre la ventana de edición,
// donde se puede editar o eliminar. El historial ya no es accesible desde esta
// pantalla (el modal queda solo para refrescarse si ya está abierto al
// completar un mes).
function bindListInteractions(container) {
    if (!container) return;
    const items = container.querySelectorAll('.task-item[data-id]');
    items.forEach(function (item) {
        const id = Number(item.getAttribute('data-id'));
        const completable = item.getAttribute('data-completable') === 'true';
        if (completable) {
            MyTickTick.attachToggle(item, { onClick: function () { completeTask(id); } });
            MyTickTick.attachSwipe(item, { onSwipe: function () { completeTask(id); }, color: 'green' });
        }
        // Long press → ventana de edición (editar / eliminar), en cualquier
        // estado de la tarea (pendiente o ya cumplida este mes).
        MyTickTick.attachLongPress(item, { onLongPress: function () { editTask(id); } });
    });
}

// Fecha ISO → "15 agosto 2026"
function formatDate(iso) {
    const d = new Date(iso);
    if (isNaN(d.getTime())) return '';
    return d.getDate() + ' ' + MONTH_NAMES[d.getMonth()] + ' ' + d.getFullYear();
}

// ---------- 8.2 - Formulario de creación ----------

// Mostrar/ocultar el formulario de creación.
function openCreateForm() {
    const section = document.getElementById('create-form');
    if (!section) return;
    section.hidden = false;
    closeEditForm();
    clearCreateFormFeedback();
    const nameInput = document.getElementById('task-name');
    if (nameInput && nameInput.focus) nameInput.focus();
}

function closeCreateForm() {
    const section = document.getElementById('create-form');
    if (!section) return;
    section.hidden = true;
}

// Feedback (éxito/error) dentro del formulario de creación.
function setCreateFormFeedback(message, kind) {
    const el = document.getElementById('form-feedback');
    if (!el) return;
    el.textContent = message;
    el.hidden = false;
    el.className = 'form-feedback ' + (kind === 'success' ? 'success' : (kind === 'error' ? 'error' : ''));
}

function clearCreateFormFeedback() {
    const el = document.getElementById('form-feedback');
    if (!el) return;
    el.textContent = '';
    el.hidden = true;
    el.className = 'form-feedback';
}

// Validación: el nombre es obligatorio (spec data-model: name requerido).
// Devuelve el nombre limpio o null si falta. `formPrefix` permite reutilizar
// la validación en el formulario de edición (8.3).
function validateName(formPrefix) {
    const nameInput = document.getElementById(formPrefix + '-name');
    const nameError = document.getElementById(formPrefix + '-name-error');
    const name = nameInput ? nameInput.value.trim() : '';
    if (nameError) nameError.hidden = true;
    if (!name) {
        if (nameError) nameError.hidden = false;
        if (nameInput && nameInput.focus) nameInput.focus();
        return null;
    }
    return name;
}

// Submit del formulario de creación: valida, crea la tarea y refresca la lista.
async function submitCreateTask(event) {
    if (event && event.preventDefault) event.preventDefault();

    const name = validateName('task');
    if (name === null) return;

    const descriptionInput = document.getElementById('task-description');
    const description = descriptionInput ? descriptionInput.value.trim() : '';

    const submitBtn = document.getElementById('submit-create-btn');
    if (submitBtn) submitBtn.disabled = true;
    setCreateFormFeedback('Creando tarea…', '');

    try {
        const task = await MyTickTick.fetchAPI('/monthly-tasks', {
            method: 'POST',
            json: { name: name, description: description }
        });
        console.log('Tarea creada:', task);
        const form = document.getElementById('create-task-form');
        if (form && form.reset) form.reset();
        setCreateFormFeedback('Tarea creada. Se repite cada mes.', 'success');
        closeCreateForm();
        loadTasks();
    } catch (error) {
        console.error('Error creating task:', error);
        setCreateFormFeedback('Error al crear la tarea: ' + error.message, 'error');
    } finally {
        if (submitBtn) submitBtn.disabled = false;
    }
}

// Botón "Nueva Tarea" + botón cancelar del formulario de creación.
function setupCreateTaskButton() {
    const btn = document.getElementById('create-task-btn');
    if (btn) {
        btn.addEventListener('click', function () {
            const section = document.getElementById('create-form');
            if (section && !section.hidden) closeCreateForm();
            else openCreateForm();
        });
    }

    const cancelBtn = document.getElementById('cancel-create-btn');
    if (cancelBtn) cancelBtn.addEventListener('click', closeCreateForm);
}

function setupCreateForm() {
    const form = document.getElementById('create-task-form');
    if (form) form.addEventListener('submit', submitCreateTask);
}

// ---------- 8.3 - Formulario de edición ----------

// Abrir el formulario de edición precargado con los datos de la tarea.
// (spec api-rest: PUT /api/monthly-tasks/{id} devuelve el recurso actualizado)
async function editTask(id) {
    const section = document.getElementById('edit-form');
    if (!section) return;

    let task;
    try {
        task = await MyTickTick.fetchAPI('/monthly-tasks/' + id);
    } catch (e) {
        console.error('Error loading task for edit:', e);
        closeEditForm(); // no dejar el formulario abierto con datos de otra tarea
        alert('No se pudo cargar la tarea: ' + e.message);
        return;
    }
    if (!task) return;

    document.getElementById('edit-task-name').value = task.name || '';
    const descInput = document.getElementById('edit-task-description');
    if (descInput) descInput.value = task.description || '';

    clearEditFormFeedback();
    closeCreateForm();
    section.hidden = false;
    section.dataset.taskId = String(id);

    const nameInput = document.getElementById('edit-task-name');
    if (nameInput && nameInput.focus) nameInput.focus();
}

function closeEditForm() {
    const section = document.getElementById('edit-form');
    if (section) section.hidden = true;
}

function setEditFormFeedback(message, kind) {
    const el = document.getElementById('edit-form-feedback');
    if (!el) return;
    el.textContent = message;
    el.hidden = false;
    el.className = 'form-feedback ' + (kind === 'success' ? 'success' : (kind === 'error' ? 'error' : ''));
}

function clearEditFormFeedback() {
    const el = document.getElementById('edit-form-feedback');
    if (!el) return;
    el.textContent = '';
    el.hidden = true;
    el.className = 'form-feedback';
}

// Submit del formulario de edición: valida y actualiza (PUT).
async function submitEditTask(event) {
    if (event && event.preventDefault) event.preventDefault();

    const name = validateName('edit-task');
    if (name === null) return;

    const descriptionInput = document.getElementById('edit-task-description');
    const description = descriptionInput ? descriptionInput.value.trim() : '';

    const section = document.getElementById('edit-form');
    const id = section ? section.dataset.taskId : null;
    if (!id) return;

    const submitBtn = document.getElementById('submit-edit-btn');
    if (submitBtn) submitBtn.disabled = true;
    setEditFormFeedback('Guardando cambios…', '');

    try {
        const updated = await MyTickTick.fetchAPI('/monthly-tasks/' + id, {
            method: 'PUT',
            json: { name: name, description: description }
        });
        console.log('Tarea actualizada:', updated);
        setEditFormFeedback('Cambios guardados.', 'success');
        closeEditForm();
        loadTasks();
    } catch (error) {
        console.error('Error updating task:', error);
        setEditFormFeedback('Error al actualizar la tarea: ' + error.message, 'error');
    } finally {
        if (submitBtn) submitBtn.disabled = false;
    }
}

function setupEditForm() {
    const form = document.getElementById('edit-task-form');
    if (form) form.addEventListener('submit', submitEditTask);

    const cancelBtn = document.getElementById('cancel-edit-btn');
    if (cancelBtn) cancelBtn.addEventListener('click', closeEditForm);

    const deleteBtn = document.getElementById('delete-edit-btn');
    if (deleteBtn) {
        deleteBtn.addEventListener('click', function () {
            const section = document.getElementById('edit-form');
            const id = section ? section.dataset.taskId : null;
            if (id) deleteTask(id);
        });
    }
}

// ---------- 8.5 - Marcar como completada (toggle / swipe / clic, sin confirm) ----------

// Marcar la tarea como cumplida para el mes en curso.
// (spec monthly-tasks: "toggle o swipe, sin confirmación" + toast "Deshacer")
// El backend es idempotente: volver a completarla solo refresca la fecha.
// Sin confirm(): la acción se aplica de inmediato.
// "Deshacer": recarga la lista (la tarea vuelve a verse pendiente visualmente
// en la lista, aunque el backend mantiene el registro).
async function completeTask(taskId) {
    const p = currentPeriod();
    const label = monthYearLabel(p.month, p.year);

    try {
        const record = await MyTickTick.fetchAPI('/monthly-tasks/' + taskId + '/completion', {
            method: 'PUT',
            json: { month: p.month, year: p.year }
        });
        console.log('Tarea completada:', record);

        // Toast con "Deshacer" (visual: recarga la lista).
        MyTickTick.toast(
            'Cumplida este mes (' + label + ')',
            {
                kind: 'success',
                actionLabel: 'Deshacer',
                onUndo: function () { loadTasks(); },
                timeout: 6000
            }
        );

        // Si el modal de historial está abierto, refrescarlo
        const modal = document.getElementById('history-modal');
        if (modal && !modal.hidden && modal.dataset.taskId === String(taskId)) {
            openHistory(taskId);
        }
        loadTasks();
    } catch (error) {
        console.error('Error completing task:', error);
        MyTickTick.toast('No se pudo completar: ' + error.message, { kind: 'error', timeout: 5000 });
    }
}

// ---------- 8.4 - Vista de historial ----------

// Mes/año → "agosto 2026"
function monthYearLabel(month, year) {
    return (MONTH_NAMES[month - 1] || '') + ' ' + year;
}

// Abrir el historial de una tarea: estado de cumplimiento mes a mes,
// más recientes primero (spec monthly-tasks: "ver el historial de
// cumplimiento mes a mes (más recientes primero)").
async function openHistory(taskId) {
    const modal = document.getElementById('history-modal');
    if (!modal) return;

    const list = document.getElementById('history-list');
    list.innerHTML = '<p class="loading-state">Cargando historial…</p>';

    let history = [];
    let task = null;
    try {
        [task, history] = await Promise.all([
            MyTickTick.fetchAPI('/monthly-tasks/' + taskId),
            MyTickTick.fetchAPI('/monthly-tasks/' + taskId + '/history')
        ]);
    } catch (e) {
        console.error('Error loading history:', e);
        list.innerHTML = '<div class="error-state">Error al cargar el historial: ' + escapeHtml(e.message) + '</div>';
    }

    const nameEl = document.getElementById('history-task-name');
    if (nameEl) nameEl.textContent = 'Historial de ' + (task ? task.name : 'la tarea');

    history = history || [];
    if (history.length === 0) {
        list.innerHTML = '<div class="empty-state">Todavía no hay registros de cumplimiento. ' +
            'Marcala como completada este mes y quedará en el historial.</div>';
    } else {
        const p = currentPeriod();
        list.innerHTML = history.map((h) => {
            const done = h.completed;
            const badge = done
                ? '<span class="status-badge status-badge--done">Cumplida' +
                  (h.completedAt ? ' · ' + formatDate(h.completedAt) : '') + '</span>'
                : '<span class="status-badge status-badge--pending">Pendiente</span>';
            // 8.5: acción de completar solo para el mes en curso y si aún no está cumplida
            const completeBtn = (h.month === p.month && h.year === p.year && !done)
                ? '<button class="btn-complete" onclick="completeTask(' + taskId + ')">Completar este mes</button>'
                : '';
            return `
            <div class="history-entry ${done ? 'done' : 'pending'}">
                <span class="history-month-year">${monthYearLabel(h.month, h.year)}</span>
                <span class="history-entry-actions">${badge} ${completeBtn}</span>
            </div>`;
        }).join('');
    }

    closeCreateForm();
    closeEditForm();
    modal.dataset.taskId = String(taskId);
    modal.hidden = false;
}

function closeHistory() {
    const modal = document.getElementById('history-modal');
    if (modal) modal.hidden = true;
}

function setupHistoryModal() {
    const modal = document.getElementById('history-modal');
    if (!modal) return;

    const closeBtn = document.getElementById('close-history');
    if (closeBtn) closeBtn.addEventListener('click', closeHistory);

    // Click en el fondo (fuera del .modal) cierra
    modal.addEventListener('click', (e) => { if (e.target === modal) closeHistory(); });

    // Tecla Escape cierra
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && !modal.hidden) closeHistory();
    });
}

// Eliminar tarea (con confirmación + undo): se pierde todo el historial de
// cumplimientos, por lo que se pide confirmación antes de borrar.
async function deleteTask(id) {
    let task = null;
    try {
        task = await MyTickTick.fetchAPI('/monthly-tasks/' + id);
    } catch (e) {
        console.error('No se pudo cargar la tarea para eliminar:', e);
    }

    // Confirmar: evita borrar por accidente (se pierde el historial).
    const ok = await MyTickTick.confirm({
        title: '¿Eliminar tarea' + (task && task.name ? ': ' + task.name : '') + '?',
        message: 'Se eliminará la tarea mensual y TODO su historial de cumplimientos. Esta acción no se puede deshacer por el sistema (el botón "Deshacer" solo recrea la tarea, no el historial).',
        confirmLabel: 'Eliminar',
        cancelLabel: 'Cancelar',
        kind: 'danger'
    });
    if (!ok) return;

    try {
        await MyTickTick.fetchAPI(`/monthly-tasks/${id}`, { method: 'DELETE' });
        closeEditForm();
        closeCreateForm();

        MyTickTick.toast('Tarea eliminada', {
            kind: 'error',
            actionLabel: 'Deshacer',
            onUndo: function () {
                if (!task) { loadTasks(); return; }
                MyTickTick.fetchAPI('/monthly-tasks', {
                    method: 'POST',
                    json: { name: task.name, description: task.description || '' }
                }).then(loadTasks).catch(function (e) {
                    console.error('No se pudo restaurar la tarea:', e);
                    MyTickTick.toast('No se pudo restaurar la tarea', { kind: 'error' });
                });
            },
            timeout: 8000
        });

        loadTasks();
    } catch (error) {
        console.error('Error deleting task:', error);
        MyTickTick.toast('Error al eliminar la tarea: ' + error.message, { kind: 'error' });
    }
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
