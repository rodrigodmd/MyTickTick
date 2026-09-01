// MyTickTick - Layout compartido (7.3)
//
// Provee el "layout base" común a todas las secciones:
//   - Header con marca + navegación + logout
//   - Barra de navegación inferior con iconos (mobile): navegar entre
//     secciones con un solo toque
//   - Enlace activo según la ruta actual
//   - Menú colapsable en pantallas chicas (responsive)
//
// Cada página declara los contenedores vacíos (#site-header)
// y este módulo los rellena. El contenido de cada sección vive en <main id="site-main">.
(function () {
  'use strict';

  // Secciones disponibles en la app.
  var SECTIONS = [
    { path: '/monthly-tasks',   label: 'Tareas Mensuales' },
    { path: '/immediate-tasks', label: 'Tareas Inmediatas' },
    { path: '/trackers',        label: 'Trackers' },
    { path: '/dashboard',       label: 'Dashboard' }
  ];

  // Iconos (trazo, estilo feather) para la barra inferior.
  var ICONS = {
    home: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>',
    calendar: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>',
    bolt: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon></svg>',
    chart: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><line x1="18" y1="20" x2="18" y2="10"></line><line x1="12" y1="20" x2="12" y2="4"></line><line x1="6" y1="20" x2="6" y2="14"></line></svg>',
    pie: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M21.21 15.89A10 10 0 1 1 8 2.83"></path><path d="M22 12A10 10 0 0 0 12 2v10z"></path></svg>'
  };

  // Pestañas de la barra inferior (mobile): las 4 secciones (sin "Inicio":
  // se accede desde la marca del header). Orden pedido por Rodrigo:
  // Inmediatas → Trackers → Mensuales → Dashboard.
  var TABS = [
    { path: '/immediate-tasks', label: 'Inmediatas', icon: 'bolt' },
    { path: '/trackers',        label: 'Trackers',   icon: 'chart' },
    { path: '/monthly-tasks',   label: 'Mensuales',  icon: 'calendar' },
    { path: '/dashboard',       label: 'Dashboard',  icon: 'pie' }
  ];

  function headerHTML() {
    var links = SECTIONS.map(function (s) {
      return '<a href="' + s.path + '" data-nav="' + s.path + '">' + s.label + '</a>' +
             '<span class="nav-sep"></span>';
    }).join('');

    return '' +
      '<div class="site-header__bar">' +
        '<a href="/" class="brand" aria-label="MyTickTick — inicio">MyTickTick</a>' +
        '<button type="button" class="nav-toggle" aria-label="Abrir menú" aria-expanded="false" aria-controls="site-nav">&#9776;</button>' +
      '</div>' +
      '<nav id="site-nav" class="site-nav" aria-label="Principal">' + links +
        '<a href="#" class="nav-logout" id="logout-btn" title="Cerrar sesión">Salir</a>' +
      '</nav>';
  }

  // Footer: vacio por pedido de Rodrigo (2026-08-28) — se elimino el texto
  // "MyTickTick — Gestión de tareas personalizadas" de todas las pantallas.
  // Si vuelve a pedirse un footer, reponer el <footer id="site-footer"> en los
  // HTML y el <p> aqui.
  function footerHTML() {
    return '';
  }

  function tabbarHTML() {
    var links = TABS.map(function (t) {
      return '<a href="' + t.path + '" data-tab="' + t.path + '" aria-label="' + t.label + '">' +
        '<span class="tabbar__icon">' + ICONS[t.icon] + '</span>' +
        '<span class="tabbar__label">' + t.label + '</span>' +
      '</a>';
    }).join('');
    return '<nav id="site-tabbar" class="site-tabbar" aria-label="Navegación principal">' + links + '</nav>';
  }

  // Marca el enlace activo: navegación del header y barra inferior.
  function setActive() {
    var path = window.location.pathname;
    var links = document.querySelectorAll('.site-nav a[data-nav], .site-tabbar a[data-tab]');
    links.forEach(function (a) {
      var key = a.getAttribute('data-nav') || a.getAttribute('data-tab');
      if (key === path) a.classList.add('active');
    });
  }

  // Menú móvil: el botón toggle muestra/oculta la navegación.
  function initToggle() {
    var header = document.getElementById('site-header');
    if (!header) return;
    var toggle = header.querySelector('.nav-toggle');
    var nav = header.querySelector('.site-nav');
    if (!toggle || !nav) return;

    toggle.addEventListener('click', function () {
      var open = nav.classList.toggle('open');
      header.classList.toggle('nav-open', open);
      toggle.setAttribute('aria-expanded', open ? 'true' : 'false');
    });

    // Cierra el menú al navegar (cambio de página en multi-página).
    nav.addEventListener('click', function (e) {
      if (e.target.closest('a')) {
        nav.classList.remove('open');
        header.classList.remove('nav-open');
        toggle.setAttribute('aria-expanded', 'false');
      }
    });
  }

  function boot() {
    var h = document.getElementById('site-header');
    if (h) h.innerHTML = headerHTML();
    var f = document.getElementById('site-footer');
    if (f) f.innerHTML = footerHTML();
    // Barra inferior (mobile): se agrega una vez, al final de <body>.
    if (!document.getElementById('site-tabbar')) {
      document.body.insertAdjacentHTML('beforeend', tabbarHTML());
    }
    setActive();
    initToggle();
    initLogout();
  }

  // Cerrar sesión: POST /api/logout (borra la cookie) y volver al login.
  function initLogout() {
    var btn = document.getElementById('logout-btn');
    if (!btn) return;
    btn.addEventListener('click', function (e) {
      e.preventDefault();
      var api = (window.API_BASE || '/api') + '/logout';
      fetch(api, { method: 'POST', credentials: 'include' })
        .then(function () { window.location.href = '/login'; })
        .catch(function () { window.location.href = '/login'; });
    });
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', boot);
  } else {
    boot();
  }

  // API pública mínima para otras secciones.
  window.Layout = { SECTIONS: SECTIONS, setActive: setActive };
})();
