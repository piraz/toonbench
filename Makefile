# Makefile

.DEFAULT_GOAL := help

.PHONY: help prepare go-bench go-bench-short js-bench

define PIRATE_BANNER
	@echo "@@@@@@@@@@@@@@@@@@@@@@@"
	@echo "@@@@@@@@@@@@@@@@@@@@@@@"
	@echo "          @@@@@@@@     JSON"
	@echo "          @@@@@@@@          Go      Protobuf"
	@echo "          @@@@@@@@                                 Python"
	@echo "          @@@@@@@@        JavaScript"
	@echo "                       "
	@echo "    @@@@              TOON TURUN TOON - TOOM TOOM... "
	@echo "    @@@@               "
	@echo "    @@@@@@@@@@@@@@@  TOON comparado a formatos de interoperabilidade."
	@echo "    @@@@@@@@@@@@@@@  Conta de padaria monstrando que TOON é só para LLMs..."
	@echo "                                                 by Piraz do Cafofo P)"
endef

.DEFAULT:
	@echo ""
	@echo "Unknown target: $(MAKECMDGOALS)"
	@echo ""
	@$(MAKE) help

help: ## Mostra esse help
	$(PIRATE_BANNER)
	@printf "\nUsage: make [target]\n\n"
	@printf "Targets:\n\n"
	@grep -E '^[a-zA-Z0-9._-]+:.*?## ' $(MAKEFILE_LIST) | sort | awk -F':.*?## ' '{ printf "  %-20s %s\n", $$1, $$2 }'
	@printf "\n"

prepare: ## Prepara o ambiente para o benchmark
	@printf "\n"
	@printf "\033[1;34m=========================================\033[0m\n"
	@printf "\033[1;32m Generating protobuf artifacts (prepare)\033[0m\n"
	@printf "\033[1;34m=========================================\033[0m\n\n"
	@echo "-> protoc go_out -> ."
	@protoc --go_out=. user.proto
	@echo "-> go mod tidy -> prepare go env"
	@go mod tidy
	@echo "-> pbjs (static commonjs) -> ./src/compiled_protos.js"
	@echo "-> npm run yarn -> prepare JavaScript env"
	@npm run yarn
	@echo "-> pbjs (static commonjs) -> ./src/compiled_protos.js"
	@npx pbjs -t static-module -w commonjs -o ./src/compiled_protos.js user.proto
	@mv ./src/compiled_protos.js ./src/compiled_protos.cjs
	@echo "-> protoc js_out (import_style=commonjs,binary) -> ./src"
	@protoc --js_out=import_style=commonjs,binary:./src user.proto
	@mv ./src/user_pb.js ./src/user_pb.cjs
	@printf "\n\033[1;34mDone.\033[0m\n"

go-bench: ## Executar benchmarks em go(gobench)
	@printf "\n"
	@printf "\033[1;34m========================================\033[0m\n"
	@printf "\033[1;33m Running Go benchmarks (gobench) \033[0m\n"
	@printf "\033[1;34m========================================\033[0m\n\n"
	@go test -bench=. -benchmem -benchtime=1s ./gobench
	@printf "\n\033[1;34mGo bench finished.\033[0m\n"

go-bench-parallel: ## Executar parallel benchmarks em go(gobench)
	@printf "\n"
	@printf "\033[1;34m===============================================\033[0m\n"
	@printf "\033[1;33m Running With Go Parallel benchmarks (gobench) \033[0m\n"
	@printf "\033[1;34m===============================================\033[0m\n\n"
	@go test -tags=parallel -bench=. -benchmem -benchtime=1s ./gobench
	@printf "\n\033[1;34mGo bench finished.\033[0m\n"

go-bench-short: ## Executar benchmarks em go(curto, benchtime=2x)
	@printf "\n"
	@printf "\033[1;34m========================================\033[0m\n"
	@printf "\033[1;33m Running short Go benchmarks (gobench) \033[0m\n"
	@printf "\033[1;34m========================================\033[0m\n\n"
	@go test -bench=. -benchmem -benchtime=2x ./gobench
	@printf "\n\033[1;34mGo short bench finished.\033[0m\n"

js-bench: ## Executar benchmarks em JavaScript (node)
	@printf "\n"
	@printf "\033[1;34m============================================\033[0m\n"
	@printf "\033[1;36m Running JS benchmark (Node + benchmark.js) \033[0m\n"
	@printf "\033[1;34m============================================\033[0m\n\n"
	@node --expose-gc src/bench.js
	@printf "\n\033[1;34mJS bench finished.\033[0m\n"
