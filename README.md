# voting

## Requisitos
Para rodar o sistema e seus microsserviços, você precisará de:
- Golang Build Tools 1.22.6+ 
- Docker
- Arquivo `config.json` próprio (presente em cada diretório raiz)



### No momento para rodar todo o projeto é necessário seguir os passos abaixo:
### (01) Rodar Kafka
- Acessar a pasta kafka na raiz do projeto;
- Abrir um terminal e executar o comando:
    - 'docker compose up -d'
- Entrar no container kafika-kafka-1-1 para criar o tópico:
    - 'docker exec -it kafika-kafka-1-1 bash'
- Criar o tópico "votacaoVotoBBB":
    - 'kafka-topics --create --bootstrap-server localhost:19092 --replication-factor 2 --partitions 3 --topic votacaoVotoBBB --config min.insync.replicas=2 acks=all enable.idempotence=true'

### (02) Rodar Memcache
- Executar o comando abaixo para rodar o Memcache:
    - 'docker run --name memcached -d -p 11211:11211 memcached:alpine3.21'

### (03) Rodar Banco de Dados
- Acessar o diretório raiz do sistema e executar o comando abaixo: 
    - 'docker compose up -d'
- Acessar o banco de dados com as configurações presentes no arquivo `docker-compose.yml` (voting/docker-compose.yml) e aplicar os scrips presentes no arquivo `db_migration_01.sql` (voting/migrationdb_migration_01.sql).

### (04) Rodar o Server de Votação (Voting)
- Acessar o diretório raiz do sistema e executar o arquivo `main.go`:
    - 'go run main.go'

### (05) Rodar Microsserviço Summary
- Acessar o diretório summary (voting/cmd/summay) e executar o arquivo `main.go`:
    - 'go run main.go'

### (06) Rodar Microsserviço Register
- Acessar o diretório register (voting/cmd/register) e executar o arquivo `main.go`:
    - 'go run main.go'
