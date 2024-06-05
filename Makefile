.PHONY: run-backend build-script run-demo lint-script

run-backend:
	@cd backend && encore run
build-script:
	@cd script && npm run build && rm ../local-demo/cursors.min.js && cp dist/cursors.min.js ../local-demo/cursors.min.js
run-demo:
	@cd local-demo && python3 -m http.server 8000
lint-script:
	@cd script && npm run lint