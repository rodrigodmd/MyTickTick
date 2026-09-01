// MyTickTick - API base + cliente de peticiones (7.4)
//
// Todas las secciones usan window.fetchAPI(...) para hablar con el backend.
// El layout (header/nav/footer) lo maneja layout.js; aquí solo vivimos la capa de datos.
(function () {
  'use strict';

  // Base URL de la API (relativa al origen: el backend Go sirve el frontend
  // y expone los endpoints bajo /api).
  window.API_BASE = window.API_BASE || '/api';

  // fetchAPI(endpoint, options)
  //   - endpoint: p. ej. '/monthly-tasks' o '/trackers/3/history'
  //   - options:  opciones de fetch; si `options.json` es un objeto se
  //               serializa a JSON y se fija el header Content-Type.
  // Devuelve el cuerpo ya parseado (JSON) o null si la respuesta no tiene cuerpo.
  // Lanza un Error legible si la respuesta no es 2xx o no es JSON.
  // Incluye credentials: 'include' para enviar la cookie de sesión (httpOnly).
  async function fetchAPI(endpoint, options) {
    var opts = Object.assign({ credentials: 'include' }, options || {});
    var headers = Object.assign({}, opts.headers || {});
    var body = opts.body;

    if (opts.json !== undefined) {
      body = JSON.stringify(opts.json);
      headers['Content-Type'] = 'application/json';
    } else if (typeof body === 'object' && body !== null && !(body instanceof FormData) && body.constructor === Object) {
      body = JSON.stringify(body);
      headers['Content-Type'] = 'application/json';
    }

    var url = (endpoint.startsWith('http') ? '' : window.API_BASE) + endpoint;
    var response;
    try {
      response = await fetch(url, Object.assign({}, opts, { headers: headers, body: body }));
    } catch (networkErr) {
      throw new Error('No se pudo conectar con el servidor: ' + networkErr.message);
    }

    // 401 Unauthorized: sesión expirada o no iniciada → redirigir a login.
    if (response.status === 401) {
      var path = window.location.pathname;
      // No redirigir si ya estamos en login o si la ruta es pública.
      if (path !== '/login' && path !== '/api/register' && path !== '/api/login' && path !== '/api/health') {
        window.location.href = '/login?next=' + encodeURIComponent(path);
      }
      var err401 = new Error('Sesión expirada. Iniciá sesión para continuar.');
      err401.status = 401;
      err401.unauthorized = true;
      throw err401;
    }

    // 204 No Content no tiene cuerpo.
    if (response.status === 204) {
      return null;
    }

    // Intentar leer JSON; si no es JSON (texto/error), devolver el texto.
    var text = await response.text();
    var parsed = null;
    if (text) {
      try {
        parsed = JSON.parse(text);
      } catch (e) {
        parsed = text; // respuesta no-JSON (p. ej. error de texto plano)
      }
    }

    if (!response.ok) {
      var message = (parsed && parsed.error) ? parsed.error : (response.status + ' ' + response.statusText);
      var err = new Error('API error: ' + message);
      err.status = response.status;
      err.body = parsed;
      throw err;
    }

    return parsed;
  }

  window.fetchAPI = fetchAPI;
  window.MyTickTick = {
    fetchAPI: fetchAPI,
    API_BASE: window.API_BASE
  };
})();
