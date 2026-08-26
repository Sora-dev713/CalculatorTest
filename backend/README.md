# Backend de MintCalc

Este documento explica cómo está organizado el backend, qué responsabilidad tiene cada carpeta y cómo una expresión enviada por el frontend termina convirtiéndose en llamadas a los microservicios matemáticos.

## Vista general

El backend está escrito en Go y contiene ocho aplicaciones ejecutables:

- Un **orquestador**, que recibe y valida expresiones completas como `2 + 3 * 4`.
- Siete **servicios de operaciones**: suma, resta, multiplicación, división, potencia, raíz cuadrada y porcentaje.

El orquestador no realiza directamente las operaciones binarias o funciones matemáticas. Primero convierte el texto en un árbol de expresión y luego delega cada cálculo al microservicio correspondiente mediante HTTP.

```text
Frontend
   │ POST /api/calculate
   ▼
Orquestador
   ├── Tokenizer y parser
   ├── Árbol de expresión
   └── Evaluador HTTP
          ├── addition
          ├── subtraction
          ├── multiplication
          ├── division
          ├── power
          ├── sqrt
          └── percentage
```

## Estructura de carpetas

```text
backend/
├── cmd/
│   ├── orchestrator/       Punto de entrada del orquestador
│   ├── addition/           Ejecutable y tests del servicio de suma
│   ├── subtraction/        Ejecutable del servicio de resta
│   ├── multiplication/     Ejecutable del servicio de multiplicación
│   ├── division/           Ejecutable del servicio de división
│   ├── power/              Ejecutable del servicio de potencia
│   ├── sqrt/               Ejecutable del servicio de raíz cuadrada
│   └── percentage/         Ejecutable del servicio de porcentaje
├── internal/
│   ├── api/                Contratos JSON y utilidades HTTP compartidas
│   ├── expression/         Tokenizer, parser y árbol de sintaxis
│   ├── operation/          Handlers y funciones matemáticas puras
│   ├── orchestrator/       Endpoint público y evaluación distribuida
│   └── server/             Configuración HTTP común de los servicios
├── Dockerfile              Imagen parametrizable para los ocho binarios
└── go.mod                  Módulo Go del backend
```

En Go, una carpeta llamada `internal` solo puede importarse desde el módulo que la contiene. Esto evita que los detalles internos se conviertan accidentalmente en una API pública.

### `cmd`: aplicaciones ejecutables

Cada subcarpeta de `cmd` contiene un `package main` y representa un proceso independiente. Los archivos son intencionalmente pequeños: seleccionan el handler correspondiente y levantan el servidor.

Por ejemplo, `cmd/addition/main.go` conecta `operation.Add` con un handler binario. Aunque las operaciones reutilizan infraestructura, cada comando se compila como un binario y se ejecuta en su propio contenedor.

### `internal/api`: contratos HTTP

Define los tipos compartidos:

- `BinaryRequest`: `valor1` y `valor2`.
- `UnaryRequest`: únicamente `valor1`.
- `Response`: `Status`, `resultado` o `Error`.

También contiene:

- Decodificación JSON estricta mediante `DisallowUnknownFields`.
- Límite para el body de la solicitud.
- Rechazo de múltiples objetos JSON en un mismo body.
- Validación de números finitos.
- Funciones comunes para escribir respuestas JSON.

Los campos numéricos del request son punteros para distinguir correctamente entre un campo ausente y un valor válido igual a cero.

### `internal/expression`: tokenizer y parser

El orquestador no usa `eval`. La expresión pasa por dos etapas:

1. El **tokenizer** recorre el texto y genera tokens: números, operadores, nombres de funciones, paréntesis y comas.
2. El **parser** consume esos tokens y construye un árbol de sintaxis abstracta o AST.

Los nodos del árbol pueden ser:

- `Number`: número literal.
- `Unary`: signo positivo o negativo.
- `Binary`: suma, resta, multiplicación, división o potencia.
- `Function`: `sqrt` o `percent`.

La estructura de funciones del parser representa la precedencia:

```text
expression
└── addition          + y -
    └── multiplication  * y /
        └── unary       +x y -x
            └── power   ^, asociativa a la derecha
                └── primary  números, paréntesis y funciones
```

Por eso `2 + 3 * 4` se interpreta como `2 + (3 * 4)`, mientras que `2^3^2` se interpreta como `2^(3^2)`.

El parser valida toda la expresión antes de realizar una llamada HTTP. Texto desconocido, paréntesis desbalanceados, números inválidos, argumentos faltantes o caracteres como `%` detienen inmediatamente la solicitud.

### `internal/orchestrator`: coordinación

Este paquete tiene dos responsabilidades separadas:

- `handler.go` implementa `POST /api/calculate`, valida el JSON, invoca el parser y devuelve la respuesta pública.
- `evaluator.go` recorre recursivamente el árbol y llama al servicio asociado a cada nodo.

El mapa `URLs` desacopla el evaluador de direcciones concretas. En Docker contiene nombres como `http://addition:8080`; durante tests contiene URLs de servidores HTTP temporales.

El cliente HTTP:

- Propaga el contexto de la petición.
- Aplica el timeout configurado en `SERVICE_TIMEOUT`.
- Limita el tamaño de las respuestas.
- Valida status HTTP, estructura JSON y resultados finitos.
- Detiene la evaluación en el primer error.

### `internal/operation`: operaciones y handlers

Las funciones matemáticas son puras y pequeñas, por ejemplo `Add`, `Divide` y `SquareRoot`. Esto permite probarlas sin levantar un servidor.

`BinaryHandler` y `UnaryHandler` adaptan esas funciones a HTTP. Ambos validan método, JSON, campos obligatorios, números y resultado. Las validaciones específicas se mantienen junto a cada operación:

- División rechaza divisor cero.
- Raíz cuadrada rechaza valores negativos.
- Potencia rechaza dominios que produzcan `NaN` o infinito.
- Porcentaje calcula `valor1 * valor2 / 100` y permite porcentajes negativos.

### `internal/server`: servidor común

Configura `POST /calculate`, `GET /health`, puerto y timeouts HTTP. Centralizar esto mantiene consistentes los siete servicios y evita repetir configuración de infraestructura en cada `main`.

## Flujo de un cálculo

Para la expresión `2 + 3 * 4` ocurre lo siguiente:

1. El frontend envía `{"expression":"2 + 3 * 4"}` al orquestador.
2. El parser construye un árbol cuya raíz es suma y cuyo lado derecho es multiplicación.
3. El evaluador resuelve primero `3 * 4` llamando a `multiplication` con `{"valor1":3,"valor2":4}`.
4. `multiplication` responde `{"Status":"ok","resultado":12}`.
5. El evaluador llama a `addition` con `{"valor1":2,"valor2":12}`.
6. `addition` responde con `14` y el orquestador entrega `{"Status":"ok","resultado":14}` al frontend.

Si la multiplicación falla, el paso 5 nunca se ejecuta.

## Contratos HTTP

### API pública del orquestador

```http
POST /api/calculate
Content-Type: application/json

{"expression":"sqrt(16) + 2^3"}
```

Respuesta exitosa:

```json
{"Status":"ok","resultado":12}
```

### Servicios binarios

```http
POST /calculate
Content-Type: application/json

{"valor1":10,"valor2":2}
```

### Raíz cuadrada

```http
POST /calculate
Content-Type: application/json

{"valor1":16}
```

Respuesta de error:

```json
{"Status":"ERROR","Error":"division by zero is not allowed"}
```

Los servicios de operaciones retornan `400` ante datos o dominios inválidos. El orquestador retorna `400` para una expresión inválida y `502` cuando un servicio dependiente falla, no responde o rechaza una operación.

## Configuración

| Variable | Predeterminado | Uso |
|---|---:|---|
| `PORT` | `8080` | Puerto del proceso actual |
| `SERVICE_TIMEOUT` | `3s` | Timeout del cliente del orquestador |
| `ADDITION_URL` | `http://localhost:8081` | Servicio de suma |
| `SUBTRACTION_URL` | `http://localhost:8082` | Servicio de resta |
| `MULTIPLICATION_URL` | `http://localhost:8083` | Servicio de multiplicación |
| `DIVISION_URL` | `http://localhost:8084` | Servicio de división |
| `POWER_URL` | `http://localhost:8085` | Servicio de potencia |
| `SQRT_URL` | `http://localhost:8086` | Servicio de raíz |
| `PERCENTAGE_URL` | `http://localhost:8087` | Servicio de porcentaje |

Docker Compose reemplaza las URLs localhost por los nombres DNS internos de los contenedores.

## Ejecución y pruebas

Desde la raíz del repositorio:

```bash
go -C backend test ./...
go -C backend vet ./...
```

Para ejecutar un servicio individual:

```bash
PORT=8081 go -C backend run ./cmd/addition
```

En PowerShell:

```powershell
$env:PORT="8081"
go -C backend run ./cmd/addition
```

Los tests están agrupados junto al servicio o componente que verifican. Los archivos `*_test.go` no forman parte de los binarios ni se copian a las imágenes Docker:

- `expression/parser_test.go`: gramática, precedencia y expresiones inválidas.
- `cmd/<servicio>/main_test.go`: comportamiento y casos de dominio propios de cada microservicio.
- `operation/handler_test.go`: validaciones comunes del adaptador HTTP reutilizado por los servicios.
- `orchestrator/handler_test.go`: validación previa, propagación de errores, timeout y evaluación integrada usando handlers reales.

## English summary

The backend contains one orchestrator and seven independently runnable Go operation services. `cmd` contains executable entry points; `internal/api` owns strict JSON contracts; `internal/expression` tokenizes and parses expressions into an AST; `internal/orchestrator` evaluates that AST through HTTP; `internal/operation` contains pure math functions and HTTP adapters; and `internal/server` provides common server configuration.

The complete expression is validated before any downstream request. Each AST operation is delegated to its corresponding service, and evaluation stops at the first failure. Run `go -C backend test ./...` and `go -C backend vet ./...` from the repository root.
