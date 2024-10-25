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

export function ShowError(texte) {
  // Créer un nouvel élément div
  const div = document.createElement("div");
  div.id = "response";

  // Définir le contenu du div
  div.textContent = texte;

  // Styles pour le rectangle flottant
  div.style.position = "fixed";
  div.style.left = "50%";
  div.style.top = "50%";
  div.style.transform = "translate(-50%, -50%)";
  div.style.padding = "15px";
  div.style.backgroundColor = "rgba(0, 0, 0, 0.7)";
  div.style.color = "white";
  div.style.borderRadius = "5px";
  div.style.boxShadow = "0 2px 10px rgba(0, 0, 0, 0.2)";
  div.style.zIndex = "1000";
  div.style.maxWidth = "80%";
  div.style.textAlign = "center";
  div.style.opacity = "0";
  div.style.transition = "opacity 0.5s ease-in-out";

  // Ajouter le div au corps du document
  document.body.appendChild(div);

  // Animation d'apparition
  setTimeout(() => {
    div.style.opacity = "1";
  }, 100);

  // Fonction pour faire disparaître le rectangle
  const fermer = () => {
    if (div && div.parentNode) {
      div.style.opacity = "0";
      setTimeout(() => {
        if (div.parentNode) {
          div.parentNode.removeChild(div);
        }
      }, 500);
    }
  };

  // Ajouter un bouton de fermeture
  const btnFermer = document.createElement("button");
  btnFermer.textContent = "X";
  btnFermer.style.position = "absolute";
  btnFermer.style.top = "5px";
  btnFermer.style.right = "5px";
  btnFermer.style.background = "none";
  btnFermer.style.border = "none";
  btnFermer.style.color = "white";
  btnFermer.style.cursor = "pointer";
  btnFermer.onclick = fermer;
  div.appendChild(btnFermer);

  // Fermeture automatique après 5 secondes
  setTimeout(fermer, 5000);
}
