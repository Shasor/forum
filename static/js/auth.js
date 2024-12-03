const div = document.querySelector(".posts-container");
const content = div.outerHTML;
const login_link = document.querySelector("#header-login-link");
const signup_link = document.querySelector("#header-signup-link");

export function GetLogin() {
  div.className = "auth-container";
  div.innerHTML = `
  <div class="auth">
      <buttons id="close">x</buttons>
      <h2>Sign in</h2>
      <form action="/login" method="post">
          <label for="username">Username</label>
          <input type="text" id="username" name="username" required autocomplete="username">
          <label for="password">Password</label>
          <input type="password" id="password" name="password" required autocomplete="password">
          <button type="submit">Login</button>
      </form>
      <div id="redirect">
          <p>Not registered yet ?</p>
          <button type="button" id="signup-link">Create an account</button>
      </div>
      <div id="auths">
        <a href="auth/google/login">Google</a>
        <a href="">Github</a>
        <a href="">Discord</a>
      </div>
  </div>`;
  const signup_bttn = document.getElementById("signup-link");
  signup_bttn.addEventListener("click", GetSignup);

  const close = document.getElementById("close");
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
        <h2>Sign up</h2>
        <buttons id="close">x</buttons>
        <form action="/signup" method="post">
            <label for="email">Email</label>
            <input type="email" id="email" name="email" required autocomplete="email" ${formValue.email}>
            <label for="username">Enter an Username</label>
            <input type="text" id="username" maxlength="15" name="username" required autocomplete="username" ${formValue.username}>
            <label for="password">Enter a password </label>
            <input type="password" id="password" name="password" required autocomplete="new-password">
            <button type="submit">Register</button>
        </form>
        <div id="auths">
          <a href="auth/google/login">Google</a>
          <a href="">Github</a>
          <a href="">Discord</a>
      </div>
    </div>`;
  const close = document.getElementById("close");
  close.addEventListener("click", CloseLogin);
}

function CloseSignup(e) {
  div.innerHTML = content;
  div.className = "posts-container";
  document.removeEventListener("click", CloseSignup);
}
