# Documentación del Servicio de Órdenes (@order-service)

## Descripción General

Este proyecto es un microservicio de gestión de órdenes desarrollado en Go, utilizando el framework Echo para la creación de APIs RESTful. El servicio proporciona funcionalidades para la creación y consulta de órdenes, interactuando con otros microservicios como el servicio de productos y el servicio de usuarios.

## Estructura del Proyecto

El proyecto sigue una estructura modular típica de aplicaciones Go:

- `config/`: Contiene la configuración de la aplicación.
- `db/`: Maneja la conexión y operaciones con la base de datos.
- `handlers/`: Contiene los manejadores de las rutas HTTP.
- `middleware/`: Incluye middlewares personalizados.
- `models/`: Define las estructuras de datos utilizadas en la aplicación.
- `main.go`: Punto de entrada de la aplicación.

## Configuración

La configuración se maneja a través de variables de entorno y un archivo `.env` para entornos de desarrollo. Las principales configuraciones incluyen:

- `DB_CONNECTION_STRING`: Cadena de conexión a la base de datos PostgreSQL.
- `APP_ENV`: Entorno de la aplicación (development, production, etc.).
- `APP_USER_SERVICE_URL`: URL del servicio de usuarios.
- `APP_PRODUCT_SERVICE_URL`: URL del servicio de productos.

## Rutas Principales

1. Crear Orden:
   - Ruta: `POST /orders`
   - Funcionalidad: Permite crear una nueva orden en el sistema.

2. Obtener Órdenes del Usuario:
   - Ruta: `GET /orders`
   - Funcionalidad: Obtiene la lista de órdenes del usuario autenticado.

## Modelos de Datos

El modelo principal es `Order`, que incluye campos como:

```5:13:user-service/models/user.go
type Order struct {
    gorm.Model
    UserID    uint    `json:"user_id"`
    ProductID uint    `json:"product_id"`
    Quantity  int     `json:"quantity"`
    Total     float64 `json:"total"`
}
```

## Seguridad

- Se utiliza un middleware de autenticación para proteger las rutas, el cual valida el token con el servicio de usuarios.

## Despliegue

El proyecto incluye configuraciones para despliegue utilizando Docker y Docker Compose. Los archivos relevantes son:

- `Dockerfile`: Para la construcción de la imagen Docker del servicio.
- `docker-compose.yml`: Para orquestar el servicio junto con otros microservicios y la base de datos PostgreSQL.

## Desarrollo y Pruebas

Para el desarrollo local:

1. Clonar el repositorio.
2. Copiar el archivo `.env.example` a `.env` y configurar las variables de entorno.
3. Ejecutar `go mod download` para instalar las dependencias.
4. Usar `go run main.go` para iniciar el servidor de desarrollo.

Para pruebas, se incluye un flujo de CI/CD en GitHub Actions que ejecuta pruebas automáticas en cada pull request.

## Notas Adicionales

Este proyecto es parte de una arquitectura de microservicios y demuestra conceptos como:
- Desarrollo de microservicios en Go
- Uso de frameworks web como Echo
- Implementación de autenticación y autorización
- Manejo de bases de datos con GORM
- Configuración de CI/CD con GitHub Actions
- Comunicación entre microservicios

Se recomienda revisar y mejorar las prácticas de seguridad antes de utilizar en un entorno de producción real.


