package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ===================== STRUCT =====================

type Bahan struct {
	Nama   string
	Jumlah string
}

type Resep struct {
	ID           int
	Judul        string
	Kategori     string
	BahanUtama   string
	Bahan        []Bahan
	Langkah      []string
	DurasiMasak  int // dalam menit
	SeringDicari int
}

// ===================== DATA GLOBAL =====================

var daftarResep []Resep
var nextID int = 1
var scanner = bufio.NewScanner(os.Stdin)

// ===================== HELPER INPUT =====================

func input(prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func inputInt(prompt string) int {
	for {
		s := input(prompt)
		n, err := strconv.Atoi(s)
		if err == nil {
			return n
		}
		fmt.Println("  [!] Masukkan angka yang valid.")
	}
}

func cariIndexByID(id int) int {
	for i, r := range daftarResep {
		if r.ID == id {
			return i
		}
	}
	return -1
}

// ===================== CRUD =====================

func tambahResep() {
	fmt.Println("\n──────────── TAMBAH RESEP ────────────")
	var r Resep
	r.ID = nextID
	nextID++

	r.Judul = input("Judul resep       : ")
	r.Kategori = input("Kategori          : ")
	r.BahanUtama = input("Bahan utama       : ")
	r.DurasiMasak = inputInt("Durasi masak (mnt): ")

	// Input bahan
	jmlBahan := inputInt("Jumlah bahan      : ")
	for i := 1; i <= jmlBahan; i++ {
		var b Bahan
		b.Nama = input(fmt.Sprintf("  Bahan ke-%d nama  : ", i))
		b.Jumlah = input(fmt.Sprintf("  Bahan ke-%d jumlah: ", i))
		r.Bahan = append(r.Bahan, b)
	}

	// Input langkah
	jmlLangkah := inputInt("Jumlah langkah    : ")
	for i := 1; i <= jmlLangkah; i++ {
		l := input(fmt.Sprintf("  Langkah ke-%d: ", i))
		r.Langkah = append(r.Langkah, l)
	}

	daftarResep = append(daftarResep, r)
	fmt.Printf("\n  [✓] Resep \"%s\" berhasil ditambahkan (ID: %d).\n", r.Judul, r.ID)
}

func tampilkanResep(r Resep) {
	fmt.Printf("\n  ID          : %d\n", r.ID)
	fmt.Printf("  Judul       : %s\n", r.Judul)
	fmt.Printf("  Kategori    : %s\n", r.Kategori)
	fmt.Printf("  Bahan Utama : %s\n", r.BahanUtama)
	fmt.Printf("  Durasi      : %d menit\n", r.DurasiMasak)
	fmt.Println("  Bahan-bahan :")
	for i, b := range r.Bahan {
		fmt.Printf("    %d. %s - %s\n", i+1, b.Nama, b.Jumlah)
	}
	fmt.Println("  Langkah-langkah:")
	for i, l := range r.Langkah {
		fmt.Printf("    %d. %s\n", i+1, l)
	}
	fmt.Printf("  Sering Dicari: %d kali\n", r.SeringDicari)
}

func lihatSemuaResep() {
	fmt.Println("\n──────────── DAFTAR SEMUA RESEP ────────────")
	if len(daftarResep) == 0 {
		fmt.Println("  [!] Belum ada resep.")
		return
	}
	for _, r := range daftarResep {
		fmt.Printf("  [%d] %-25s | Kategori: %-15s | Durasi: %d mnt\n",
			r.ID, r.Judul, r.Kategori, r.DurasiMasak)
	}
}

func ubahResep() {
	fmt.Println("\n──────────── UBAH RESEP ────────────")
	lihatSemuaResep()
	if len(daftarResep) == 0 {
		return
	}
	id := inputInt("\nMasukkan ID resep yang ingin diubah: ")
	idx := cariIndexByID(id)
	if idx == -1 {
		fmt.Println("  [!] ID tidak ditemukan.")
		return
	}

	r := &daftarResep[idx]
	fmt.Printf("\n  Mengubah resep: %s\n", r.Judul)
	fmt.Println("  (Kosongkan untuk tidak mengubah)")

	if v := input("  Judul baru       : "); v != "" {
		r.Judul = v
	}
	if v := input("  Kategori baru    : "); v != "" {
		r.Kategori = v
	}
	if v := input("  Bahan utama baru : "); v != "" {
		r.BahanUtama = v
	}
	if v := input("  Durasi baru (mnt): "); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			r.DurasiMasak = n
		}
	}

	fmt.Println("  [✓] Resep berhasil diubah.")
}

func hapusResep() {
	fmt.Println("\n──────────── HAPUS RESEP ────────────")
	lihatSemuaResep()
	if len(daftarResep) == 0 {
		return
	}
	id := inputInt("\nMasukkan ID resep yang ingin dihapus: ")
	idx := cariIndexByID(id)
	if idx == -1 {
		fmt.Println("  [!] ID tidak ditemukan.")
		return
	}

	judul := daftarResep[idx].Judul
	daftarResep = append(daftarResep[:idx], daftarResep[idx+1:]...)
	fmt.Printf("  [✓] Resep \"%s\" berhasil dihapus.\n", judul)
}

// ===================== SEARCH =====================

func sequentialSearch() {
	fmt.Println("\n──────────── SEQUENTIAL SEARCH ────────────")
	if len(daftarResep) == 0 {
		fmt.Println("  [!] Belum ada resep.")
		return
	}
	kata := strings.ToLower(input("Cari berdasarkan bahan utama: "))
	ditemukan := false
	for i := range daftarResep {
		if strings.Contains(strings.ToLower(daftarResep[i].BahanUtama), kata) {
			tampilkanResep(daftarResep[i])
			daftarResep[i].SeringDicari++
			ditemukan = true
		}
	}
	if !ditemukan {
		fmt.Println("  [!] Resep dengan bahan utama tersebut tidak ditemukan.")
	}
}

// Binary search mengharuskan data terurut berdasarkan BahanUtama
func binarySearch() {
	fmt.Println("\n──────────── BINARY SEARCH ────────────")
	if len(daftarResep) == 0 {
		fmt.Println("  [!] Belum ada resep.")
		return
	}

	// Buat salinan terurut berdasarkan BahanUtama (tidak mengubah data asli)
	sorted := make([]Resep, len(daftarResep))
	copy(sorted, daftarResep)
	// Urutkan dengan insertion sort (bahan utama abjad)
	for i := 1; i < len(sorted); i++ {
		key := sorted[i]
		j := i - 1
		for j >= 0 && strings.ToLower(sorted[j].BahanUtama) > strings.ToLower(key.BahanUtama) {
			sorted[j+1] = sorted[j]
			j--
		}
		sorted[j+1] = key
	}

	kata := strings.ToLower(input("Cari berdasarkan bahan utama (exact): "))

	lo, hi := 0, len(sorted)-1
	ditemukan := false
	for lo <= hi {
		mid := (lo + hi) / 2
		cmp := strings.ToLower(sorted[mid].BahanUtama)
		if cmp == kata {
			// Cari semua yang cocok di sekitar mid
			// kiri
			l := mid
			for l > 0 && strings.ToLower(sorted[l-1].BahanUtama) == kata {
				l--
			}
			// kanan
			r := mid
			for r < len(sorted)-1 && strings.ToLower(sorted[r+1].BahanUtama) == kata {
				r++
			}
			for k := l; k <= r; k++ {
				tampilkanResep(sorted[k])
				// update SeringDicari di data asli
				for idx := range daftarResep {
					if daftarResep[idx].ID == sorted[k].ID {
						daftarResep[idx].SeringDicari++
					}
				}
			}
			ditemukan = true
			break
		} else if cmp < kata {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	if !ditemukan {
		fmt.Println("  [!] Resep tidak ditemukan.")
	}
}

// ===================== SORT =====================

func selectionSortDurasi() {
	if len(daftarResep) == 0 {
		fmt.Println("  [!] Belum ada resep.")
		return
	}
	sorted := make([]Resep, len(daftarResep))
	copy(sorted, daftarResep)

	n := len(sorted)
	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if sorted[j].DurasiMasak < sorted[minIdx].DurasiMasak {
				minIdx = j
			}
		}
		sorted[i], sorted[minIdx] = sorted[minIdx], sorted[i]
	}

	fmt.Println("\n  Resep diurutkan berdasarkan durasi masak (tercepat → terlama):")
	for _, r := range sorted {
		fmt.Printf("  [%d] %-25s | %d menit\n", r.ID, r.Judul, r.DurasiMasak)
	}
}

func insertionSortAbjad() {
	if len(daftarResep) == 0 {
		fmt.Println("  [!] Belum ada resep.")
		return
	}
	sorted := make([]Resep, len(daftarResep))
	copy(sorted, daftarResep)

	for i := 1; i < len(sorted); i++ {
		key := sorted[i]
		j := i - 1
		for j >= 0 && strings.ToLower(sorted[j].Judul) > strings.ToLower(key.Judul) {
			sorted[j+1] = sorted[j]
			j--
		}
		sorted[j+1] = key
	}

	fmt.Println("\n  Resep diurutkan berdasarkan judul (A → Z):")
	for _, r := range sorted {
		fmt.Printf("  [%d] %-25s | %s\n", r.ID, r.Judul, r.Kategori)
	}
}

func menuUrutan() {
	fmt.Println("\n──────────── URUTKAN RESEP ────────────")
	fmt.Println("  1. Selection Sort  — berdasarkan durasi masak")
	fmt.Println("  2. Insertion Sort  — berdasarkan judul (abjad)")
	pilihan := input("\nPilih metode pengurutan: ")
	switch pilihan {
	case "1":
		selectionSortDurasi()
	case "2":
		insertionSortAbjad()
	default:
		fmt.Println("  [!] Pilihan tidak valid.")
	}
}

// ===================== STATISTIK =====================

func statistik() {
	fmt.Println("\n──────────── STATISTIK ────────────")
	if len(daftarResep) == 0 {
		fmt.Println("  [!] Belum ada resep.")
		return
	}

	// Hitung jumlah resep per kategori
	katMap := make(map[string]int)
	for _, r := range daftarResep {
		katMap[r.Kategori]++
	}

	fmt.Println("\n  Jumlah Resep per Kategori:")
	for kat, jml := range katMap {
		fmt.Printf("    %-20s : %d resep\n", kat, jml)
	}

	// Daftar menu paling sering dicari (top 5)
	sorted := make([]Resep, len(daftarResep))
	copy(sorted, daftarResep)

	// Selection sort berdasarkan SeringDicari (descending)
	n := len(sorted)
	for i := 0; i < n-1; i++ {
		maxIdx := i
		for j := i + 1; j < n; j++ {
			if sorted[j].SeringDicari > sorted[maxIdx].SeringDicari {
				maxIdx = j
			}
		}
		sorted[i], sorted[maxIdx] = sorted[maxIdx], sorted[i]
	}

	fmt.Println("\n  Top 5 Menu Paling Sering Dicari:")
	limit := 5
	if len(sorted) < limit {
		limit = len(sorted)
	}
	for i := 0; i < limit; i++ {
		fmt.Printf("    %d. %-25s — %d kali dicari\n",
			i+1, sorted[i].Judul, sorted[i].SeringDicari)
	}
}

// ===================== DATA CONTOH =====================

func muatDataContoh() {
	contoh := []Resep{
		{
			ID:          nextID,
			Judul:       "Ayam Goreng Kremes",
			Kategori:    "Ayam",
			BahanUtama:  "ayam",
			DurasiMasak: 45,
			Bahan: []Bahan{
				{"Ayam potong", "1 ekor"},
				{"Bawang putih", "5 siung"},
				{"Kunyit", "1 ruas"},
			},
			Langkah:      []string{"Marinasi ayam", "Goreng hingga keemasan", "Sajikan dengan kremes"},
			SeringDicari: 3,
		},
		{
			ID:          nextID + 1,
			Judul:       "Nasi Goreng Spesial",
			Kategori:    "Nasi",
			BahanUtama:  "nasi",
			DurasiMasak: 20,
			Bahan: []Bahan{
				{"Nasi putih", "2 piring"},
				{"Telur", "2 butir"},
				{"Kecap manis", "2 sdm"},
			},
			Langkah:      []string{"Tumis bumbu", "Masukkan nasi", "Tambahkan kecap"},
			SeringDicari: 7,
		},
		{
			ID:          nextID + 2,
			Judul:       "Soto Ayam Bening",
			Kategori:    "Sup",
			BahanUtama:  "ayam",
			DurasiMasak: 60,
			Bahan: []Bahan{
				{"Ayam kampung", "1/2 ekor"},
				{"Serai", "2 batang"},
				{"Daun salam", "3 lembar"},
			},
			Langkah:      []string{"Rebus ayam", "Tumis bumbu halus", "Gabungkan dan didihkan"},
			SeringDicari: 5,
		},
		{
			ID:          nextID + 3,
			Judul:       "Tempe Orek Kering",
			Kategori:    "Sayur",
			BahanUtama:  "tempe",
			DurasiMasak: 30,
			Bahan: []Bahan{
				{"Tempe", "1 papan"},
				{"Cabe merah", "5 buah"},
				{"Gula merah", "1 sdm"},
			},
			Langkah:      []string{"Potong tempe", "Goreng tempe", "Orek dengan bumbu"},
			SeringDicari: 2,
		},
		{
			ID:          nextID + 4,
			Judul:       "Bakso Kuah",
			Kategori:    "Sup",
			BahanUtama:  "daging sapi",
			DurasiMasak: 90,
			Bahan: []Bahan{
				{"Daging sapi giling", "500 gr"},
				{"Tepung tapioka", "100 gr"},
				{"Bawang putih", "4 siung"},
			},
			Langkah:      []string{"Buat adonan bakso", "Bentuk bulat", "Rebus dalam kaldu"},
			SeringDicari: 9,
		},
	}
	daftarResep = append(daftarResep, contoh...)
	nextID += 5
}

// ===================== MENU UTAMA =====================

func menuUtama() {
	for {
		fmt.Println("\n╔══════════════════════════════════════╗")
		fmt.Println("║        ResepKu — Menu Utama          ║")
		fmt.Println("╠══════════════════════════════════════╣")
		fmt.Println("║  1. Tambah Resep                     ║")
		fmt.Println("║  2. Lihat Semua Resep                ║")
		fmt.Println("║  3. Ubah Resep                       ║")
		fmt.Println("║  4. Hapus Resep                      ║")
		fmt.Println("║  5. Cari Resep (Sequential Search)   ║")
		fmt.Println("║  6. Cari Resep (Binary Search)       ║")
		fmt.Println("║  7. Urutkan Resep                    ║")
		fmt.Println("║  8. Statistik                        ║")
		fmt.Println("║  9. Lihat Detail Resep               ║")
		fmt.Println("║  0. Keluar                           ║")
		fmt.Println("╚══════════════════════════════════════╝")
		pilihan := input("Pilih menu: ")

		switch pilihan {
		case "1":
			tambahResep()
		case "2":
			lihatSemuaResep()
		case "3":
			ubahResep()
		case "4":
			hapusResep()
		case "5":
			sequentialSearch()
		case "6":
			binarySearch()
		case "7":
			menuUrutan()
		case "8":
			statistik()
		case "9":
			lihatSemuaResep()
			if len(daftarResep) > 0 {
				id := inputInt("\nMasukkan ID resep: ")
				idx := cariIndexByID(id)
				if idx == -1 {
					fmt.Println("  [!] ID tidak ditemukan.")
				} else {
					tampilkanResep(daftarResep[idx])
				}
			}
		case "0":
			fmt.Println("\n  Terima kasih telah menggunakan ResepKu. Selamat memasak! 🍳")
			return
		default:
			fmt.Println("  [!] Pilihan tidak valid.")
		}
	}
}

// ===================== MAIN =====================

func main() {
	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║     Selamat Datang di ResepKu!       ║")
	fmt.Println("║  Aplikasi Manajemen Resep Kuliner    ║")
	fmt.Println("╚══════════════════════════════════════╝")

	muatDataContoh()
	fmt.Println("\n  [✓] 5 data resep contoh telah dimuat.")

	menuUtama()
}