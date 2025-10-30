function parseFlexible(dateStr) {
  if (!dateStr) return null
  let d = new Date(dateStr)
  if (!isNaN(d)) return d
  const reDate = /^(\d{4}-\d{2}-\d{2})$/
  const m = dateStr.match(reDate)
  if (m) return new Date(m[1] + 'T00:00:00')
  return null
}

function isoDate(d) {
  if (!d) return null
  try {
    return new Date(d).toISOString()
  } catch (e) { return null }
}

function fmtLocal(d) {
  if (!d) return ''
  const dt = new Date(d)
  return dt.toLocaleString()
}

// theme
const THEME_KEY = 'tm_theme'
function applyTheme(mode) {
  if (mode === 'dark') document.documentElement.classList.add('dark')
  else if (mode === 'light') document.documentElement.classList.remove('dark')
  else {
    // auto: follow system
    const prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
    if (prefersDark) document.documentElement.classList.add('dark')
    else document.documentElement.classList.remove('dark')
  }
  localStorage.setItem(THEME_KEY, mode)
}
function loadTheme() {
  const saved = localStorage.getItem(THEME_KEY) || 'auto'
  const sel = document.getElementById('theme-select')
  if (sel) { sel.value = saved; sel.onchange = () => applyTheme(sel.value) }
  applyTheme(saved)
}
document.addEventListener('DOMContentLoaded', () => {
  try { loadTheme() } catch (e) { }
})
