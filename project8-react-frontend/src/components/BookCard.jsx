// A functional component: receives data and callbacks through props
export default function BookCard({ book, favorite, onToggleFavorite }) {
  return (
    <li className={favorite ? 'book favorite' : 'book'}>
      <div>
        <div className="book-title">{book.title}</div>
        <div className="book-meta">
          {book.author}, {book.year}
        </div>
      </div>
      <button onClick={onToggleFavorite} aria-pressed={favorite}>
        {favorite ? '★ Favorite' : '☆ Add to favorites'}
      </button>
    </li>
  )
}
