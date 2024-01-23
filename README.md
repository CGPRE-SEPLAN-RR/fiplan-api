# API do FIPLAN 

API para consulta dos dados do FIPLAN.

## Objetivos

- Tornar a consulta de relatórios específicios fácil, reproduzível e em um formato amplamente utilizado e suportado para modificações via script (`json`)

## Requisitos

- Ser extremamente intuitivo na sua utilização, de forma a gerar o menor número de dúvidas possível para os usuários
- Atender às [especificações da OpenAPI](https://swagger.io/specification/), visando fornecer uma documentação simples, porém extensiva, via [Swagger](https://swagger.io/)

## Como Iniciar o Ambiente de Desenvolvimento

1. Clone o repositório

```bash
git clone git@github.com:CGPRE-SEPLAN-RR/fiplan-api.git
```

2. Vá para a pasta raiz do repositório

```bash
cd fiplan-api
```

3. Baixe as dependências do projeto

```bash
go get .
```

4. Copie o arquivo com as variáveis de ambiente

```bash
cp .env.example .env
```

5. Complete as variáveis de ambiente como disposto abaixo

```bash
PORT=8080
APP_ENV=local

DB_DATABASE= # SID do banco de dados Oracle
DB_PASSWORD= # Senha do banco de dados
DB_USERNAME= # Usuário do banco de dados
DB_PORT=1521 # Porta do banco de dados (A porta padrão para bancos de dado Oracle é a 1521)
DB_HOST= # Endereço IP do banco de dados
```

6. Execute o código

```bash
make run
```

7. Acesse a documentação no [servidor local](http://localhost:8080/swagger/index.html)

## Comandos `make`

- Rodar todos os comandos, incluindo testes

```bash
make all build
```

- Gerar o executável da aplicação

```bash
make build
```

- Rodar a aplicação

```bash
make run
```

- Rodar a aplicação com *hot reload*

```bash
make watch
```

- Rodar testes

```bash
make test
```

- Limpar o executável antigo

```bash
make clean
```

