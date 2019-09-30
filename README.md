# RedCrawler

Para executar o crawler antes deve-se antes instalar os seguintes packages:
- golang.org/x/net/html/charset
- github.com/PuerkitoBio/goquery
- github.com/lib/pq

Utilizando o comando: 'go get github.com/{nome do package}' no terminal e criar um banco de dados vazio Postgres chamado 'redcrawler'.
Após isso para executar o crawler deve-se rodar o seguinte comando 'go build main.go' e em seguida 'go ./main'.

