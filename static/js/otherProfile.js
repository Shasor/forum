const div = document.querySelector(".posts-container");

export function GetOtherProfile(){

  
  const username = event.target.getAttribute('user');
  const usernameID = event.target.getAttribute('userID')
  const userPicture = event.target.getAttribute('src')
  const userRole = event.target.getAttribute('userRole')
  console.log('Nom d\'utilisateur récupéré :', username);
  console.log('ID utilisateur : ', usernameID );
  console.log(getUserInfo(usernameID));

    if (!div){
        console.error("Error: .posts-container");
        return;
    }

    div.className = "profile-container";
    div.innerHTML = `
    <header>
      <h1>Profile of ${username}</h1>
      <h1> Bute <span id="userUsername"></span> </h1>
    </header>
  <div class="profile-menu">
      <nav>
        <ul>
          <li><a href="#myposts" id="header-profile-posts">Posts of ${username} </a></li>
          <li>Hihihihi</li>
          <li><a href="#likedposts" id="header-profile-liked">My Liked Posts</a></li>
          </ul>
      </nav>
      <div>
          <img src="${userPicture}" alt="Profile Picture of ${username}">
          <p>Role: ${userRole}</p>
      </div>
  </div>
  <div id="posts-content"> 
    <p>Activité de ${username}</p>
  
  </div>
  `;
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
      document.getElementById('userUsername').textContent = data.userData.Username || 'Non disponible';
      console.log('Informations utilisateur :', data.userPosts);
      displayPosts(data.userPosts)
    } else {
      console.error('Erreur de récupération des données:', response.status);
    }
  } catch (error) {
    console.error('Erreur lors de la requête :', error);
  }
}

async function displayPosts(posts) {

  console.log(posts.length)

  const container = document.getElementById('posts-content');

  console.log(container)
  for (let i = 1; i < posts.length; i++) {
    let categoriesHTML = '';
    if(posts[i].Categories){
    for (let k = 0; k < posts[i].Categories.length; k++) {
        categoriesHTML += `<a href="/?catID=${posts[i].Categories[k].ID}">#${posts[i].Categories[k].Name}</a>`;
      }
    }
  

    const postDiv = document.createElement('div');

    postDiv.innerHTML= `
    <div class="post">
    <div class="post-header">
    <div class="sender">${posts[i].Sender.Username}</div>
    <div class="category"><a href="?catID=${posts[i].Categories[0].ID}"> ${posts[i].Categories[0].Name} </a></div>
    <div class="date">${posts[i].Date}</div>
  </div>
  <div class="title">${posts[i].Title}</div>
  <div class="post-content">
    <div class="content">${posts[i].Content}</div>
    <div class="post-image"><img src="data:image/jpeg;base64,${posts[i].Picture}" /></div>
  </div>
  <div class="reactions">
    <div class="post-categories">
        ${categoriesHTML}
    </div>
      <div class="reaction">
          <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="currentcolor"><path d="m480-120-58-52q-101-91-167-157T150-447.5Q111-500 95.5-544T80-634q0-94 63-157t157-63q52 0 99 22t81 62q34-40 81-62t99-22q94 0 157 63t63 157q0 46-15.5 90T810-447.5Q771-395 705-329T538-172l-58 52Zm0-108q96-86 158-147.5t98-107q36-45.5 50-81t14-70.5q0-60-40-100t-100-40q-47 0-87 26.5T518-680h-76q-15-41-55-67.5T300-774q-60 0-100 40t-40 100q0 35 14 70.5t50 81q36 45.5 98 107T480-228Zm0-273Z" /></svg>
        <span class="like-count" >${posts[i].Likes}</span>
      </div>
      <div class="reaction">
          <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="currentcolor"><path d="M240-840h440v520L400-40l-50-50q-7-7-11.5-19t-4.5-23v-14l44-174H120q-32 0-56-24t-24-56v-80q0-7 2-15t4-15l120-282q9-20 30-34t44-14Zm360 80H240L120-480v80h360l-54 220 174-174v-406Zm0 406v-406 406Zm80 34v-80h120v-360H680v-80h200v520H680Z" /></svg>
        <span class="dislike-count" >${posts[i].Dislikes}</span>
      </div>
    <div class="reaction">
        <button class="reaction-button" id="comments">
          <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="currentcolor"><path d="M880-80 720-240H320q-33 0-56.5-23.5T240-320v-40h440q33 0 56.5-23.5T760-440v-280h40q33 0 56.5 23.5T880-640v560ZM160-473l47-47h393v-280H160v327ZM80-280v-520q0-33 23.5-56.5T160-880h440q33 0 56.5 23.5T680-800v280q0 33-23.5 56.5T600-440H240L80-280Zm80-240v-280 280Z" /></svg>
        </button>
    </div>
  </div>
  </div>
</div>`;

    container.appendChild(postDiv);
  }
}