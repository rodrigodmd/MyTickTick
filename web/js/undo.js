// MyTickTick - Interacción de lista (toast "Deshacer" + swipe)
//
// Implementa el patrón pedido en las specs (immediate-tasks / monthly-tasks):
//   - Sin confirm(): el usuario completa o elimina y el cambio se aplica ya.
//   - Toast "Deshacer": aparece tras el cambio y permite revertirlo.
//   - Swipe: gesto horizontal (táctil o mouse) sobre un item dispara la acción.
//
// API pública:
//   MyTickTick.toast(message, { actionLabel, onUndo, timeout, kind })
//     -> muestra un toast inferior con un botón opcional (p. ej. "Deshacer").
//        Devuelve la función que lo cierra.
//   MyTickTick.attachSwipe(el, { onSwipe, threshold, color })
//     -> habilita el gesto swipe horizontal sobre `el`. onSwipe() se llama
//        una sola vez cuando el deslizamiento supera `threshold` px.
//   MyTickTick.attachLongPress(el, { onLongPress, duration, ignoreButtons })
//     -> habilita el gesto long press (mantener presionado) sobre `el`.
//        onLongPress() se llama una sola vez al sostener `duration` ms
//        (default 500) en cualquier punto del item.
(function () {
  'use strict';

  var TOAST_ID = 'mtt-toast-host';

  // ---------- Toast ----------

  function ensureHost() {
    var host = document.getElementById(TOAST_ID);
    if (!host) {
      host = document.createElement('div');
      host.id = TOAST_ID;
      host.className = 'mtt-toast-host';
      host.setAttribute('aria-live', 'polite');
      document.body.appendChild(host);
    }
    return host;
  }

  // toast(message, opts) -> close()
  //   opts.actionLabel : texto del botón (p. ej. "Deshacer"). Opcional.
  //   opts.onUndo      : callback que revierte la acción. Opcional.
  //   opts.timeout     : ms antes de auto-cerrar (default 5000). 0 = no cierra.
  //   opts.kind        : 'success' | 'error' | '' (default '').
  function toast(message, opts) {
    opts = opts || {};
    var host = ensureHost();

    // Reemplaza cualquier toast anterior (un solo toast a la vez).
    host.innerHTML = '';

    var el = document.createElement('div');
    el.className = 'mtt-toast mtt-toast--' + (opts.kind || 'default');

    var text = document.createElement('span');
    text.className = 'mtt-toast__msg';
    text.textContent = message || '';
    el.appendChild(text);

    var closeTimer = null;
    var closed = false;

    function close() {
      if (closed) return;
      closed = true;
      if (closeTimer) clearTimeout(closeTimer);
      if (el.parentNode) el.parentNode.removeChild(el);
    }

    if (opts.actionLabel) {
      var btn = document.createElement('button');
      btn.type = 'button';
      btn.className = 'mtt-toast__action';
      btn.textContent = opts.actionLabel;
      btn.addEventListener('click', function () {
        if (typeof opts.onUndo === 'function') {
          try {
            opts.onUndo();
          } catch (e) {
            console.error('onUndo falló:', e);
          }
        }
        close();
      });
      el.appendChild(btn);
    }

    host.appendChild(el);

    // Entrada animada (rAF para forzar layout).
    requestAnimationFrame(function () {
      if (el.parentNode) el.classList.add('mtt-toast--show');
    });

    var timeout = opts.timeout !== undefined ? opts.timeout : 5000;
    if (timeout > 0) {
      closeTimer = setTimeout(close, timeout);
    }

    return close;
  }

  // ---------- Swipe ----------

  // attachSwipe(el, opts)
  //   opts.onSwipe    : callback al completar el gesto.
  //   opts.threshold  : px mínimos de desplazamiento (default 64).
  //   opts.color      : 'green' | 'red' | 'accent' (fondo revelado, default 'accent').
  function attachSwipe(el, opts) {
    if (!el || el.__swipeBound) return;
    opts = opts || {};
    var threshold = opts.threshold || 64;
    var color = opts.color || 'accent';

    el.classList.add('swipeable', 'swipe-' + color);
    el.__swipeBound = true;

    var startX = 0;
    var startY = 0;
    var active = false;
    var triggered = false;
    var width = el.getBoundingClientRect().width || 1;

    function setTranslate(px) {
      el.style.transform = 'translateX(' + px + 'px)';
    }

    function onDown(e) {
      if (e.button !== undefined && e.button !== 0) return; // solo botón primario
      // No iniciar swipe desde botones/links internos (que tienen su propio click).
      if (e.target.closest('button, a, input, select, textarea')) return;

      active = true;
      triggered = false;
      startX = e.clientX;
      startY = e.clientY;
      width = el.getBoundingClientRect().width || 1;
      if (el.setPointerCapture) el.setPointerCapture(e.pointerId);
    }

    function onMove(e) {
      if (!active) return;
      var dx = e.clientX - startX;
      var dy = e.clientY - startY;

      // Solo si el gesto es claramente horizontal (|dx| > |dy|).
      if (Math.abs(dy) > Math.abs(dx)) return;
      if (Math.abs(dx) < 4) return;

      if (Math.abs(dx) >= threshold) {
        if (!triggered) {
          triggered = true;
          // Revelar el fondo de acción antes de disparar.
          setTranslate(dx < 0 ? -width : width);
          setTimeout(function () {
            if (typeof opts.onSwipe === 'function') opts.onSwipe();
            reset();
          }, 120);
        }
        return;
      }
      setTranslate(dx);
    }

    function reset() {
      active = false;
      triggered = false;
      el.style.transition = 'transform 0.18s ease';
      setTranslate(0);
      setTimeout(function () { el.style.transition = ''; }, 200);
    }

    function onUp() {
      if (active && !triggered) reset();
    }

    el.addEventListener('pointerdown', onDown, { passive: true });
    el.addEventListener('pointermove', onMove, { passive: true });
    el.addEventListener('pointerup', onUp, { passive: true });
    el.addEventListener('pointercancel', onUp, { passive: true });
  }

  // ---------- Toggle (clic en el ítem) ----------

  // attachToggle(el, opts)
  //   opts.onClick : callback al hacer clic en el ítem (no en sus controles).
  function attachToggle(el, opts) {
    if (!el || el.__toggleBound) return;
    opts = opts || {};
    el.__toggleBound = true;
    el.classList.add('clickable');
    el.addEventListener('click', function (e) {
      // Ignorar clics en botones/links/inputs internos.
      if (e.target.closest('button, a, input, select, textarea')) return;
      // Si un long press ya consumio esta interaccion (el gesto de "sostener"
      // ya abrio la accion), no disparar el toggle con el click residual.
      if (el.__longPressConsumed) {
        el.__longPressConsumed = false;
        return;
      }
      if (typeof opts.onClick === 'function') opts.onClick();
    });
  }

  // ---------- Long press ----------

  // attachLongPress(el, opts)
  //   opts.onLongPress : callback al mantener presionado el item (default 500ms).
  //   opts.duration    : ms de presion sostenida (default 500).
  //   opts.ignoreButtons: si es true, no dispara desde botones/links/inputs
  //                       internos (default false: cualquier parte del item).
  //   Dispara UNA sola vez por presion. Cancela el click sintetico para que no
  //   se dispare a la vez el toggle (attachToggle) del mismo item.
  function attachLongPress(el, opts) {
    if (!el || el.__longPressBound) return;
    opts = opts || {};
    var duration = opts.duration || 500;
    var ignoreButtons = !!opts.ignoreButtons;
    el.__longPressBound = true;
    el.classList.add('longpress-target');

    var timer = null;
    var fired = false;

    // Suprime el click que el navegador sintetiza tras soltar el dedo, para
    // que no se dispare a la vez el toggle (attachToggle) del mismo item.
    function suppressClick() {
      fired = true;
      // Bandera leida por attachToggle: su handler de click se registra ANTES
      // que el capture de abajo, asi que solo esta bandera evita el toggle.
      el.__longPressConsumed = true;
      var h = function (e) {
        el.__longPressConsumed = false;
        e.stopPropagation();
        if (e.cancelable) e.preventDefault();
      };
      el.addEventListener('click', h, true);
      // Se retira solo a los 500ms por si no llega ningun click (p. ej. el
      // punto sale del elemento); en caso contrario el primer click ya lo
      // suprimio y este retiro es inocuo.
      setTimeout(function () { el.removeEventListener('click', h, true); }, 500);
    }

    function clearHold() {
      if (timer) { clearTimeout(timer); timer = null; }
      el.classList.remove('longpress-hold');
    }

    function onDown(e) {
      if (e.button !== undefined && e.button !== 0) return; // solo boton primario
      // Reset de la presion anterior (si un long press ya se consumio, el
      // handler de click lo suprimio; el proximo tap vuelve a ser normal).
      fired = false;
      if (ignoreButtons && e.target.closest('button, a, input, select, textarea')) return;

      timer = setTimeout(function () {
        timer = null;
        el.classList.remove('longpress-hold');
        suppressClick();
        if (typeof opts.onLongPress === 'function') opts.onLongPress();
      }, duration);
      // Feedback visual de "sostener" (solo si sigue activo).
      requestAnimationFrame(function () {
        if (timer) el.classList.add('longpress-hold');
      });
    }

    function onUp(e) {
      clearHold();
      if (fired && e.cancelable) e.preventDefault();
    }

    function onCancel() {
      clearHold();
    }

    el.addEventListener('pointerdown', onDown, { passive: true });
    el.addEventListener('pointerup', onUp, { passive: false });
    el.addEventListener('pointercancel', onCancel, { passive: true });
  }

  // ---------- Confirmación (modal propio, no window.confirm) ----------
  //
  // confirm({ title, message, confirmLabel, cancelLabel, kind }) -> Promise<bool>
  //   Devuelve true si el usuario confirma, false si cancela (o cierra con
  //   Escape / clic fuera). Se usa antes de acciones destructivas que borran
  //   historial (tracker, tareas mensuales) para evitar borrados accidentales.
  function confirmDialog(opts) {
    opts = opts || {};
    return new Promise(function (resolve) {
      var overlay = document.createElement('div');
      overlay.className = 'mtt-confirm-overlay';
      // Icono de advertencia (triángulo) inline, se pinta con el color del
      // contexto (var(--danger) para kind='danger', var(--accent) por defecto).
      var icon =
        '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" ' +
        'stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">' +
        '<path d="M10.29 3.86 1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>' +
        '<line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/>' +
        '</svg>';
      overlay.innerHTML =
        '<div class="mtt-confirm" role="alertdialog" aria-modal="true">' +
          '<div class="mtt-confirm__head">' +
            '<span class="mtt-confirm__icon">' + icon + '</span>' +
            '<h3 class="mtt-confirm__title"></h3>' +
          '</div>' +
          '<p class="mtt-confirm__msg"></p>' +
          '<div class="mtt-confirm__actions">' +
            '<button type="button" class="btn btn-ghost mtt-confirm__cancel"></button>' +
            '<button type="button" class="btn mtt-confirm__ok"></button>' +
          '</div>' +
        '</div>';

      overlay.querySelector('.mtt-confirm__title').textContent =
        opts.title || '¿Estás seguro?';
      overlay.querySelector('.mtt-confirm__msg').textContent =
        opts.message || '';

      var cancelBtn = overlay.querySelector('.mtt-confirm__cancel');
      var okBtn = overlay.querySelector('.mtt-confirm__ok');
      cancelBtn.textContent = opts.cancelLabel || 'Cancelar';
      okBtn.textContent = opts.confirmLabel || 'Confirmar';
      if (opts.kind === 'danger') okBtn.classList.add('mtt-confirm__ok--danger');

      function finish(value) {
        overlay.remove();
        document.removeEventListener('keydown', onKey, true);
        resolve(value);
      }

      function onKey(e) {
        if (e.key === 'Escape') { e.preventDefault(); e.stopPropagation(); finish(false); }
        if (e.key === 'Enter')  { e.preventDefault(); e.stopPropagation(); finish(true); }
      }

      cancelBtn.addEventListener('click', function () { finish(false); });
      okBtn.addEventListener('click', function () { finish(true); });
      overlay.addEventListener('click', function (e) {
        if (e.target === overlay) finish(false);
      });
      document.addEventListener('keydown', onKey, true);

      document.body.appendChild(overlay);
      requestAnimationFrame(function () {
        overlay.classList.add('mtt-confirm--show');
        okBtn.focus();
      });
    });
  }

  window.MyTickTick = window.MyTickTick || {};
  window.MyTickTick.toast = toast;
  window.MyTickTick.attachSwipe = attachSwipe;
  window.MyTickTick.attachToggle = attachToggle;
  window.MyTickTick.attachLongPress = attachLongPress;
  window.MyTickTick.confirm = confirmDialog;
})();
