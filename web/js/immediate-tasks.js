// MyTickTick - Tareas Inmediatas (usa el cliente fetchAPI de app.js)
//
// Sección 9 del plan:
//   9.1 Lista de tareas inmediatas
//   9.2 Formulario de creación
//   9.3 Selector de fecha límite
//   9.4 Selector de prioridad
//   9.5 Ordenamiento por fecha límite
//
// Orden por defecto (spec immediate-tasks "Asignar prioridad alta → la muestra
// primero"): prioridad (alta > media > baja) y, dentro de cada prioridad,
// fecha límite ascendente. El selector "Ordenar por" permite cambiar a
// fecha límite pura (9.5 / spec "Ver tareas próximas").
document.addEventListener('DOMContentLoaded', function () {
    loadTasks();
    setupQuickAdd();
    setupCreateTaskButton();
    setupCreateForm();
    setupEditForm();
    setupFilters();
});

// Nombres de meses en español para mostrar fechas.
const MONTH_NAMES = ['enero', 'febrero', 'marzo', 'abril', 'mayo', 'junio',
    'julio', 'agosto', 'septiembre', 'octubre', 'noviembre', 'diciembre'];

// Pesos para ordenar por prioridad (menor = más urgente).
const PRIORITY_RANK = { high: 0, medium: 1, low: 2 };

// Estado de la vista (filtros + orden), conservado entre recargas de la lista.
const viewState = {
    status: 'open',     // 'open' | 'all' | 'completed'
    sort: 'priority'    // 'priority' | 'dueDate'
};

// ---------- 9.1 - Lista de tareas ----------

async function loadTasks() {
    const container = document.getElementById('task-list');
    if (container) container.innerHTML = '<p class="loading-state">Cargando tareas…</p>';
    try {
        const tasks = (await MyTickTick.fetchAPI('/immediate-tasks')) || [];
        displayTasks(tasks);
    } catch (error) {
        console.error('Error loading tasks:', error);
        if (container) {
            container.innerHTML = '<div class="error-state">Error al cargar las tareas: ' +
                escapeHtml(error.message) + '</div>';
        }
    }
}

// Filtro de estado + ordenamiento elegido por el usuario (9.4 / 9.5).
function applyView(tasks) {
    let list = tasks.slice();

    if (viewState.status === 'open') {
        list = list.filter((t) => !t.isCompleted);
    } else if (viewState.status === 'completed') {
        list = list.filter((t) => t.isCompleted);
    }

    list.sort(function (a, b) {
        // Completadas siempre al final (dentro del grupo, el orden elegido).
        if (!!a.isCompleted !== !!b.isCompleted) return a.isCompleted ? 1 : -1;
        if (viewState.sort === 'dueDate') {
            return dueTime(a) - dueTime(b);
        }
        const ra = PRIORITY_RANK[a.priority] !== undefined ? PRIORITY_RANK[a.priority] : 1;
        const rb = PRIORITY_RANK[b.priority] !== undefined ? PRIORITY_RANK[b.priority] : 1;
        if (ra !== rb) return ra - rb;
        return dueTime(a) - dueTime(b);
    });
    return list;
}

function dueTime(task) {
    const t = task.dueDate ? new Date(task.dueDate).getTime() : NaN;
    return isNaN(t) ? Infinity : t;
}

function displayTasks(tasks) {
    const container = document.getElementById('task-list');
    const list = applyView(tasks);

    if (!list || list.length === 0) {
        const hint = tasks && tasks.length > 0
            ? 'Ninguna tarea coincide con el filtro actual.'
            : 'No hay tareas inmediatas. Creá una para los próximos días.';
        container.innerHTML = '<div class="empty-state">' + hint + '</div>';
        return;
    }

    const now = Date.now();
    const html = list.map(function (task) {
        const due = task.dueDate ? new Date(task.dueDate) : null;
        const dueLabel = due ? formatDueLabel(due) : 'sin fecha';
        const priority = task.priority || 'medium';
        let dueBadge = '';
        let itemClass = 'task-item';
        if (task.isCompleted) {
            itemClass += ' completed';
        } else if (due) {
            const diff = due.getTime() - now;
            if (diff < 0) {
                itemClass += ' overdue';
                dueBadge = '<span class="status-badge" style="background:var(--danger);color:#fff">Vencida</span>';
            } else if (diff <= 24 * 3600 * 1000) {
                itemClass += ' due-soon';
                dueBadge = '<span class="status-badge" style="background:var(--warning-soft);color:var(--warning)">Vence pronto</span>';
            }
        }
        const checkbox = task.isCompleted
            ? '<span class="task-toggle task-toggle--on" aria-hidden="true">✓</span>'
            : '<span class="task-toggle" aria-hidden="true"></span>';
        return `
        <div class="${itemClass}" data-id="${task.id}">
            <div class="task-item__head">
                ${checkbox}
                <h3>${escapeHtml(task.name)}</h3>
            </div>
            <p>${escapeHtml(task.description || '')}</p>
            <div class="task-meta">
                <span>Vence: ${dueLabel}</span>
                <span class="priority priority-${priority}">${priority}</span>
                ${dueBadge}
            </div>
            <div class="task-actions">
                <button class="btn btn-ghost" onclick="openEditForm(${task.id})">Editar</button>
                <button class="btn btn-ghost btn-danger" onclick="deleteTask(${task.id})">Eliminar</button>
            </div>
        </div>`;
    }).join('');

    container.innerHTML = html;
    bindListInteractions(container);
}

// Vincular toggle (clic) + swipe sobre cada item: cambiar el estado de
// completada (toggle o swipe, sin confirmación) con toast "Deshacer".
function bindListInteractions(container) {
    if (!container) return;
    const items = container.querySelectorAll('.task-item[data-id]');
    items.forEach(function (item) {
        const id = Number(item.getAttribute('data-id'));
        MyTickTick.attachToggle(item, { onClick: function () { toggleTask(id); } });
        MyTickTick.attachSwipe(item, { onSwipe: function () { toggleTask(id); }, color: 'green' });
    });
}

// Cambiar el estado de completada (el estado opuesto al actual).
function toggleTask(id) {
    const item = document.querySelector('.task-item[data-id="' + id + '"]');
    const done = item ? item.classList.contains('completed') : false;
    completeTask(id, !done);
}

// "15 ago 2026" + marca "hoy" / "mañana" cuando corresponde.
function formatDueLabel(due) {
    if (isNaN(due.getTime())) return 'sin fecha';
    const base = due.getDate() + ' ' + MONTH_NAMES[due.getMonth()].slice(0, 3) + ' ' + due.getFullYear();
    const today = new Date();
    const startOf = (x) => new Date(x.getFullYear(), x.getMonth(), x.getDate()).getTime();
    const diffDays = Math.round((startOf(due) - startOf(today)) / 86400000);
    if (diffDays === 0) return base + ' (hoy)';
    if (diffDays === 1) return base + ' (mañana)';
    return base;
}

// ---------- 9.2 - Formulario de creación ----------

function openCreateForm() {
    const section = document.getElementById('create-form');
    if (!section) return;
    section.hidden = false;
    closeEditForm();
    clearCreateFormFeedback();
    // 9.3 - fecha límite por defecto: +7 días (tareas para "próximos días")
    const dueInput = document.getElementById('task-due-date');
    if (dueInput) dueInput.value = toISODate(new Date(Date.now() + 7 * 86400 * 1000));
    const nameInput = document.getElementById('task-name');
    if (nameInput && nameInput.focus) nameInput.focus();
}

function closeCreateForm() {
    const section = document.getElementById('create-form');
    if (section) section.hidden = true;
}

function setCreateFormFeedback(message, kind) {
    const el = document.getElementById('form-feedback');
    if (!el) return;
    el.textContent = message;
    el.hidden = false;
    el.className = 'form-feedback ' + (kind === 'success' ? 'success' : (kind === 'error' ? 'error' : ''));
}

function clearCreateFormFeedback() {
    const el = document.getElementById('form-feedback');
    if (el) {
        el.textContent = '';
        el.hidden = true;
        el.className = 'form-feedback';
    }
}

// Validación de nombre, reutilizada por creación y edición.
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

// 9.3 - Validación de fecha límite (input type=date → YYYY-MM-DD).
// Devuelve "YYYY-MM-DD" o null si falta.
function validateDueDate(formPrefix) {
    const dueInput = document.getElementById(formPrefix + '-due-date');
    const dueError = document.getElementById(formPrefix + '-due-date-error');
    const value = dueInput ? dueInput.value.trim() : '';
    if (dueError) dueError.hidden = true;
    if (!value || isNaN(new Date(value + 'T00:00:00').getTime())) {
        if (dueError) dueError.hidden = false;
        if (dueInput && dueInput.focus) dueInput.focus();
        return null;
    }
    return value;
}

// 9.4 - Prioridad del select (high | medium | low), con "medium" de respaldo.
function selectedPriority(formPrefix) {
    const sel = document.getElementById(formPrefix + '-priority');
    const value = sel ? sel.value : 'medium';
    return PRIORITY_RANK[value] !== undefined ? value : 'medium';
}

async function submitCreateTask(event) {
    if (event && event.preventDefault) event.preventDefault();

    const name = validateName('task');
    if (name === null) return;
    const dueDate = validateDueDate('task');
    if (dueDate === null) return;
    const priority = selectedPriority('task');

    const descriptionInput = document.getElementById('task-description');
    const description = descriptionInput ? descriptionInput.value.trim() : '';

    const submitBtn = document.getElementById('submit-create-btn');
    if (submitBtn) submitBtn.disabled = true;
    setCreateFormFeedback('Creando tarea…', '');

    try {
        const task = await MyTickTick.fetchAPI('/immediate-tasks', {
            method: 'POST',
            json: {
                name: name,
                description: description,
                dueDate: dueDate + 'T23:59:59Z', // fin del día, sin ambigüedad de zona
                priority: priority
            }
        });
        console.log('Tarea creada:', task);
        const form = document.getElementById('create-task-form');
        if (form && form.reset) form.reset();
        setCreateFormFeedback('Tarea creada.', 'success');
        closeCreateForm();
        loadTasks();
    } catch (error) {
        console.error('Error creating task:', error);
        setCreateFormFeedback('Error al crear la tarea: ' + error.message, 'error');
    } finally {
        if (submitBtn) submitBtn.disabled = false;
    }
}

// ---------- Carga rápida (campo superior, solo nombre + Enter) ----------
// Permite agregar una tarea inmediata escribiendo solo el título y
// presionando Enter. Usa los mismos defaults que el formulario completo:
// fecha límite +7 días y prioridad "medium".

async function submitQuickAdd(event) {
    if (event && event.preventDefault) event.preventDefault();

    const input = document.getElementById('quick-add-input');
    const btn = document.getElementById('quick-add-btn');
    const name = input ? input.value.trim() : '';
    if (!name) {
        if (input && input.focus) input.focus();
        return;
    }
    if (btn) btn.disabled = true;

    try {
        const task = await MyTickTick.fetchAPI('/immediate-tasks', {
            method: 'POST',
            json: {
                name: name,
                description: '',
                dueDate: toISODate(new Date(Date.now() + 7 * 86400 * 1000)) + 'T23:59:59Z',
                priority: 'medium'
            }
        });
        console.log('Tarea creada (rápida):', task);
        if (input) input.value = '';
        MyTickTick.toast('Tarea agregada: ' + name, { kind: 'success' });
        loadTasks();
    } catch (error) {
        console.error('Error creating task (quick):', error);
        MyTickTick.toast('Error al agregar la tarea: ' + error.message, { kind: 'error' });
    } finally {
        if (btn) btn.disabled = false;
    }
}

function setupQuickAdd() {
    const form = document.getElementById('quick-add-form');
    if (form) {
        form.addEventListener('submit', submitQuickAdd);
    }
}

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

// ---------- Formulario de edición ----------

async function openEditForm(id) {
    const section = document.getElementById('edit-form');
    if (!section) return;

    let task;
    try {
        task = await MyTickTick.fetchAPI('/immediate-tasks/' + id);
    } catch (e) {
        console.error('Error loading task for edit:', e);
        closeEditForm();
        alert('No se pudo cargar la tarea: ' + e.message);
        return;
    }
    if (!task) return;

    document.getElementById('edit-task-name').value = task.name || '';
    const descInput = document.getElementById('edit-task-description');
    if (descInput) descInput.value = task.description || '';
    const dueInput = document.getElementById('edit-task-due-date');
    if (dueInput) dueInput.value = task.dueDate ? task.dueDate.slice(0, 10) : '';
    const prioInput = document.getElementById('edit-task-priority');
    if (prioInput) {
        prioInput.value = PRIORITY_RANK[task.priority] !== undefined ? task.priority : 'medium';
    }

    clearEditFormFeedback();
    closeCreateForm();
    section.hidden = false;
    section.dataset.taskId = String(id);
    section.dataset.isCompleted = task.isCompleted ? 'true' : 'false';
    const nameInput = document.getElementById('edit-task-name');
    if (nameInput && nameInput.focus) nameInput.focus();
}

function closeEditForm() {
    const section = document.getElementById('edit-form');
    if (section) section.hidden = true;
}

function setEditFormFeedback(message, kind) {
    const el = document.getElementById('edit-form-feedback');
    if (el) {
        el.textContent = message;
        el.hidden = false;
        el.className = 'form-feedback ' + (kind === 'success' ? 'success' : (kind === 'error' ? 'error' : ''));
    }
}

function clearEditFormFeedback() {
    const el = document.getElementById('edit-form-feedback');
    if (el) {
        el.textContent = '';
        el.hidden = true;
        el.className = 'form-feedback';
    }
}

async function submitEditTask(event) {
    if (event && event.preventDefault) event.preventDefault();

    const name = validateName('edit-task');
    if (name === null) return;
    const dueDate = validateDueDate('edit-task');
    if (dueDate === null) return;
    const priority = selectedPriority('edit-task');

    const descriptionInput = document.getElementById('edit-task-description');
    const description = descriptionInput ? descriptionInput.value.trim() : '';

    const section = document.getElementById('edit-form');
    const id = section ? section.dataset.taskId : null;
    if (!id) return;
    // El estado de completada de la tarea se conserva (no se cambia al editar).
    const isCompleted = section ? section.dataset.isCompleted === 'true' : false;

    const submitBtn = document.getElementById('submit-edit-btn');
    if (submitBtn) submitBtn.disabled = true;
    setEditFormFeedback('Guardando cambios…', '');

    try {
        const updated = await MyTickTick.fetchAPI('/immediate-tasks/' + id, {
            method: 'PUT',
            json: {
                name: name,
                description: description,
                dueDate: dueDate + 'T23:59:59Z',
                priority: priority,
                isCompleted: isCompleted
            }
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

// ---------- Completar / reabrir (toggle o swipe, sin confirm) ----------

// PUT completo: el backend espera todos los campos; se envían los de la
// tarea cargada con isCompleted actualizado.
// Sin confirm(): el cambio se aplica de inmediato y el toast ofrece "Deshacer".
async function completeTask(id, value) {
    let task;
    try {
        task = await MyTickTick.fetchAPI('/immediate-tasks/' + id);
    } catch (e) {
        MyTickTick.toast('No se pudo cargar la tarea: ' + e.message, { kind: 'error' });
        return;
    }
    if (!task) return;

    try {
        await MyTickTick.fetchAPI('/immediate-tasks/' + id, {
            method: 'PUT',
            json: {
                name: task.name,
                description: task.description,
                dueDate: task.dueDate,
                priority: task.priority || 'medium',
                isCompleted: value
            }
        });
        MyTickTick.toast(
            value ? 'Tarea completada' : 'Tarea reabierta',
            {
                kind: value ? 'success' : 'default',
                actionLabel: 'Deshacer',
                onUndo: function () { completeTask(id, !value); },
                timeout: 6000
            }
        );
        loadTasks();
    } catch (error) {
        console.error('Error completing task:', error);
        MyTickTick.toast('Error al actualizar el estado: ' + error.message, { kind: 'error' });
    }
}

// ---------- Eliminar (sin confirm, con undo) ----------

async function deleteTask(id) {
    let task = null;
    try {
        task = await MyTickTick.fetchAPI('/immediate-tasks/' + id);
    } catch (e) {
        console.error('No se pudo cargar la tarea para eliminar:', e);
    }

    try {
        await MyTickTick.fetchAPI('/immediate-tasks/' + id, { method: 'DELETE' });
        closeEditForm();
        closeCreateForm();

        MyTickTick.toast('Tarea eliminada', {
            kind: 'error',
            actionLabel: 'Deshacer',
            onUndo: function () {
                if (!task) { loadTasks(); return; }
                MyTickTick.fetchAPI('/immediate-tasks', {
                    method: 'POST',
                    json: {
                        name: task.name,
                        description: task.description || '',
                        dueDate: task.dueDate,
                        priority: task.priority || 'medium'
                    }
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

// ---------- 9.1 / 9.5 - Filtros y ordenamiento ----------

function setupFilters() {
    const status = document.getElementById('status-filter');
    if (status) {
        status.value = viewState.status;
        status.addEventListener('change', function () {
            viewState.status = status.value;
            loadTasks();
        });
    }
    const sort = document.getElementById('sort-order');
    if (sort) {
        sort.value = viewState.sort;
        sort.addEventListener('change', function () {
            viewState.sort = sort.value;
            loadTasks();
        });
    }
}

// ---------- Utilidades ----------

function toISODate(d) {
    return d.toISOString().split('T')[0];
}

function escapeHtml(value) {
    return String(value)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}
