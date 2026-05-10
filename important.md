# REAL EXISTING PROJECT STRUCTURE

Project menggunakan structure berikut:

```text
internal/
 ├── domain
 ├── handler
 ├── repository
 └── service
```

Routing dipusatkan di:

```text
routes/routes.go
```

Database setup ada di:

```text
pkg/postgres/db.go
```

Shared response helper:

```text
pkg/response/response.go
```

Middleware auth:

```text
middleware/auth.go
```

AI integration:

```text
pkg/n8n/client.go
```

---

# IMPORTANT IMPLEMENTATION RULE

Jangan membuat:

- module folder baru
- feature-based architecture baru
- clean architecture yang terlalu kompleks
- hexagonal architecture
- cqrs
- unnecessary abstraction

Karena project saat ini menggunakan architecture sederhana berbasis:

- domain
- repository
- service
- handler

Maka fitur quiz HARUS mengikuti pattern existing.

---

# EXPECTED QUIZ FILE STRUCTURE

Tambahkan file baru seperti berikut:

```text
internal/
 ├── domain/
 │    ├── quiz.go
 │    ├── quiz_attempt.go
 │    └── quiz_answer.go
 │
 ├── repository/
 │    └── quiz_repo.go
 │
 ├── service/
 │    └── quiz_service.go
 │
 ├── handler/
 │    └── quiz_handler.go
```

Jika perlu DTO:

```text
internal/dto/
```

Tetapi hanya jika memang diperlukan.

Jika existing project belum memakai dto folder:

- response/request struct boleh tetap di handler atau service
- prioritaskan consistency

---

# ROUTING INTEGRATION

Tambahkan route quiz ke existing:

```text
routes/routes.go
```

Jangan membuat routing system baru.

---

# DATABASE MIGRATION

Tambahkan quiz migration ke existing auto migration flow di:

```text
pkg/postgres/db.go
```

Jangan membuat migration framework baru.

---

# RESPONSE FORMAT

Gunakan helper existing:

```text
pkg/response/response.go
```

Jangan membuat custom response formatter baru.

---

# AUTHENTICATION

Gunakan middleware existing:

```text
middleware/auth.go
```

Semua endpoint quiz harus protected.

---

# REPOSITORY STYLE

Ikuti pattern existing repository:

- material_repo.go
- chat_repo.go
- user_repo.go

Jangan membuat generic repository abstraction.

---

# SERVICE STYLE

Ikuti pattern existing service:

- material_service.go
- chat_service.go

Business logic tetap di service layer.

---

# HANDLER STYLE

Ikuti style existing handler:

- auth_handler.go
- material_handler.go

Handler hanya:

- parse request
- call service
- return response

Jangan taruh business logic besar di handler.

---

# FINAL GOAL

Quiz system harus terasa seperti:

- bagian natural dari existing project
- bukan project baru yang ditempel

Prioritaskan:

- consistency
- maintainability
- low complexity
- scalable enough
