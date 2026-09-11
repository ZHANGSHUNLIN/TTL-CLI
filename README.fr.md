<div align="center">

# ttl

### Votre Archive Personnelle de Connaissances

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

[English](README.md) | [简体中文](README.zh-CN.md) | [日本語](README.ja.md) | [Español](README.es.md) | [Português](README.pt.md)

---


</div>

---

## 📖 L'Histoire

Chaque développeur y a déjà été confronté :

> "Attends, c'était quoi cette commande Docker que j'ai utilisée le mois dernier ?"
> "Où mon collègue a-t-il partagé ce fichier de configuration ?"
> "Je sais que j'ai lu un article là-dessus... mais je ne le retrouve pas."

Nous stockons nos connaissances partout — onglets du navigateur, messages Slack, emails, dossiers de favoris, applications de notes. Quand nous en avons vraiment besoin, nous perdons du temps à fouiller dans des onglets sans fin et à défiler l'historique des conversations.

**C'était aussi mon problème.**

C'est pourquoi j'ai créé **ttl**.

Le nom vient de "Time to Live" — mais avec une signification différente. Au lieu d'expirer, il s'agit de donner à vos connaissances un **temps pour vivre** éternellement.

- Stockez tout au même endroit sous forme de paires clé-valeur
- Étiquetez pour une organisation facile
- Recherchez instantanément par mot-clé

Plus besoin de chercher dans les vieux emails ou de défiler l'historique Slack. Juste `ttl get <mot-clé>` et c'est là.

**ttl est votre archive personnelle de connaissances — tout ce dont vous avez besoin, quand vous en avez besoin.**

---

## ✨ Fonctionnalités

| Fonctionnalité | Description |
|------|------|
| 🗄️ **Stockage KV Local** | Base de données embarquée rapide, sans configuration (bbolt) |
| 🏷️ **Système de Tags** | Organisez les ressources avec des tags flexibles et recherchables |
| 🔍 **Recherche Floue** | Trouvez ce dont vous avez besoin instantanément parmi clés et tags |
| ☁️ **Sync Cloud** | Auto-hébergez un serveur multi-tenant avec isolation des données par utilisateur |
| 🚀 **Ouverture Intelligente** | Ouvrez URLs et fichiers avec les programmes par défaut du système |
| 📤 **Export** | Exportez les données en JSON ou CSV |

---

## 🚀 Démarrage Rapide

### Installation

#### Linux / macOS

```bash
# Installer depuis GitHub releases
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.sh)"

# Ou compiler depuis les sources
go build -o ttl .
sudo mv ttl /usr/local/bin/
```

#### Windows

```powershell
# Installer depuis GitHub releases
irm https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.ps1 | iex
```

#### URL de Téléchargement Personnalisé

Pour les réseaux internes ou mirrors personnalisés :

```bash
# Linux/macOS
TTL_DOWNLOAD_URL="https://your-mirror.com/ttl-cli-v1.0.0-linux-amd64" /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.sh)"
```

```powershell
# Windows
$env:TTL_DOWNLOAD_URL="https://your-mirror.com/ttl-cli-v1.0.0-windows-amd64.zip"; irm https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.ps1 | iex
```

### Utilisation de Base

```bash
# Ajouter une ressource
ttl add my-link https://example.com

# Ajouter avec des tags
ttl add docker-cmd "docker run -d -p 8080:80 nginx"
ttl tag docker-cmd dev ops

# Rechercher des ressources
ttl get docker

# Ouvrir dans le navigateur
ttl open my-link

# Supprimer
ttl del old-key
```

---


## 📝 Journal de Travail

Suivez et filtrez votre travail quotidien.

```bash
# Écrire une entrée de journal
ttl log write "Refactoring du module utilisateur terminé" --tags "projetA,dev"

# Voir les journaux
ttl log list                    # Journaux d'aujourd'hui
ttl log list --range week       # Cette semaine
ttl log list --range month      # Ce mois

```

---


## ☁️ Serveur Cloud & Synchronisation

### Démarrer Votre Serveur

```bash
# Créer un utilisateur
ttl server user add alice

# Démarrer le serveur multi-tenant
ttl server start --port 8080
```

### Synchroniser Vos Données

```bash
# Configurer le serveur distant
ttl config
# Éditez la section server avec votre endpoint et clé API

# Synchroniser local <-> distant
ttl sync
```

**Architecture :**
- Conception multi-tenant avec bases de données isolées par utilisateur
- Authentification par clé API
- API REST pour accès programmatique

---

## ⚙️ Configuration

Fichier de configuration : `~/.ttl/ttl.ini`

```ini
[default]
db_path = ~/.ttl/data.db

[server]
endpoint  = https://votre-serveur.com
api_key   = votre-user-api-key
```

```bash
# Voir la configuration actuelle
ttl config

```

---

## 📤 Exporter Vos Données

```bash
# Exporter en JSON
ttl export --format json

# Exporter en CSV
ttl export --format csv

# Exporter vers un fichier spécifique
ttl export --format json --output backup.json
```

---

## 🏗️ Structure du Projet

```
ttl
├── main.go              # Point d'entrée, configuration CLI (cobra)
├── command/             # Définitions des commandes CLI
│   ├── commands.go      # Commandes principales (get/add/del/tag/open...)
│   ├── log.go           # Commandes journal de travail
│   ├── export.go        # Commande d'export
│   └── server.go        # Commandes serveur
├── db/                  # Couche de stockage (bbolt)
│   ├── db.go            # Initialisation de la base de données
│   ├── storage.go       # Implémentation du stockage local
│   ├── tenant_storage.go# Routeur de stockage multi-tenant
│   ├── user_store.go    # CRUD utilisateurs (users.json)
│   └── context.go       # Stockage par requête
├── api/                 # Serveur HTTP
│   ├── server.go        # Démarrage du serveur
│   ├── handlers.go      # Handlers API REST
│   └── middleware.go     # Middleware d'authentification
├── sync/                # Logique de synchronisation des données
├── models/              # Modèles de données partagés
├── conf/                # Gestion du fichier de configuration (INI)
└── util/                # Fonctions utilitaires
```

---

## 🔧 Stack Technique

| Composant | Technologie |
|------|------|
| Langage | [Go 1.23](https://golang.org) |
| Framework CLI | [cobra](https://github.com/spf13/cobra) |
| Stockage | [bbolt](https://github.com/etcd-io/bbolt) |
| Configuration | [ini.v1](https://gopkg.in/ini.v1) |

---

## 🌐 Traductions

- [English](README.md)
- [简体中文](README.zh-CN.md)
- [日本語](README.ja.md)
- [Español](README.es.md)
- [Português](README.pt.md)

---

## 🤝 Contribuer

Les contributions sont les bienvenues ! Voici comment vous pouvez aider :

1. Forkez le dépôt
2. Créez une branche de fonctionnalité (`git checkout -b feature/amazing-feature`)
3. Commitez vos changements (`git commit -m 'Add amazing feature'`)
4. Poussez vers la branche (`git push origin feature/amazing-feature`)
5. Ouvrez une Pull Request

Pour les changements majeurs, veuillez d'abord ouvrir un issue pour discuter de ce que vous souhaitez modifier.

---

## 📄 Licence

Ce projet est sous licence Apache License 2.0 - voir le fichier [LICENSE](LICENSE) pour plus de détails.

---

## 🙏 Remerciements

- [cobra](https://github.com/spf13/cobra) pour l'excellent framework CLI
- [bbolt](https://github.com/etcd-io/bbolt) pour le stockage key-value embarqué fiable
- La communauté open-source

---

<div align="center">

**Fait avec ❤️ par des développeurs qui détestent chercher des connaissances perdues**

</div>
