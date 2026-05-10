# TutorKu Backend

TutorKu adalah aplikasi **AI-powered learning assistant** berbasis platform, yang membantu mahasiswa/pelajar belajar dari materi PDF mereka sendiri. Backend ini adalah REST API yang ditulis dalam **Go (Golang)**, bertindak sebagai orchestrator antara client (mobile/web) dan layanan AI yang dikelola lewat **n8n** (workflow automation).

---

## 🧠 Konsep Utama

Pengguna mengunggah file PDF (materi kuliah/buku/dokumen), sistem memproses file tersebut menggunakan AI pipeline di n8n (parsing, embedding, vector store), lalu pengguna dapat:

1. **Chat** dengan materi mereka (RAG / Retrieval-Augmented Generation)
2. **Merangkum** materi dalam 3 mode: singkat, detail, atau mindmap
3. **Generate quiz** (pilihan ganda, esai, benar/salah) dari materi
4. **Evaluasi esai** yang ditulis oleh pengguna terhadap materi

---

## 🏗️ Tech Stack

| Komponen         | Teknologi                              |
| ---------------- | -------------------------------------- |
| Language         | Go 1.25                                |
| HTTP Framework   | Gin                                    |
| ORM              | GORM                                   |
| Database         | PostgreSQL (with `pgcrypto` extension) |
| Auth             | JWT (HS256, `golang-jwt/jwt/v5`)       |
| Password Hashing | bcrypt                                 |
| AI Orchestration | n8n (via webhook HTTP calls)           |
| UUID             | `google/uuid`                          |
| Config           | `godotenv`                             |
| Deployment       | Docker / Docker Compose                |

---

## 📁 Struktur Project

```
tutorku-backend/
├── main.go                    # Entry point: init config, DB, DI, server
├── config/
│   └── config.go              # Load env variables (APP, DB, JWT, N8N webhooks)
├── routes/
│   └── routes.go              # Semua API route definitions
├── middleware/
│   └── auth.go                # JWT auth middleware
├── internal/
│   ├── domain/                # Struct model + request/response types
│   │   ├── user.go
│   │   ├── material.go
│   │   └── chat.go
│   ├── repository/            # Database access layer (GORM queries)
│   │   ├── user_repo.go
│   │   ├── material_repo.go
│   │   └── chat_repo.go
│   ├── service/               # Business logic layer
│   │   ├── auth_service.go
│   │   ├── material_service.go
│   │   ├── chat_service.go
│   │   └── feature_service.go
│   └── handler/               # HTTP handler layer (gin context)
│       ├── auth_handler.go
│       ├── material_handler.go
│       ├── chat_handler.go
│       └── feature_handler.go
└── pkg/
    ├── n8n/
    │   └── client.go          # HTTP client untuk memanggil n8n webhooks
    ├── postgres/
    │   └── ...                # Koneksi PostgreSQL
    └── response/
        └── ...                # Helper untuk response JSON standar
```

---

## 🗄️ Database Schema

### `users`

| Field        | Type            | Keterangan                             |
| ------------ | --------------- | -------------------------------------- |
| `id`         | UUID (PK)       | Auto-generated via `gen_random_uuid()` |
| `name`       | string          | Nama pengguna                          |
| `email`      | string (unique) | Email untuk login                      |
| `password`   | string          | bcrypt hash                            |
| `created_at` | timestamp       |                                        |
| `deleted_at` | timestamp       | Soft delete                            |

### `materials`

| Field        | Type              | Keterangan                        |
| ------------ | ----------------- | --------------------------------- |
| `id`         | UUID (PK)         |                                   |
| `user_id`    | UUID (FK → users) | Owner materi                      |
| `title`      | string            | Judul materi                      |
| `category`   | string            | Kategori (opsional)               |
| `filename`   | string            | Nama file PDF                     |
| `file_size`  | int64             | Ukuran file dalam bytes           |
| `status`     | enum              | `processing` → `ready` / `failed` |
| `created_at` | timestamp         |                                   |
| `deleted_at` | timestamp         | Soft delete                       |

### `chat_sessions`

| Field         | Type                  | Keterangan               |
| ------------- | --------------------- | ------------------------ |
| `id`          | UUID (PK)             |                          |
| `user_id`     | UUID (FK → users)     |                          |
| `material_id` | UUID (FK → materials) | Sesi terkait materi mana |
| `created_at`  | timestamp             |                          |

### `chat_messages`

| Field        | Type                      | Keterangan              |
| ------------ | ------------------------- | ----------------------- |
| `id`         | UUID (PK)                 |                         |
| `session_id` | UUID (FK → chat_sessions) |                         |
| `role`       | string                    | `user` atau `assistant` |
| `content`    | text                      | Isi pesan               |
| `created_at` | timestamp                 |                         |

---

## 🔗 API Endpoints

Base URL: `/api/v1`

### 🔓 Public (tanpa auth)

| Method | Endpoint              | Deskripsi                                    |
| ------ | --------------------- | -------------------------------------------- |
| `POST` | `/auth/register`      | Daftar akun baru                             |
| `POST` | `/auth/login`         | Login, mendapatkan JWT token                 |
| `POST` | `/callback/ingestion` | Callback dari n8n setelah proses PDF selesai |

### 🔒 Protected (wajib `Authorization: Bearer <token>`)

#### Materials

| Method   | Endpoint                  | Deskripsi                             |
| -------- | ------------------------- | ------------------------------------- |
| `POST`   | `/materials/upload`       | Upload file PDF (multipart/form-data) |
| `GET`    | `/materials`              | Ambil semua materi milik user         |
| `GET`    | `/materials/:id`          | Ambil detail materi                   |
| `GET`    | `/materials/:id/status`   | Cek status processing materi          |
| `GET`    | `/materials/:id/download` | Download materi PDF                   |
| `DELETE` | `/materials/:id`          | Hapus materi                          |

#### Chat

| Method | Endpoint                    | Deskripsi                           |
| ------ | --------------------------- | ----------------------------------- |
| `POST` | `/chat`                     | Kirim pertanyaan ke AI (RAG)        |
| `GET`  | `/chat/sessions`            | Ambil semua sesi chat milik user    |
| `GET`  | `/chat/history/:session_id` | Ambil riwayat pesan dalam satu sesi |

#### Features (AI)

| Method | Endpoint              | Deskripsi              |
| ------ | --------------------- | ---------------------- |
| `POST` | `/features/summarize` | Rangkum materi         |
| `POST` | `/features/quiz`      | Generate soal kuis     |
| `POST` | `/features/essay`     | Evaluasi esai pengguna |

---

## 🔄 Alur Sistem (Flow)

### 1. Registrasi & Login

```
Client → POST /auth/register → AuthHandler → AuthService
  → bcrypt(password) → UserRepo.Create() → DB
  → generateToken(userID) → return { token, user }

Client → POST /auth/login → AuthHandler → AuthService
  → UserRepo.FindByEmail() → bcrypt.CompareHash() → generateToken() → return { token, user }
```

### 2. Upload Materi (Async Processing)

```
Client → POST /materials/upload (multipart PDF)
  → [AuthMiddleware: validasi JWT]
  → MaterialHandler.Upload()
  → MaterialService.Upload()
    ├─ Simpan record Material ke DB (status: "processing")
    ├─ Return response ke client (langsung, tidak menunggu)
    └─ [goroutine async] → n8nClient.TriggerIngestion()
         → POST ke n8n webhook "ingest"
           Payload: { material_id, user_id, file_base64, filename }
         → n8n memproses: parse PDF → chunk text → embed → simpan ke vector store
         → n8n memanggil callback: POST /api/v1/callback/ingestion
           Payload: { material_id, status: "ready"/"failed" }
         → MaterialService.UpdateStatus() → DB (status diperbarui)
```

### 3. Chat dengan Materi (RAG)

```
Client → POST /chat { material_id, session_id (optional), query }
  → [AuthMiddleware]
  → ChatHandler.Chat()
  → ChatService.Chat()
    ├─ Validasi: material milik user & status == "ready"
    ├─ Jika session_id kosong → buat ChatSession baru di DB
    ├─ Ambil 10 pesan terakhir dari sesi (chat history)
    ├─ Simpan pesan user ke DB (role: "user")
    ├─ n8nClient.QueryRAG()
    │    → POST ke n8n webhook "chat"
    │      Payload: { material_id, user_id, query, chat_history }
    │    → n8n: retrieve context dari vector store → LLM → return answer
    └─ Simpan jawaban AI ke DB (role: "assistant")
    → Return { session_id, answer }
```

### 4. Ringkasan Materi

```
Client → POST /features/summarize { material_id, mode: "short"|"detailed"|"mindmap" }
  → [AuthMiddleware]
  → FeatureHandler.Summarize()
  → FeatureService.Summarize()
    ├─ Validasi: material milik user & status == "ready"
    └─ n8nClient.Summarize()
         → POST ke n8n webhook "summarize"
           Payload: { material_id, user_id, mode }
         → n8n: retrieve content → LLM summarize → return result
```

### 5. Generate Kuis

```
Client → POST /features/quiz { material_id, type, count, difficulty }
  → [AuthMiddleware]
  → FeatureHandler.GenerateQuiz()
  → FeatureService.GenerateQuiz()
    ├─ Validasi: material milik user & status == "ready"
    └─ n8nClient.GenerateQuiz()
         → POST ke n8n webhook "quiz"
           Payload: { material_id, user_id, type, count, difficulty }
         → n8n: generate questions via LLM → return quiz data
```

**type**: `multiple_choice` | `essay` | `true_false`  
**difficulty**: `easy` | `medium` | `hard`  
**count**: 1–20

### 6. Evaluasi Esai

```
Client → POST /features/essay { material_id, title, content (min 50 chars) }
  → [AuthMiddleware]
  → FeatureHandler.EvaluateEssay()
  → FeatureService.EvaluateEssay()
    ├─ Validasi: material milik user & status == "ready"
    └─ n8nClient.EvaluateEssay()
         → POST ke n8n webhook "essay"
           Payload: { material_id, user_id, title, content }
         → n8n: compare essay vs materi via LLM → return evaluation/feedback
```

---

## ⚙️ Integrasi n8n

Backend ini tidak menjalankan AI sendiri. Semua AI processing di-delegate ke **n8n** melalui webhook HTTP calls.

| Webhook Key | Fungsi                                           |
| ----------- | ------------------------------------------------ |
| `ingest`    | Proses PDF: parse → chunk → embed → vector store |
| `chat`      | RAG query: retrieve context → LLM answer         |
| `summarize` | Summarize materi dengan mode tertentu            |
| `quiz`      | Generate soal kuis                               |
| `essay`     | Evaluasi/feedback esai pengguna                  |

Config webhook diset via environment variables:

```
N8N_BASE_URL=http://your-n8n-instance
N8N_WEBHOOK_INGEST=/webhook/xxx
N8N_WEBHOOK_CHAT=/webhook/yyy
N8N_WEBHOOK_SUMMARIZE=/webhook/zzz
N8N_WEBHOOK_QUIZ=/webhook/aaa
N8N_WEBHOOK_ESSAY=/webhook/bbb
```

---

## 🔑 Environment Variables

```env
# App
APP_PORT=8080
APP_ENV=development

# Database
DB_HOST=
DB_PORT=
DB_NAME=
DB_USER=
DB_PASS=

# JWT
JWT_SECRET=your-super-secret-key-here
JWT_EXPIRES_HOUR=72

# n8n
N8N_BASE_URL=
N8N_WEBHOOK_INGEST=/pdf-ingest
N8N_WEBHOOK_CHAT=/rag-chat
N8N_WEBHOOK_SUMMARIZE=/summarize
N8N_WEBHOOK_QUIZ=/quiz
N8N_WEBHOOK_ESSAY=/essay
```

---

## 🚀 Menjalankan

### Lokal

```bash
go mod tidy
cp .env.example .env   # isi konfigurasi
go run main.go
```

### Docker Compose

```bash
docker-compose up -d
```

---

## 🏛️ Arsitektur Layering

```
HTTP Request
     ↓
[Middleware] → AuthMiddleware (JWT validation, set user_id ke context)
     ↓
[Handler] → Parse request, validasi binding, call service, return response
     ↓
[Service] → Business logic: validasi ownership, orchestrate repo + n8n calls
     ↓
[Repository] → GORM queries ke PostgreSQL
     ↓          ↕
[pkg/n8n]   → HTTP calls ke n8n webhooks (AI processing)
```

**Prinsip:**

- Handler hanya Handle HTTP (parse, respond)
- Service berisi semua business logic
- Repository hanya database queries
- n8n Client hanya HTTP calls ke n8n
- Domain hanya definisi struct

---

## 📌 Catatan Penting untuk AI / Developer

- **Material harus berstatus `ready`** sebelum bisa digunakan untuk chat/summarize/quiz/essay. Ini divalidasi di semua service terkait.
- **Upload PDF bersifat async**: response diberikan langsung setelah record disimpan ke DB, sedangkan proses n8n berjalan di goroutine terpisah.
- **Chat session**: jika `session_id` tidak dikirim, sistem otomatis membuat sesi baru. Session digunakan untuk menyimpan history percakapan.
- **Callback ingestion** (`POST /callback/ingestion`) adalah endpoint **tanpa auth**, dipanggil oleh n8n setelah selesai memproses PDF untuk memperbarui status material.
- Semua endpoint protected menggunakan `user_id` dari JWT claims, bukan dari request body, sehingga resource isolation per-user dijamin di level service.
