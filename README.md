# Forum

This project is a web-based login and registration application, built with **Go**, **SQLite** for the database, and using **HTML/CSS** for the user interface. The application allows users to create an account, log in and access protected pages via session management.

## Run

```bash
$ make          # to build docker image and run* container
$ make start    # to just start the docker container
$ make delete   # to delete docker container and image
```

###### \*run = create and start a docker container

## Features

- Account creation with secure password hashing
- Login with password verification
- Session management to maintain user login status
- User logout
- Redirects users not logged in to the login page

## Technologies used

- **Go** (Golang)
- **SQLite** for the database
- **HTML/CSS** for user interface
- **bcrypt** for password hashing

## Author(s)

- [Clement NUTTENS](https://github.com/ClemNTTS)
- [Antoine MORLAY](https://github.com/Yssnogood)
- [Nathan SANNIER](https://github.com/Naofumi76)
- [Adam GONÇALVES](https://github.com/Shasor) (aka [Shasor#3755](https://discordapp.com/users/282816260075683841))
