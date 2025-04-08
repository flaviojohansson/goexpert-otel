# goexpert-otel
Desafio Fullcycle - Pós GoExpert - Labs - Observabilidade e Open Telemetry

## Testar o projeto

1. Baixar e subir os serviços
```

git clone https://github.com/flaviojohansson/goexpert-otel

cd goexpert-otel

# crie um arquivo .env com a chave para o serviço ex:
# echo "WEATHER_API_KEY=b193f39823d249099e4141815252303" > .env

docker compose up -d
```
2. Chamar a API
```
curl -X POST localhost:8080 -d '{"cep": "80530000"}'
```
3. Abrir Zipkin no browser e conferir os tracings e os spans
```
http://localhost:9411
```
4. interromper os serviços e fazer a limpeza final, tal qual O Lobo do PulpFiction
```
docker compose down -v
```
