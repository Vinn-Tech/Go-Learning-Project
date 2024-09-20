package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	var err error

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	DBusername := os.Getenv("DB_USERNAME")
	DBpassword := os.Getenv("DB_PASSWORD")

	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/db_bts", DBusername, DBpassword))
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	defer db.Close()

	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/dataBTS", btsHandler)

	fmt.Println("Server started on port http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "API is running"})
}

func btsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var bts_lists []struct {
		ID_BTS      int    `json:"id"`
		KEL_DES     string `json:"kel_des"`
		KAB_KOTA    string `json:"kab_kota"`
		PROV        string `json:"prov"`
		LUAS_DESA   string `json:"luas_desa"`
		TOTAL_NE    string `json:"total_ne"`
		RASIO_NE    string `json:"rasio_ne"`
		TOTAL_NE_4G string `json:"total_ne_4g"`
		RASIO_NE_4G string `json:"rasio_ne_4g"`
		KEC         string `json:"kec"`
	}

	rows, err := db.Query("SELECT * FROM tb_ne")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var bts_list struct {
			ID_BTS      int    `json:"id"`
			KEL_DES     string `json:"kel_des"`
			KAB_KOTA    string `json:"kab_kota"`
			PROV        string `json:"prov"`
			LUAS_DESA   string `json:"luas_desa"`
			TOTAL_NE    string `json:"total_ne"`
			RASIO_NE    string `json:"rasio_ne"`
			TOTAL_NE_4G string `json:"total_ne_4g"`
			RASIO_NE_4G string `json:"rasio_ne_4g"`
			KEC         string `json:"kec"`
		}
		err := rows.Scan(
			&bts_list.ID_BTS,
			&bts_list.KEL_DES,
			&bts_list.KAB_KOTA,
			&bts_list.PROV,
			&bts_list.LUAS_DESA,
			&bts_list.TOTAL_NE,
			&bts_list.RASIO_NE,
			&bts_list.TOTAL_NE_4G,
			&bts_list.RASIO_NE_4G,
			&bts_list.KEC,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bts_lists = append(bts_lists, bts_list)
	}

	json.NewEncoder(w).Encode(bts_lists)
}
