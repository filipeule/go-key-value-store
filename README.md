# Go Simple Key-Value Store

Este projeto foi desenvolvido durante a leitura do livro Cloud Native Go, de Matthew A. Titmus. Trata-se de uma store para armazenamento de informações no formato chave-valor, com suporte a um Transaction Log que registra todas as transações realizadas, garantindo a integridade dos dados.

Além disso, foi criada uma imagem Docker para facilitar a execução da aplicação em ambientes conteinerizados, permitindo uma implantação simples e portátil.

## Instalação

1. Crie as chaves cert.pem e key.pem, de acordo com este [tutorial](https://www.suse.com/pt-br/support/kb/doc/?id=000018152) e adicione na pasta "cert" do projeto

2. Crie a imagem do container utilizando o Dockerfile e execute o container

```bash
docker image build -t go-key-value-store:1.0.0

docker container run -d -p 8080:8080 --name kvs go-key-value-store:1.0.0
```

## Modo de uso

1. Criar ou atualizar uma chave
```bash
PUT /v1/key/{key}
```
Exemplo:
```bash
curl -X PUT https://localhost:8080/v1/key/message \
     -d "Hello from KVS!"
```

2. Ler o valor de uma chave
```bash
GET /v1/key/{key}
```
Exemplo:
```bash
curl https://localhost:8080/v1/key/message
```

3. Deletar uma chave
```bash
DELETE /v1/key/{key}
```
Exemplo:
```bash
curl -X DELETE https://localhost:8080/v1/key/message
```