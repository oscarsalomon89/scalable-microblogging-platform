# 📄 Business Rules & Assumptions

## 🧠 Supuestos funcionales

### 1. **Longitud de tweet**

Cada tweet tiene un límite máximo de 280 caracteres. Este valor se asume por convención y no está definido en el enunciado.

### 2. **Relaciones entre usuarios**

- Un usuario puede seguir a múltiples otros usuarios.
- Un usuario **no puede seguirse a sí mismo**.
- No se permite duplicación de follows.

### 3. **Eliminación de relaciones**

Se permite dejar de seguir a un usuario (`unfollow`).

### 4. **Visualización de timeline**

- El timeline muestra los tweets de los usuarios a los que el usuario sigue.
- No incluye los tweets propios del usuario (aunque podría ser modificado).
- Los tweets se ordenan de más nuevo a más antiguo.
- Se devuelve una cantidad limitada (por ejemplo, los últimos 50).

### 5. **Edición y eliminación de tweets**

- No se implementa la edición ni la eliminación de tweets, ya que no fue solicitado en el enunciado.
- Solo se permite la creación de nuevos tweets.

### 6. **Usuarios y autenticación**

- Se asume que los IDs de usuario que llegan por la API son válidos.
- No se implementa autenticación ni autorización en esta versión.

---

## ⚙️ Supuestos técnicos

### 7. Arquitectura

- Se adoptó una **arquitectura hexagonal (Clean Architecture)**, respetando los principios de separación de responsabilidades entre:

  - **Application** (entidades, reglas de negocio y casos de uso),
  - **Adapters** (handlers HTTP, cache Redis, repositorios, etc).
  - **platform/**: componentes de infraestructura (configuración, enviroment, herramientas internas). Proporciona soporte a la aplicación sin acoplarla directamente.
  - **pkg/**: bibliotecas genéricas y reutilizables, independientes del dominio del negocio. Contiene helpers, validadores, contextos, etc.

- La organización del código sigue el enfoque de **Package Oriented Design**, donde cada paquete representa una unidad funcional autocontenida, evitando estructuras artificiales basadas únicamente en capas técnicas.

- Para la **inyección de dependencias**, se utiliza **Uber FX**, lo que permite definir claramente los componentes del sistema, sus relaciones y ciclos de vida. Esto mejora la escalabilidad del código y facilita el testing, habilitando el uso de **mocks y fakes** en tests sin acoplamientos innecesarios.

### 8. Contenerización y Despliegue

- Se utiliza **Docker Compose** para levantar el entorno de desarrollo completo.
- Los servicios involucrados son:
  - PostgreSQL (persistencia de usuarios y tweets)
  - Redis (cache de timelines)
- La solución está preparada para ser desplegada en entornos orquestados (Kubernetes, AWS, Cloud Run, etc.) con pequeñas modificaciones de configuración.

### 9. Testing

- Se implementaron **pruebas unitarias enfocadas en la capa de casos de uso**, donde reside la lógica de negocio y se coordinan las reglas entre entidades, repositorios y servicios externos. Este enfoque permite validar el comportamiento del sistema de forma aislada, sin necesidad de duplicar tests en capas inferiores como entidades puras o adaptadores.

- Se utilizan herramientas como `testify` para aserciones y `mockery` para la generación automática de mocks a partir de interfaces. En algunos casos también se emplean fakes escritos a mano para mayor control.

### 10. Cache de timeline con Redis

- El timeline de cada usuario se almacena temporalmente en Redis para lecturas rápidas.
- Se aplica una política de TTL (ej. 1 minuto) para evitar inconsistencias prolongadas.
- Redis contiene una lista ordenada por tiempo y se invalida ante ciertos eventos.

### 11. Estrategia de timeline: invalidación (fan-out-on-read)

- Cuando un usuario publica un nuevo tweet, se invalidan los timelines cacheados de todos sus seguidores.
- En la próxima lectura, se reconstruye desde base de datos y se vuelve a cachear.

### 12. Alternativas consideradas\*\*

- El timeline no se actualiza al momento de publicar un tweet.
- Cuando un usuario publica un tweet, se lanza una goroutine que obtiene a todos sus seguidores e invalida (borra) sus timelines cacheados en Redis.
- Cuando un usuario hace follow o unfollow, también se invalida su timeline para asegurar consistencia.
- En la próxima lectura (`GET /timeline/:user_id`), si no existe cache, se reconstruye consultando los últimos tweets de los usuarios que sigue, se ordena y se cachea nuevamente en Redis con TTL.
- Por simplicidad y claridad, se optó por invalidación + reconstrucción.

### 13. Escalabilidad futura

- Actualmente se usan goroutines para invalidar timelines de seguidores.
- Para escalar correctamente, se sugiere migrar esta lógica a una cola asincrónica (ej. Pub/Sub, Kafka, SQS) que procese los eventos fuera del request/response.

---

## 🚫 Alcance limitado

### 14. Endpoints omitidos por simplicidad

No se implementaron endpoints para eliminar tweets, listar seguidores o seguir múltiples usuarios en lote, ya que no eran requeridos directamente.

### 15. Manejo de errores y validaciones

- Se hace validación de existencia de usuarios y contenido de tweets.
- Errores técnicos y de negocio se diferencian para facilitar su manejo en la API.
