
# Go + Gin + MongoDB + New Relic APM Example

Bu proje, Go ile yazılmış, Gin framework'ü ve MongoDB kullanan, New Relic APM ile izlenebilir bir REST API örneğidir.

## Özellikler
- **/api/data**: MongoDB'den item listeler (GET)
- **/api/add**: MongoDB'ye yeni item ekler (POST, JSON: `{ "name": "değer" }`)
- New Relic ile HTTP ve DB işlemleri izlenir
- Başlangıçta koleksiyon boşsa otomatik örnek veri ekler

## Kurulum
1. **Go ve MongoDB kurulu olmalı**
2. New Relic hesabınızdan bir lisans anahtarı alın
3. `main.go` içindeki `ConfigLicense` kısmına anahtarı girin
4. Bağımlılıkları yükleyin:
   ```sh
   go mod tidy
   ```
5. MongoDB başlatın:
   ```sh
   mongod
   ```
6. Uygulamayı başlatın:
   ```sh
   go run main.go
   ```

## API Kullanımı
- **GET /api/data**: Tüm item'ları listeler
- **POST /api/add**: Yeni item ekler
  ```sh
  curl -X POST -H "Content-Type: application/json" -d '{"name":"örnek"}' http://localhost:8080/api/add
  ```

## New Relic
- Uygulama başlatıldığında New Relic'e otomatik olarak veri gönderir
- DB işlemleri (find, insert) New Relic APM UI'da Datastore olarak görünür

## Notlar
- `main.go` içindeki lisans anahtarınızı paylaşmayın!
- Demo amaçlıdır, production için ek güvenlik ve yapılandırma gereklidir.
# newrelic-integration
