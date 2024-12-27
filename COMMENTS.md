# Estratégia de Solução Para o Problema Votação do BBB

O sistema para registro de votos e acompanhamento da votação de um paredão contará com os serviços/estrutura:

    - API em Go pra receber os votos, enviar votos para mensageria e fornecer dados para acompanhamento da votação (Voting);
    - Microsserviço em GO para receber os votos da mensageria e registrar no banco de dados relacional (Register);
    - Microsserviço em GO para ler as informações no banco de dados relacional, montar relatórios, e registrá-los no banco de dados de chave valor (Summary);
    - Banco de Dados relacional em PostgreSQL;
    - Serviço de Mensageria Kafka;
    - Banco de Dados de estrutura de dados de chave-valor Memcached.

## API em Go pra Receber os Votos e Mostrar Acompanhamento (Voting)
### A API será responsável por:

    - Receber as solicitação do front-end para listar as opções de participantes do paredão;
    - Ler no banco de dados do Memcached as opções de participantes e responder ao front-end;
    - Registrar o voto: criar o voto com suas informações (identificação, participante escolhido, IP do usuário e horário do voto) e enviar para a mensageria;
    - Ler no banco de dados do Memcached o resultado parcial da votação até o momento e responder para o usuário junto com a confirmação do voto;
    - Receber as solicitação da Produção do programa para listar as informações da votação do paredão;
    - Ler no banco de dados do Memcached as informações da votação do paredão e responder para a Produção.

### Evitar o Uso de Máquinas/bots na Votação

    - Podemos aplicar uso de um CAPTCHA no front-end para cada voto;
    - O banck-end também pode gerar um desafio que deve ser respondido pelo usuário;
    - Aplicar regras de Firewall para bloquear solicitações vindas de um mesmo endereço de IP em um curto espaço de tempo.

## Microsserviço Para Receber os Votos da Mensageria e Registrar no Banco de Dados Relacional (Register)
Microsserviço responsável por receber um voto enviado pelo Apache Kafka e realizar a inserção do voto no banco de dados relacional.

## Microsserviço Para Ler as Informações do Banco de Dados Relacional (Summary)
Microsserviço responsável por:

    - Ler as informações do paredão ativo no momento e cadastrar os dados no banco de dados do Memcached.
    - Fazer a contagem dos votos de cada participante no banco de dados relacional, montar as informações e inserir no banco de dados do Memcached. Esta contagem deve ser feita em um intervalo de tempo, em estudo, de forma que o microsserviço não sobrecarregue o banco de dados com leituras excessivas.

## Banco de Dados Relacional
    - PostgreSQL (https://www.postgresql.org/);
    - Banco de dados que terá as tabelas para registro das informações de cada paredão, participantes do paredão e votos de cada paredão. Será feito uso de índices na tabela de votos para permitir que a recuperação de dados seja mais rápida;
    - Escolhi o PostgreSQL pois é o banco de dados que atualmente eu trabalho, é open-source e tem um desempenho satisfatório para escrita de dados e consultas mais complexas.

## Banco de Dados de Estrutura de Dados de chave-valor
    - Memcached (https://memcached.org/);
    - O Memcached será responsável por receber e fornecer os dados que serão acessados com mais frequência pelo API de Registro de Votos e Acompanhamento da Votação, sendo eles:
        - os dados dos participantes para serem enviados ao front-end e posteriormente enviados para a mensageria; 
        - os dados de acompanhamento da votação.
    - Escolhi o Memcached pois é uma solução open-source.

## Serviço de Mensageria
    - Apache Kafka (https://kafka.apache.org/);
    - O serviço de mensageria será responsável por receber os votos enviados pela API, armazenar e enviar cada voto para o microsserviço responsável por registrar os votos no banco de dados relacional.
    - Escolhi o Kafka pois encontrei mais documentação e exemplos, solução open-source.


# Solução Implementada Até o Momento

## API Voting

API desevolvida em golang seguindo o padrão de arquitetura limpa. A API se conecta aos bancos de dados PostgreSQL (relacional) e Memcache (chave-valor). A API também se conecta com o Apache Kafka para enviar os dados de um voto para regsitro.

Todas as interações com os usuários são feitas atravez das rotas HTTP explicadas a seguir:

### Cadastro de Novo Participante
[POST] host:8000/participant
- Corpo da requisição:

    `{
	"name": "Lucia Maria"
}`

A rota retorna o ID do participante cadastrado.


### Listagen dos Participantes Cadastrados
[GET] host:8000/all-participants

- Resposta da requisição:
    
    `{
	"participants": [
		{
			"id": "2a466427-e8f0-480b-b280-030f9e0bfbd7",
			"name": "Patricia Lima",
			"created_at": "2024-12-23T12:44:51.268216-03:00"
		},
		{
			"id": "56079f38-38f1-496e-98fe-05ea0477c523",
			"name": "Lucia Maria",
			"created_at": "2024-12-23T12:39:11.570243-03:00"
		}
	]
}`

A rota lista todos os participantes cadastrados no banco de dados do PostgreSQL.


### Cadastro de Novo Paredão
[POST] host:8000/wall

- Corpo da requisição:

    `{
	"name": "Paredão 01",
	"start_time": "2024-12-25T00:00:00.000000-03:00",
	"end_time": "2024-12-25T23:59:59.999999-03:00",
	"participants_id": ["ce985aee-1be0-488a-92aa-fb849f40060a","56079f38-38f1-496e-98fe-05ea0477c523"]
}`

A rota faz as validações básicas, registra o paredão no banco de dados do PostgreSQL e retorna o id do paredão cadastrado.


### Listagem dos Participantes do Paredão Ativo
[GET] host:8000/participants-active-wall

- Resposta da requisição:
    
    `{
	"wall_id": "103423d2-9d15-4f4b-b759-03bf50ec70a5",
	"wall_name": "Paredão de Teste 1",
	"participants": [
		{
			"id": "ce985aee-1be0-488a-92aa-fb849f40060a",
			"name": "Pedro André F"
		},
		{
			"id": "de52c486-0f42-4784-8d31-545c44215f85",
			"name": "João Almeida"
		},
		{
			"id": "5c1a59e5-d0f6-420f-ad07-1daa8ad7f036",
			"name": "André Martins"
		}
	]
}`

A rota retorna os participantes do paredão ativo encontrados no Memecache.


### Regitro do Voto
[POST] host:8000/vote

- Corpo da requisição:

    `{
	"participant_id": "5c1a59e5-d0f6-420f-ad07-1daa8ad7f036",
	"wall_id": "103423d2-9d15-4f4b-b759-03bf50ec70a5"
}`

A rota faz as validações necessárias, registra o voto e retorna o ID do voto registrado e o resultado da votação até o momento.

- Corpo da resposta:

    `{
	"register_vote_id": "c9fe875e-effd-4d5d-a615-159a683fa7cc",
	"partial_result": [
		{
			"id": "ce985aee-1be0-488a-92aa-fb849f40060a",
			"name": "Pedro André F",
			"votes": 0,
			"votes_percentage": 0
		},
		{
			"id": "de52c486-0f42-4784-8d31-545c44215f85",
			"name": "João Almeida",
			"votes": 30,
			"votes_percentage": 0.08670520231213873
		},
		{
			"id": "5c1a59e5-d0f6-420f-ad07-1daa8ad7f036",
			"name": "André Martins",
			"votes": 316,
			"votes_percentage": 0.9132947976878613
		}
	]
}`


### Listagem dos Paredrões Cadastrados 
[GET] host:8000/all-walls

- Corpo da resposta:

    `{
	"walls": [
		{
			"id": "103423d2-9d15-4f4b-b759-03bf50ec70a5",
			"name_wall": "Paredão de Teste 1",
			"created_at": "2024-12-17T14:07:53.711973-03:00",
			"start_time": "2024-12-22T00:00:00-03:00",
			"end_time": "2024-12-23T23:59:59.999-03:00"
		},
		{
			"id": "f823174f-805a-4100-a1e5-1e56c7615f41",
			"name_wall": "Paredão de Teste 2",
			"created_at": "2024-12-17T14:07:53.711973-03:00",
			"start_time": "2024-12-16T00:00:00-03:00",
			"end_time": "2024-12-17T23:59:59.999-03:00"
		}
	]
}`

A rota lista todos os paredões cadastrados no banco de dados do PostgreSQL.


### Busca todos os dados de um Paredão pelo ID
[GET] host:8000/wall/:id

- Corpo da resposta:

    `{
	"id": "103423d2-9d15-4f4b-b759-03bf50ec70a5",
	"name_wall": "Paredão de Teste 1",
	"start_time": "2024-12-22T00:00:00-03:00",
	"end_time": "2024-12-23T23:59:59.999-03:00",
	"total_votes": 348,
	"participants": [
		{
			"id": "ce985aee-1be0-488a-92aa-fb849f40060a",
			"name_participant": "Pedro André F",
			"votes": 0,
			"votes_percentage": 0
		},
		{
			"id": "de52c486-0f42-4784-8d31-545c44215f85",
			"name_participant": "João Almeida",
			"votes": 30,
			"votes_percentage": 0.08620689655172414
		},
		{
			"id": "5c1a59e5-d0f6-420f-ad07-1daa8ad7f036",
			"name_participant": "André Martins",
			"votes": 318,
			"votes_percentage": 0.9137931034482759
		}
	]
}`

A rota lista todos os dados, incluindo dados de votação, de um paredão encontrado no banco de dados do PostgreSQL.


### Listagem do Resultado Parcial do Paredão Ativo
[GET] host:8000/partial-result-active-wall

- Corpo da resposta:

    `{
	"id": "103423d2-9d15-4f4b-b759-03bf50ec70a5",
	"name_wall": "Paredão de Teste 1",
	"start_time": "2024-12-22T00:00:00-03:00",
	"end_time": "2024-12-23T23:59:59.999-03:00",
	"total_votes": 348,
	"participants": [
		{
			"id": "ce985aee-1be0-488a-92aa-fb849f40060a",
			"name_participant": "Pedro André F",
			"votes": 0,
			"votes_percentage": 0
		},
		{
			"id": "de52c486-0f42-4784-8d31-545c44215f85",
			"name_participant": "João Almeida",
			"votes": 30,
			"votes_percentage": 0.08620689655172414
		},
		{
			"id": "5c1a59e5-d0f6-420f-ad07-1daa8ad7f036",
			"name_participant": "André Martins",
			"votes": 318,
			"votes_percentage": 0.9137931034482759
		}
	]
}`

A rota lista todos os dados, incluindo dados de votação, do paredão ativo encontrado no banco de dados do Memecache.
