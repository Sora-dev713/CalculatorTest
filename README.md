# MintCalc — Calculadora distribuida / Distributed calculator

## Español

MintCalc es una calculadora responsive creada con React, TypeScript y Go. Interpreta expresiones sin usar `eval`, aplica precedencia PEMDAS y delega cada operación matemática a un microservicio HTTP independiente.

### Inicio rápido con Docker

Requisitos: Docker Engine con Docker Compose v2.

```bash
cp .env.example .env
docker compose up --build
```

Abre `http://localhost:3000`. El puerto se puede cambiar con `FRONTEND_PORT`. Solo Nginx/frontend se publica al host; el orquestador y los siete servicios permanecen en la red bridge `calculator-app-network`.

Para detener la aplicación:

```bash
docker compose down
```

### Ejecución local separada

Requisitos: Go 1.24+ y Node.js 24+.

En terminales distintas, ejecuta cada servicio con su puerto:

```bash
PORT=8081 go -C backend run ./cmd/addition
PORT=8082 go -C backend run ./cmd/subtraction
PORT=8083 go -C backend run ./cmd/multiplication
PORT=8084 go -C backend run ./cmd/division
PORT=8085 go -C backend run ./cmd/power
PORT=8086 go -C backend run ./cmd/sqrt
PORT=8087 go -C backend run ./cmd/percentage
go -C backend run ./cmd/orchestrator
```

En PowerShell, usa `$env:PORT="8081"; go -C backend run ./cmd/addition` y cambia el puerto/servicio en cada terminal. El orquestador usa por defecto esas direcciones localhost; se pueden sobrescribir mediante las variables `ADDITION_URL`, `SUBTRACTION_URL`, `MULTIPLICATION_URL`, `DIVISION_URL`, `POWER_URL`, `SQRT_URL`, `PERCENTAGE_URL` y `SERVICE_TIMEOUT`.

Finalmente:

```bash
cd frontend
npm ci
npm run dev
```

Vite abre `http://localhost:5173` y envía `/api` a `localhost:8080`.

### Uso de la API

Endpoint público:

```bash
curl -X POST http://localhost:3000/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression":"sqrt(16) + 2^3"}'
```

```json
{"Status":"ok","resultado":12}
```

La gramática admite decimales, signos unarios, `+`, `-`, `*`, `/`, `^`, paréntesis, `sqrt(x)` y `percent(valor, porcentaje)`. Por ejemplo, `percent(200,10)` devuelve `20`. `%` no es un operador válido. La potencia asocia a la derecha: `2^3^2 = 512`.

Los servicios internos aceptan `POST /calculate`:

```json
{"valor1":10,"valor2":2}
```

Raíz cuadrada usa únicamente `{"valor1":16}`. Una respuesta inválida tiene la forma:

```json
{"Status":"ERROR","Error":"division by zero is not allowed"}
```

Consulta [la referencia de servicios](docs/SERVICES.md) para contratos y validaciones individuales.

Para una explicación detallada de la estructura, el parser y el flujo interno, consulta el [README del backend](backend/README.md).

### Pruebas y calidad

```bash
go -C backend test ./...
go -C backend vet ./...
go -C backend test '-coverpkg=./...' '-coverprofile=coverage.out' ./...
go -C backend tool cover '-func=coverage.out'
cd frontend
npm test
npm run test:coverage
npm run build
npm run lint
```

“Código idiomático” significa seguir las convenciones naturales del lenguaje: errores explícitos, `gofmt`, `net/http` y tablas de pruebas en Go; componentes funcionales, estado predecible, TypeScript estricto y elementos semánticos en React.

### Decisiones de diseño

- El orquestador tokeniza y construye un árbol sintáctico completo antes de llamar servicios. Esto impide ejecuciones parciales para expresiones inválidas y evita los riesgos de `eval`.
- Cada nodo se evalúa secuencialmente y el primer error detiene el cálculo. Los clientes internos tienen timeout.
- Los binarios comparten solo contratos JSON, validación HTTP y funciones puras; continúan siendo procesos y contenedores desplegables de forma independiente.
- Se usa `float64`, rechazando resultados no finitos. La UI limita el ruido visual sin alterar la respuesta de la API.
- Nginx concentra el acceso público y evita CORS o la exposición accidental de servicios internos.

Los prompts entregables se encuentran en [docs/PROMPTS.md](docs/PROMPTS.md). El cumplimiento del correo original está detallado en [docs/REQUIREMENTS.md](docs/REQUIREMENTS.md) y los resultados de cobertura en [docs/COVERAGE.md](docs/COVERAGE.md).

---

## English

MintCalc is a responsive React, TypeScript, and Go calculator. It parses expressions without `eval`, follows PEMDAS precedence, and delegates every mathematical operation to an independent HTTP microservice.

### Quick start

Requirements: Docker Engine and Docker Compose v2.

```bash
cp .env.example .env
docker compose up --build
```

Open `http://localhost:3000`. Only the Nginx/frontend container publishes a host port. The orchestrator and operation services remain inside `calculator-app-network`. Use `docker compose down` to stop the stack.

For local development, start the seven Go commands under `backend/cmd` on ports 8081–8087, then `go -C backend run ./cmd/orchestrator`. Run `npm ci && npm run dev` inside `frontend`; Vite proxies API traffic to port 8080. The exact commands and environment variables are documented in the Spanish section above.

### API and design

Send `POST /api/calculate` with `{"expression":"sqrt(16) + 2^3"}`. Successful responses use `{"Status":"ok","resultado":12}`; errors use `{"Status":"ERROR","Error":"message"}`. Supported syntax includes decimal numbers, unary signs, arithmetic operators, parentheses, `sqrt(x)`, and `percent(value, percentage)`.

The orchestrator validates and parses the complete input before any downstream request. Evaluation stops at the first error, internal calls have timeouts, and non-finite inputs/results are rejected. Shared Go packages contain only HTTP contracts, validation, and pure operation logic; every operation remains an independently deployable process.

Run backend checks with `go -C backend test ./...` and `go -C backend vet ./...`. Run frontend checks with `npm test`, `npm run build`, and `npm run lint` inside `frontend`. See the [backend guide](backend/README.md), [service reference](docs/SERVICES.md), and [shared prompts](docs/PROMPTS.md) for more detail.
