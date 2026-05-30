# CryptoTracker Alert Service
Chega de ter o TradingView aberto o dia todo. Define o teu preço-alvo e o serviço avisa quando lá chegar.

---

## Requisitos Funcionais
* Criar um alerta para um par símbolo de Crypto e um preço-alvo.
* Cancelar um alerta pendente.
* Listar alertas existentes.
* Monitorizar o preço do par em tempo real via Binance WebSocket (`wss://stream.binance.com:9443/ws`) — pares disponíveis via [Exchange Info API](https://api.binance.com/api/v3/exchangeInfo).
* Disparar o alerta quando o preço-alvo é atingido.
* Persistir o estado dos alertas na DB (`PENDING`, `TRIGGERED`, `CANCELLED`).
* No arranque, retomar a monitorização de todos os alertas `PENDING` existentes.

---

## A magia do CryptoTracker - Alert!
### Scout (Price Prime Monitor)
Componente responsável por gerir as ligações WebSocket à Binance e distribuir
eventos de preço pelo Sniper.

#### Responsabilidades
- Manter uma ligação WebSocket ativa por símbolos a serem monitorizados.
- Abrir/fechar streams dinamicamente consoante os alertas activos (Existe simbolo a ser monitorizado ou não?)
- Passar eventos de preço ao AlertDispatcher para avaliação

#### Estado interno
| Campo | Tipo | Descrição |
|---|---|---|
| `streams` | `map[string]context.CancelFunc` | Símbolo → função de cancelamento da goroutine |
| `mu` | `sync.Mutex` | Protege o mapa de streams contra race conditions |

#### Métodos
| Método | Descrição |
|---|---|
| `Subscribe(symbol)` | Abre stream WS para o símbolo, se ainda não existir |
| `Unsubscribe(symbol)` | Fecha stream WS se não houver mais alertas pendentes |

#### Ciclo de vida de um stream
1. `Subscribe` lança uma goroutine com um `context.Context` próprio.
2. A goroutine conecta ao WS da Binance (`<symbol>@aggTrade`).
3. Para cada evento recebido → `AlertDispatcher.Check(symbol, price)`.
4. `Unsubscribe` cancela o context → goroutine termina.

#### Resiliência
- No arranque do serviço, todos os símbolos com alertas `PENDING` na DB são subscritos automaticamente.
- Em caso de desconexão do WS, a goroutine faz reconnect com backoff exponencial.

### Sniper (Alert Dispatcher)
Componente responsável por avaliar os eventos de preço recebidos pelo Scout
e disparar os alertas quando o preço-alvo é atingido.

#### Responsabilidades
- Manter em memória os alertas `PENDING` agrupados por símbolo.
- Avaliar cada evento de preço contra os alertas ativos.
- Marcar alertas como `TRIGGERED` na DB quando o preço-alvo é atingido.
- Remover alertas da memória quando são disparados ou cancelados.

#### Estado interno
| Campo | Tipo | Descrição |
|---|---|---|
| `alerts` | `map[string][]Alert` | Símbolo → lista de alertas PENDING |
| `mu` | `sync.Mutex` | Protege o mapa contra race conditions |

#### Métodos
| Método | Descrição |
|---|---|
| `Load(alerts []Alert)` | Carrega alertas `PENDING` da DB no arranque. |
| `Add(alert Alert)` | Adiciona um novo alerta à memória. |
| `Remove(id)` | Remove um alerta da memória (cancel). |
| `Check(symbol, price)` | Avalia o preço contra os alertas activos do símbolo. |
| `Done(id)` | Marca o alerta como `TRIGGERED` na DB + remove da memória. |

---

## gRPC
| RPC | Request | Response | Descrição |
|---|---|---|---|
| `CreateAlert` | `symbol`, `target_price` | `Alert` | Cria um alerta → persiste na DB + adiciona ao Sniper + subscreve o Scout |
| `CancelAlert` | `id` | `Empty` | Cancela um alerta → marca `CANCELLED` na DB + remove do Sniper |
| `ListAlerts` | `Empty` | `[]Alert` | Devolve todos os alertas no estado `PENDING` |

---

## Alertas

#### Estados
| Estado | Descrição |
|---|---|
| `PENDING` | Alerta ativo, à espera do preço-alvo |
| `TRIGGERED` | Preço-alvo atingido |
| `CANCELLED` | Cancelado pelo utilizador |

#### Schema
| Campo | Tipo | Descrição |
|---|---|---|
| `id` | `BIGINT PK` | Identificador único |
| `symbol` | `VARCHAR(20)` | Símbolo (ex: `BTCUSDT`) |
| `target_price` | `DECIMAL(20,8)` | Preço-alvo |
| `status` | `ENUM` | `PENDING`, `TRIGGERED`, `CANCELLED` |
| `created_at` | `DATETIME` | Data de criação |
| `triggered_at` | `DATETIME` | Data de disparo (nullable) |

## Notificação

Quando um alerta é atingido, o evento é publicado num tópico Kafka.
O `notification-service` consome esse evento e envia a notificação via Telegram ou Webhook.

## Tecnologias
* **Linguagem**: Go
* **Comunicação**: gRPC
* **Base de Dados**: MySQL
* **Infraestrutura**: Dockerfile & Docker Compose
* **Ferramentas**: Visual Studio Code, Docker Desktop, Heidisql (Everything installed with winget <3)