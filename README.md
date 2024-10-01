# Projet Go : Système de Connexion et d'Inscription

Ce projet est une application web de gestion de connexion et d'inscription, construite avec **Go**, **SQLite** pour la base de données, et utilisant **HTML/CSS** pour l'interface utilisateur. L'application permet aux utilisateurs de créer un compte, de se connecter et d'accéder à des pages protégées via la gestion des sessions.

## Fonctionnalités

- Création de compte avec hashage sécurisé des mots de passe.
- Connexion avec vérification du mot de passe.
- Gestion des sessions pour maintenir l'état de connexion des utilisateurs.
- Déconnexion des utilisateurs.
- Redirection des utilisateurs non connectés vers la page de connexion.

## Technologies Utilisées

- **Go** (Golang)
- **SQLite** pour la base de données
- **HTML/CSS** pour l'interface utilisateur
- **bcrypt** pour le hashage des mots de passe
- **gorilla/sessions** pour la gestion des sessions