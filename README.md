# Client-Server API - Cotação do Dólar

Este projeto implementa um sistema cliente-servidor em Go para consulta de cotação do dólar americano (USD/BRL).

## Funcionalidades

- **Server**: Consome API externa de cotação, persiste dados em SQLite e retorna cotação para o cliente
- **Client**: Faz requisição ao servidor e salva cotação em arquivo texto

## Requisitos

- Go 1.24.5 ou superior
- SQLite (incluído via driver Go)

## Como executar

### 1. Baixar dependências

```bash
go mod tidy
```

### 2. Executar o servidor

```bash
# Na raiz do projeto
go run server/server.go
```

O servidor será iniciado na porta 8080 e estará disponível no endpoint `/cotacao`.

### 3. Executar o cliente (em outro terminal)

```bash
# Na raiz do projeto  
go run client/client.go
```

O cliente fará uma requisição ao servidor e salvará a cotação no arquivo `cotacao.txt`.

## Estrutura do projeto

```
├── client/
│   └── client.go          # Cliente HTTP
├── server/
│   ├── server.go          # Servidor HTTP
│   ├── exchange/
│   │   └── exchange.go    # Integração com API externa
│   └── database/
│       └── database.go    # Operações SQLite
├── model/
│   └── model.go           # Estruturas de dados
├── go.mod                 # Dependências
└── README.md
```

## Timeouts configurados

- **API externa**: 200ms
- **Banco de dados**: 10ms  
- **Cliente**: 300ms

## Arquivos gerados

- `cotacao.db`: Banco SQLite com histórico de cotações
- `cotacao.txt`: Arquivo com a cotação atual no formato "Dólar: {valor}"

## API Externa utilizada

- **URL**: https://economia.awesomeapi.com.br/json/last/USD-BRL

## Exemplo de uso

1. Execute o servidor:
```bash
go run server/server.go
```

2. Em outro terminal, execute o cliente:
```bash
go run client/client.go
```

3. Verifique o arquivo gerado:
```bash
cat cotacao.txt
```

4. No banco cotacao.db
```sql
select * from exchange e 
```