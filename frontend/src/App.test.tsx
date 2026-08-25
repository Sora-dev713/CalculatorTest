import { cleanup, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import App from './App'

afterEach(() => { cleanup(); vi.restoreAllMocks() })

describe('MintCalc', () => {
  it('builds an expression with buttons and clears it', async () => {
    const user = userEvent.setup(); render(<App />)
    await user.click(screen.getByRole('button', { name: '7' })); await user.click(screen.getByRole('button', { name: 'Sumar' })); await user.click(screen.getByRole('button', { name: '3' }))
    expect(screen.getByLabelText('Expresión')).toHaveValue('7+3')
    await user.click(screen.getByRole('button', { name: 'Limpiar' })); expect(screen.getByLabelText('Expresión')).toHaveValue('')
  })
  it('submits from the keyboard and shows the result', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({ Status: 'ok', resultado: 12 }), { status: 200 }))
    const user = userEvent.setup(); render(<App />); await user.type(screen.getByLabelText('Expresión'), 'sqrt(16)+2^3{Enter}')
    await waitFor(() => expect(screen.getByLabelText('Resultado')).toHaveTextContent('12'))
    expect(fetch).toHaveBeenCalledWith('/api/calculate', expect.objectContaining({ method: 'POST', body: JSON.stringify({ expression: 'sqrt(16)+2^3' }) }))
  })
  it('shows API errors in an accessible toast', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({ Status: 'ERROR', Error: 'division by zero is not allowed' }), { status: 400 }))
    const user = userEvent.setup(); render(<App />); await user.type(screen.getByLabelText('Expresión'), '1/0'); await user.click(screen.getByRole('button', { name: 'Calcular' }))
    expect(await screen.findByRole('alert')).toHaveTextContent('division by zero')
  })
  it('rejects an empty calculation without a request', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch'); const user = userEvent.setup(); render(<App />); await user.click(screen.getByRole('button', { name: 'Calcular' }))
    expect(screen.getByRole('alert')).toBeInTheDocument(); expect(fetchMock).not.toHaveBeenCalled()
  })
})
