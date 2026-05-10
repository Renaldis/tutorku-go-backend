# Dokumentasi API - Flow Mengerjakan Kuis

Dokumentasi ini berisi panduan untuk menguji flow pengerjaan kuis (Quiz) melalui Postman, mulai dari men-generate kuis hingga melakukan submit (penilaian).

Semua endpoint di bawah ini memerlukan otorisasi (Protected Endpoint).
Pastikan menambahkan Header berikut di setiap request:
`Authorization: Bearer <token_jwt_anda>`

---

## 1. Generate Kuis (AI)

Menggunakan AI untuk men-generate kuis dari materi yang sudah berstatus `ready`.

* **Endpoint:** `POST /api/v1/materials/:id/generate-quiz`
* **Path Variable:**
  * `id`: ID dari Material (Materi PDF)

**Request Body (JSON):**
```json
{
  "type": "multiple_choice",  // Pilihan: "multiple_choice", "true_false", "essay"
  "difficulty": "medium",     // Pilihan: "easy", "medium", "hard"
  "count": 5                  // Jumlah soal (1 - 20)
}
```

*Response sukses (201 Created) akan mengembalikan objek Kuis yang telah tersimpan di Database beserta ID Kuis tersebut (`quiz_id`). Catat `id` kuis ini untuk langkah selanjutnya.*

---

## 2. Melihat Detail Kuis (Opsional)

Sebelum memulai kuis, Anda bisa melihat detail kuis beserta soal-soalnya (tanpa kunci jawaban, tergantung logic backend).

* **Endpoint:** `GET /api/v1/quizzes/:id`
* **Path Variable:**
  * `id`: ID dari Kuis (`quiz_id` yang didapat dari langkah 1)

---

## 3. Memulai Kuis (Start Attempt)

Endpoint ini digunakan ketika user menekan tombol "Mulai Kuis" di aplikasi. Sistem akan mencatat waktu mulai (`started_at`) dan membuat sesi pengerjaan (Attempt).

* **Endpoint:** `POST /api/v1/quizzes/:id/start`
* **Path Variable:**
  * `id`: ID dari Kuis yang akan dikerjakan

**Request Body:**
*(Tidak membutuhkan request body, karena ID kuis sudah ada di URL)*

**Contoh Response Sukses (201 Created):**
Perhatikan `id` yang dikembalikan di dalam data response. Itu adalah **Attempt ID** (`attempt_id`). Anda akan membutuhkannya untuk men-submit jawaban.
```json
{
  "status": "success",
  "message": "Quiz dimulai",
  "data": {
    "id": "abc123xx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", // INI ADALAH ATTEMPT ID
    "quiz_id": "...",
    "user_id": "...",
    "started_at": "2023-10-25T10:00:00Z"
  }
}
```

---

## 4. Submit Jawaban Kuis

Endpoint ini digunakan di akhir, ketika user sudah menjawab semua soal dan menekan tombol "Selesai" atau "Kumpulkan". Sistem akan menghitung skor akhir berdasarkan kunci jawaban di database.

* **Endpoint:** `POST /api/v1/quizzes/attempt/:attempt_id/submit`
* **Path Variable:**
  * `attempt_id`: ID Attempt yang didapatkan dari langkah 3.

**Request Body (JSON):**
Kirimkan *array* dari jawaban user yang berisi `question_id` (ID dari setiap soal) dan `user_answer` (Jawaban yang dipilih/diketik user).

```json
{
  "answers": [
    {
      "question_id": "id-soal-pertama",
      "user_answer": "A"  // Jika pilihan ganda, isi dengan option_key
    },
    {
      "question_id": "id-soal-kedua",
      "user_answer": "B"
    },
    {
      "question_id": "id-soal-ketiga",
      "user_answer": "True" // Jika true_false
    }
  ]
}
```

**Contoh Response Sukses (200 OK):**
Anda akan mendapatkan hasil skor, total soal yang benar, dan waktu selesai (`finished_at`).
```json
{
  "status": "success",
  "message": "Quiz disubmit",
  "data": {
    "id": "abc123xx-...",
    "score": 100,
    "total_correct": 3,
    "total_questions": 3,
    "started_at": "2023-10-25T10:00:00Z",
    "finished_at": "2023-10-25T10:05:00Z"
  }
}
```

---

## 5. Melihat Riwayat Pengerjaan Kuis (History)

Untuk melihat daftar riwayat pengerjaan user pada suatu kuis.

* **Endpoint:** `GET /api/v1/quizzes/:id/attempts`
* **Path Variable:**
  * `id`: ID dari Kuis
