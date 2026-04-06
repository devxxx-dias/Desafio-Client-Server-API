# Desafio-Client-Server-API
Desafio técnico pós-graduação Full Cycle

## Como executar

> **Importante:** todos os comandos devem ser executados a partir da **raiz do projeto**.

### 1. Instalar dependências

```bash
go mod tidy
```

### 2. Iniciar o servidor (Terminal 1)

A partir da raiz do projeto:
```bash
go run src/server/server.go
```

Ou a partir da pasta do servidor:
```bash
cd src/server && go run server.go
```

Saída esperada:
```
Server started on port 8080
```

### 3. Executar o cliente (Terminal 2)

Com o servidor rodando, em outro terminal.

A partir da raiz do projeto:
```bash
go run src/client/client.go
```

Ou a partir da pasta do cliente:
```bash
cd src/client && go run client.go
```

O cliente irá criar o arquivo `cotacao.txt` (ou adicionar uma linha se já existir) com o seguinte conteúdo:

```
Dólar: 5.XXXX
```

Cada execução do cliente adiciona uma nova linha ao arquivo.

## Estrutura do projeto

```
.
├── src/
│   ├── client/
│   │   └── client.go   # programa cliente (package main)
│   ├── database/
│   │   └── database.db # banco SQLite (criado automaticamente)
│   └── server/
│       └── server.go   # programa servidor (package main)
├── cotacao.txt          # gerado pelo cliente
├── go.mod
└── go.sum
```

## Acessando o banco de dados localmente

O banco de dados SQLite fica em `src/database/database.db` e é criado automaticamente quando o servidor é iniciado pela primeira vez.

### Usando o SQLite CLI

```bash
sqlite3 src/database/database.db
```

Dentro do shell do SQLite, alguns comandos úteis:

```sql
-- Listar todas as cotações salvas
SELECT * FROM cotacoes;

-- Ver as N cotações mais recentes
SELECT * FROM cotacoes ORDER BY id DESC LIMIT 10;

-- Sair
.quit
```

> **Instalação do SQLite CLI:**
> - macOS: `brew install sqlite`
> - Ubuntu/Debian: `sudo apt install sqlite3`
> - Windows: baixe em https://sqlite.org/download.html

## Timeouts

| Operação                        | Timeout |
|---------------------------------|---------|
| Cliente → Servidor              | 300 ms  |
| Servidor → API externa (AwesomeAPI) | 200 ms  |
| Servidor → Banco de dados       | 10 ms   |
