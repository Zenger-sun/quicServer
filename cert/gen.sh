openssl genrsa -out ca.key 4096

openssl req \
   -new \
   -sha256 \
   -out ca.csr \
   -key ca.key \
   -config ca.conf

openssl x509 \
     -req \
     -days 3650 \
     -in ca.csr \
     -signkey ca.key \
     -out ca.crt

openssl genrsa -out server.key 2048

openssl req \
   -new \
   -sha256 \
   -out server.csr \
   -key server.key \
   -config server.conf

openssl x509 \
   -req \
   -days 3650 \
   -CA ca.crt \
   -CAkey ca.key \
   -CAcreateserial \
   -in server.csr \
   -out server.pem\
   -extensions req_ext \
   -extfile server.conf
