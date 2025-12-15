# 🌤️ Weather Agent

Un agente inteligente de clima desarrollado en Go que utiliza OpenAI GPT para analizar condiciones meteorológicas y proporcionar recomendaciones prácticas con evaluación de riesgo.

- Análisis detallado del clima
- Recomendaciones prácticas basadas en las condiciones
- Evaluación de nivel de riesgo (bajo, medio, alto)
- Arquitectura limpia con inyección de dependencias
- Manejo robusto de errores y timeouts
- Validación de entrada y sanitización

## Requisitos

- Go 1.23 o superior
- API Key de [OpenWeatherMap](https://openweathermap.org/api)
- API Key de [OpenAI](https://platform.openai.com/)

## Configuración

Configura las siguientes variables de entorno:

```bash
export OPENWEATHER_API_KEY="tu_api_key_de_openweather"
export OPENAI_API_KEY="tu_api_key_de_openai"
export PORT="8080"  # Opcional, por defecto usa 8080
export LOG_LEVEL="info"  # Opcional, por defecto usa info (debug, info, warn, error)
```

O crea un archivo `.env`

## Compilación y Ejecución

```bash
# Compilar desde la raíz del proyecto
go build -o weather-agent ./cmd/server

# Ejecutar
./weather-agent
```

O ejecutar directamente:

```bash
go run ./cmd/server/main.go
```

O desde el directorio `cmd/server`:

```bash
cd cmd/server
go run main.go
```

## Uso del API

### Endpoints

#### Weather Agent
```
GET http://localhost:8080/agent/weather?city=<nombre_ciudad>
```

#### Health Check
```
GET http://localhost:8080/health
```

### Parámetros

- `city` (opcional): Nombre de la ciudad. Si no se proporciona, usa "Managua,NI" por defecto.
  - Máximo 100 caracteres
  - Se sanitiza automáticamente

### Ejemplo de uso

Consulta el clima de León:

```bash
curl "http://localhost:8080/agent/weather?city=Leon"
```

O desde tu navegador:
```
http://localhost:8080/agent/weather?city=Leon
```

### Respuesta de ejemplo

```json
{
  "analysis": "El clima está cálido y soleado...",
  "recommendation": "Usa protector solar y mantente hidratado...",
  "risk_level": "medio",
  "city": "Leon",
  "temp": 32.5,
  "condition": "cielo claro"
}
```

## Arquitectura

El proyecto sigue los principios de Clean Architecture y mejores prácticas de Go:

- **cmd/**: Puntos de entrada de la aplicación
- **internal/**: Lógica de aplicación (no expuesta externamente)
  - **agent/**: Lógica del agente de clima
  - **api/**: Handlers HTTP
  - **config/**: Configuración y carga de variables de entorno
  - **llm/**: Integración con OpenAI
  - **weather/**: Integración con OpenWeatherMap
- **pkg/**: Utilidades compartidas (si aplica)

### Características implementadas

- ✅ Context propagation para cancelación y timeouts
- ✅ Inyección de dependencias mediante interfaces
- ✅ Timeouts en todas las llamadas HTTP
- ✅ Validación de códigos de estado HTTP
- ✅ Manejo de errores con contexto (error wrapping)
- ✅ Validación y sanitización de entrada
- ✅ Graceful shutdown del servidor
- ✅ Documentación GoDoc en funciones públicas
- ✅ Configuración centralizada
- ✅ Uso del paquete oficial de OpenAI ([openai-go](https://github.com/openai/openai-go)) con retries automáticos
- ✅ Logging estructurado con [logrus](https://github.com/sirupsen/logrus) en formato JSON para producción

## Ejemplos

### Consulta para León

![Consulta del clima en León](demo/leon.png)

### Consulta para Managua

![Consulta del clima en Managua](demo/managua.png)
