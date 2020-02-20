# Executando
No momento executa-se o aplicativo nativamente enquanto a database pode ser executado usando o docker. Suponho ser possível fazer os dois funcionarem usando o docker.

Para as funcionalidades que usam database funcionar (e o app inicializar) é necessário uma database mongoDb com a porta padrão aberta. Você pode iniciar um container docker com mongoDb com o comando
```
docker run -d -p 27017-27019:27017-27019 --name mongodb mongo
```

## Autocomplete
O serviço de autocomplete completa palavras enviadas em /complete/{palavra} usando os eventos já coletados pela api coletora em /coletar


Você pode usar a api coletora usando o curl:
```
curl -H "Content-Type: application/json" -d '{"event":"buy", "timestamp":"2016-09-22T13:57:31.2311892-04:00"}' http://localhost:8080/coletar
```
E ver sua entrada em localhost:8080/eventos
Os parâmetros podem ser quaisquer que caibam no padrão da struct evento.

> Autocomplete: 
Para localhost:8080/complete/bu
O retorno em json será:
``` json
{
  "matchs": [
    "buy"
  ]
}
```
Para /complete/com
``` json
{
  "matchs": [
    "comprou-camisa",
    "comprou-produto",
    "comprou"
  ]
}
```

## Manipulação de dados
Em /timeline se vê o resultado do agrupamento e ordenação do endpoint. Nenhum dado é armazenado na database.
