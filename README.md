# PT GEU Waste Management API

Backend RESTful API untuk sistem manajemen pengambilan sampah, dikembangkan sebagai bagian dari technical test PT Green Energi Utama. Proyek ini dibangun menggunakan bahasa pemrograman Golang, database PostgreSQL, dan Docker Desktop untuk containerize.

## Teknologi yang Digunakan

* **Bahasa:** Go (Golang)

* **Framework Web:** Gin

* **ORM:** GORM

* **Database:** PostgreSQL

* **Containerization:** Docker & Docker Compose

## Prasyarat (Prerequisites)

Pastikan sistem Anda sudah terinstal:

* [Docker Desktop](https://www.docker.com/products/docker-desktop/) atau Docker Engine & Docker Compose.

*(Aplikasi sudah di-containerize sepenuhnya, sehingga tidak memerlukan instalasi Go atau PostgreSQL secara manual di komputer lokal).*

## Cara Menjalankan Aplikasi

1. **Clone repositori**

   ```
   git clone https://github.com/jooharver/GEU-Test.git
   cd geu-test
   
   ```

2. **Siapkan Environment Variables**
   Ubah nama file `.env.example` menjadi `.env` (konfigurasi default sudah disesuaikan untuk Docker, Anda bisa mengganti nama DB dll sesuai dengan keinginan).

   ```
   cp .env.example .env
   
   ```

3. **Jalankan dengan Docker Compose**
   Di terminal, gunakan perintah berikut untuk melakukan build dan menjalankan seluruh container (API dan Database):

   ```
   docker compose up -d --build
   
   ```

   *Catatan: Sistem secara otomatis akan menjalankan database seeder saat pertama kali dihidupkan untuk memasukkan data warga (household) awal. Anda dapat mengecek container yang berjalan di Docker Desktop. Pastikan bahwa service `db` dan `api` sudah berjalan.*

4. **Akses API**
   Setelah langkah ke-3 selesai, Server akan berjalan secara lokal pada `http://localhost:8080` dan langsung bisa digunakan di Postman/Insomnia.

## Dokumentasi API & Testing

API ini dilengkapi dengan dokumentasi Postman yang sudah mencakup seluruh endpoint (Create, Read, Update, Delete) beserta penjelasan dan contoh payload JSON.

* **Postman Published Docs:** [Lihat Dokumentasi Lengkap API PT GEU di Sini](https://documenter.getpostman.com/view/51813315/2sBYHQ1hwf)

* **Postman Collection (Offline):** Anda juga dapat mengimpor file `geu_waste_management.postman_collection.json` yang terdapat pada root direktori proyek ini ke dalam workspace Postman atau Insomnia Anda.

## Contoh Penggunaan API

Berikut adalah beberapa contoh penggunaan endpoint utama pada aplikasi ini:

### 1. Menambahkan Warga Baru

**Request (POST `/api/households`):**

```
{
    "owner_name": "Joko Anwar",
    "email": "joko@example.com",
    "phone": "085712345678",
    "address": "Jl. Anggrek No. 99, Malang"
}

```

**Response (201 Created):**

```
{
    "data": {
        "id": "e4215b2e-a2f0-4a8b-8a8b-1c5c0a3e8b0a",
        "owner_name": "Joko Anwar",
        ...
    },
    "message": "Data warga berhasil ditambahkan"
}

```

### 2. Meminta Pengambilan Sampah

**Request (POST `/api/pickups`):**

```
{
    "household_id": "MASUKAN_UUID_WARGA_DISINI",
    "type": "electronic",
    "safety_check": true
}

```

**Response (201 Created):**

```
{
    "data": {
        "id": "f5b21c3d...",
        "household_id": "...",
        "type": "electronic",
        "status": "pending",
        "safety_check": true
    },
    "message": "Permintaan pickup berhasil dibuat"
}

```

### 3. Menjadwalkan Pengambilan Sampah

**Request (PUT `/api/pickups/:id/schedule`):**

```
{
    "pickup_date": "2026-10-15 10:00"
}

```

**Response (200 OK):**

```
{
    "message": "Status pickup berhasil diubah jadi scheduled"
}

```

### 4. Mengonfirmasi Pembayaran Tagihan

**Request (PUT `/api/payments/:id/confirm`):**
Endpoint ini menggunakan `multipart/form-data` untuk menerima file gambar bukti transfer.

* **Key:** `receipt`

* **Type:** `File`

* **Value:** `[Pilih file gambar Anda, misal: transfer_geu.jpg]`

**Response (200 OK):**

```
{
    "message": "Pembayaran berhasil dikonfirmasi"
}

```