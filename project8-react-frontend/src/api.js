// Mock API: books.json lives in public/ and is served by the dev server.
// The artificial delay makes the loading state visible.
export async function fetchBooks() {
  await new Promise((resolve) => setTimeout(resolve, 500))
  const res = await fetch('/books.json')
  if (!res.ok) {
    throw new Error(`Failed to load books (status ${res.status})`)
  }
  return res.json()
}
