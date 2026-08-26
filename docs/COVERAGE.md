# Reporte de cobertura / Coverage report

Reporte generado el 25 de agosto de 2026. Los porcentajes deben regenerarse después de cada cambio relevante.

## Backend

```bash
go -C backend test '-coverpkg=./...' '-coverprofile=coverage.out' ./...
go -C backend tool cover '-func=coverage.out'
```

| Métrica | Cobertura |
|---|---:|
| Statements | 75.7% |

El perfil detallado queda temporalmente en `backend/coverage.out` y no se versiona. Los `main` y el arranque bloqueante de servidores se verifican mediante build/vet, no invocándolos directamente desde tests unitarios.

## Frontend

```bash
cd frontend
npm run test:coverage
```

| Métrica | Cubierto | Porcentaje |
|---|---:|---:|
| Statements | 97 / 126 | 76.98% |
| Branches | 52 / 82 | 63.41% |
| Functions | 23 / 25 | 92.00% |
| Lines | 78 / 84 | 92.85% |

Vitest genera el reporte navegable en `frontend/coverage/`, también excluido de Git. Las pruebas cubren entrada por teclado y botones, edición en el cursor, editor expandido, resultados, errores del API y rechazo de expresiones vacías.

## English

This report was generated on August 25, 2026. Backend statement coverage is 75.7%. Frontend coverage is 76.98% statements, 63.41% branches, 92.00% functions, and 92.85% lines. Use the commands above to regenerate the reports after meaningful changes.
