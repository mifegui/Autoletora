
Você pode rodar o aplicativo como um container docker:
docker build -t autoletora . && docker run --name autoletora -d --rm -p 8080:8080 autoletora

Para as funcionalidades que usam database funcionar elas precisam de uma database mongoDb com a porta padrão aberta. Você pode iniciar um container docker com mongoDb com o comando
```
docker run -d --name autoletora-mongo -d mongo
docker run -d -p 27017-27019:27017-27019 --name mongodb mongo
```
Sendo "autoletora-mongo" o nome do container.


Você pode testar a api coletora usando o curl:
```
curl -H "Content-Type: application/json" -d '{"event":"buy", "timestamp":"2016-09-22T13:57:31.2311892-04:00"}' http://localhost:8080/coletar
```
E ver sua entrada em localhost:8080/eventos
Os parâmetros podem ser quaisquer que caibam no padrão da struct evento.

