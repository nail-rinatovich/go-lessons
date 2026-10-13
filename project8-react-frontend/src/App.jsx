import { useEffect, useState } from 'react'
import { fetchBooks } from './api.js'
import BookCard from './components/BookCard.jsx'

export default function App() {
  const [books, setBooks] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [search, setSearch] = useState('')
  const [favorites, setFavorites] = useState([])

  // Load the books once, when the component mounts
  useEffect(() => {
    fetchBooks()
      .then(setBooks)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  function toggleFavorite(id) {
    setFavorites((prev) =>
      prev.includes(id) ? prev.filter((f) => f !== id) : [...prev, id],
    )
  }

  // Derived data: computed on every render, not stored in state
  const query = search.trim().toLowerCase()
  const visible = books.filter(
    (b) => b.title.toLowerCase().includes(query) || b.author.toLowerCase().includes(query),
  )

  return (
    <main className="container">
      <h1>Book Shelf</h1>
      <p className="subtitle">Favorites: {favorites.length}</p>

      <input
        className="search"
        placeholder="Search by title or author"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
      />

      {loading && <p className="muted">Loading…</p>}
      {error && <p className="error">{error}</p>}
      {!loading && !error && visible.length === 0 && (
        <p className="muted">Nothing found.</p>
      )}

      <ul className="book-list">
        {visible.map((book) => (
          <BookCard
            key={book.id}
            book={book}
            favorite={favorites.includes(book.id)}
            onToggleFavorite={() => toggleFavorite(book.id)}
          />
        ))}
      </ul>
    </main>
  )
}
