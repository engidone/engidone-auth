# Engidone Auth Service

Servicio de autenticación microservicio construido con **Go** y **gRPC**, implementando inyección de dependencias con **fx** y arquitectura limpia. Proporciona endpoints de autenticación y un servicio de saludo para demostración.

## 🚀 Características

✅ **Inyección de Dependencias con fx**: Gestión automática de dependencias
✅ **Arquitectura Limpia**: Separación clara de dominios, casos de uso, infraestructura y transporte
✅ **gRPC**: Comunicación eficiente con Protocol Buffers
✅ **Modular por Diseño**: Providers organizados por módulo
✅ **Go-Kit**: Estructura de microservicios estándar
✅ **Ciclo de Vida Automático**: Gestión fx de start/stop
✅ **Testing Ready**: Estructura optimizada para pruebas
✅ **Kubernetes Ready**: Ideal para despliegue en contenedores

## 📁 Estructura del Proyecto

```
engidone-auth/
├── internal/
│   ├── di/                      # 🔥 Dependency Injection con fx
│   │   ├── app_providers.go     # Providers de aplicación (logger, config)
│   │   ├── hello_providers.go   # Providers del servicio Hello
│   │   ├── signin_providers.go  # Providers del servicio Signin
│   │   └── grpc_providers.go    # Providers de transporte gRPC
│   ├── hello/                   # 📝 Servicio Hello
│   │   ├── domain/              # Interfaces y modelos de dominio
│   │   ├── infrastructure/      # Implementaciones concretas
│   │   ├── usecase/            # Casos de uso de negocio
│   │   ├── transport/          # Capa de transporte gRPC
│   │   └── proto/              # Definiciones Protocol Buffers
│   └── signin/                  # 🔐 Servicio de Autenticación
│       ├── domain/              # Interfaces y modelos de dominio
│       ├── infrastructure/      # Implementaciones (repositorios, tokens)
│       ├── usecase/            # Casos de uso (signin, validate, refresh)
│       ├── transport/          # Capa de transporte gRPC
│       └── proto/              # Definiciones Protocol Buffers
├── cmd/
│   └── server/
│       └── main.go             # 🚀 Aplicación principal con fx
├── bin/                        # Binarios compilados
├── go.mod                      # Módulo Go con dependencias fx
├── Makefile                   # Automatización de builds
└── README.md                  # Este archivo
```

## 🏗️ Arquitectura con Inyección de Dependencias

### **Estructura de Providers fx:**

```go
// main.go - Aplicación simplificada con fx
func main() {
    app := fx.New(
        di.LoggerModule,      // Logger y configuración
        di.ConfigModule,
        di.HelloModule,       // Servicio Hello completo
        di.SigninModule,      // Servicio Signin completo
        di.GRPCModule,        // Transporte gRPC y server
    )
    app.Run()
}
```

### **Flujo de Inyección:**

1. **App Layer**: `LoggerModule` + `ConfigModule`
2. **Domain Layer**: `HelloModule` + `SigninModule`
3. **Transport Layer**: `GRPCModule`
4. **Lifecycle**: fx gestiona start/stop automáticamente

## 🏃‍♂️ Ejecutar

### 1. Instalar dependencias

```bash
go mod tidy
```

### 2. Compilar la aplicación

```bash
go build -o bin/server ./cmd/server
```

### 3. Iniciar el servidor

```bash
./bin/server
```

**Salida esperada con fx:**
```
[Fx] PROVIDE  *zap.Logger <= engidone-auth/internal/di.NewZapLogger()
[Fx] PROVIDE  log.Logger <= engidone-auth/internal/di.NewGoKitLogger()
[Fx] PROVIDE  *di.AppConfig <= engidone-auth/internal/di.NewAppConfig()
[Fx] PROVIDE  domain.HelloService <= engidone-auth/internal/di.NewHelloService()
[Fx] PROVIDE  domain.HelloUseCase <= engidone-auth/internal/di.NewHelloUseCase()
[Fx] PROVIDE  domain.UserRepository <= engidone-auth/internal/di.NewUserRepository()
[Fx] PROVIDE  domain.TokenService <= engidone-auth/internal/di.NewTokenService()
[Fx] PROVIDE  domain.SigninUseCase <= engidone-auth/internal/di.NewSigninUseCase()
[Fx] PROVIDE  domain.ValidateTokenUseCase <= engidone-auth/internal/di.NewValidateTokenUseCase()
[Fx] PROVIDE  domain.RefreshTokenUseCase <= engidone-auth/internal/di.NewRefreshTokenUseCase()
[Fx] PROVIDE  domain.GetUserUseCase <= engidone-auth/internal/di.NewGetUserUseCase()
[Fx] PROVIDE  *grpc.Server <= engidone-auth/internal/di.NewGRPCServer()
[Fx] PROVIDE  endpoints.Set <= engidone-auth/internal/di.NewHelloEndpoints()
[Fx] PROVIDE  endpoints.Set <= engidone-auth/internal/di.NewSigninEndpoints()
[Fx] PROVIDE  proto.HelloServiceServer <= engidone-auth/internal/di.NewHelloGRPCServer()
[Fx] PROVIDE  proto.SigninServiceServer <= engidone-auth/internal/di.NewSigninGRPCServer()
[Fx] PROVIDE  net.Listener <= engidone-auth/internal/di.NewTCPListener()
[Fx] INVOKE   engidone-auth/internal/di.RegisterGRPCServices()

msg==== Engidone Auth Service ===
transport=gRPC addr=:9000
msg=Servidor iniciado en :9000
msg=Servicios disponibles:
msg=  - Signin Service
msg=  - Hello Service
msg==== Usuarios disponibles para testing ===
msg=Username: admin, Password: password123
msg=Username: testuser, Password: test123
msg=Username: john, Password: john123
msg======================================
```

### 4. Variables de Entorno

```bash
# Configurar puerto del servidor (default: 9000)
export SERVER_PORT=9000

# Configurar secreto JWT (default: your-secret-key)
export JWT_SECRET=your-production-secret
```

## 🔌 API gRPC

### HelloService

#### `Hello`
Genera un saludo personalizado.

```protobuf
message HelloRequest {
  string name = 1;
}

message HelloResponse {
  string message = 1;
  bool success = 2;
}
```

### SigninService

#### `Signin`
Autentica usuario y genera token.

```protobuf
message SigninRequest {
  string username = 1;
  string password = 2;
}

message SigninResponse {
  bool success = 1;
  string message = 2;
  string user_id = 3;
  string username = 4;
  string email = 5;
  string token = 6;
  int64 expires_at = 7;
}
```

#### `ValidateToken`
Valida un token JWT y retorna información del usuario.

```protobuf
message ValidateTokenRequest {
  string token = 1;
}

message ValidateTokenResponse {
  bool valid = 1;
  string message = 2;
  string user_id = 3;
  string username = 4;
  string email = 5;
}
```

#### `RefreshToken`
Refresca un token existente.

```protobuf
message RefreshTokenRequest {
  string user_id = 1;
  string token = 2;
}
```

#### `GetUser`
Obtiene información de un usuario por ID.

```protobuf
message GetUserRequest {
  string user_id = 1;
}

message GetUserResponse {
  bool success = 1;
  string message = 2;
  string user_id = 3;
  string username = 4;
  string email = 5;
  int64 created_at = 6;
  int64 updated_at = 7;
}
```

## 👥 Usuarios de Prueba

| Username | Password | Rol |
|----------|----------|-----|
| admin | password123 | Administrador |
| testuser | test123 | Usuario de prueba |
| john | john123 | Usuario regular |

## 🛠️ Tecnologías

- **Go**: 1.21+
- **fx**: v1.24.0 (Inyección de dependencias)
- **gRPC**: Comunicación via Protocol Buffers
- **Go-Kit**: Estructura de microservicios
- **Clean Architecture**: Separación de capas
- **Zap**: Logger estructurado

## 🔧 Desarrollo

### Agregar Nuevo Servicio

1. **Crear estructura de dominio:**
```
internal/newservice/
├── domain/        # Interfaces y modelos
├── infrastructure/ # Implementaciones
├── usecase/       # Casos de uso
├── transport/     # gRPC handlers
└── proto/         # Protocol Buffers
```

2. **Crear providers fx:**
```go
// internal/di/newservice_providers.go
var NewServiceModule = fx.Options(
    fx.Provide(
        NewNewServiceRepository,
        NewNewServiceUseCase,
        NewNewServiceGRPCServer,
    ),
)
```

3. **Agregar al main.go:**
```go
app := fx.New(
    // ... módulos existentes
    di.NewServiceModule,  // ← Agregar nuevo módulo
)
```

### Debug con fx

```bash
# Ver grafo de dependencias
export FX_GRAPH=1
./bin/server > dependency-graph.dot
```

### Logs de fx

Fx proporciona logs detallados del proceso de inyección:

```
[Fx] PROVIDE  Type <= package.Function()
[Fx] RUN      provide: package.Function() in 1.234µs
[Fx] INVOKE   package.Function()
```

## ☸️ Despliegue

### Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/server .
EXPOSE 9000
CMD ["./server"]
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: engidone-auth
spec:
  replicas: 3
  selector:
    matchLabels:
      app: engidone-auth
  template:
    metadata:
      labels:
        app: engidone-auth
    spec:
      containers:
      - name: auth
        image: engidone-auth:latest
        ports:
        - containerPort: 9000
        env:
        - name: SERVER_PORT
          value: "9000"
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: auth-secrets
              key: jwt-secret
```

## 🐛 Troubleshooting

### Errores Comunes

1. **Puerto en uso:**
```bash
lsof -i :9000
kill -9 <PID>
```

2. **Errores de fx:**
```bash
# Revisa que todos los providers tengan las dependencias correctas
go build ./...
```

3. **Problemas de módulos:**
```bash
go mod tidy
go mod verify
```

### Verificar Configuración

```bash
# Test de compilación
go build -o bin/server ./cmd/server

# Test de dependencias fx
go test ./internal/di/...
```

## 🧪 Testing

### Unit Tests con fx

```go
func TestSigninUseCase(t *testing.T) {
    app := fx.New(
        di.TestModule,  // módulo de pruebas
        di.SigninModule,
    )

    startCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := app.Start(startCtx); err != nil {
        t.Fatal(err)
    }
    defer app.Stop(context.Background())

    // ... pruebas
}
```

## 📚 Arquitectura Limpia

El proyecto sigue los principios de Clean Architecture:

- **Domain**: Interfaces y reglas de negocio (sin dependencias externas)
- **Use Cases**: Lógica de aplicación (dependen solo de domain)
- **Infrastructure**: Implementaciones concretas (DB, APIs externas)
- **Transport**: gRPC, HTTP, etc (adaptadores externos)
- **DI**: fx gestiona el ensamblaje de todas las capas

## 🤝 Contribuir

1. Fork el proyecto
2. Crear feature branch: `git checkout -b feature/amazing-feature`
3. Commit cambios: `git commit -m 'Add amazing feature'`
4. Push: `git push origin feature/amazing-feature`
5. Pull Request

## 📄 Licencia

MIT License - ver archivo [LICENSE](LICENSE) para detalles