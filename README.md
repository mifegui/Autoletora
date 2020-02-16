
Você pode testar a api coletora usando o curl:
'''
curl -H "Content-Type: application/json" -d '{"event":"buy", "timestamp":"2016-09-22T13:57:31.2311892-04:00"}' http://localhost:8080/coletar
'''
E ver sua entrada em localhost:8080/eventos

