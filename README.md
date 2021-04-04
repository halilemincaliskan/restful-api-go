# restful-api-go
Kartaca Görev Task

## Özel Anahtar Kodu
gAAAAABgZG9Pc-hOPraRo4MPchmmdn4Y2ckQnoX9DcDTHhS1b-AGqtEej6yRBflweiu3fk-uHcfQTTfBy8smpicJurfacck-Wc9EtdYdaVeECzwzd_qLXaaVNV1eY_GJH5ZO1_ZZ
OSnRTc2K9Fn_g_6-FnWziHV8G0eUIk9h0v7lYHovk8HROBh8qfVCcCplAs9YVi87m-0mvtpvvhNYj57n_rZfYHz4-w==


## Çalıştırmak İçin

Projenin api tarafı Go ile yazdım. Log bilgisini okuyup database'e yazması için Kafka kullandım. Kafkanın logları yazacağı database olarak ise mongoDB kullandım.
Projeyi çalıştırmak için aşağıdaki yolları takip edebilirsiniz.

### İlk olarak docker-compose aşağıdaki gibi ayağa kaldırılır.
```go
docker-compose build
```

### Log oluşturmak için aşağıdaki rootlara istek atılır
İster Postman kullanarak localhost:10000 url'ine GET,POST,PUT,DELETE istekleri atabilirsiniz isterseniz de aşağıdaki kodlardaki gibi Curl ile istekleri atabilirsiniz.
```go
curl -X GET localhost:10000
curl -X POST localhost:10000
curl -X PUT localhost:10000
curl -X DELETE localhost:10000
```

### Verilerin Chart Üzerinde Gösterilmesi
Aşağıdaki link üzerinden chart'ı görüntüleyebilirsiniz.
[localhost:10000/public](http://localhost:10000/public)