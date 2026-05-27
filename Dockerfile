# ==============================================================================
# Porquê Multi-stage build? Para seguir boas práticas.
# 1. Otimização de Tamanho: Separamos o ambiente de Desenvolvimento (Stage 1) do ambiente de Produção (Stage 2), pois só vamos usar Alpine Linux puro e o binário final compilado em Go.
# 2. Segurança: A imagem de produção não contém o compilador do Go, código-fonte ou ferramentas de desenvolvimento. Se o container for comprometido, o atacante não tem ferramentas para compilar ou executar scripts maliciosos.
# 3. Eficiência de CI/CD: Imagens ultra-leves são transferidas e inicializadas na Cloud (AWS/GCP) em milissegundos, reduzindo custos de storage e tráfego.
# ==============================================================================


# ==============================================================================
# STAGE 1: Builder (Ambiente de Desenvolvimento e Compilação)
# ==============================================================================
# Utilizamos a imagem oficial do Go 1.26 baseada em Alpine para termos as ferramentas de compilação necessárias de forma ligeiramente mais otimizada.
FROM golang:1.26.3-alpine AS builder 

# Define o diretório de trabalho dentro do container onde o código será processado.
WORKDIR /app

# Copia os ficheiros que definem as dependências e as suas assinaturas de segurança (hashes).
# Fazer isto ANTES de copiar o resto do código permite aproveitar a cache do Docker.
COPY go.mod go.sum ./

# Descarrega e valida todas as dependências listadas no go.mod e go.sum. Se não mudarem, este passo é ignorado pelo Docker.
RUN go mod download

# Copia todo o resto do código-fonte do projeto para dentro do diretório /app do container
COPY . .

# Compila o código-fonte num binário executável único e estático.
# - CGO_ENABLED=0: Desativa dependências de bibliotecas C dinâmicas (essencial para rodar em Alpine/Scratch).
# - GOOS=linux: Garante que o binário é compilado especificamente para o sistema operativo Linux.
# - -o cryptotracker-alert: Define o nome do ficheiro binário de saída.
RUN CGO_ENABLED=0 GOOS=linux go build -o cryptotracker-alert ./cmd/cryptotracker-alert-server/main.go


# ==============================================================================
# STAGE 2: A Imagem de Produção (Ambiente de Execução Mínimo)
# ==============================================================================
# Utilizamos uma imagem Alpine Linux totalmente limpa e minimalista.
FROM alpine:latest

# Instala os certificados CA (Root Certificates) necessários para que a nossa aplicação consiga fazer pedidos HTTPS seguros (ex: falar com as APIs da Binance ou CoinGecko).
# Cria também um utilizador dedicado sem privilégios (lgcarvalho) para correr a aplicação, evitando que o processo tenha acesso root ao sistema de ficheiros do container.
RUN apk --no-cache add ca-certificates && \
    adduser -D -g '' lgcarvalho

# Define o diretório de trabalho no home do utilizador dedicado.
WORKDIR /home/lgcarvalho/

# A magia do Multi-stage: Vai ao Stage 1 (builder) e copia apenas o binário compilado.
# Todo o SDK do Go, compilador e código-fonte original são descartados aqui.
COPY --from=builder /app/cryptotracker-alert .

# Garante que o utilizador lgcarvalho tem permissão de execução sobre o binário.
RUN chown lgcarvalho:lgcarvalho ./cryptotracker-alert

# A partir deste ponto, todos os comandos (incluindo o CMD) correm como lgcarvalho e não como root.
USER lgcarvalho

# Documenta que o container vai escutar na porta 50051 (porta padrão que usamos para o gRPC)
EXPOSE 50051

# Comando executado quando o container inicia. Corre o nosso binário de forma nativa e direta.
CMD ["./cryptotracker-alert"]