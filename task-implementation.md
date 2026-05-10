# TASK: Implement AI Quiz System (Backend + Database Architecture)

Saya sedang membangun platform AI learning dengan stack:

- Golang
- GORM
- PostgreSQL
- Existing entities:
  - User
  - Material
  - ChatSession
  - ChatMessage

Saat ini sistem sudah bisa:

- upload material
- AI chat berdasarkan material

Sekarang saya ingin menambahkan fitur:

# AI Quiz Generator System

---

# EXISTING MODELS

```go
type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type Material struct {
	ID        string
	UserID    string
	Title     string
	Category  string
	Filename  string
	FileSize  int64
	Status    MaterialStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type ChatSession struct {
	ID         string
	UserID     string
	MaterialID string
	CreatedAt  time.Time
}

type ChatMessage struct {
	ID        string
	SessionID string
	Role      string
	Content   string
	CreatedAt time.Time
}
```

---

# GOAL

Buatkan implementation lengkap untuk sistem quiz yang scalable dan production-ready.

Quiz di-generate dari Material menggunakan AI.

---

# QUIZ FEATURES

Support tipe soal:

- multiple_choice
- true_false
- essay

---

# REQUIRED DATABASE MODELS

Buat model GORM lengkap beserta relation:

## Quiz

Relasi:

- belongs to User
- belongs to Material
- has many QuizQuestions
- has many QuizAttempts

Fields:

- id
- user_id
- material_id
- title
- description
- generated_by
- created_at
- updated_at
- deleted_at

---

## QuizQuestion

Relasi:

- belongs to Quiz
- has many QuizOptions

Fields:

- id
- quiz_id
- question
- type
- correct_answer
- explanation
- order_no
- created_at
- updated_at
- deleted_at

Question type:

- multiple_choice
- true_false
- essay

---

## QuizOption

Relasi:

- belongs to QuizQuestion

Fields:

- id
- question_id
- option_key
- option_text
- created_at

option_key example:

- A
- B
- C
- D

---

## QuizAttempt

Relasi:

- belongs to User
- belongs to Quiz
- has many QuizAnswers

Fields:

- id
- user_id
- quiz_id
- score
- total_correct
- total_questions
- started_at
- finished_at
- created_at

---

## QuizAnswer

Relasi:

- belongs to QuizAttempt
- belongs to QuizQuestion

Fields:

- id
- attempt_id
- question_id
- user_answer
- is_correct
- earned_points
- created_at

---

# REQUIREMENTS

Implement:

## 1. GORM Models

- full struct
- relations
- json tags
- gorm tags
- indexes
- UUID support

---

## 2. Migration Setup

Buat auto migration setup untuk semua quiz entities.

---

## 3. Quiz Generation Flow

Flow:

Material
→ AI generate quiz JSON
→ backend validate JSON
→ save quiz
→ save questions
→ save options

---

## 4. API Endpoints

Buat REST API structure:

### Generate Quiz

POST /materials/:id/generate-quiz

### Get Quiz

GET /quizzes/:id

### Start Attempt

POST /quizzes/:id/start

### Submit Quiz

POST /quizzes/:id/submit

### Get Attempt History

GET /quizzes/history

---

## 5. DTO / Request Response

Buat request & response struct untuk:

- generate quiz
- submit answers
- get quiz detail
- attempt result

---

## 6. Quiz JSON Structure

Gunakan format berikut:

```json
{
  "title": "DBMS Quiz",
  "description": "Quiz dasar DBMS",
  "questions": [
    {
      "question": "Apa itu Primary Key?",
      "type": "multiple_choice",
      "options": [
        {
          "key": "A",
          "text": "Kunci utama"
        },
        {
          "key": "B",
          "text": "Foreign key"
        }
      ],
      "correct_answer": "A",
      "explanation": "Primary key adalah identitas unik."
    }
  ]
}
```

---

# IMPORTANT ARCHITECTURE RULES

- Backend is source of truth
- Correct answer jangan hanya di frontend
- Quiz harus persistent di database
- Attempt user harus tersimpan
- Questions harus reusable
- Support future scalability:
  - retry quiz
  - review answers
  - analytics
  - adaptive learning
  - leaderboard

---

# OUTPUT EXPECTATION

Saya ingin output:

- folder structure
- models
- migrations
- repositories
- services
- handlers/controllers
- DTOs
- route examples
- clean architecture recommendation
- scalable project structure
- example implementation snippets

Gunakan best practice Golang + GORM + PostgreSQL.
