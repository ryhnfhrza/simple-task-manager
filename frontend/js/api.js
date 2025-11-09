const API_BASE = (function () {
  return `${window.APP_CONFIG.API_BASE_URL}/api/`;
})();

function extractErrorMessage(parsed, rawText, status) {
  if (!parsed && rawText) return rawText;
  if (!parsed) return `HTTP ${status}`;

  if (typeof parsed === "string") return parsed;
  if (parsed.error) return parsed.error;
  if (parsed.message) return parsed.message;
  if (parsed.data && parsed.data.message) return parsed.data.message;
  if (parsed.data && typeof parsed.data === "string") return parsed.data;
  if (parsed.errors) {
    if (Array.isArray(parsed.errors)) return parsed.errors.join(", ");
    if (typeof parsed.errors === "object") {
      try {
        return Object.values(parsed.errors).flat().join(", ");
      } catch (e) {}
    }
  }
  try {
    return JSON.stringify(parsed);
  } catch (e) {
    return `HTTP ${status}`;
  }
}

async function apiFetch(path, opts = {}) {
  const headers = opts.headers ? { ...opts.headers } : {};
  if (opts.body && !(opts.body instanceof FormData)) {
    headers["Content-Type"] = "application/json";
  }

  const token = authGetToken && authGetToken();
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const fetchOpts = { ...opts, headers };

  try {
    const res = await fetch(API_BASE + path, fetchOpts);

    const raw = await res.clone().text();

    let parsed = null;
    const contentType = res.headers.get("content-type") || "";
    if (contentType.includes("application/json")) {
      try {
        parsed = JSON.parse(raw);
      } catch (e) {
        parsed = null;
      }
    }

    if (!res.ok) {
      const msg = extractErrorMessage(parsed, raw, res.status);
      console.error(`[apiFetch] ERROR ${res.status} ${res.statusText} ->`, {
        path: API_BASE + path,
        request: fetchOpts,
        responseText: raw,
        responseJson: parsed,
      });
      if (res.status === 401) {
      authClear(); 
      window.location.replace("./index.html"); 
      return; 
  }

      throw new Error(msg);
    }

    if (parsed !== null) return parsed;
    if (raw) return raw;
    return {};
  } catch (err) {
    // network error
    if (err.name === "TypeError" || err.message === "Failed to fetch") {
      console.error("[apiFetch] Network/Fetch error ->", err);
      throw new Error(
        "Server tidak dapat dihubungi. Periksa koneksi Anda atau CORS."
      );
    }
    throw err;
  }
}

/* Auth */
async function apiRegister(payload) {
  return apiFetch("/register", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}
async function apiLogin(payload) {
  return apiFetch("/login", { method: "POST", body: JSON.stringify(payload) });
}

/* Tasks */
async function apiGetTasks(query = "") {
  return apiFetch("/tasks" + (query ? "?" + query : ""), { method: "GET" });
}
async function apiGetTask(taskId) {
  return apiFetch("/task/" + taskId, { method: "GET" });
}
async function apiCreateTask(payload) {
  return apiFetch("/tasks", { method: "POST", body: JSON.stringify(payload) });
}
async function apiUpdateTask(payload) {
  return apiFetch("/tasks/" + payload.id, {
    method: "PATCH",
    body: JSON.stringify(payload),
  });
}
async function apiDeleteTask(taskId) {
  return apiFetch("/tasks/" + taskId, { method: "DELETE" });
}
