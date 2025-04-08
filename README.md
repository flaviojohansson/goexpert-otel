# goexpert-otel
Desafio Fullcycle - Pós GoExpert - Labs - Observabilidade e Open Telemetry

## Testar o projeto

1. Baixar e subir os serviços
```

git clone https://github.com/flaviojohansson/goexpert-otel

cd goexpert-otel

# Para efeito de aprendizado, este repositório já possui um arquivo .env com a chave da Weather API

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
