FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    postgresql-client \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

# [Optional] Uncomment the next lines to use go get to install anything else you need
USER vscode

# Download the latest installer, run it and remove it.
RUN curl -LsSf https://astral.sh/uv/install.sh -o ~/uv-installer.sh && sh ~/uv-installer.sh && rm ~/uv-installer.sh

# Ensure the installed binary is on the `PATH`
ENV PATH="~/.cargo/bin/:$PATH"

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.1 \
    && go install github.com/vektra/mockery/v2@v2.46.3 \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0 \
    && go install github.com/magefile/mage@v1.15.0 \
    && go clean -cache -modcache
USER root

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1
