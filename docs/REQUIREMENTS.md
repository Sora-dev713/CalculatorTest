# Matriz de cumplimiento / Requirements checklist

Esta matriz contrasta la implementación con el correo original de la prueba técnica.

| Requerimiento | Estado | Evidencia |
|---|---|---|
| React frontend y backend por API | Cumple | React/TypeScript en `frontend`; orquestador y servicios Go en `backend` |
| Suma, resta, multiplicación y división | Cumple | Servicios independientes y tests en `backend/cmd/<operación>` |
| Potencia, raíz y porcentaje opcionales | Cumple | Servicios `power`, `sqrt` y `percentage` |
| UI intuitiva para input y resultado | Cumple | Teclado visual, editor directo/expandido, resultado y ayuda de sintaxis |
| Validación y manejo de errores frontend | Cumple | Expresión vacía, límite de 512 caracteres, loading y toast accesible |
| Diseño responsive | Cumple | CSS mobile-first, media queries y controles táctiles |
| Endpoints REST para operaciones | Cumple | `POST /api/calculate` y `POST /calculate` por servicio |
| Casos borde e input inválido | Cumple | JSON estricto, división por cero, raíz negativa, dominio de potencia y finitud |
| Resultados JSON | Cumple | Contrato uniforme `Status`, `resultado` o `Error` |
| Código limpio, legible e idiomático | Cumple | `gofmt`, `go vet`, errores explícitos, handlers compuestos, TypeScript estricto, ESLint y comentarios selectivos |
| Tests unitarios frontend/backend | Cumple | Vitest/Testing Library y tests Go agrupados por componente/servicio |
| Coverage report | Cumple | [Reporte reproducible](COVERAGE.md) |
| Setup, API y decisiones de diseño | Cumple | README principal, README de backend y referencia de servicios |
| Dockerfile full-stack opcional | Cumple mediante Compose | Dockerfiles multi-stage coordinados por Compose, apropiado para varios procesos |
| Repositorio Git con ambos componentes | Cumple localmente | Monorepo con frontend, backend, documentación y configuración Docker |
| Compartir prompts | Cumple | [Prompts utilizados](PROMPTS.md) |
| Publicar y compartir enlace del repositorio | Pendiente del candidato | Requiere hacer push y enviar la URL a Sezzle |

## Observaciones

- El correo habla de “a backend microservice”; esta solución va más allá y separa cada operación debido al alcance aclarado posteriormente.
- El único entregable externo pendiente es publicar el estado final y compartir el enlace del repositorio.
- Antes de enviar, conviene repetir todos los comandos de calidad y revisar manualmente desktop y mobile.

## English summary

All code-based functional and non-functional requirements are implemented, including optional operations, Docker deployment, shared prompts, and a reproducible coverage report. Publishing the final repository and sharing its URL remains an external candidate-owned deliverable.
