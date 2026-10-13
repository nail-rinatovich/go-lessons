import { useState } from 'react'

// Client-side validation mirrors the backend rules in handlers.go
export function validate({ title, description }) {
  const errors = {}
  if (!title.trim()) errors.title = 'Title is required'
  else if (title.trim().length > 100) errors.title = 'Title must be at most 100 characters'
  if (description.trim().length > 500) errors.description = 'Description must be at most 500 characters'
  return errors
}

export default function TaskForm({ onSubmit }) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [errors, setErrors] = useState({})
  const [submitting, setSubmitting] = useState(false)

  async function handleSubmit(e) {
    e.preventDefault()
    const input = { title, description }
    const found = validate(input)
    setErrors(found)
    if (Object.keys(found).length > 0) return

    setSubmitting(true)
    try {
      await onSubmit({ title: title.trim(), description: description.trim() })
      // Clear only what was submitted: the user may have typed a new task
      // while the request was in flight
      setTitle((current) => (current === title ? '' : current))
      setDescription((current) => (current === description ? '' : current))
    } catch (err) {
      // Show server-side validation errors too
      setErrors({ form: err.message })
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form className="task-form" onSubmit={handleSubmit} noValidate>
      <input
        placeholder="What needs to be done?"
        value={title}
        onChange={(e) => { setTitle(e.target.value); setErrors({}) }}
        aria-invalid={!!errors.title}
      />
      {errors.title && <span className="field-error">{errors.title}</span>}

      <textarea
        placeholder="Description (optional)"
        rows={2}
        value={description}
        onChange={(e) => { setDescription(e.target.value); setErrors({}) }}
        aria-invalid={!!errors.description}
      />
      {errors.description && <span className="field-error">{errors.description}</span>}
      {errors.form && <span className="field-error">{errors.form}</span>}

      <button type="submit" disabled={submitting}>
        {submitting ? 'Adding…' : 'Add task'}
      </button>
    </form>
  )
}
