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
  const signup_bttn = document.getElementById('signup-link');
  signup_bttn.addEventListener('click',GetSignup);
  
  const close = document.getElementById('close')
  close.addEventListener("click", CloseLogin);
}

function CloseLogin(e) {
  div.innerHTML = content;
  div.className = "posts-container";
  document.removeEventListener("click", CloseLogin);
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
    const close = document.getElementById('close');  
    close.addEventListener("click", CloseLogin);
}

function CloseSignup(e) {
  div.innerHTML = content;
  div.className = "posts-container";
  document.removeEventListener("click", CloseSignup);
}
