const content = document.getElementById("posts").textContent;
const div = document.querySelector(".posts");
const login_link = document.querySelector("#login-link");
const signup_link = document.querySelector("#header-signup-link");

export function GetLogin() {
  document.getElementById("posts").innerHTML = `
  <div class="login">
      <buttons id="close">x</buttons>
      <h2>Connexion</h2>
      <form action="/login" method="post">
          <label for="username">Nom d'utilisateur</label>
          <input type="text" id="username" name="username" required>
          <label for="password">Mot de passe</label>
          <input type="password" id="password" name="password" required>
          <button type="submit">Se connecter</button>
      </form>
      <div id="redirect">
          <p>Pas encore de compte ?</p>
            <button type="button" id="signup-link">Créer un compte</button>
      </div>
  </div>`;
  document.addEventListener("click", CloseLogin);
}

function CloseLogin(e) {
  const login = document.querySelector(".login");
  if (login && !(login.contains(e.target) || login_link.contains(e.target))) {
    div.innerHTML = content;
    document.removeEventListener("click", CloseLogin);
  }
}

export function GetSignup() {
  document.getElementById("posts").innerHTML = `
    <div class="signup">
        <h2>Créer un compte</h2>
        <buttons id="close">x</buttons>
        <form action="/signup" method="post">
            <label for="email">Email</label>
            <input type="email" id="email" name="email" required autocomplete="email">
            <label for="pseudo">Nom d'utilisateur</label>
            <input type="text" id="pseudo" name="pseudo" required autocomplete="username">
            <label for="password">Mot de passe</label>
            <input type="password" id="password" name="password" required autocomplete="new-password">
            <button type="submit">S'inscrire</button>
        </form>
    </div>`;
  document.addEventListener("click", CloseSignup);
}

function CloseSignup(e) {
  const signup = document.querySelector(".signup");
  if (signup && !(signup.contains(e.target) || signup_link.contains(e.target))) {
    div.innerHTML = content;
    document.removeEventListener("click", CloseSignup);
  }
}
