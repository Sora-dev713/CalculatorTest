# Referencia de servicios / Service reference

## Español

Todos los servicios escuchan en `PORT` (8080 por defecto), ofrecen `GET /health` y aceptan únicamente `POST /calculate`. Rechazan JSON malformado, campos desconocidos, valores ausentes y resultados no finitos. Los tests compartidos en `backend/internal/operation/handler_test.go` ejecutan el contrato y los casos de dominio de cada operación.

| Servicio | Body | Cálculo y validación |
|---|---|---|
| addition | `{"valor1":a,"valor2":b}` | `a + b`; resultado finito |
| subtraction | `{"valor1":a,"valor2":b}` | `a - b`; resultado finito |
| multiplication | `{"valor1":a,"valor2":b}` | `a * b`; resultado finito |
| division | `{"valor1":a,"valor2":b}` | `a / b`; rechaza divisor cero |
| power | `{"valor1":a,"valor2":b}` | `a^b`; rechaza dominios y resultados no finitos |
| sqrt | `{"valor1":a}` | `sqrt(a)`; rechaza números negativos y `valor2` |
| percentage | `{"valor1":a,"valor2":b}` | `a * b / 100`; permite porcentajes negativos |

El orquestador expone `POST /api/calculate`, valida un máximo de 512 caracteres y llama los servicios configurados mediante `*_URL`. Un error de sintaxis retorna 400. Un error o indisponibilidad aguas abajo retorna 502 y detiene la evaluación.

## English

Every operation service listens on `PORT` (8080 by default), provides `GET /health`, and only accepts `POST /calculate`. Malformed JSON, unknown or missing fields, and non-finite results are rejected. Shared table-driven tests exercise every operation's contract and domain rules.

The table above is also the complete English contract: binary services receive `valor1` and `valor2`; sqrt receives only `valor1`. The orchestrator exposes `POST /api/calculate`, limits input to 512 characters, and resolves service locations through `*_URL` variables. Syntax errors return 400; downstream failures return 502 and immediately stop evaluation.
