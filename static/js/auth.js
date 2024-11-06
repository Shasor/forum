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


export function GetProfile() {
  if (!div) {
    console.error("Error: .posts-container not found!");
    return;
  }

  div.className = "profile-container";
  div.innerHTML = `
      <div class="profile">
          <header>
              <h1>Profile of ${window.userData.pseudo}</h1>
          </header>
          <nav>
            <ul>
              <li><a href="#" id="header-profile-link">Voir mon profil</a></li>
              <li><a href="#" id="header-profile-posts">My Posts</a></li>
              <li><a href="#" id="header-profile-liked">My Liked Posts</a></li>
              <li><a href="#" id="header-profile-edit">Edit Profile</a></li>
              <form method="post" action="/logout">
                <li><button type="submit">Se déconnecter</button></li>
              </form>
            </ul>
          </nav>
          <div>
              <img src="${window.userData.profilePicture}" alt="Profile Picture of ${window.userData.pseudo}">
              <p>Email: ${window.userData.email}</p>
              <p>Role: ${window.userData.role}</p>
          </div>
          <footer>
              <a href="/dashboard">Back to Dashboard</a>
          </footer>
      </div>`;

  // Attach event listener for "Edit Profile" after updating inner HTML
  const editProfileLink = document.getElementById("header-profile-edit");
  editProfileLink?.addEventListener("click", (event) => {
    event.preventDefault();
    GetEditProfile(); // Call the function to load the edit profile view
  });

  // Optionally, re-attach other event listeners if needed
  const postsLink = document.getElementById("header-profile-posts");
  postsLink?.addEventListener("click", (event) => {
    event.preventDefault();
    GetMyPosts(); // Load user's posts when clicking "My Posts"
  });

  const likedPostsLink = document.getElementById("header-profile-liked");
  likedPostsLink?.addEventListener("click", (event) => {
    event.preventDefault();
    GetLikedPosts(); // Load liked posts when clicking "My Liked Posts"
  });
}





function CloseProfile(e) {
  const profile = document.querySelector(".profile");
  const response = document.querySelector("#response");
  const close = document.querySelector("#close");
  if (!((profile?.contains(e.target) && e.target !== close) || profile_link.contains(e.target) || signup_link.contains(e.target) || response?.contains(e.target))) {
    div.innerHTML = content;
    div.className = "posts-container";
    document.removeEventListener("click", CloseProfile);
  }
}



export function GetEditProfile() {
  div.className = "edit-profile-container";
  div.innerHTML = `
    <header>
        <h1>Edit Profile of ${window.userData.pseudo}</h1>
    </header>

    <nav>
        <a href="#" id="my-profile-link">My Profile</a>
        <a href="#" id="my-posts-link">My Posts</a> 
        <a href="#" id="my-liked-posts-link">My Liked Posts</a>
        <a href="#" id="edit-profile-link">Edit Profile</a>
    </nav>

    <form action="/profile/edit" method="POST" enctype="multipart/form-data">
        <div>
            <label for="email">Email:</label>
            <input type="email" id="email" name="email" placeholder="${window.userData.email}">
        </div>

        <div>
            <label for="profile_picture">Profile Picture:</label>
            <input type="file" id="profile_picture" name="profile_picture" accept=".jpeg, .png, .gif">
            <p>Current profile picture:</p>
            <img src="${window.userData.profilePicture}" alt="Current Profile Picture" style="max-width: 150px; max-height: 150px;">
        </div>

        <button type="submit">Update Profile</button>
    </form>

    <form action="/delete-account" method="POST">
        <input type="submit" value="Supprimer mon compte" onclick="return confirm('Êtes-vous sûr de vouloir supprimer votre compte ? Cela est irréversible.');">
    </form>

    <footer>
        <a href="/dashboard">Back to Dashboard</a>
    </footer>
  `;

  // Attach click events to navigation links
  document.getElementById("my-profile-link").addEventListener("click", GetProfile);
  document.getElementById("my-posts-link").addEventListener("click", GetMyPosts); // Placeholder function
  document.getElementById("my-liked-posts-link").addEventListener("click", GetLikedPosts); // Placeholder function
}

// Function to display the My Posts section
export function GetMyPosts() {
  div.className = "posts-container";
  div.innerHTML = `
    <header>
        <h1>Your Posts</h1>
    </header>

    <nav>
        <a href="#" id="my-profile-link">My Profile</a>
        <a href="#" id="my-posts-link">My Posts</a>
        <a href="#" id="my-liked-posts-link">My Liked Posts</a>
        <a href="#" id="edit-profile-link">Edit Profile</a>
    </nav>


    <div id="posts-content">
      ${window.userData.myposts}
    </div>

    <footer>
        <a href="/dashboard">Back to Dashboard</a>
    </footer>
  `;
  
  // Add event listeners for nav links
  document.getElementById("my-profile-link").addEventListener("click", GetProfile);
  document.getElementById("my-liked-posts-link").addEventListener("click", GetLikedPosts);
  document.getElementById("edit-profile-link").addEventListener("click", GetEditProfile);
}

// Function to display the My Liked Posts section
export function GetLikedPosts() {
  div.className = "posts-container";
  div.innerHTML = `
    <header>
        <h1>Your Liked Posts</h1>
    </header>

    <nav>
        <a href="#" id="my-profile-link">My Profile</a>
        <a href="#" id="my-posts-link">My Posts</a>
        <a href="#" id="my-liked-posts-link">My Liked Posts</a>
        <a href="#" id="edit-profile-link">Edit Profile</a>
    </nav>

    <div id="liked-posts-content">
        ${window.userData.likedPosts.map(post => `
          <div class="post">
              <h2>${post.title}</h2>
              <p class="author">Posted by ${post.author} on ${post.date}</p>
              <div class="content">${post.content}</div>
          </div>
        `).join('') || '<p>No liked posts found.</p>'}
    </div>

    <footer>
        <a href="/dashboard">Back to Dashboard</a>
    </footer>
  `;
  
  // Add event listeners for nav links
  document.getElementById("my-profile-link").addEventListener("click", GetProfile);
  document.getElementById("my-posts-link").addEventListener("click", GetMyPosts);
  document.getElementById("edit-profile-link").addEventListener("click", GetEditProfile);
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
