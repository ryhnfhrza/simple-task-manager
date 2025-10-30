let currentPage = 0
let pageLimit = 20

function toast(rawMsg, type = 'info') {
  const t = document.getElementById('toast');
  const id = 't' + Date.now();
  const el = document.createElement('div');

  let msg = rawMsg;

  try {
    if (typeof rawMsg === 'string' && rawMsg.includes('Key:') && rawMsg.includes('Error:')) {
      const cleaned = rawMsg.replace(/\r?\n|\r/g, ' ');
      const parts = cleaned.split(/Key:\s*/g).map((x) => x.trim()).filter(Boolean);

      let messages = [];

      const requiredWithoutFields = [];

      parts.forEach((p) => {
        const fieldMatch = p.match(/'(?:\w+Request\.)?(\w+)'/);
        const tagMatch = p.match(/'(\w+)' tag/);
        let field = fieldMatch ? fieldMatch[1] : 'Field';
        const tag = tagMatch ? tagMatch[1] : '';

        field = field.charAt(0).toUpperCase() + field.slice(1);

        if (tag === 'required_without') {
          requiredWithoutFields.push(field);
        } else {
          switch (tag) {
            case 'required':
              messages.push(`${field} wajib diisi.`);
              break;
            case 'min':
              messages.push(`${field} terlalu pendek.`);
              break;
            case 'max':
              messages.push(`${field} terlalu panjang.`);
              break;
            case 'email':
              messages.push(`${field} harus berupa alamat email yang valid.`);
              break;
            case 'passComplex':
              messages.push(`${field} harus mengandung huruf besar, huruf kecil, angka, dan simbol.`);
              break;
            case 'alphanum':
              messages.push(`${field} hanya boleh berisi huruf dan angka.`);
              break;
            case 'lte':
              messages.push(`${field} harus lebih kecil atau sama dengan nilai yang diizinkan.`);
              break;
            case 'gte':
              messages.push(`${field} harus lebih besar atau sama dengan nilai yang diizinkan.`);
              break;
            default:
              messages.push(`${field} tidak valid.`);
          }
        }
      });

      if (requiredWithoutFields.length >= 2) {
        messages.unshift(`Harus diisi salah satu antara ${requiredWithoutFields.join(' dan ')}.`);
      } else if (requiredWithoutFields.length === 1) {
        messages.unshift(`${requiredWithoutFields[0]} wajib diisi jika kolom lain kosong.`);
      }

      msg = messages.join('\n');
    }
  } catch {
    msg = rawMsg;
  }

  el.id = id;
  el.className =
    'p-3 rounded-md mb-2 ' +
    (type === 'error'
      ? 'bg-red-100 text-red-700 border border-red-300'
      : 'bg-emerald-100 text-emerald-700 border border-emerald-300') +
    ' shadow transition-all duration-300 whitespace-pre-line';
  el.innerText = msg;

  t.appendChild(el);
  setTimeout(() => el.remove(), 4000);
}



async function loadUserInfo() {
  const el = document.getElementById('user-info')
  el.textContent = 'You'
}

function buildQuery() {
  const dueBefore = document.getElementById('filter-due-before').value.trim();
  const dueAfter = document.getElementById('filter-due-after').value.trim();
  const completed = document.getElementById('filter-completed').value;
  const sortBy = document.getElementById('sort-by').value;
  const order = document.getElementById('order').value;
  const limitInput = document.getElementById('limit').value;
  const limit = limitInput ? parseInt(limitInput) : pageLimit;
  const offset = currentPage * limit;

  const params = new URLSearchParams();

  const dateRegex = /^\d{4}-\d{2}-\d{2}(T\d{2}:\d{2}(:\d{2})?)?$/;

  if (completed !== '') params.append('completed', completed);

  if (dueBefore && dateRegex.test(dueBefore)) {
    params.append('due_before', dueBefore);
  }
  if (dueAfter && dateRegex.test(dueAfter)) {
    params.append('due_after', dueAfter);
  }

  if (sortBy) params.append('sort_by', sortBy);
  if (order) params.append('order', order);
  if (limit) params.append('limit', limit);
  if (offset) params.append('offset', offset);

  return params.toString();
}


async function fetchAndRender() {
  try {
    const q = buildQuery()
    const body = await apiGetTasks(q)
    const tasks = body.data || []
    renderTasks(tasks)
    document.getElementById('paging-info').innerText = `Showing ${tasks.length} items (page ${currentPage + 1})`
  } catch (e) {
    toast(e.message, 'error')
    if (e.message.toLowerCase().includes('unauthorized')) {
      authClear()
      location.href = './index.html'
    }
  }
}

function renderTasks(tasks) {
  const container = document.getElementById('task-list')
  container.innerHTML = ''
  if (!tasks.length) {
    container.innerHTML = `<div class="p-4 bg-white dark:bg-slate-800 rounded-md text-slate-500">No tasks</div>`
    return
  }
  tasks.forEach(t => {
    const due = t.due_date ? new Date(t.due_date) : null;
    const now = new Date();
    const isOverdue = due && due < now && !t.completed;

    const card = document.createElement('div');
    card.classList.add(
      'task-card',
      'animate-fadeInUp',
      'p-4',
      'rounded-md',
      'shadow',
      'flex',
      'justify-between',
      'items-start'
    );

    if (t.completed) {
      card.classList.add('task-completed');
    } else if (isOverdue) {
      card.classList.add('task-overdue');
    } else {
      card.classList.add('task-default');
    }

    const titleClass = t.completed
      ? 'font-semibold line-through text-slate-400 flex items-center gap-1'
      : 'font-semibold flex items-center gap-1';

    const statusIcon = t.completed
      ? '<span class="text-emerald-600 dark:text-emerald-400">✅</span>'
      : '';

    card.innerHTML = `
  <div class="flex-1">
    <div class="flex items-center gap-3">
      <h3 class="font-semibold ${t.completed ? "line-through text-slate-400" : ""}">
        ${escapeHtml(t.title || "(no title)")}
      </h3>
    </div>
    <p class="mt-1 text-sm text-slate-300">${escapeHtml(t.description || "No Desc")}</p>
    <div class="mt-2 text-xs text-slate-500">
      ${due
        ? `<span>Due: ${new Date(due).toLocaleString()}</span>`
        : "<span>No due date</span>"
      }
      <span class="ml-3">Created: ${new Date(t.created_at).toLocaleString()}</span>
    </div>
  </div>

  <div class="flex flex-col items-end gap-2 ml-4">
    <div class="flex gap-2">
      <button data-id="${t.id}" class="btn-edit px-2 py-1 bg-sky-600 text-white rounded-md">Edit</button>
      <button data-id="${t.id}" class="btn-delete px-2 py-1 bg-red-500 text-white rounded-md">Delete</button>
    </div>
    <div>
      <label class="inline-flex items-center gap-2 text-sm">
        <input type="checkbox" data-id="${t.id}" class="toggle-complete" ${t.completed ? "checked" : ""} />
        <span>${t.completed ? "Completed" : "Mark done"}</span>
      </label>
    </div>
  </div>
`;


    card.style.animationDelay = `${tasks.indexOf(t) * 0.05}s`;
    container.appendChild(card);


    container.appendChild(card);
  });


  document.querySelectorAll('.btn-delete').forEach(btn => {
    btn.onclick = async (e) => {
      const id = btn.getAttribute('data-id')
      if (!confirm('Delete this task?')) return
      try {
        await apiDeleteTask(id)
        toast('Deleted', 'info')
        await fetchAndRender()
      } catch (err) { toast(err.message, 'error') }
    }
  })

  document.querySelectorAll('.btn-edit').forEach(btn => {
    btn.onclick = async () => {
      const id = btn.getAttribute('data-id')
      try {
        const res = await apiGetTask(id)
        const t = res.data
        openModal('Edit Task', t)
      } catch (err) { toast(err.message, 'error') }
    }
  })

  document.querySelectorAll('.toggle-complete').forEach(cb => {
    cb.onchange = async () => {
      const id = cb.getAttribute('data-id');
      const card = cb.closest('.task-card');
      try {
        card.style.opacity = "0.5";
        await apiUpdateTask({ id: parseInt(id), completed: cb.checked ? 1 : 0 });
        card.style.opacity = "1";
        toast('Status updated', 'success');
        await fetchAndRender();
      } catch (err) {
        toast(err.message, 'error');
        card.style.opacity = "1";
      }
    };
  });
}

// modal logic
const modal = document.getElementById('modal')
const formTask = document.getElementById('form-task')
let editingId = null

function openModal(title, task = null) {
  document.getElementById('modal-title').innerText = title
  document.getElementById('task-title').value = task?.title || ''
  document.getElementById('task-desc').value = task?.description || ''
  document.getElementById('task-due').value = task?.due_date ? new Date(task.due_date).toISOString().slice(0, 19) : ''
  editingId = task?.id || null
  modal.classList.remove('hidden')
}

document.getElementById('btn-open-create').onclick = () => openModal('Create Task')

document.getElementById('modal-cancel').onclick = () => modal.classList.add('hidden')

formTask.addEventListener('submit', async (e) => {
  e.preventDefault()
  const submitBtn = document.getElementById('modal-submit')
  submitBtn.disabled = true
  submitBtn.innerText = 'Saving...'
  try {
    const title = document.getElementById('task-title').value.trim()
    const description = document.getElementById('task-desc').value.trim()
    const due = document.getElementById('task-due').value.trim()
    const payload = { title, description }
    if (due) payload.due_date = due
    if (editingId) {
      payload.id = editingId
      await apiUpdateTask(payload)
      toast('Updated task', 'success')
    } else {
      await apiCreateTask(payload)
      toast('Created task', 'success')
    }
    modal.classList.add('hidden')
    formTask.reset()
    editingId = null
    await fetchAndRender()
  } catch (err) {
    toast(err.message, 'error')
  } finally {
    submitBtn.disabled = false
    submitBtn.innerText = 'Save'
  }
})

// pagination
document.getElementById('btn-refresh').onclick = () => { currentPage = 0; fetchAndRender() }
document.getElementById('prev-page').onclick = () => { if (currentPage > 0) { currentPage--; fetchAndRender() } }
document.getElementById('next-page').onclick = () => { currentPage++; fetchAndRender() }

// helper debounce
function debounce(fn, delay = 700) {
  let timer;
  return (...args) => {
    clearTimeout(timer);
    timer = setTimeout(() => fn(...args), delay);
  };
}

function attachFilterListeners() {
  const ids = [
    'filter-due-before',
    'filter-due-after',
    'filter-completed',
    'sort-by',
    'order',
    'limit'
  ];

  const debouncedFetch = debounce(() => {
    currentPage = 0;
    fetchAndRender();
  }, 700);

  ids.forEach(id => {
    const el = document.getElementById(id);
    if (!el) return;

    if (el.tagName.toLowerCase() === 'input') {
      el.addEventListener('input', debouncedFetch);
    } else {
      el.addEventListener('change', debouncedFetch);
    }
  });
}


// init
document.addEventListener('DOMContentLoaded', () => {
  console.log("Dashboard loaded, fetching tasks...");
  loadTheme()
  loadUserInfo()
  document.getElementById('btn-logout').onclick = () => logoutAndBack()
  attachFilterListeners()
  fetchAndRender()
})
// helper
function escapeHtml(s) { if (!s) return ''; return s.replaceAll('<', '&lt;').replaceAll('>', '&gt;') }
