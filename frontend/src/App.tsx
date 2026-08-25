import { useCallback, useEffect, useRef, useState } from 'react'

type ApiResponse = { Status: 'ok'; resultado: number } | { Status: 'ERROR'; Error: string }
type Key = { label: string; value: string; kind?: 'action' | 'operator' | 'equals'; wide?: boolean; aria?: string }

const keys: Key[] = [
  { label: 'AC', value: 'clear', kind: 'action', aria: 'Limpiar' }, { label: '⌫', value: 'backspace', kind: 'action', aria: 'Borrar último carácter' },
  { label: '(', value: '(' }, { label: ')', value: ')' }, { label: '÷', value: '/', kind: 'operator', aria: 'Dividir' },
  { label: '7', value: '7' }, { label: '8', value: '8' }, { label: '9', value: '9' }, { label: '×', value: '*', kind: 'operator', aria: 'Multiplicar' },
  { label: '√', value: 'sqrt(', kind: 'operator', aria: 'Raíz cuadrada' }, { label: '4', value: '4' }, { label: '5', value: '5' }, { label: '6', value: '6' },
  { label: '−', value: '-', kind: 'operator', aria: 'Restar' }, { label: 'xʸ', value: '^', kind: 'operator', aria: 'Potencia' },
  { label: '1', value: '1' }, { label: '2', value: '2' }, { label: '3', value: '3' }, { label: '+', value: '+', kind: 'operator', aria: 'Sumar' },
  { label: '% de', value: 'percent(', kind: 'operator', aria: 'Porcentaje' }, { label: '0', value: '0', wide: true }, { label: '.', value: '.' },
  { label: ',', value: ',' }, { label: '=', value: 'equals', kind: 'equals', aria: 'Calcular' },
]

function formatResult(value: number): string {
  if (!Number.isFinite(value)) return String(value)
  return Number.parseFloat(value.toPrecision(15)).toString()
}

export default function App() {
  const [expression, setExpression] = useState('')
  const [result, setResult] = useState('0')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const toastTimer = useRef<number | undefined>(undefined)

  const showError = useCallback((message: string) => {
    setError(message)
    window.clearTimeout(toastTimer.current)
    toastTimer.current = window.setTimeout(() => setError(''), 3000)
  }, [])

  useEffect(() => () => window.clearTimeout(toastTimer.current), [])

  const calculate = useCallback(async () => {
    if (loading) return
    if (!expression.trim()) { showError('Ingresa una expresión para calcular.'); return }
    setLoading(true)
    try {
      const response = await fetch('/api/calculate', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ expression }) })
      const data = await response.json() as ApiResponse
      if (!response.ok || data.Status === 'ERROR') throw new Error(data.Status === 'ERROR' ? data.Error : 'No se pudo completar el cálculo.')
      setResult(formatResult(data.resultado))
    } catch (caught) {
      showError(caught instanceof Error ? caught.message : 'No se pudo conectar con la calculadora.')
    } finally { setLoading(false) }
  }, [expression, loading, showError])

  const press = useCallback((value: string) => {
    if (value === 'clear') { setExpression(''); setResult('0'); return }
    if (value === 'backspace') { setExpression(current => current.slice(0, -1)); return }
    if (value === 'equals') { void calculate(); return }
    setExpression(current => (current + value).slice(0, 512))
  }, [calculate])

  useEffect(() => {
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.ctrlKey || event.metaKey || event.altKey) return
      if (event.key === 'Enter' || event.key === '=') { event.preventDefault(); void calculate(); return }
      if (event.key === 'Backspace') { event.preventDefault(); press('backspace'); return }
      if (event.key === 'Escape') { event.preventDefault(); press('clear'); return }
      if (/^[0-9+\-*/^().,]$/.test(event.key)) { event.preventDefault(); press(event.key); return }
      if (/^[a-zA-Z]$/.test(event.key)) { event.preventDefault(); press(event.key.toLowerCase()) }
    }
    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [calculate, press])

  return <main className="app-shell">
    {error && <div className="toast" role="alert" aria-live="assertive"><span aria-hidden="true">!</span>{error}</div>}
    <section className="calculator" aria-label="Calculadora">
      <header className="brand"><span className="brand-mark">M</span><div><h1>MintCalc</h1><p>Distributed calculator</p></div></header>
      <div className="display" aria-live="polite">
        <label htmlFor="expression">Expresión</label>
        <input id="expression" value={expression} onChange={event => setExpression(event.target.value.slice(0, 512))} placeholder="sqrt(16) + 2^3" autoComplete="off" spellCheck={false} />
        <div className="result-row"><span>Resultado</span><output aria-label="Resultado">{loading ? 'Calculando…' : result}</output></div>
      </div>
      <div className="keypad">
        {keys.map((key, index) => <button key={`${key.value}-${index}`} type="button" className={`key ${key.kind ?? ''} ${key.wide ? 'wide' : ''}`} onClick={() => press(key.value)} disabled={loading && key.value === 'equals'} aria-label={key.aria}>{key.label}</button>)}
      </div>
      <aside className="help"><strong>Sintaxis avanzada</strong><span><code>sqrt(16)</code> raíz · <code>2^8</code> potencia · <code>percent(200,10)</code> porcentaje</span></aside>
    </section>
    <footer>Usa los botones o escribe con tu teclado.</footer>
  </main>
}
