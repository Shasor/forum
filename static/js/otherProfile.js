const div = document.querySelector(".posts-container");

// Fonction pour récupérer les informations d'un utilisateur
export async function GetOtherProfile() {
  try {
    // Récupération des attributs depuis l'événement
    const otherID = event.target.getAttribute("other_id");
    const userID = event.target.getAttribute("user_id");
    const userRole = event.target.getAttribute("user_role");

    console.log(otherID)

    if (!otherID) {
      console.error("L'attribut 'other_id' est manquant !");
      return;
    }

    const user = {
      id: userID,
      role: userRole,
    };

    // Envoi de la requête au serveur
    const response = await fetch(`https://localhost:8080/users`, {
      method: "POST", // Changez en GET si nécessaire

      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ otherID }),
    });

    // Traitement de la réponse
    if (response.ok) {
      const data = await response.json();
      console.log(data.otherActivity);
      displayProfile(data.otherData, data.otherActivity, user);
    } else {
      const errorDetails = await response.text();
      console.error("Erreur de récupération des données:", response.status, errorDetails);
    }
  } catch (error) {
    console.error("Erreur lors de la requête :", error.message, error.stack);
  }
}


async function displayProfile(otherData, otherActivity, user) {
  let picture;
  if (otherData.Picture) {
    picture = `data:image/jpeg;base64,${otherData.Picture}`;
  } else {
    picture = "/static/assets/img/default_profile_picture.png";
  }

  let selected;
  if (otherData.Role === "user") {
    selected = `selected`;
  }
  let role;
  if (user.role === "admin" && otherData.Role != "admin") {
    role = `
    <p>Role : 
      <form action="/role" method="POST">
        <select id="role" name="role">
          <option value="moderator">Moderator</option>
          <option value="user" ${selected}>User</option>
        </select>
        <input hidden id="otherID" name="otherID" value="${otherData.ID}"/>
        <input type="submit" value="Send!"/>
      </form>
    </p>`;
  } else {
    role = `<p>Role : ${otherData.Role}</p>`;
  }

  div.innerHTML = `
    <header>
      <h1>Profile of ${otherData.Username}</h1>
    </header>
    <div class="profile-menu">
      <div>
        <img src="${picture}" alt="Profile Picture of ${otherData.Username}">
        ${role}
      </div>
    </div>
    <div id="posts-content"> 
      <p><h2> Latest activities of ${otherData.Username} : </h2></p>
    </div>`;

  const container = document.getElementById("posts-content");
  for (let i = 0; i < otherActivity.length; i++) {
    //Check if the post has multiples categories or not
    let categoriesHTML = "";
    if (otherActivity[i].Post.Categories) {
      for (let k = 0; k < otherActivity[i].Post.Categories.length; k++) {
        categoriesHTML += `<a href="/?catID=${otherActivity[i].Post.Categories[k].ID}">#${otherActivity[i].Post.Categories[k].Name}</a>`;
      }
    }

    let typeHTML;
    if (otherActivity[i].Action === "post") {
      typeHTML = `<p>${otherActivity[i].Post.Sender.Username} posted on ${categoriesHTML} :</p>`;
    }
    if (otherActivity[i].Action === "comment") {
      typeHTML = `<p>${otherData.Username} commented to someone post <a href="/?postID=${otherActivity[i].Post.ParentID}">here</a>:</p>`;
    }
    if (otherActivity[i].Action === "LIKE") {
      typeHTML = `<p>${otherData.Username} liked:</p>`;
    }
    if (otherActivity[i].Action === "DISLIKE") {
      typeHTML = `<p>${otherData.Username} disliked:</p>`;
    }

    // Set the correct Avatar on the post
    let avatarHTML;
    if (otherActivity[i].Post.Sender.Picture) {
      avatarHTML = `<img src="data:image/jpeg;base64,${otherActivity[i].Post.Sender.Picture}" alt="Profile Picture" id="avatar-post" style="max-width: 150px; max-height: 150px" />`;
    } else {
      avatarHTML = `<img src="/static/assets/img/default_profile_picture.png" alt="Default Profile Picture" id="avatar-post" style="max-width: 150px; max-height: 150px" />`;
    }

    // Check if a picture is present inside the post
    let pictureHTML;
    if (otherActivity[i].Post.Picture) {
      pictureHTML = `<div class="post-image"><img src="data:image/jpeg;base64,${otherActivity[i].Post.Picture}" /></div>`;
    } else {
      pictureHTML = ``;
    }

    const postDiv = document.createElement("div");
    postDiv.className = "post";

    postDiv.innerHTML = `
    ${typeHTML}
      <div class="post-header"> 
        <div class="sender"> 
          ${avatarHTML}
          <p id = "user_name"   > ${otherActivity[i].Post.Sender.Username} </p>
        </div>
        <div class="category"><a href="?catID=${otherActivity[i].Post.Categories[0].ID}"> ${otherActivity[i].Post.Categories[0].Name} </a></div>
        <div class="date">${otherActivity[i].Post.Date}</div>
      </div>
      <div class="title">${otherActivity[i].Post.Title}</div>
      <div class="post-content">
        <div class="content">${otherActivity[i].Post.Content}</div>
        ${pictureHTML}
        </div>
      <div class="reactions">
        <div class="post-categories">
          ${categoriesHTML}
        </div>
        <div class="reaction">
          <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="currentcolor"><path d="m480-120-58-52q-101-91-167-157T150-447.5Q111-500 95.5-544T80-634q0-94 63-157t157-63q52 0 99 22t81 62q34-40 81-62t99-22q94 0 157 63t63 157q0 46-15.5 90T810-447.5Q771-395 705-329T538-172l-58 52Zm0-108q96-86 158-147.5t98-107q36-45.5 50-81t14-70.5q0-60-40-100t-100-40q-47 0-87 26.5T518-680h-76q-15-41-55-67.5T300-774q-60 0-100 40t-40 100q0 35 14 70.5t50 81q36 45.5 98 107T480-228Zm0-273Z" /></svg>
          <span class="like-count" >${otherActivity[i].Post.Likes}</span>
        </div>
        <div class="reaction">
          <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="currentcolor"><path d="M240-840h440v520L400-40l-50-50q-7-7-11.5-19t-4.5-23v-14l44-174H120q-32 0-56-24t-24-56v-80q0-7 2-15t4-15l120-282q9-20 30-34t44-14Zm360 80H240L120-480v80h360l-54 220 174-174v-406Zm0 406v-406 406Zm80 34v-80h120v-360H680v-80h200v520H680Z" /></svg>
          <span class="dislike-count" >${otherActivity[i].Post.Dislikes}</span>
        </div>
        <div class="reaction">
          <button class="reaction-button" id="comments">
          <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="currentcolor"><path d="M880-80 720-240H320q-33 0-56.5-23.5T240-320v-40h440q33 0 56.5-23.5T760-440v-280h40q33 0 56.5 23.5T880-640v560ZM160-473l47-47h393v-280H160v327ZM80-280v-520q0-33 23.5-56.5T160-880h440q33 0 56.5 23.5T680-800v280q0 33-23.5 56.5T600-440H240L80-280Zm80-240v-280 280Z" /></svg>
          </button>
        </div>
      </div>
    `;
    container.appendChild(postDiv);
  }
}
