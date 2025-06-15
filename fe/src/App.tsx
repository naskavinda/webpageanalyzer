import { useState } from 'react'
import './App.css'

function App() {
  const [url, setUrl] = useState('')
  const [result, setResult] = useState<null | {
    URL: string;
    HTMLVersion: string;
    Title: string;
    HeadingCounts: { h1: number; h2: number; h3: number };
    InternalLinks: number;
    ExternalLinks: number;
    InaccessibleLinks: number;
    HasLoginForm: boolean;
  }>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError(null)
    setResult(null)
    try {
      const response = await fetch('http://localhost:8080/analyzer', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ webpageUrl: url }),
      })
      
      const data = await response.json()
      if (!response.ok) throw new Error(data.error || 'Failed to analyze webpage')
      setResult(data.content)
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message)
      } else {
        setError('An unknown error occurred')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="App">
      <h1>Webpage Analyzer</h1>
      <form onSubmit={handleSubmit} style={{ marginBottom: 24 }}>
        <input
          type="text"
          placeholder="Enter webpage URL"
          value={url}
          onChange={e => setUrl(e.target.value)}
          style={{ width: '350px', marginRight: 8 }}
        />
        <button type="submit" disabled={loading || !url}>
          {loading ? 'Analyzing...' : 'Analyze'}
        </button>
      </form>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      {result && (
        <div style={{ marginTop: 24, textAlign: 'left', background: '#f8f8f8', padding: 16, borderRadius: 8 }}>
          <h2>Analysis Result</h2>
          <p><strong>URL:</strong> {result.URL}</p>
          <p><strong>HTML Version:</strong> {result.HTMLVersion}</p>
          <p><strong>Title:</strong> {result.Title}</p>
          <p><strong>Headings:</strong> h1: {result.HeadingCounts.h1}, h2: {result.HeadingCounts.h2}, h3: {result.HeadingCounts.h3}</p>
          <p><strong>Internal Links:</strong> {result.InternalLinks}</p>
          <p><strong>External Links:</strong> {result.ExternalLinks}</p>
          <p><strong>Inaccessible Links:</strong> {result.InaccessibleLinks}</p>
          <p><strong>Has Login Form:</strong> {result.HasLoginForm ? 'Yes' : 'No'}</p>
        </div>
      )}
    </div>
  )
}

export default App
