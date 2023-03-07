openssl genrsa -out example.key 2048
openssl req -new -key example.key -out example.csr
openssl x509 -req -days 365 -in example.csr -signkey example.key -out example.crt
