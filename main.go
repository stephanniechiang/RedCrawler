package main

import (
	"crawl"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	results := []crawl.Acao{}

	results = crawl.Crawl("https://www.fundamentus.com.br/detalhes.php")

	//fmt.Println(results)

	psqlInfo := fmt.Sprintf("postgres://postgres:postgres@localhost/redcrawler?sslmode=disable")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlClear := `DELETE FROM acoes`
	sqlStatement :=`INSERT INTO acoes (posicao, papel, empresa, oscildia, valormerc) VALUES ($1, $2, $3, $4, $5) RETURNING papel`

	_, err = db.Exec(sqlClear)
	if err != nil {
		panic(err)
	}

	var i = 0
	for i < len(results){
		var papel string
		empresa := string(results[i].Empresa)

		err = db.QueryRow(sqlStatement, results[i].Posicao, results[i].Papel, empresa, results[i].Oscil_dia, results[i].Valor_merc).Scan(&papel)

		if err != nil {
			panic(err)
		}

		i++
	}
}



