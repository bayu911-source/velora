# Velora Development TODO

Ini adalah daftar tugas untuk pengembangan proyek Velora di masa depan.

## Fase 1: Peningkatan Inti & Stabilitas

- [ ] **Integrasikan `MemoryManager`**:
  - Modifikasi `Runner` untuk menampung sebuah `MemoryManager`.
  - Izinkan `Agent` untuk membaca dari dan menulis ke `MemoryManager` selama eksekusi `Run`. Ini akan memungkinkan agen untuk berbagi status dan memiliki konteks percakapan yang lebih panjang.
  - Perbarui `Agent` interface untuk menyertakan parameter `context` atau `memory`.

- [ ] **Terapkan Unit Testing**:
  - Tulis unit test untuk `workflow.Runner` untuk memastikan pipeline dieksekusi dengan benar.
  - Tulis unit test untuk setiap `Agent` untuk memverifikasi perilakunya secara terpisah.
  - Gunakan table-driven tests untuk mencakup berbagai skenario input.

- [ ] **File Konfigurasi Terpusat**:
  - Buat file konfigurasi (misalnya, `config.yaml` atau `config.json`).
  - Pindahkan nilai-nilai yang di-hardcode seperti nama model AI (`gemini-1.5-pro-latest`) ke dalam file konfigurasi ini.
  - Muat konfigurasi ini saat aplikasi dimulai.

- [ ] **Logging Terstruktur**:
  - Ganti `log.Fatalf` dan `fmt.Println` dengan library logging terstruktur (misalnya, `slog` dari Go 1.21+).
  - Tambahkan ID korelasi untuk melacak satu eksekusi pipeline di seluruh log.

## Fase 2: Kemampuan Agen & Pipeline

- [ ] **Sistem "Tools" untuk Agen**:
  - Definisikan `Tool` interface yang dapat digunakan oleh agen (misalnya, `ReadFileTool`, `ExecuteCommandTool`).
  - Modifikasi `Agent` untuk dapat diberikan daftar `Tool` yang tersedia.
  - Refactor `WebResearchAgent` untuk menggunakan `WebSearchTool` sebagai contoh.

- [ ] **Agen Perencana (Planner Agent)**:
  - Buat `PlannerAgent` baru yang tugasnya adalah menerima tujuan tingkat tinggi dari pengguna.
  - `PlannerAgent` ini akan secara dinamis menghasilkan `Pipeline` (daftar langkah dan agen) untuk mencapai tujuan tersebut.

- [ ] **Output Terstruktur dari Agen**:
  - Alih-alih hanya mengembalikan string, refactor `Agent.Run` untuk mengembalikan struct yang lebih kaya, yang bisa mencakup pemikiran agen, alat yang digunakan, dan hasil akhirnya.

- [ ] **Eksekusi Pipeline Secara Paralel**:
  - Analisis dependensi antar langkah pipeline.
  - Modifikasi `Runner` untuk menjalankan langkah-langkah yang tidak saling bergantung secara bersamaan (concurrently) menggunakan goroutine.

## Fase 3: Pengalaman Pengguna & Developer

- [ ] **Peningkatan Antarmuka Baris Perintah (CLI)**:
  - Gunakan library seperti `cobra` atau `urfave/cli` untuk membuat CLI yang lebih kuat.
  - Tambahkan sub-perintah seperti `velora run <pipeline.json>`, `velora agent list`, `velora plan "my goal is..."`.

- [ ] **Dokumentasi Proyek**:
  - Tulis dokumentasi yang jelas tentang cara membuat `Agent` baru.
  - Dokumentasikan format file `pipeline.json` dan file konfigurasi.
  - Tambahkan contoh penggunaan yang lebih banyak.

- [ ] **Manajemen Paket yang Lebih Baik**:
  - Pindahkan definisi `Agent` interface ke paket `agent` atau `pkg/agent` sendiri untuk kejelasan.
  - Evaluasi kembali struktur paket saat fungsionalitas bertambah.
