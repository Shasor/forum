const div = document.querySelector(".posts-container");

export function GetOtherProfile(){

  
  const username = event.target.getAttribute('user');
  const usernameID = event.target.getAttribute('userID')
  const userPicture = event.target.getAttribute('src')
  const userRole = event.target.getAttribute('userRole')
  console.log('Nom d\'utilisateur récupéré :', username);
  console.log('ID utilisateur : ', usernameID );
  console.log(getUserInfo(usernameID))

    if (!div){
        console.error("Error: .posts-container");
        return;
    }

    div.className = "profile-container";
    div.innerHTML = `
    <header>
      <h1>Profile of ${username}</h1>
    </header>
  <div class="profile-menu">
      <nav>
        <ul>
          <li><a href="#myposts" id="header-profile-posts">Posts of ${username} </a></li>
          <li>aled</li>
          <li><a href="#likedposts" id="header-profile-liked">My Liked Posts</a></li>
          </ul>
      </nav>
      <div>
          <img src="${userPicture}" alt="Profile Picture of ${username}">
          <p>Role: ${userRole}</p>
      </div>
  </div>
  
  ${window.otheruserData.myposts}

  `;

  


  // Optionally, re-attach other event listeners if needed
  const postsLink = document.getElementById("header-profile-posts");
  postsLink?.addEventListener("click", (event) => {
    event.preventDefault();
    GetOtherPosts(); // Load user's posts when clicking "My Posts"
  });


}


export function GetOtherPosts() {
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
        ${window.otheruserData.myposts}
      </div>
    `;

  // Add event listeners for nav links
  document.getElementById("my-profile-link").addEventListener("click", GetOtherProfile);
}


// Fonction pour récupérer les informations d'un utilisateur
async function getUserInfo(userId) {
  try {
    // Envoi de la requête GET avec l'ID de l'utilisateur
    const response = await fetch(`http://localhost:8080/users/${userId}`, {
      method: 'GET', // Méthode GET pour récupérer les données
      headers: {
        'Content-Type': 'application/json', // Type de contenu JSON
      },
    });

    // Si la réponse est ok, on la transforme en JSON
    if (response.ok) {
      const data = await response.json();
      console.log('Informations utilisateur :', data);
    } else {
      console.error('Erreur de récupération des données:', response.status);
    }
  } catch (error) {
    console.error('Erreur lors de la requête :', error);
  }
}