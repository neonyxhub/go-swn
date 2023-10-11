define protoc_model
	protoc --go_out=. \
		--go_opt=paths=source_relative \
		$(1)
endef

define protoc_service
	protoc --go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		$(1)
endef

define protoc_model_arr
	@for proto in $^; do \
		$(call protoc_model,$$proto); \
	done
endef

define protoc_service_arr
	@for proto in $^; do \
		$(call protoc_service,$$proto); \
	done
endef

pb_pkg_model_gen: $(shell find pkg/ -name '*_model.proto')
	$(protoc_model_arr)

pb_pkg_service_gen: $(shell find pkg/ -name '*_api.proto')
	$(protoc_service_arr)

pb_internal_model_gen: $(shell find internal/ -name '*_model.proto')
	$(protoc_model_arr)

pb_internal_service_gen: $(shell find internal/ -name '*_api.proto')
	$(protoc_service_arr)