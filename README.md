# 🌤️ Weather Agent

Un agente inteligente de clima desarrollado en Go que utiliza OpenAI GPT para analizar condiciones meteorológicas y proporcionar recomendaciones prácticas con evaluación de riesgo.

- Análisis detallado del clima
- Recomendaciones prácticas basadas en las condiciones
- Evaluación de nivel de riesgo (bajo, medio, alto)

## 🚀 Requisitos

- API Key de [OpenWeatherMap](https://openweathermap.org/api)
- API Key de [OpenAI](https://platform.openai.com/)

## 🔧 Configuración

Configura las siguientes variables de entorno:

```bash
export OPENWEATHER_API_KEY="tu_api_key_de_openweather"
export OPENAI_API_KEY="tu_api_key_de_openai"
export PORT="8080"  # Opcional, por defecto usa 8080
```

O crea un archivo `.env`

## 🏃 Ejecución

Ejecuta el servidor:

```bash
go run main.go
```

Verás el mensaje:
```
🤖 Weather Agent corriendo en puerto 8080
```

## 📡 Uso del API

### Endpoint

```
GET http://localhost:8080/agent/weather?city=<nombre_ciudad>
```

### Parámetros

- `city` (opcional): Nombre de la ciudad. Si no se proporciona, usa "Managua,NI" por defecto.

### Ejemplo de uso

Consulta el clima de León:

```bash
curl "http://localhost:8080/agent/weather?city=Leon"
```

O desde tu navegador:
```
http://localhost:8080/agent/weather?city=Leon
```

### Respuesta

El endpoint retorna un JSON con la siguiente estructura:

```json
{
  "city": "Leon",
  "temp": 28.5,
  "condition": "cielo despejado",
  "analysis": "Análisis detallado del clima...",
  "recommendation": "Recomendación práctica...",
  "risk_level": "bajo"
}
```

## 📸 Ejemplos

### Consulta para León

![Consulta del clima en León](demo/leon.png)

### Consulta para Managua

![Consulta del clima en Managua](demo/managua.png)