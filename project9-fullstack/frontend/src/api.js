// Relative URL: works with the Vite proxy in development
// and when Gin serves the build in production.
const BASE = '/api'

async function request(path, options = {}) {
  const res = await fetch(BASE + path, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (res.status === 204) return null

  const data = await res.json().catch(() => null)
  if (!res.ok) {
    throw new Error(data?.error || `Request failed with status ${res.status}`)
  }
  return data
}

export const api = {
  listTasks: (status = 'all') => request(`/tasks?status=${status}`),
  createTask: (task) => request('/tasks', { method: 'POST', body: JSON.stringify(task) }),
  updateTask: (id, task) => request(`/tasks/${id}`, { method: 'PUT', body: JSON.stringify(task) }),
  deleteTask: (id) => request(`/tasks/${id}`, { method: 'DELETE' }),
}
