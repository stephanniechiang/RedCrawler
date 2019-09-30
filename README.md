# RedCrawler

### Prepare os packages e banco de dados
Para executar o crawler antes deve-se antes instalar os seguintes packages:
- golang.org/x/net/html/charset
- github.com/PuerkitoBio/goquery
- github.com/lib/pq

Utilizando o comando: `go get {nome do package)` no terminal.
Criar um banco de dados vazio Postgres chamado **redcrawler**.
E executar a seguinte query no banco de dados postgres:

    CREATE TABLE acoes(
    	posicao INT,
    	papel VARCHAR(10) PRIMARY KEY,
    	empresa VARCHAR(256),
    	oscilDia VARCHAR(256),
    	valorMerc FLOAT
	);

### Execute o crawler
Em seguida, para executar o crawler deve-se rodar o seguinte comando `go build main.go` e em seguida `./main`.
O crawler irá rodar por alguns minutos e ao final populará o banco de dados com os dados requisitados.

