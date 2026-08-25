# Prompts utilizados / Prompts used

## Español — requerimiento original resumido

> Construye una calculadora responsive para una prueba técnica de Software Engineer II. Usa React y TypeScript en frontend y Go en backend. Implementa suma, resta, multiplicación, división, potencia, raíz y porcentaje como microservicios independientes. Añade un orquestador que reciba expresiones, valide sintaxis y paréntesis, respete PEMDAS y detenga la evaluación al primer error. Usa JSON con `Status`, `resultado` o `Error`, tests unitarios, documentación bilingüe y Docker Compose con una red llamada `calculator-app-network`, exponiendo únicamente el frontend. La UI debe aceptar teclado, mostrar errores mediante un toast superior de tres segundos y utilizar una paleta Solarized Dark con acentos menta/turquesa.

## Español — prompt reproducible

```text
Crea un monorepo de una calculadora distribuida con estas decisiones cerradas:
- React + TypeScript + Vite, diseño mobile-first Solarized Dark/menta, botones y teclado.
- Go con biblioteca estándar: orquestador y siete binarios independientes para +, -, *, /, ^, sqrt y percent.
- POST /api/calculate recibe {"expression":"..."}. Implementa tokenizer y parser propio, sin eval, con PEMDAS, paréntesis, signos unarios, potencia asociativa a la derecha, sqrt(x) y percent(a,b)=a*b/100.
- Cada operación interna expone POST /calculate. Binarias reciben valor1/valor2; sqrt solo valor1. Respuestas: {"Status":"ok","resultado":n} o {"Status":"ERROR","Error":"..."}.
- Valida JSON estricto, división por cero, raíz negativa, overflow/NaN/infinito, longitud y tokens. Detén la evaluación ante el primer error y usa timeouts HTTP.
- Incluye tests Go y React, TypeScript estricto, ESLint, Dockerfiles multi-stage, Nginx proxy y Compose. Solo frontend publica puerto y todos usan calculator-app-network.
- Documenta setup Docker/local, REST, arquitectura y decisiones en español e inglés.
Entrega código ejecutable, lockfile y comandos de verificación.
```

## English — reproducible prompt

```text
Create a distributed calculator monorepo using React, TypeScript, Vite, and Go. Build one orchestrator plus seven independently deployable operation services. Parse the full expression without eval, support PEMDAS, parentheses, unary signs, right-associative powers, sqrt(x), and percent(a,b)=a*b/100. Use strict JSON contracts with Status/result/Error, validate all input and numeric edge cases, stop on the first error, and set downstream HTTP timeouts. Add a responsive Solarized Dark mint UI with keyboard support and a three-second top toast. Include Go and React tests, strict TypeScript, ESLint, multi-stage Dockerfiles, an Nginx API proxy, and Docker Compose where only the frontend publishes a port and the network is named calculator-app-network. Document Docker/local setup, REST usage, architecture, and design decisions in Spanish first and English second.
```

Nota: el texto completo de la conversación puede compartirse adicionalmente mediante la función de exportación de la plataforma usada. Este archivo conserva una versión autocontenida, revisable y reproducible sin incluir metadatos personales.
