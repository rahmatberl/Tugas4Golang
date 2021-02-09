package main

import (
	"database/sql"
	"encoding/xml"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

type mahasiswa struct {
	Idmahasiswa string `json:"id_mahasiswa"`
	Nama        string `json:"nama"`
	Alamat      struct {
		Jalan     string `json:"jalan"`
		Kelurahan string `json:"kelurahan"`
		Kecamatan string `json:"kecamatan"`
		Kabupaten string `json:"kabupaten"`
		Provinsi  string `json:"provinsi"`
	} `json:"alamat"`
	Fakultas string  `json:"fakultas"`
	Jurusan  string  `json:"jurusan"`
	Nilai    []nilai `json:"Nilai"`
}

type nilai struct {
	Idmahasiswa string  `json:"id_mahasiswa"`
	Idmatkul    string  `json:"id_matkul"`
	Mkuliah     string  `json:"m_kuliah"`
	Nilai       float32 `json:"nilai"`
	Semester    int8    `json:"semester"`
}

func getMahasiswa(w http.ResponseWriter, r *http.Request) {

	var mhs mahasiswa
	var ni nilai
	params := mux.Vars(r)

	sql := `SELECT
				id_mahasiswa,
				IFNULL(nama,'') nama,
				IFNULL(jalan,'') jalan,
				IFNULL(kelurahan,'') kelurahan,
				IFNULL(kecamatan,'') kecamatan,
				IFNULL(kabupaten,'') kabupaten,
				IFNULL(provinsi,'') provinsi,
				IFNULL(fakultas,'') fakultas,
				IFNULL(jurusan,'') jurusan				
				FROM mahasiswa WHERE id_mahasiswa IN (?)`

	result, err := db.Query(sql, params["id"])

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		err := result.Scan(&mhs.Idmahasiswa, &mhs.Nama, &mhs.Alamat.Jalan, &mhs.Alamat.Kelurahan, &mhs.Alamat.Kecamatan, &mhs.Alamat.Kabupaten, &mhs.Alamat.Provinsi, &mhs.Fakultas, &mhs.Jurusan)

		if err != nil {
			panic(err.Error())
		}

		sqlNilai := `SELECT
						id_mahasiswa		
						, matkul.id_matkul
						, matkul.m_kuliah
						, nilai
						, semester
					FROM
						nilai INNER JOIN matkul
						ON (nilai.id_matkul = matkul.id_matkul)
						WHERE id_mahasiswa = ?`

		id_mahasiswa := &mhs.Idmahasiswa

		resultDetail, errDet := db.Query(sqlNilai, *id_mahasiswa)

		defer resultDetail.Close()

		if errDet != nil {
			panic(err.Error())
		}

		for resultDetail.Next() {

			err := resultDetail.Scan(&ni.Idmahasiswa, &ni.Idmatkul, &ni.Mkuliah, &ni.Nilai, &ni.Semester)

			if err != nil {
				panic(err.Error())
			}

			mhs.Nilai = append(mhs.Nilai, ni)

		}

	}
	//header for exml
	w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"))
	xml.NewEncoder(w).Encode(mhs)
}
func main() {

	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_akademik")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/mahasiswa/{id}", getMahasiswa).Methods("GET")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}
