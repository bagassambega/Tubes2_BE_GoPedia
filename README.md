# Tugas Besar Strategi Algoritma 2
### Pemanfaatan Algoritma IDS dan BFS dalam Permainan WikiRace

### Kelompok GoPedia:
| No. | Nama                     |   NIM    |
|:---:|:-------------------------|:--------:|
| 1.  | Rafiki Prawhira Harianto | 13522065 |
| 2.  | Bagas Sambega Rosyada    | 13522071 |
| 3.  | Abdullah Mubarak         | 13522101 |

## Daftar Isi
1. [Deskripsi Program](#deskripsi-program)
2. [Implementasi Algoritma](#implementasi-algoritma)
3. [Cara Penggunaan Program](#cara-penggunaan-program)

## Deskripsi Program
Program ini adalah program untuk menyelesaikan permainan Wikirace menggunakan algoritma IDS dan BFS. Permainan Wikirace adalah permainan untuk mencari 
jalur terpendek dari satu artikel Wikipedia ke artikel lainnya. Program ini akan mencari jalur terpendek dari artikel awal ke artikel tujuan dengan melakukan 
_scraping_ pada artikel-artikel yang dikunjungi, dan mengunjungi artikel-artikel tersebut untuk mencari artikel tujuan. Program akan menghasilkan jalur terpendek dari artikel awal ke artikel tujuan.

_Repository_ ini adalah bagian _backend_ dari program yang berisi _script_ dalam bahasa Go untuk menjalankan fungsi pemrosesan permainan Wikirace dan mengirimkan datanya melalui API ke 
<a href="github.com/bagassambega/Tubes2_FE_GoPedia">_frontend_</a> yang dibuat menggunakan _framework_ ReactJS. Kedua _repository_ perlu dijalankan bersamaan untuk
menjalankan program Wikirace. Link kedua repository:
1. <a href="github.com/bagassambega/Tubes2_BE_GoPedia">Backend</a>
2. <a href="github.com/bagassambega/Tubes2_FE_GoPedia">Frontend</a>

## Implementasi Algoritma
Program ini menggunakan dua algoritma untuk menyelesaikan permainan Wikirace, yaitu:
1. Algoritma IDS (_Iterative Deepening Search_): Algoritma ini adalah algoritma _search_ yang melakukan _depth-first search_ dengan level _depth_ yang bertambah secara iteratif. Implementasi algoritma ini terdapat pada file src/IDS.go, yang berisi fungsi utama IDS yang akan melakukan pemanggilan fungsi DLS/_Depth Limited Search_ sampai dengan level tertentu. Jika pada level tersebut tidak ditemukan solusi, maka level kedalaman akan ditingkatkan dan fungsi DLS akan
dipanggil kembali dengan level kedalaman yang baru.
2. Algoritma BFS (_Breadth-First Search_): Algoritma ini adalah algoritma _search_ yang melakukan pencarian pada satu level terlebih dahulu secara keseluruhan sebelum mencari di level kedalaman berikutnya.
Implementasi algoritma ini terdapat pada file src/BFS.go, yang berisi fungsi utama BFS yang akan melakukan pencarian dengan melakukan _scraping_ dan menyimpan seluruh tautan pada suatu level ke dalam _queue_, dan mengunjungi setiap artikel pada _queue_
tersebut untuk mencari artikel tujuan. Jika tidak ditemukan artikel tujuan pada level tersebut, fungsi akan melakukan _scraping_ pada level kedalaman selanjutnya dan mengunjunginya secara keseluruhan.

## Cara Penggunaan Program
Program memerlukan <a href="github.com/bagassambega/Tubes2_FE_GoPedia">_frontend_</a> untuk menjalankan program Wikirace. Langkah instalasi terdapat pada _repository_ _frontend_ tersebut.
### Requirement
1. Go terinstal di perangkat
2. _Framework_ Gin dan Gocolly
3. Docker Desktop

### Instalasi
1. Clone _repository_ ini
```bash
git clone https://github.com/bagassambega/Tubes2_BE_GoPedia.git
```
2. Masuk ke direktori _repository_ yang telah di-_clone_ dan ke folder src
```bash
cd Tubes2_BE_GoPedia/src
```
3. Jalankan _command_ berikut untuk memastikan seluruh _dependency_ terinstal
```bash
go mod tidy
```
4. Jalankan Docker Desktop atau Docker, lalu build _docker image_ dari _Dockerfile_ yang telah disediakan
```bash
docker build -t gopedia-backend .
```
5. Jalankan _docker container_ dari _docker image_ yang telah dibuat
```bash
docker run -p 8080:8080 gopedia-backend
```
6. Untuk menghentikan program, jalankan _command_ berikut. Lihat _container_id_ dengan menjalankan _command_ `docker ps` atau pada Docker Desktop
```bash
docker stop [container_id]
```


Setelah program _backend_ berjalan dan langkah-langkah menjalankan _frontend_ selesai, program dapat diakses pada _browser_ dengan membuka alamat http://localhost:5173/
dan pengguna dapat memasukkan artikel awal dan artikel tujuan untuk mencari jalur terpendek dari artikel awal ke artikel tujuan.



**_NOTE:_**
Setelah langkah-langkah di atas dilakukan, bagian _backend_ Wikirace akan berjalan pada _port_ 8080. Data sudah dapat diakses melalui _browser_ dengan membuka alamat http://localhost:8080/gopedia/?method=[metode]&source=[awal]&target=[akhir], dengan mengganti [metode] menjadi IDS/BFS, [awal] menjadi artikel awal, dan [akhir] menjadi artikel tujuan. Contoh: http://localhost:8080/gopedia/?method=IDS&source=Indonesia&target=Jepang. 
Jika bagian _frontend_ tidak berjalan dengan baik, data dapat diakses langsung dengan langkah di atas.

**_NOTE:_**
Jika _build docker_ gagal atau menjalankan dengan _docker_ tidak berhasil, program dapat dijalankan dengan menjalankan _command_ berikut pada folder src:
```bash
go run .
```