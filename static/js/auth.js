const div = document.querySelector(".posts-container");
const content = div.outerHTML;
const login_link = document.querySelector("#header-login-link");
const signup_link = document.querySelector("#header-signup-link");

export function GetLogin() {
  div.className = "auth-container";
  div.innerHTML = `
  <div class="auth">
      <buttons id="close">x</buttons>
      <h2>Connexion</h2>
      <form action="/login" method="post">
          <label for="username">Nom d'utilisateur</label>
          <input type="text" id="username" name="username" required autocomplete="username">
          <label for="password">Mot de passe</label>
          <input type="password" id="password" name="password" required autocomplete="password">
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
  const login = document.querySelector(".auth");
  const response = document.querySelector("#response");
  const close = document.querySelector("#close");
  if (!((login?.contains(e.target) && e.target !== close) || login_link.contains(e.target) || signup_link.contains(e.target) || response?.contains(e.target))) {
    div.innerHTML = content;
    div.className = "posts-container";
    document.removeEventListener("click", CloseLogin);
  }
}

export function GetSignup(formValue) {
  div.className = "auth-container";
  div.innerHTML = `
    <div class="auth">
        <h2>Créer un compte</h2>
        <buttons id="close">x</buttons>
        <form action="/signup" method="post">
            <label for="email">Email</label>
            <input type="email" id="email" name="email" required autocomplete="email" ${formValue.email}">
            <label for="username">Nom d'utilisateur</label>
            <input type="text" id="username" name="username" required autocomplete="username" ${formValue.username}">
            <label for="password">Mot de passe</label>
            <input type="password" id="password" name="password" required autocomplete="new-password">
            <button type="submit">S'inscrire</button>
        </form>
    </div>`;
  document.addEventListener("click", CloseSignup);
}

function CloseSignup(e) {
  const signup = document.querySelector(".auth");
  const response = document.querySelector("#response");
  const close = document.querySelector("#close");
  if (!((signup?.contains(e.target) && e.target !== close) || signup_link.contains(e.target) || login_link.contains(e.target) || response?.contains(e.target))) {
    div.innerHTML = content;
    div.className = "posts-container";
    document.removeEventListener("click", CloseSignup);
  }
}
