package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ===================== ANSI COLORS =====================

const (
	Reset = "\033[0m"
	Bold  = "\033[1m"
	Dim   = "\033[2m"

	FgBlack   = "\033[30m"
	FgRed     = "\033[31m"
	FgGreen   = "\033[32m"
	FgYellow  = "\033[33m"
	FgBlue    = "\033[34m"
	FgMagenta = "\033[35m"
	FgCyan    = "\033[36m"
	FgWhite   = "\033[37m"

	BgGreen = "\033[42m"
	BgBlue  = "\033[44m"
	BgCyan  = "\033[46m"

	// Bright variants
	FgBrightGreen  = "\033[92m"
	FgBrightYellow = "\033[93m"
	FgBrightCyan   = "\033[96m"
	FgBrightWhite  = "\033[97m"
)

func clr(color, text string) string { return color + text + Reset }
func bold(text string) string       { return Bold + text + Reset }
func dim(text string) string        { return Dim + text + Reset }

// ===================== TABLE DRAWING =====================

// Box drawing characters
const (
	TL  = "╔"
	TR  = "╗"
	BL  = "╚"
	BR  = "╝"
	H   = "═"
	V   = "║"
	ML  = "╠"
	MR  = "╣"
	MT  = "╦"
	MB  = "╩"
	MC  = "╬"
	SH  = "─"
	SV  = "│"
	STL = "┌"
	STR = "┐"
	SBL = "└"
	SBR = "┘"
	SML = "├"
	SMR = "┤"
	SMT = "┬"
	SMB = "┴"
	SMC = "┼"
)

func repeatStr(s string, n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteString(s)
	}
	return b.String()
}

func padRight(s string, width int) string {
	r := []rune(s)
	if len(r) >= width {
		if len(r) > width-1 {
			return string(r[:width-3]) + "..."
		}
		return s
	}
	return s + strings.Repeat(" ", width-len(r))
}

func center(s string, width int) string {
	r := []rune(s)
	if len(r) >= width {
		return s
	}
	pad := width - len(r)
	left := pad / 2
	right := pad - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

// printDoubleBox draws a double-line box (for menus/headers)
func printDoubleBox(lines []string, width int) {
	top := TL + repeatStr(H, width) + TR
	bottom := BL + repeatStr(H, width) + BR
	sep := ML + repeatStr(H, width) + MR

	fmt.Println(clr(FgCyan, top))
	for i, line := range lines {
		if line == "---" {
			fmt.Println(clr(FgCyan, sep))
		} else {
			inner := center(line, width)
			fmt.Printf("%s%s%s\n", clr(FgCyan, V), inner, clr(FgCyan, V))
			_ = i
		}
	}
	fmt.Println(clr(FgCyan, bottom))
}

// printTable draws a simple table with single-line borders
// cols: column widths; headers: header labels; rows: data rows
func printTable(cols []int, headers []string, rows [][]string) {
	total := len(cols)

	// Build separator line
	buildSep := func(left, mid, right, h string) string {
		var b strings.Builder
		b.WriteString(left)
		for i, w := range cols {
			b.WriteString(repeatStr(h, w+2))
			if i < total-1 {
				b.WriteString(mid)
			}
		}
		b.WriteString(right)
		return b.String()
	}

	topLine := buildSep(STL, SMT, STR, SH)
	midLine := buildSep(SML, SMC, SMR, SH)
	botLine := buildSep(SBL, SMB, SBR, SH)

	fmt.Println(clr(FgBrightCyan, topLine))

	// Header row
	fmt.Print(clr(FgBrightCyan, SV))
	for i, h := range headers {
		cell := " " + padRight(h, cols[i]) + " "
		fmt.Print(clr(FgBrightYellow, Bold+cell+Reset) + clr(FgBrightCyan, SV))
	}
	fmt.Println()
	fmt.Println(clr(FgBrightCyan, midLine))

	// Data rows
	for ri, row := range rows {
		rowColor := FgWhite
		if ri%2 == 1 {
			rowColor = Dim + FgWhite
		}
		fmt.Print(clr(FgBrightCyan, SV))
		for i, cell := range row {
			if i >= len(cols) {
				break
			}
			padded := " " + padRight(cell, cols[i]) + " "
			if i == 0 {
				fmt.Print(clr(FgBrightGreen, padded))
			} else {
				fmt.Print(rowColor + padded + Reset)
			}
			fmt.Print(clr(FgBrightCyan, SV))
		}
		fmt.Println()
	}
	fmt.Println(clr(FgBrightCyan, botLine))
}

// printInfoBox draws a single info/highlight box
func printInfoBox(label, value, color string) {
	line := fmt.Sprintf("  %s%-20s%s %s", Bold, label+":", Reset, clr(color, value))
	fmt.Println(line)
}

func printSuccess(msg string) {
	fmt.Printf("\n  %s %s\n", clr(FgBrightGreen, "✔"), clr(FgBrightGreen, msg))
}

func printError(msg string) {
	fmt.Printf("\n  %s %s\n", clr(FgRed, "✘"), clr(FgRed, msg))
}

func printSectionTitle(title string) {
	bar := repeatStr("─", 44)
	fmt.Printf("\n%s\n  %s\n%s\n", clr(FgCyan, bar), clr(FgBrightCyan, bold(title)), clr(FgCyan, bar))
}

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
	DurasiMasak  int
	SeringDicari int
}

// ===================== DATA GLOBAL =====================

var daftarResep []Resep
var nextID int = 1
var scanner = bufio.NewScanner(os.Stdin)

// ===================== HELPER INPUT =====================

func input(prompt string) string {
	fmt.Print(clr(FgBrightYellow, "  ❯ ") + prompt)
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
		printError("Masukkan angka yang valid.")
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
	printSectionTitle("TAMBAH RESEP BARU")
	var r Resep
	r.ID = nextID
	nextID++

	r.Judul = input("Judul resep        : ")
	r.Kategori = input("Kategori           : ")
	r.BahanUtama = input("Bahan utama        : ")
	r.DurasiMasak = inputInt("Durasi masak (mnt) : ")

	jmlBahan := inputInt("Jumlah bahan       : ")
	for i := 1; i <= jmlBahan; i++ {
		var b Bahan
		b.Nama = input(fmt.Sprintf("  Bahan ke-%d nama   : ", i))
		b.Jumlah = input(fmt.Sprintf("  Bahan ke-%d jumlah : ", i))
		r.Bahan = append(r.Bahan, b)
	}

	jmlLangkah := inputInt("Jumlah langkah     : ")
	for i := 1; i <= jmlLangkah; i++ {
		l := input(fmt.Sprintf("  Langkah ke-%d      : ", i))
		r.Langkah = append(r.Langkah, l)
	}

	daftarResep = append(daftarResep, r)
	printSuccess(fmt.Sprintf("Resep \"%s\" berhasil ditambahkan (ID: %d)", r.Judul, r.ID))
}

func tampilkanResep(r Resep) {
	printSectionTitle(fmt.Sprintf("DETAIL RESEP — %s", strings.ToUpper(r.Judul)))
	printInfoBox("ID", fmt.Sprintf("%d", r.ID), FgBrightCyan)
	printInfoBox("Judul", r.Judul, FgBrightWhite)
	printInfoBox("Kategori", r.Kategori, FgBrightYellow)
	printInfoBox("Bahan Utama", r.BahanUtama, FgBrightGreen)
	printInfoBox("Durasi", fmt.Sprintf("%d menit", r.DurasiMasak), FgBrightCyan)
	printInfoBox("Dicari", fmt.Sprintf("%d kali", r.SeringDicari), FgMagenta)

	fmt.Printf("\n  %s\n", clr(FgBrightYellow, bold("Bahan-bahan:")))
	for i, b := range r.Bahan {
		fmt.Printf("    %s %s — %s\n",
			clr(FgCyan, fmt.Sprintf("%d.", i+1)),
			clr(FgWhite, b.Nama),
			clr(FgBrightGreen, b.Jumlah))
	}

	fmt.Printf("\n  %s\n", clr(FgBrightYellow, bold("Langkah-langkah:")))
	for i, l := range r.Langkah {
		fmt.Printf("    %s %s\n",
			clr(FgCyan, fmt.Sprintf("%d.", i+1)),
			clr(FgWhite, l))
	}
}

func lihatSemuaResep() {
	printSectionTitle("DAFTAR SEMUA RESEP")
	if len(daftarResep) == 0 {
		printError("Belum ada resep.")
		return
	}
	cols := []int{4, 26, 16, 10}
	headers := []string{"ID", "Judul", "Kategori", "Durasi"}
	rows := [][]string{}
	for _, r := range daftarResep {
		rows = append(rows, []string{
			fmt.Sprintf("%d", r.ID),
			r.Judul,
			r.Kategori,
			fmt.Sprintf("%d mnt", r.DurasiMasak),
		})
	}
	printTable(cols, headers, rows)
}

func ubahResep() {
	printSectionTitle("UBAH RESEP")
	lihatSemuaResep()
	if len(daftarResep) == 0 {
		return
	}

	id := inputInt("ID resep yang ingin diubah : ")
	idx := cariIndexByID(id)
	if idx == -1 {
		printError("ID tidak ditemukan.")
		return
	}

	r := &daftarResep[idx]
	fmt.Printf("\n  %s %s\n", clr(FgBrightCyan, "✦"), clr(FgBrightWhite, bold("Mengubah: "+r.Judul)))
	fmt.Println(dim("  (Kosongkan untuk tidak mengubah)"))

	if v := input("Judul baru        : "); v != "" {
		r.Judul = v
	}
	if v := input("Kategori baru     : "); v != "" {
		r.Kategori = v
	}
	if v := input("Bahan utama baru  : "); v != "" {
		r.BahanUtama = v
	}
	if v := input("Durasi baru (mnt) : "); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			r.DurasiMasak = n
		}
	}
	printSuccess("Resep berhasil diubah.")
}

func hapusResep() {
	printSectionTitle("HAPUS RESEP")
	lihatSemuaResep()
	if len(daftarResep) == 0 {
		return
	}

	id := inputInt("ID resep yang ingin dihapus : ")
	idx := cariIndexByID(id)
	if idx == -1 {
		printError("ID tidak ditemukan.")
		return
	}

	judul := daftarResep[idx].Judul
	konfirmasi := input(fmt.Sprintf("Hapus \"%s\"? (y/n) : ", judul))
	if strings.ToLower(konfirmasi) != "y" {
		fmt.Println(clr(FgYellow, "  Penghapusan dibatalkan."))
		return
	}
	daftarResep = append(daftarResep[:idx], daftarResep[idx+1:]...)
	printSuccess(fmt.Sprintf("Resep \"%s\" berhasil dihapus.", judul))
}

// ===================== SEARCH =====================

func sequentialSearch() {
	printSectionTitle("SEQUENTIAL SEARCH")
	if len(daftarResep) == 0 {
		printError("Belum ada resep.")
		return
	}

	kata := strings.ToLower(input("Cari bahan utama (partial) : "))
	ditemukan := false
	fmt.Println()
	for i := range daftarResep {
		if strings.Contains(strings.ToLower(daftarResep[i].BahanUtama), kata) {
			tampilkanResep(daftarResep[i])
			daftarResep[i].SeringDicari++
			ditemukan = true
		}
	}
	if !ditemukan {
		printError("Resep dengan bahan utama tersebut tidak ditemukan.")
	}
}

func binarySearch() {
	printSectionTitle("BINARY SEARCH")
	if len(daftarResep) == 0 {
		printError("Belum ada resep.")
		return
	}

	sorted := make([]Resep, len(daftarResep))
	copy(sorted, daftarResep)
	for i := 1; i < len(sorted); i++ {
		key := sorted[i]
		j := i - 1
		for j >= 0 && strings.ToLower(sorted[j].BahanUtama) > strings.ToLower(key.BahanUtama) {
			sorted[j+1] = sorted[j]
			j--
		}
		sorted[j+1] = key
	}

	kata := strings.ToLower(input("Cari bahan utama (exact) : "))
	lo, hi := 0, len(sorted)-1
	ditemukan := false
	for lo <= hi {
		mid := (lo + hi) / 2
		cmp := strings.ToLower(sorted[mid].BahanUtama)
		if cmp == kata {
			l, r2 := mid, mid
			for l > 0 && strings.ToLower(sorted[l-1].BahanUtama) == kata {
				l--
			}
			for r2 < len(sorted)-1 && strings.ToLower(sorted[r2+1].BahanUtama) == kata {
				r2++
			}
			for k := l; k <= r2; k++ {
				tampilkanResep(sorted[k])
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
		printError("Resep tidak ditemukan.")
	}
}

// ===================== SORT =====================

func selectionSortDurasi() {
	if len(daftarResep) == 0 {
		printError("Belum ada resep.")
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
	fmt.Printf("\n  %s\n\n", clr(FgBrightYellow, bold("Diurutkan berdasarkan durasi masak (tercepat → terlama):")))
	cols := []int{4, 26, 16, 10}
	headers := []string{"ID", "Judul", "Kategori", "Durasi"}
	rows := [][]string{}
	for _, r := range sorted {
		rows = append(rows, []string{
			fmt.Sprintf("%d", r.ID), r.Judul, r.Kategori,
			fmt.Sprintf("%d mnt", r.DurasiMasak),
		})
	}
	printTable(cols, headers, rows)
}

func insertionSortAbjad() {
	if len(daftarResep) == 0 {
		printError("Belum ada resep.")
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
	fmt.Printf("\n  %s\n\n", clr(FgBrightYellow, bold("Diurutkan berdasarkan judul (A → Z):")))
	cols := []int{4, 26, 16, 10}
	headers := []string{"ID", "Judul", "Kategori", "Durasi"}
	rows := [][]string{}
	for _, r := range sorted {
		rows = append(rows, []string{
			fmt.Sprintf("%d", r.ID), r.Judul, r.Kategori,
			fmt.Sprintf("%d mnt", r.DurasiMasak),
		})
	}
	printTable(cols, headers, rows)
}

func menuUrutan() {
	printSectionTitle("URUTKAN RESEP")
	fmt.Printf("  %s Selection Sort  — berdasarkan durasi masak\n", clr(FgBrightCyan, "[1]"))
	fmt.Printf("  %s Insertion Sort  — berdasarkan judul (abjad)\n", clr(FgBrightCyan, "[2]"))
	pilihan := input("Pilih metode : ")
	switch pilihan {
	case "1":
		selectionSortDurasi()
	case "2":
		insertionSortAbjad()
	default:
		printError("Pilihan tidak valid.")
	}
}

// ===================== STATISTIK =====================

func statistik() {
	printSectionTitle("STATISTIK RESEPKU")
	if len(daftarResep) == 0 {
		printError("Belum ada resep.")
		return
	}

	// Per kategori
	katMap := make(map[string]int)
	for _, r := range daftarResep {
		katMap[r.Kategori]++
	}

	fmt.Printf("  %s\n\n", clr(FgBrightYellow, bold("Jumlah Resep per Kategori:")))
	cols1 := []int{20, 8}
	headers1 := []string{"Kategori", "Jumlah"}
	rows1 := [][]string{}
	for kat, jml := range katMap {
		rows1 = append(rows1, []string{kat, fmt.Sprintf("%d resep", jml)})
	}
	printTable(cols1, headers1, rows1)

	// Top 5
	sorted := make([]Resep, len(daftarResep))
	copy(sorted, daftarResep)
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

	fmt.Printf("\n  %s\n\n", clr(FgBrightYellow, bold("Top 5 Menu Paling Sering Dicari:")))
	limit := 5
	if len(sorted) < limit {
		limit = len(sorted)
	}
	cols2 := []int{3, 26, 10}
	headers2 := []string{"#", "Judul", "Dicari"}
	rows2 := [][]string{}
	for i := 0; i < limit; i++ {
		rows2 = append(rows2, []string{
			fmt.Sprintf("%d", i+1),
			sorted[i].Judul,
			fmt.Sprintf("%d kali", sorted[i].SeringDicari),
		})
	}
	printTable(cols2, headers2, rows2)
}

// ===================== DATA CONTOH =====================

func muatDataContoh() {
	contoh := []Resep{
		{ID: nextID, Judul: "Ayam Goreng Kremes", Kategori: "Ayam", BahanUtama: "ayam", DurasiMasak: 45, SeringDicari: 3,
			Bahan:   []Bahan{{"Ayam potong", "1 ekor"}, {"Bawang putih", "5 siung"}, {"Kunyit", "1 ruas"}},
			Langkah: []string{"Marinasi ayam", "Goreng hingga keemasan", "Sajikan dengan kremes"}},
		{ID: nextID + 1, Judul: "Nasi Goreng Spesial", Kategori: "Nasi", BahanUtama: "nasi", DurasiMasak: 20, SeringDicari: 7,
			Bahan:   []Bahan{{"Nasi putih", "2 piring"}, {"Telur", "2 butir"}, {"Kecap manis", "2 sdm"}},
			Langkah: []string{"Tumis bumbu", "Masukkan nasi", "Tambahkan kecap"}},
		{ID: nextID + 2, Judul: "Soto Ayam Bening", Kategori: "Sup", BahanUtama: "ayam", DurasiMasak: 60, SeringDicari: 5,
			Bahan:   []Bahan{{"Ayam kampung", "1/2 ekor"}, {"Serai", "2 batang"}, {"Daun salam", "3 lembar"}},
			Langkah: []string{"Rebus ayam", "Tumis bumbu halus", "Gabungkan dan didihkan"}},
		{ID: nextID + 3, Judul: "Tempe Orek Kering", Kategori: "Sayur", BahanUtama: "tempe", DurasiMasak: 30, SeringDicari: 2,
			Bahan:   []Bahan{{"Tempe", "1 papan"}, {"Cabe merah", "5 buah"}, {"Gula merah", "1 sdm"}},
			Langkah: []string{"Potong tempe", "Goreng tempe", "Orek dengan bumbu"}},
		{ID: nextID + 4, Judul: "Bakso Kuah", Kategori: "Sup", BahanUtama: "daging sapi", DurasiMasak: 90, SeringDicari: 9,
			Bahan:   []Bahan{{"Daging sapi giling", "500 gr"}, {"Tepung tapioka", "100 gr"}, {"Bawang putih", "4 siung"}},
			Langkah: []string{"Buat adonan bakso", "Bentuk bulat", "Rebus dalam kaldu"}},
	}
	daftarResep = append(daftarResep, contoh...)
	nextID += 5
}

// ===================== MENU UTAMA =====================

func printMenu() {
	lines := []string{
		clr(FgBrightGreen, bold("R E S E P K U")),
		clr(FgBrightCyan, "Aplikasi Manajemen Resep Kuliner"),
		"---",
		clr(FgBrightCyan, "[1]") + "  Tambah Resep",
		clr(FgBrightCyan, "[2]") + "  Lihat Semua Resep",
		clr(FgBrightCyan, "[3]") + "  Ubah Resep",
		clr(FgBrightCyan, "[4]") + "  Hapus Resep",
		"---",
		clr(FgBrightCyan, "[5]") + "  Cari — Sequential Search",
		clr(FgBrightCyan, "[6]") + "  Cari — Binary Search",
		"---",
		clr(FgBrightCyan, "[7]") + "  Urutkan Resep",
		clr(FgBrightCyan, "[8]") + "  Statistik",
		clr(FgBrightCyan, "[9]") + "  Lihat Detail Resep",
		"---",
		clr(FgRed, "[0]") + "  Keluar",
	}
	printDoubleBox(lines, 40)
}

func menuUtama() {
	for {
		printMenu()
		pilihan := input("Pilih menu : ")
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
				id := inputInt("Masukkan ID resep : ")
				idx := cariIndexByID(id)
				if idx == -1 {
					printError("ID tidak ditemukan.")
				} else {
					tampilkanResep(daftarResep[idx])
				}
			}
		case "0":
			fmt.Printf("\n  %s\n\n", clr(FgBrightGreen, bold("Terima kasih telah menggunakan ResepKu. Selamat memasak! 🍳")))
			return
		default:
			printError("Pilihan tidak valid.")
		}
	}
}

// ===================== MAIN =====================

func main() {
	// Clear screen
	fmt.Print("\033[H\033[2J")

	printDoubleBox([]string{
		"",
		clr(FgBrightGreen, bold("  ╭─────────────────────╮  ")),
		clr(FgBrightGreen, bold("  │  🍳  R E S E P K U  │  ")),
		clr(FgBrightGreen, bold("  ╰─────────────────────╯  ")),
		"",
		clr(FgBrightCyan, "Aplikasi Manajemen Resep Kuliner"),
		clr(FgBrightYellow, "Algoritma Pemrograman 2 — 2026"),
		"",
	}, 40)

	muatDataContoh()
	printSuccess("5 data resep contoh telah dimuat.")
	fmt.Println()

	menuUtama()
}
