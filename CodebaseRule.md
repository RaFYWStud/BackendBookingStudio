# 💡 Aturan dan Konvensi Umum (The Rules)

Ini adalah seperangkat aturan yang Wajib Di ikuti oleh developer Go.

## **1. Penamaan (Naming Convention)**

Aturan penamaan di Go sangat penting karena ia menentukan visibilitas (*export* atau *unexport*).

* **Gunakan `camelCase`**: Untuk nama variabel, fungsi, dan parameter yang bersifat lokal (tidak di-*export*). Ini persis seperti yang Anda inginkan.

    ```go
    var httpClient string
    func calculateTotal(price int, quantity int) int { ... }
    ```

* **Gunakan `PascalCase`**: Untuk nama yang perlu di-*export* (dapat diakses dari package lain). Ini berlaku untuk variabel, konstanta, struct, interface, dan fungsi.

    ```go
    // Konstanta yang diekspor
    const DefaultTimeout = 30 

    // Struct yang diekspor
    type User struct {
        FirstName string // Field ini juga diekspor
        lastName  string // Field ini tidak diekspor
    }

    // Fungsi yang diekspor
    func NewClient(apiKey string) *Client { ... }
    ```

* **Nama Package**: Gunakan nama yang pendek, ringkas, dan semua huruf kecil. Hindari `snake_case` atau `kebab-case`.
  * ✅ Bagus: `httpclient`, `controller`, `models`
  * ❌ Kurang bagus: `http_client`, `userController`
* **Interface**: Interface yang hanya memiliki satu metode seringkali diberi nama dengan menambahkan akhiran "-er".

    ```go
    type Reader interface {
        Read(p []byte) (n int, err error)
    }
    ```

* **Akronim**: Perlakukan akronim (seperti URL, API, ID) sebagai satu kata. Tulis semuanya dalam huruf besar jika di awal nama (PascalCase) atau semuanya huruf kecil (camelCase).
  * ✅ Bagus: `apiClient`, `UserID`, `serveHTTP`, `URL`
  * ❌ Kurang bagus: `ApiClient`, `UserId`, `ServeHttp`, `Url`

-----

## ⚙️ Tools untuk Otomatisasi Aturan

### Linting

```bash
golangci-lint run
```

Jalankan secara lokal sebelum commit penting, dan wajib di CI.
