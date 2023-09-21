install-pre-commit-config:
	@pre-commit autoupdate || sudo pre-commit autoupdate
	@pre-commit install --hook-type pre-commit --hook-type pre-push --hook-type post-commit || sudo pre-commit install --hook-type pre-commit --hook-type pre-push --hook-type post-commit
