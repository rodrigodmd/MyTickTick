// MyTickTick - Dashboard (usa el cliente fetchAPI de app.js)
//
// Fuente de datos:
//   GET /api/metrics?month=M&year=Y  → { monthly, monthlySeries, trackers }
//   GET /api/immediate-tasks         → para el contador de pendientes
//   GET /api/trackers                → para el contador de trackers activos
//   GET /api/trackers/{id}/history   → para las series de los gráficos por tracker
let monthlyChart = null;
let trackerCharts = []; // un Chart por tracker (carousel)

const MONTH_NAMES = ['Enero', 'Febrero', 'Marzo', 'Abril', 'Mayo', 'Junio',
    'Julio', 'Agosto', 'Septiembre', 'Octubre', 'Noviembre', 'Diciembre'];

// Colores del gráfico de trackers (alineados con --success / --danger).
// El umbral es cian para no mezclarse con verde, rojo ni el acento naranja.
const TRACKER_LINE_MET = '#10b981';
const TRACKER_LINE_MISS = '#ef4444';
const TRACKER_LINE_LIMIT = '#22d3ee';
const TRACKER_LINE_NEUTRAL = '#60a5fa';

document.addEventListener('DOMContentLoaded', function () {
    initFilters();
    loadDashboard();
});

// Devuelve [first, last] (YYYY-MM-DD) del mes/año indicado.
function monthBounds(month, year) {
    const first = `${year}-${String(month).padStart(2, '0')}-01`;
    const last = new Date(Date.UTC(year, month, 0));
    return [first, last.toISOString().slice(0, 10)];
}

// Llena los selects de mes/año y fija el mes en curso.
function initFilters() {
    const monthSel = document.getElementById('metric-month');
    const yearSel = document.getElementById('metric-year');
    const now = new Date();

    MONTH_NAMES.forEach((name, i) => {
        const opt = document.createElement('option');
        opt.value = String(i + 1);
        opt.textContent = name;
        monthSel.appendChild(opt);
    });
    monthSel.value = String(now.getMonth() + 1);

    for (let y = now.getFullYear() - 2; y <= now.getFullYear() + 1; y++) {
        const opt = document.createElement('option');
        opt.value = String(y);
        opt.textContent = String(y);
        yearSel.appendChild(opt);
    }
    yearSel.value = String(now.getFullYear());
}

// Carga las métricas y redibuja ambos gráficos.
async function loadDashboard() {
    try {
        const month = document.getElementById('metric-month').value;
        const year = document.getElementById('metric-year').value;

        const [metrics, immediateTasks, trackers] = await Promise.all([
            fetchAPI(`/metrics?month=${month}&year=${year}`),
            fetchAPI('/immediate-tasks'),
            fetchAPI('/trackers')
        ]);

        renderMetrics(metrics, immediateTasks, trackers, month, year);
        renderMonthlyChart(metrics.monthlySeries);
        await renderTrackerChart(trackers, month, year);
    } catch (error) {
        console.error('Error cargando el dashboard:', error);
        alert('Error al cargar el dashboard: ' + error.message);
    }
}

// Tarjetas de resumen + tabla de métricas por tracker (11.6).
function renderMetrics(metrics, immediateTasks, trackers, month, year) {
    const m = metrics.monthly;
    const rate = Math.round((m.rate || 0) * 1000) / 10;
    document.getElementById('monthly-completion-rate').textContent =
        `${m.completed}/${m.total} (${rate}%)`;

    const pending = immediateTasks.filter(t => !t.isCompleted).length;
    document.getElementById('immediate-pending-count').textContent = String(pending);

    document.getElementById('active-trackers-count').textContent = String(trackers.length);

    document.getElementById('metrics-period-label').textContent =
        `${MONTH_NAMES[(month - 1) % 12]} ${year}`;
    renderMetricsTable(metrics.trackers);
}

// Tabla de métricas de cumplimiento por tracker para el período elegido:
// días con registro, días cumplidos, tasa % y desviación media.
function renderMetricsTable(trackers) {
    const wrap = document.getElementById('metrics-table-wrap');
    if (!trackers || trackers.length === 0) {
        wrap.innerHTML = '<p class="empty-state">No hay trackers registrados.</p>';
        return;
    }

    const fmtDev = v => {
        const n = Math.round(v * 10) / 10;
        return (n > 0 ? '+' : '') + n;
    };

    const rows = trackers.map(t => {
        const pct = Math.round((t.rate || 0) * 100);
        const metClass = t.total > 0 && t.met === t.total ? 'met-ok' : (t.met > 0 ? 'met-part' : 'met-bad');
        return `<tr>
            <td>${t.name}</td>
            <td>${t.total}</td>
            <td>${t.met}</td>
            <td class="${metClass}">${t.total ? pct + '%' : '&mdash;'}</td>
            <td>${t.total ? fmtDev(t.avgDeviation) : '&mdash;'}</td>
            <td>${t.total ? fmtDev(t.avgAbsDeviation) : '&mdash;'}</td>
        </tr>`;
    }).join('');

    wrap.innerHTML = `
        <table class="metrics-table">
            <thead>
                <tr>
                    <th>Tracker</th>
                    <th>Días con registro</th>
                    <th>Días cumplidos</th>
                    <th>Tasa</th>
                    <th>Desviación media</th>
                    <th>Desviación media (abs)</th>
                </tr>
            </thead>
            <tbody>${rows}</tbody>
        </table>`;
}

// Gráfico de barras: cumplimiento mensual (serie anual).
function renderMonthlyChart(series) {
    const ctx = document.getElementById('monthly-chart').getContext('2d');
    if (monthlyChart) {
        monthlyChart.destroy();
    }

    monthlyChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: (series || []).map(p => MONTH_NAMES[p.month - 1]),
            datasets: [{
                label: '% de cumplimiento',
                data: (series || []).map(p => Math.round((p.rate || 0) * 100)),
                backgroundColor: 'rgba(52, 152, 219, 0.6)',
                borderColor: 'rgba(52, 152, 219, 1)',
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            scales: { y: { beginAtZero: true, max: 100 } },
            plugins: {
                legend: { display: false },
                tooltip: {
                    callbacks: {
                        afterLabel: (item) => {
                            const p = (series || [])[item.dataIndex];
                            return p ? `${p.completed}/${p.total} tareas` : '';
                        }
                    }
                }
            }
        }
    });
}

// Escala Y "ajustada": el eje no arranca en 0 sino cerca del mínimo/máximo
// de las mediciones (con un margen ~10% de padding), para que variaciones
// chicas (ej. peso 74–76 kg) se vean sin que la línea se aplique al fondo.
// Si solo hay un valor se usa el mismo min/max alrededor de ese valor.
function computeYScale(values) {
    const nums = values.filter(v => typeof v === 'number' && isFinite(v));
    if (nums.length === 0) {
        return { beginAtZero: true };
    }
    let min = Math.min(...nums);
    let max = Math.max(...nums);
    if (min === max) {
        // Un único valor (o todos iguales): margen fijo de ±10% (mínimo 1).
        const pad = Math.max(Math.abs(min) * 0.1, 1);
        min -= pad;
        max += pad;
    } else {
        const span = max - min;
        const lo = min - span * 0.1;
        // Solo se pega al 0 si los datos lo cruzan (o son todos positivos);
        // con todos negativos el 0 no tiene sentido como piso.
        min = (lo < 0 && max >= 0) ? 0 : lo;
        max = max + span * 0.1;
    }
    return { min: roundScale(min), max: roundScale(max) };
}

// Redondea el límite a un número "redondo" para que Chart.js genere ticks
// limpios (0.1/0.5/1/2/5/10 * 10^n) en lugar de 74.285714.
function roundScale(n) {
    const abs = Math.abs(n);
    if (abs >= 1000) return Math.round(n);
    if (abs >= 100) return Math.round(n * 10) / 10;
    if (abs >= 10) return Math.round(n * 100) / 100;
    return Math.round(n * 1000) / 1000;
}

// Número sin ceros de basura: 80 → "80", 79.5 → "79.5".
function fmtLimit(v) {
    const n = Math.round(v * 1000) / 1000;
    return String(n);
}

// Usa isMet del backend; si falta, evalúa contra el límite (misma regla que Go).
// null = el tracker no tiene límite (línea neutra).
function entryIsMet(entry, tracker, hasLimit) {
    if (!hasLimit) return null;
    if (typeof entry.isMet === 'boolean') return entry.isMet;
    const v = Number(entry.value);
    const limit = Number(tracker.limitValue);
    if (!isFinite(v) || !isFinite(limit)) return null;
    if (tracker.limitType === 'min') return v >= limit;
    return v <= limit;
}

function colorForMet(isMet) {
    if (isMet === null || isMet === undefined) return TRACKER_LINE_NEUTRAL;
    return isMet ? TRACKER_LINE_MET : TRACKER_LINE_MISS;
}

// Nunca tira: Chart.js a veces llama esto sin p1DataIndex.
function segmentMetColor(ctx, mets) {
    try {
        if (!ctx) return undefined;
        const idx = (ctx.p1DataIndex != null) ? ctx.p1DataIndex : ctx.dataIndex;
        if (idx == null || idx < 0 || idx >= mets.length) return undefined;
        return colorForMet(mets[idx]);
    } catch (e) {
        return undefined;
    }
}

// Carousel: un gráfico de líneas POR TRACKER dentro del rango del mes/año
// seleccionado. Cada tracker tiene su propia pantalla (peso no se mezcla con
// sueño, etc.) y se desliza entre ellos con swipe o con las flechas.
async function renderTrackerChart(trackers, month, year) {
    const [from, to] = monthBounds(month, year);
    const carousel = document.getElementById('tracker-carousel');

    // Destruye los gráficos anteriores (cambio de mes/año o recarga).
    trackerCharts.forEach(chart => chart.destroy());
    trackerCharts = [];
    carousel.innerHTML = '';

    const list = (trackers || []);
    if (list.length === 0) {
        carousel.innerHTML = '<div class="chart-container tracker-slide"><p class="empty-state">No hay trackers registrados.</p></div>';
        return;
    }

    // Serie de cada tracker, en paralelo.
    const series = await Promise.all(list.map(async (tracker) => {
        try {
            const history = await fetchAPI(`/trackers/${tracker.id}/history`);
            const rows = history
                .slice()
                .sort((a, b) => a.entryDate.localeCompare(b.entryDate))
                .filter(h => h.entryDate >= from && h.entryDate <= to);
            return { tracker, rows };
        } catch (error) {
            console.error(`Error cargando historial de ${tracker.name}:`, error);
            return { tracker, rows: [] };
        }
    }));

    series.forEach(({ tracker, rows }) => {
        const slide = document.createElement('div');
        slide.className = 'tracker-slide chart-container';

        const h3 = document.createElement('h3');
        h3.textContent = tracker.name;
        slide.appendChild(h3);
        carousel.appendChild(slide);

        if (rows.length === 0) {
            const p = document.createElement('p');
            p.className = 'empty-state';
            p.textContent = 'Sin registros en este período.';
            slide.appendChild(p);
            return;
        }

        const canvas = document.createElement('canvas');
        slide.appendChild(canvas);
        const yScale = computeYScale(rows.map(r => r.value));
        const limitValue = Number(tracker.limitValue);
        const hasLimit = isFinite(limitValue) && tracker.limitValue !== '';
        const unit = tracker.unit ? ' ' + tracker.unit : '';
        const limitLabel = hasLimit
            ? ((tracker.limitType === 'min' ? 'Mínimo' : 'Máximo') + ': ' + fmtLimit(limitValue) + unit)
            : null;

        const mets = rows.map(r => entryIsMet(r, tracker, hasLimit));
        const pointColors = mets.map(colorForMet);
        const datasets = [{
            label: tracker.name,
            data: rows.map(r => r.value),
            fill: false,
            tension: 0.1,
            spanGaps: true,
            pointRadius: 4,
            pointHoverRadius: 6,
            pointBackgroundColor: pointColors,
            pointBorderColor: pointColors,
            borderColor: hasLimit ? TRACKER_LINE_MET : TRACKER_LINE_NEUTRAL,
            borderWidth: 2,
            segment: {
                borderColor: (ctx) => segmentMetColor(ctx, mets)
            }
        }];

        if (hasLimit) {
            // Límite actual del tracker: línea horizontal punteada sobre
            // todo el período. La escala Y se expande para que siempre sea
            // visible, aunque el valor esté fuera del rango de los datos.
            delete yScale.beginAtZero;
            if (yScale.min === undefined || limitValue < yScale.min) yScale.min = limitValue;
            if (yScale.max === undefined || limitValue > yScale.max) yScale.max = limitValue;
            if (yScale.min !== undefined && yScale.max !== undefined && yScale.min === yScale.max) {
                const pad = Math.max(Math.abs(yScale.min) * 0.1, 1);
                yScale.min -= pad;
                yScale.max += pad;
            }
            datasets.push({
                label: limitLabel,
                data: rows.map(() => limitValue),
                borderDash: [6, 4],
                borderWidth: 2,
                borderColor: TRACKER_LINE_LIMIT,
                backgroundColor: TRACKER_LINE_LIMIT,
                pointRadius: 0,
                fill: false
            });
        }

        try {
            const chart = new Chart(canvas.getContext('2d'), {
                type: 'line',
                data: {
                    labels: rows.map(r => r.entryDate),
                    datasets: datasets
                },
                options: {
                    responsive: true,
                    scales: {
                        x: {
                            ticks: {
                                maxRotation: 45,
                                autoSkip: true,
                                maxTicksLimit: 12
                            }
                        },
                        y: yScale
                    },
                    plugins: {
                        legend: { display: false },
                        tooltip: {
                            filter: function (item) {
                                return item.datasetIndex === 0;
                            },
                            callbacks: {
                                label: (item) => {
                                    const base = ` ${item.parsed.y}${unit}`;
                                    const met = mets[item.dataIndex];
                                    if (met === null || met === undefined) return base;
                                    return base + (met ? ' · cumplió' : ' · no cumplió');
                                }
                            }
                        }
                    }
                }
            });
            trackerCharts.push(chart);
        } catch (err) {
            console.error('Error dibujando gráfico de ' + tracker.name + ':', err);
            const p = document.createElement('p');
            p.className = 'empty-state';
            p.textContent = 'No se pudo dibujar el gráfico.';
            slide.appendChild(p);
        }
    });
}

// Desliza el carousel de trackers (flechas). El swipe táctil funciona sin
// más: el scroll-snap del CSS hace el resto.
function scrollTrackerCarousel(direction) {
    const carousel = document.getElementById('tracker-carousel');
    if (!carousel) return;
    const slide = carousel.querySelector('.tracker-slide');
    const step = slide ? slide.clientWidth : carousel.clientWidth;
    carousel.scrollBy({ left: direction * step, behavior: 'smooth' });
}
