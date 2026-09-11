<div align="center">

# ttl

### Seu Arquivo Pessoal de Conhecimento

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

[English](README.md) | [简体中文](README.zh-CN.md) | [日本語](README.ja.md) | [Español](README.es.md) | [Français](README.fr.md)

---


</div>

---

## 📖 A História

Todo desenvolvedor já passou por isso:

> "Espera, qual era aquele comando Docker que usei mês passado?"
> "Onde meu colega compartilhou aquele arquivo de configuração?"
> "Eu sei que li um artigo sobre isso... mas não consigo encontrar."

Armazenamos conhecimento em todo lugar — abas do navegador, mensagens do Slack, emails, pastas de favoritos, aplicativos de notas. Quando realmente precisamos, perdemos tempo cavando em abas infinitas e rolando o histórico de chat.

**Essa também era minha dor.**

Por isso criei **ttl**.

O nome vem de "Time to Live" — mas com um significado diferente. Em vez de expirar, é sobre dar ao seu conhecimento um **tempo para viver** para sempre.

- Armazene tudo em um só lugar como pares chave-valor
- Marque com tags para fácil organização
- Busque instantaneamente por palavra-chave

Chega de procurar em emails antigos ou rolar o histórico do Slack. Apenas `ttl get <palavra-chave>` e você tem.

**ttl é seu arquivo pessoal de conhecimento — tudo que você precisa, quando você precisa.**

---

## ✨ Funcionalidades

| Funcionalidade | Descrição |
|------|------|
| 🗄️ **Armazenamento KV Local** | Banco de dados embarcado rápido, sem configuração (bbolt) |
| 🏷️ **Sistema de Tags** | Organize recursos com tags flexíveis e pesquisáveis |
| 🔍 **Busca Fuzzy** | Encontre o que precisa instantaneamente entre chaves e tags |
| ☁️ **Sincronização em Nuvem** | Auto-hospede um servidor multi-tenant com isolamento de dados por usuário |
| 🚀 **Abertura Inteligente** | Abra URLs e arquivos com programas padrão do sistema |
| 📤 **Exportar** | Exporte dados como JSON ou CSV |

---

## 🚀 Início Rápido

### Instalação

#### Linux / macOS

```bash
# Instalar do GitHub releases
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.sh)"

# Ou compilar a partir do código-fonte
go build -o ttl .
sudo mv ttl /usr/local/bin/
```

#### Windows

```powershell
# Instalar do GitHub releases
irm https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.ps1 | iex
```

#### URL de Download Personalizada

Para redes internas ou mirrors personalizados:

```bash
# Linux/macOS
TTL_DOWNLOAD_URL="https://your-mirror.com/ttl-cli-v1.0.0-linux-amd64" /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.sh)"
```

```powershell
# Windows
$env:TTL_DOWNLOAD_URL="https://your-mirror.com/ttl-cli-v1.0.0-windows-amd64.zip"; irm https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.ps1 | iex
```

### Uso Básico

```bash
# Adicionar um recurso
ttl add my-link https://example.com

# Adicionar com tags
ttl add docker-cmd "docker run -d -p 8080:80 nginx"
ttl tag docker-cmd dev ops

# Buscar recursos
ttl get docker

# Abrir no navegador
ttl open my-link

# Deletar
ttl del old-key
```

---


## 📝 Log de Trabalho

Acompanhe e filtre seu trabalho diário.

```bash
# Escrever uma entrada de log
ttl log write "Refatoração do módulo de usuários concluída" --tags "projetoA,dev"

# Ver logs
ttl log list                    # Logs de hoje
ttl log list --range week       # Esta semana
ttl log list --range month      # Este mês

```

---


## ☁️ Servidor em Nuvem e Sincronização

### Iniciar Seu Servidor

```bash
# Criar um usuário
ttl server user add alice

# Iniciar servidor multi-tenant
ttl server start --port 8080
```

### Sincronizar Seus Dados

```bash
# Configurar servidor remoto
ttl config
# Edite a seção server com seu endpoint e API key

# Sincronizar local <-> remoto
ttl sync
```

**Arquitetura:**
- Design multi-tenant com bancos de dados isolados por usuário
- Autenticação via API Key
- REST API para acesso programático

---

## ⚙️ Configuração

Arquivo de configuração: `~/.ttl/ttl.ini`

```ini
[default]
db_path = ~/.ttl/data.db

[server]
endpoint  = https://seu-servidor.com
api_key   = sua-user-api-key
```

```bash
# Ver configuração atual
ttl config

```

---

## 📤 Exportar Seus Dados

```bash
# Exportar como JSON
ttl export --format json

# Exportar como CSV
ttl export --format csv

# Exportar para arquivo específico
ttl export --format json --output backup.json
```

---

## 🏗️ Estrutura do Projeto

```
ttl
├── main.go              # Ponto de entrada, configuração CLI (cobra)
├── command/             # Definições de comandos CLI
│   ├── commands.go      # Comandos principais (get/add/del/tag/open...)
│   ├── log.go           # Comandos de log de trabalho
│   ├── export.go        # Comando de exportação
│   └── server.go        # Comandos do servidor
├── db/                  # Camada de armazenamento (bbolt)
│   ├── db.go            # Inicialização do banco de dados
│   ├── storage.go       # Implementação do armazenamento local
│   ├── tenant_storage.go# Roteador de armazenamento multi-tenant
│   ├── user_store.go    # CRUD de usuários (users.json)
│   └── context.go       # Armazenamento por requisição
├── api/                 # Servidor HTTP
│   ├── server.go        # Inicialização do servidor
│   ├── handlers.go      # Handlers da API REST
│   └── middleware.go     # Middleware de autenticação
├── sync/                # Lógica de sincronização de dados
├── models/              # Modelos de dados compartilhados
├── conf/                # Manipulação do arquivo de configuração (INI)
└── util/                # Funções utilitárias
```

---

## 🔧 Stack Tecnológico

| Componente | Tecnologia |
|------|------|
| Linguagem | [Go 1.23](https://golang.org) |
| Framework CLI | [cobra](https://github.com/spf13/cobra) |
| Armazenamento | [bbolt](https://github.com/etcd-io/bbolt) |
| Configuração | [ini.v1](https://gopkg.in/ini.v1) |

---

## 🌐 Traduções

- [English](README.md)
- [简体中文](README.zh-CN.md)
- [日本語](README.ja.md)
- [Español](README.es.md)
- [Français](README.fr.md)

---

## 🤝 Contribuindo

Contribuições são bem-vindas! Aqui está como você pode ajudar:

1. Faça fork do repositório
2. Crie uma branch de funcionalidade (`git checkout -b feature/amazing-feature`)
3. Commite suas mudanças (`git commit -m 'Add amazing feature'`)
4. Push para a branch (`git push origin feature/amazing-feature`)
5. Abra um Pull Request

Para mudanças grandes, por favor abra primeiro uma issue para discutir o que você gostaria de mudar.

---

## 📄 Licença

Este projeto está licenciado sob a Apache License 2.0 - veja o arquivo [LICENSE](LICENSE) para mais detalhes.

---

## 🙏 Agradecimentos

- [cobra](https://github.com/spf13/cobra) pelo excelente framework CLI
- [bbolt](https://github.com/etcd-io/bbolt) pelo armazenamento key-value embarcado confiável
- A comunidade open-source

---

<div align="center">

**Feito com ❤️ por desenvolvedores que odeiam procurar conhecimento perdido**

</div>
