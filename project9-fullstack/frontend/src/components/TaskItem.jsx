import { useState } from 'react'
import { validate } from './TaskForm.jsx'

export default function TaskItem({ task, onUpdate, onDelete }) {
  const [editing, setEditing] = useState(false)
  const [title, setTitle] = useState(task.title)
  const [description, setDescription] = useState(task.description)
  const [errors, setErrors] = useState({})
  const [busy, setBusy] = useState(false)

  async function run(action) {
    setBusy(true)
    try {
      await action()
    } catch {
      // The parent component already shows the error
    } finally {
      setBusy(false)
    }
  }

  function toggleDone() {
    run(() => onUpdate(task.id, { title: task.title, description: task.description, done: !task.done }))
  }

  function startEdit() {
    setTitle(task.title)
    setDescription(task.description)
    setErrors({})
    setEditing(true)
  }

  function save(e) {
    e.preventDefault()
    const found = validate({ title, description })
    setErrors(found)
    if (Object.keys(found).length > 0) return
    run(async () => {
      await onUpdate(task.id, { title: title.trim(), description: description.trim(), done: task.done })
      setEditing(false)
    })
  }

  if (editing) {
    return (
      <li className="task editing">
        <form onSubmit={save} noValidate>
          <input value={title} onChange={(e) => setTitle(e.target.value)} autoFocus />
          {errors.title && <span className="field-error">{errors.title}</span>}
          <textarea rows={2} value={description} onChange={(e) => setDescription(e.target.value)} />
          {errors.description && <span className="field-error">{errors.description}</span>}
          <div className="actions">
            <button type="submit" disabled={busy}>Save</button>
            <button type="button" className="secondary" onClick={() => setEditing(false)}>Cancel</button>
          </div>
        </form>
      </li>
    )
  }

  return (
    <li className={`task ${task.done ? 'done' : ''}`}>
      <input type="checkbox" checked={task.done} onChange={toggleDone} disabled={busy} />
      <div className="task-body">
        <div className="task-title">{task.title}</div>
        {task.description && <div className="task-desc">{task.description}</div>}
        <div className="task-date">{new Date(task.createdAt).toLocaleString()}</div>
      </div>
      <div className="actions">
        <button className="secondary" onClick={startEdit} disabled={busy}>Edit</button>
        <button className="danger" onClick={() => run(() => onDelete(task.id))} disabled={busy}>Delete</button>
      </div>
    </li>
  )
}
