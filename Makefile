# Makefile
.PHONY: proto-gen
proto-gen:
	@echo "A gerar código via Buf..."
	buf generate
	@echo "Código gerado com sucesso!"

.PHONY: test-alert
test-alert:
	grpcurl -plaintext -d '{"user_id": "1", "symbol": "BTC/EUR", "target_price": 50000.0}' localhost:50051 api.proto.v1.AlertService/CreateAlert