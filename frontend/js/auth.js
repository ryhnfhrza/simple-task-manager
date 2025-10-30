// auth helpers
const TOKEN_KEY = 'tm_token'
function authSaveToken(token) { localStorage.setItem(TOKEN_KEY, token) }
function authGetToken() { return localStorage.getItem(TOKEN_KEY) }
function authClear() { localStorage.removeItem(TOKEN_KEY) }
function authIsLogged() { return !!authGetToken() }

// redirect to login if not logged
document.addEventListener('DOMContentLoaded', () => {
  // If on dashboard and not logged, go back to index
  if (location.pathname.endsWith('dashboard.html') && !authIsLogged()) {
    location.href = './index.html'
  }
})
async function logoutAndBack() {
  authClear()
  location.href = './index.html'
}
