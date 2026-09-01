package static

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// StaticHandler sirve la web app (web/) con rutas limpias por sección
// (7.2) y los assets (css/js) desde la raíz.
type StaticHandler struct {
	assetPath string
	// pages mapea cada ruta de sección a su archivo HTML.
	pages map[string]string
}

// NewStaticHandler crea una nueva instancia de StaticHandler.
// assetPath es la carpeta raíz de la web app (p. ej. "web").
func NewStaticHandler(assetPath string) *StaticHandler {
	return &StaticHandler{
		assetPath: assetPath,
		pages: map[string]string{
			"/":                  "html/index.html",
			"/monthly-tasks":     "html/monthly-tasks.html",
			"/immediate-tasks":   "html/immediate-tasks.html",
			"/trackers":          "html/trackers.html",
			"/dashboard":         "html/dashboard.html",
			// Autenticación (0.3)
			"/login":         "html/login.html",
			"/auth/login":    "html/login.html",
			"/auth/register": "html/login.html",
		},
	}
}

// ServeHTTP implementa http.Handler para servir la web app.
func (h *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Rutas de sección → archivo HTML correspondiente.
	if page, ok := h.pages[path]; ok {
		http.ServeFile(w, r, h.assetPath+"/"+page)
		return
	}

	// Assets (css/, js/) y cualquier otro archivo estático.
	http.FileServer(http.Dir(h.assetPath)).ServeHTTP(w, r)
}

// MapURLs configura las rutas para servir archivos estáticos.
func (h *StaticHandler) MapURLs(mux *http.ServeMux) {
	// El patrón "/" con el método vacío captura todas las rutas restantes,
	// así que se registra al final, después de las rutas de API.
	mux.Handle("/", h)
}

// ProxyHandler crea un proxy para el frontend cuando necesita consumir la API.
type ProxyHandler struct {
	backendURL *url.URL
	proxy      *httputil.ReverseProxy
}

// NewProxyHandler crea una nueva instancia de ProxyHandler.
func NewProxyHandler(backendAddr string) *ProxyHandler {
	backendURL, _ := url.Parse("http://" + backendAddr)
	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	return &ProxyHandler{
		backendURL: backendURL,
		proxy:      proxy,
	}
}

// ServeHTTP implementa http.Handler para el proxy.
func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Cambiar el path para que apunte a la API.
	r.URL.Path = "/api" + r.URL.Path
	p.proxy.ServeHTTP(w, r)
}
