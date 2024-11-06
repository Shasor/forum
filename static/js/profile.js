const div = document.querySelector(".posts-container");

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
                <a href="/">Back to Home</a>
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
          <a href="/">Back to Home</a>
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
          <a href="/">Back to Home</a>
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
          ${
            window.userData.likedPosts
              .map(
                (post) => `
            <div class="post">
                <h2>${post.title}</h2>
                <p class="author">Posted by ${post.author} on ${post.date}</p>
                <div class="content">${post.content}</div>
            </div>
          `
              )
              .join("") || "<p>No liked posts found.</p>"
          }
      </div>
  
      <footer>
          <a href="/">Back to Home</a>
      </footer>
    `;

  // Add event listeners for nav links
  document.getElementById("my-profile-link").addEventListener("click", GetProfile);
  document.getElementById("my-posts-link").addEventListener("click", GetMyPosts);
  document.getElementById("edit-profile-link").addEventListener("click", GetEditProfile);
}
