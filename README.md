# Velora

Velora adalah agen yang didukung AI berbasis Go yang dirancang untuk menjadi pendamping pengkodean yang proaktif dan efisien. Ini memanfaatkan AI Generatif Google untuk memahami dan memenuhi permintaan pengguna dalam lingkungan pengembangan mereka.

## Fitur

*   **Agen yang Didukung AI**: Velora menggunakan serangkaian agen AI untuk melakukan berbagai tugas, termasuk:
    *   **Pembangun Aplikasi**: Membangun aplikasi berdasarkan spesifikasi pengguna.
    *   **Agen Obrolan**: Terlibat dalam interaksi percakapan dengan pengguna.
    *   **Generator Kode**: Menghasilkan cuplikan kode dan file.
    *   **Ekstraktor Data**: Mengekstrak informasi dari berbagai sumber.
*   **Arsitektur Plugin**: Fungsionalitas Velora dapat diperluas melalui sistem plugin sederhana.
*   **Berbasis Go**: Dibangun dengan Go untuk kinerja dan keandalan.
*   **Google Generative AI**: Didukung oleh model AI generatif canggih dari Google.

## Memulai

### Prasyarat

*   Go 1.18 atau lebih baru
*   Variabel lingkungan `GEMINI_API_KEY` yang valid

### Instalasi

1.  Kloning repositori:

    ```bash
    git clone https://github.com/velora-id/velora.git
    ```

2.  Arahkan ke direktori proyek:

    ```bash
    cd velora
    ```

3.  Instal dependensi:

    ```bash
    go mod tidy
    ```

## Penggunaan

Untuk menjalankan agen Velora, gunakan perintah berikut:

```bash
go run main.go
```

## Berkontribusi

Kontribusi sangat diterima! Jangan ragu untuk mengirimkan *pull request* atau membuka *issue*.
