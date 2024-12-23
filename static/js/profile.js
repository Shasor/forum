import { GetOtherProfile } from "./otherProfile.js";
const div = document.querySelector(".posts-container");

export function GetProfile() {
  if (!div) {
    console.error("Error: .posts-container not found!");
    return;
  }

  div.className = "profile-container";
  div.innerHTML = `
          <header>
            <h1>Profile of ${window.userData.pseudo}</h1>
          </header>
        <div class="profile-menu">
          <nav>
            <ul id="profile-menu-list">
              <li><a href="#profil-activity" other_id="${window.userData.id}"  user_id="${window.userData.id}" user_role="${window.userData.role}" id="header-profile-activity-link">Lastest activities</a></li>
              <li><a href="#myposts" id="header-profile-posts">My Posts</a></li>
              <li><a href="#likedposts" id="header-profile-liked">My Liked Posts</a></li>
              <li><a href="#edit" id="header-profile-edit">Edit Profile</a></li>
              <form method="post" action="/logout">
                <li><button type="submit"> Logout </button></li>
              </form>
            </ul>
          </nav>
          <div>
            <img src="${window.userData.profilePicture}" alt="Profile Picture of ${window.userData.pseudo}">
            <p>Email: ${window.userData.email}</p>
            <p>Role: ${window.userData.role}</p>
            ${window.userData.request}
          </div>
        </div>`;

  const activitiesLink = document.getElementById("header-profile-activity-link");
  activitiesLink?.addEventListener("click", (event) => {
    event.preventDefault();
    GetOtherProfile();
  });

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
  
      <form action="/edit" method="POST" enctype="multipart/form-data">
          <div>
              <label for="email">Email:</label>
              <input type="email" id="email" name="email" placeholder="${window.userData.email}">
          </div>
  
          <div>
              <label for="profile_picture">Profile Picture:</label>
              <input type="file" id="profile_picture" name="profile_picture" accept=".jpeg, .png, .gif">
              <p id="profile-picture-text">Current profile picture:</p>
              <img id="current-profile-picture" src="${window.userData.profilePicture}" alt="Current Profile Picture" style="max-width: 150px; max-height: 150px;">
          </div>
  
          <button type="submit">Update Profile</button>
      </form>
  
      <form action="/delete" method="POST">
          <input type="submit" value="Supprimer mon compte" onclick="return confirm('Êtes-vous sûr de vouloir supprimer votre compte ? Cela est irréversible.');">
      </form>

    `;

  // Attach click events to navigation links
  document.getElementById("my-profile-link").addEventListener("click", GetProfile);
  document.getElementById("my-posts-link").addEventListener("click", GetMyPosts);
  document.getElementById("my-liked-posts-link").addEventListener("click", GetLikedPosts);


  
  const profilePictureInput = document.getElementById("profile_picture");
  const profilePicturePreview = document.getElementById("current-profile-picture");
  const profilePictureText = document.getElementById("profile-picture-text");
  
  profilePictureInput.addEventListener("change", (event) => {
    const file = event.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (e) => {
        profilePicturePreview.src = e.target.result; // Update the image preview
        profilePictureText.textContent = "Your new profile picture will be:"; // Update the text
      };
      reader.readAsDataURL(file);
    }
  });
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
          ${window.userData.likedposts || "<p>No liked posts found.</p>"}
      </div>
    `;

  // Add event listeners for nav links
  document.getElementById("my-profile-link").addEventListener("click", GetProfile);
  document.getElementById("my-posts-link").addEventListener("click", GetMyPosts);
  document.getElementById("edit-profile-link").addEventListener("click", GetEditProfile);
}

export function SeeEveryUser() {
  {
    div.className = "posts-container";
    div.innerHTML = `
      <header>
          <h1>Manage users : </h1>
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
    `;

    // Add event listeners for nav links
    document.getElementById("my-profile-link").addEventListener("click", GetProfile);
    document.getElementById("my-liked-posts-link").addEventListener("click", GetLikedPosts);
    document.getElementById("edit-profile-link").addEventListener("click", GetEditProfile);
  }
}
