package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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
	http.HandleFunc("/createBTS", createBtsHandler)
	http.HandleFunc("/updateBTS", updateBtsHandler)
	http.HandleFunc("/deleteBTS", deleteBtsHandler)
	http.HandleFunc("/getBTSByID", getBtsByIDHandler)

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

func getBtsByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var bts struct {
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

	err = db.QueryRow("SELECT * FROM tb_ne WHERE ID_BTS = ?", id).Scan(
		&bts.ID_BTS,
		&bts.KEL_DES,
		&bts.KAB_KOTA,
		&bts.PROV,
		&bts.LUAS_DESA,
		&bts.TOTAL_NE,
		&bts.RASIO_NE,
		&bts.TOTAL_NE_4G,
		&bts.RASIO_NE_4G,
		&bts.KEC,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "No BTS found with that ID", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(bts)
}

func createBtsHandler(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("INSERT INTO tb_ne (KEL_DES, KAB_KOTA, PROV, LUAS_DESA, TOTAL_NE, RASIO_NE, TOTAL_NE_4G, RASIO_NE_4G, KEC) VALUES (?,?,?,?,?,?,?,?,?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var input struct {
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

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = stmt.Exec(input.KEL_DES, input.KAB_KOTA, input.PROV, input.LUAS_DESA, input.TOTAL_NE, input.RASIO_NE, input.TOTAL_NE_4G, input.RASIO_NE_4G, input.KEC)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "success",
		"data":    input,
	})
}

func updateBtsHandler(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("UPDATE tb_ne SET KEL_DES=?, KAB_KOTA=?, PROV=?, LUAS_DESA=?, TOTAL_NE=?, RASIO_NE=?, TOTAL_NE_4G=?, RASIO_NE_4G=?, KEC=? WHERE ID_BTS=?")
	if err != nil {
		http.Error(w, "Failed to prepare statement: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var update struct {
		ID_BTS      int64  `json:"id_bts"`
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

	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "Invalid input data: "+err.Error(), http.StatusBadRequest)
		return
	}

	_, err = stmt.Exec(update.KEL_DES, update.KAB_KOTA, update.PROV, update.LUAS_DESA, update.TOTAL_NE, update.RASIO_NE, update.TOTAL_NE_4G, update.RASIO_NE_4G, update.KEC, update.ID_BTS)
	if err != nil {
		http.Error(w, "Failed to execute update: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "success",
		"updated_data": update,
	})
}

func deleteBtsHandler(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("DELETE FROM tb_ne WHERE ID_BTS=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var remove struct {
		ID_BTS      int64  `json:"id_bts"`
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

	err = json.NewDecoder(r.Body).Decode(&remove)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = stmt.Exec(remove.ID_BTS)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "success",
		"removed_data": remove,
	})
}
