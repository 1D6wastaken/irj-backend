# IRJ – Backend API

## 📖 Description

Ce projet constitue l’API backend de l’application **IRJ**.
Il expose un ensemble d’API (publiques, privées et administrateur) définies dans le contrat Swagger (`swagger.yaml`) et utilisées pour la gestion :

* des **monuments et lieux**
* des **mobiliers et images**
* des **personnes morales** (entités légales)
* des **personnes physiques** (individus)
* de la **recherche avancée** (avec filtres : pays, régions, siècles, professions, etc.)
* de la **gestion des utilisateurs** (inscription, connexion, validation email, mot de passe oublié, etc.)

Le backend est développé en **Go**, utilise **sqlc** pour la génération de code lié à la base de données, et intègre un système de **migrations versionnées avec go-migrate** pour assurer l’évolution du schéma de la base.

---

## 🚀 Prérequis

* **Go >= 1.20**
* **Taskfile** (outil de gestion de tâches : [https://taskfile.dev](https://taskfile.dev))
* **Docker / PostgreSQL** (pour la base de données)
* **go-migrate** (gestion des migrations)

---

## 📂 Structure du projet

* `api/swagger.yaml` : définition OpenAPI/Swagger des endpoints REST
* `api/sqlc.yaml` : configuration de génération SQLC
* `internal/postgres/_generated` : code généré pour l’accès aux données
* `pkg/api`, `pkg/dto` : modèles générés depuis Swagger
* `migrations/` : scripts SQL versionnés avec **go-migrate**

---

## ⚡️ Commandes utiles

Les principales tâches sont définies dans le `Taskfile.yml`.

### Lancer les tests unitaires

```bash
task test
```

### Vérifier le lint

```bash
task lint
```

### Générer le code à partir de Swagger et SQLC

```bash
task generate
```

### Vérifier la couverture de tests

```bash
task cover
```

### Vérifier les dépendances

```bash
task deps-check
```

### Mettre à jour les dépendances

```bash
task deps-upgrade
```

### Nettoyer le projet

```bash
task clean
```

---

## 🗄️ Gestion des migrations

Le projet utilise [golang-migrate/migrate](https://github.com/golang-migrate/migrate) pour gérer le versioning de la base de données.

### Créer une migration

```bash
migrate create -ext sql -dir migrations -seq add_new_table
```

### Appliquer les migrations

```bash
migrate -path migrations -database "$DATABASE_URL" up
```

### Revenir en arrière

```bash
migrate -path migrations -database "$DATABASE_URL" down 1
```

---

## 🔑 Sécurité

L’API repose sur une authentification **JWT Bearer** pour les endpoints nécessitant une connexion.

* Les APIs **publiques** sont accessibles sans authentification.
* Les APIs **privées** nécessitent un utilisateur connecté.
* Les APIs **admin** sont réservées aux administrateurs.

---

## 📑 Documentation API

La documentation complète de l’API est disponible dans le fichier :

```bash
api/swagger.yaml
```

Elle peut être visualisée via [Swagger UI](https://swagger.io/tools/swagger-ui/) ou importée dans Postman/Insomnia.

---

## ✅ CI/CD & Qualité

* **golangci-lint** : vérification du code
* **gotestsum** : exécution et reporting des tests
* **govulncheck** : audit de vulnérabilités Go

---
