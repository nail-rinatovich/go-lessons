# Project 8: React Frontend

## Description
A React app that loads a list of books from a mock API and lets you search them and mark favorites (Lesson 9).

## How It Works
1. `public/books.json` is the mock API; `src/api.js` loads it with `fetch` (with an artificial 0.5 s delay so the loading state is visible).
2. `App.jsx` keeps state with `useState` and loads data in `useEffect` on mount.
3. Conditional rendering: "Loading…", an error message, or "Nothing found".
4. Lists: `books.map(...)` renders `BookCard` components with a `key`.
5. Events: typing in the search box (`onChange`) filters the list; the button (`onClick`) toggles a favorite.
6. Styles live in `src/App.css`.

## How to Run
```bash
cd path/to/project8-react-frontend
npm install
npm run dev
```
Open `http://localhost:3000`.

## Notes
- The lesson plan mentions `create-react-app`, but it is deprecated. This project uses Vite, the current standard (`npm create vite@latest`). JSX, components and hooks work exactly the same.
