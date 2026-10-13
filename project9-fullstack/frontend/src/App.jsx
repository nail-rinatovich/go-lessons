import { useEffect, useState } from 'react'
import { api } from './api.js'
import TaskForm from './components/TaskForm.jsx'
import TaskItem from './components/TaskItem.jsx'

const FILTERS = ['all', 'active', 'done']

export default function App() {
  const [tasks, setTasks] = useState([])
  const [filter, setFilter] = useState('all')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  // Reload the list whenever the filter changes
  useEffect(() => {
    let ignore = false
    setLoading(true)
    setError('')
    api
      .listTasks(filter)
      .then((data) => { if (!ignore) setTasks(data) })
      .catch((err) => { if (!ignore) setError(err.message) })
      .finally(() => { if (!ignore) setLoading(false) })
    // Ignore a stale response if the filter changed before it arrived
    return () => { ignore = true }
  }, [filter])

  async function handleCreate(input) {
    const created = await api.createTask(input)
    if (filter !== 'done') setTasks((prev) => [created, ...prev])
  }

  async function handleUpdate(id, input) {
    try {
      const updated = await api.updateTask(id, input)
      setTasks((prev) =>
        prev
          .map((t) => (t.id === id ? updated : t))
          // Drop tasks that no longer match the current filter
          .filter((t) => filter === 'all' || (filter === 'done') === t.done),
      )
    } catch (err) {
      setError(err.message)
      throw err
    }
  }

  async function handleDelete(id) {
    try {
      await api.deleteTask(id)
      setTasks((prev) => prev.filter((t) => t.id !== id))
    } catch (err) {
      setError(err.message)
    }
  }

  return (
    <main className="container">
      <h1>Task Manager</h1>
      <p className="subtitle">React + Gin + SQLite</p>

      <TaskForm onSubmit={handleCreate} />

      <div className="filters">
        {FILTERS.map((f) => (
          <button
            key={f}
            className={f === filter ? 'active' : ''}
            onClick={() => setFilter(f)}
          >
            {f}
          </button>
        ))}
      </div>

      {error && <div className="error">{error}</div>}

      {loading ? (
        <p className="muted">Loading…</p>
      ) : tasks.length === 0 ? (
        <p className="muted">No tasks here yet.</p>
      ) : (
        <ul className="task-list">
          {tasks.map((task) => (
            <TaskItem
              key={task.id}
              task={task}
              onUpdate={handleUpdate}
              onDelete={handleDelete}
            />
          ))}
        </ul>
      )}
    </main>
  )
}
