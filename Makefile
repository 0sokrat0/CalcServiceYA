build: ## 🚀 Build and start containers
	@printf "🚀 Building and starting containers...\n"
	@{ \
	  docker-compose up -d --build > /dev/null 2>&1 & \
	  pid=$$!; \
	  bar=""; \
	  total=30; \
	  i=0; \
	  finished=0; \
	  while kill -0 $$pid 2>/dev/null; do \
	    percent=$$(( i * 100 / total )); \
	    filled=$$(printf "%*s" $$i | tr " " "="); \
	    empty=$$(printf "%*s" $$(( total - i )) | tr " " " "); \
	    printf "\r\033[1;34m[%s%s]\033[0m \033[1;32m%3d%%\033[0m Building containers..." "$$filled" "$$empty" "$$percent"; \
	    sleep 0.2; \
	    i=$$(( (i + 1) % (total + 1) )); \
	  done; \
	  finished=1; \
	  # Плавное завершение до 100%, если не завершилось \
	  while [ $$i -le $$total ]; do \
	    percent=$$(( i * 100 / total )); \
	    filled=$$(printf "%*s" $$i | tr " " "="); \
	    empty=$$(printf "%*s" $$(( total - i )) | tr " " " "); \
	    printf "\r\033[1;34m[%s%s]\033[0m \033[1;32m%3d%%\033[0m Building containers..." "$$filled" "$$empty" "$$percent"; \
	    sleep 0.05; \
	    i=$$((i + 1)); \
	  done; \
	  # Финальная задержка для эффекта \
	  sleep 0.3; \
	  printf "\r\033[1;32m✅ Containers started successfully!\033[0m\n"; \
	}
