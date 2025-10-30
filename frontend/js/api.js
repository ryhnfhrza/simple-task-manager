// js/api.js
const API_BASE = (function () { return "http://localhost:3000/api" })()

/**
 * Helpers untuk mengekstrak pesan error dari berbagai bentuk respons backend
 */
function extractErrorMessage(parsed, rawText, status) {
  // parsed bisa object / array / null
  if (!parsed && rawText) return rawText;
  if (!parsed) return `HTTP ${status}`;

  // beberapa pola umum
  if (typeof parsed === 'string') return parsed;
  if (parsed.error) return parsed.error;
  if (parsed.message) return parsed.message;
  if (parsed.data && parsed.data.message) return parsed.data.message;
  if (parsed.data && typeof parsed.data === 'string') return parsed.data;
  // kadang backend mengirim errors: { field: ["msg"] } atau array
  if (parsed.errors) {
    if (Array.isArray(parsed.errors)) return parsed.errors.join(', ');
    if (typeof parsed.errors === 'object') {
      try {
        return Object.values(parsed.errors).flat().join(', ');
      } catch (e) { }
    }
  }
  // fallback: JSON stringify (singkat)
  try {
    return JSON.stringify(parsed);
  } catch (e) {
    return `HTTP ${status}`;
  }
}

async function apiFetch(path, opts = {}) {
  const headers = opts.headers ? { ...opts.headers } : {};
  // hanya set content-type bila ada body dan bukan form-data
  if (opts.body && !(opts.body instanceof FormData)) {
    headers["Content-Type"] = "application/json";
  }

  const token = authGetToken && authGetToken();
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const fetchOpts = { ...opts, headers };

  try {
    const res = await fetch(API_BASE + path, fetchOpts);

    // baca raw text dulu (untuk kasus non-json)
    const raw = await res.clone().text();

    // coba parse json, kalau gagal parsed = null
    let parsed = null;
    const contentType = res.headers.get('content-type') || '';
    if (contentType.includes('application/json')) {
      try {
        parsed = JSON.parse(raw);
      } catch (e) {
        parsed = null;
      }
    }

    if (!res.ok) {
      // build message
      const msg = extractErrorMessage(parsed, raw, res.status);
      // console logging lengkap untuk debugging
      console.error(`[apiFetch] ERROR ${res.status} ${res.statusText} ->`, {
        path: API_BASE + path,
        request: fetchOpts,
        responseText: raw,
        responseJson: parsed,
      });
      throw new Error(msg);
    }

    // sukses
    // return parsed json jika ada, else return raw text or empty object
    if (parsed !== null) return parsed;
    if (raw) return raw;
    return {}; // empty 204/empty response
  } catch (err) {
    // network error
    if (err.name === 'TypeError' || err.message === 'Failed to fetch') {
      // console.debug untuk developer, user-facing pesan dalam bahasa Indonesia
      console.error('[apiFetch] Network/Fetch error ->', err);
      throw new Error('Server tidak dapat dihubungi. Periksa koneksi Anda atau CORS.');
    }
    // already Error dari server / throw kembali
    throw err;
  }
}

/* Auth */
async function apiRegister(payload) { return apiFetch('/register', { method: 'POST', body: JSON.stringify(payload) }) }
async function apiLogin(payload) { return apiFetch('/login', { method: 'POST', body: JSON.stringify(payload) }) }

/* Tasks */
async function apiGetTasks(query = '') { return apiFetch('/tasks' + (query ? ('?' + query) : ''), { method: 'GET' }) }
async function apiGetTask(taskId) { return apiFetch('/task/' + taskId, { method: 'GET' }) }
async function apiCreateTask(payload) { return apiFetch('/tasks', { method: 'POST', body: JSON.stringify(payload) }) }
async function apiUpdateTask(payload) { return apiFetch('/tasks/' + payload.id, { method: 'PATCH', body: JSON.stringify(payload) }) }
async function apiDeleteTask(taskId) { return apiFetch('/tasks/' + taskId, { method: 'DELETE' }) }
