import { ShowError } from "./other.js";
import { GetLogin, GetSignup } from "./auth.js";
import { GetProfile, GetEditProfile, GetMyPosts, GetLikedPosts } from "./profile.js";
import { GetOtherProfile } from "./otherProfile.js";
import { PollForNotifications, ClearNotifications } from "./notifications.js";

PollForNotifications();

// ╔════════════════════ avatar ════════════════════╗
const avatar = document.getElementById("avatar");
avatar.addEventListener("click", ToggleAvatar);

function ToggleAvatar(event) {
  const popup = document.getElementById("popup");
  popup.style.display = popup.style.display === "block" ? "none" : "block";
  event.stopPropagation();

  document.addEventListener("click", function (event) {
    if (!popup.contains(event.target)) {
      popup.style.display = "none";
    }
  });
}

// ╔════════════════════ profile access ════════════════════╗

const profile_link = document.getElementById("header-profile-link");
profile_link?.addEventListener("click", (event) => {
  event.preventDefault(); // Prevent default link navigation
  GetProfile();
});

const edit_profile_link = document.getElementById("header-profile-edit");
edit_profile_link?.addEventListener("click", (event) => {
  event.preventDefault();
  GetEditProfile();
});

const posts_link = document.getElementById("header-profile-posts");
posts_link?.addEventListener("click", (event) => {
  event.preventDefault();
  GetMyPosts();
});

const liked_posts_link = document.getElementById("header-profile-liked");
liked_posts_link?.addEventListener("click", (event) => {
  GetLikedPosts();
});

// ╔════════════════════ other profile access ════════════════════╗
const other_profile_links = document.getElementsByClassName("sender");
Array.from(other_profile_links).forEach((link) => {
  link.addEventListener("click", (event) => {
    GetOtherProfile();
  });
});

const notificationList = document.getElementById("notification-ul");
if (notificationList) {
  notificationList.addEventListener("click", (event) => {
    // Vérifiez si l'élément cliqué a la classe "sender"
    if (event.target.classList.contains("sender")) {
      GetOtherProfile();
      event.preventDefault(); // Facultatif, selon votre logique
    }
  });
}

// ╔════════════════════ left bar ════════════════════╗
const leftBar = document.querySelector(".left-bar");
const toggleButton = leftBar.querySelector("#logo");
const postsDiv = document.querySelector(".posts-container");
toggleButton.addEventListener("click", () => {
  leftBar.classList.toggle("closed");
  if (leftBar.classList.contains("closed")) {
    postsDiv.style.marginRight = "-20vw";
  } else {
    postsDiv.style.marginRight = "0vw";
  }
});

// ╔════════════════════ create-post ════════════════════╗
const new_post = document.querySelector(".create-post");
const bttn_new_post = document.querySelector(".create-post-button");
bttn_new_post?.addEventListener("click", CreatePost(new_post, bttn_new_post));

function CreatePost(new_post, bttn_new_post) {
  // Vérifier si les éléments existent avant d'ajouter des événements
  if (bttn_new_post && new_post) {
    bttn_new_post.onclick = function (event) {
      new_post.style.display = new_post.style.display === "block" ? "none" : "block";
      event.stopPropagation();
    };

    new_post.onclick = function (event) {
      event.stopPropagation();
    };

    // Ajouter l'écouteur de clic global seulement si les éléments existent
    document.addEventListener("click", function (event) {
      if (!new_post.contains(event.target) && !bttn_new_post.contains(event.target)) {
        new_post.style.display = "none";
      }
    });
  }
}

// ╔════════════════════ modify post ════════════════════╗

const buttons = document.querySelectorAll(".create-post-to-modify-button");

// Ajouter des écouteurs d'événements à chaque bouton
buttons.forEach((button) => {
  const modifyPostDiv = document?.querySelector(`#${button.dataset.target}`); // Cible la div associée via un attribut data-target

  button.addEventListener("click", function (event) {
    event.stopPropagation(); // Empêche le clic de se propager au document
    if (modifyPostDiv.style.display === "none" || modifyPostDiv.style.display === "") {
      modifyPostDiv.style.display = "block";
    } else {
      modifyPostDiv.style.display = "none";
    }
  });

  // Fermer la div si un clic se produit en dehors
  document.addEventListener("click", function (event) {
    if (!modifyPostDiv.contains(event.target) && event.target !== button) {
      modifyPostDiv.style.display = "none";
    }
  });
});

// ╔════════════════════ Notification Poppup ════════════════════╗

document.addEventListener("DOMContentLoaded", () => {
  const notifIcon = document?.getElementById("notif-pic-none"); // Utilisez l'identifiant correct
  const notifContainer = document?.getElementById("notification-container"); // Utilisez l'identifiant correct

  notifIcon?.addEventListener("click", (event) => {
    // Empêcher l'événement de se propager pour éviter de masquer immédiatement
    event.stopPropagation();

    // Basculer la visibilité
    notifContainer.style.display = notifContainer.style.display === "none" || !notifContainer.style.display ? "block" : "none";
  });

  // Cacher la notification si on clique ailleurs
  document.addEventListener("click", (event) => {
    // Vérifiez si les éléments existent avant de continuer
    if (notifIcon && notifContainer) {
      // Si le clic n'est ni sur l'icône ni sur le conteneur, masquer le conteneur
      if (!notifIcon.contains(event.target) && !notifContainer.contains(event.target)) {
        notifContainer.style.display = "none";
      }
    }
  });
});

// ╔════════════════════ Notification Poppup Delete ════════════════════╗

document.addEventListener("DOMContentLoaded", function () {
  const delete_post_link = document.getElementById("delete-post");
  delete_post_link?.addEventListener("click", ClearNotifications);
});

// ╔════════════════════ Login-Signup ════════════════════╗
document.addEventListener("DOMContentLoaded", function () {
  const login_link = document.getElementById("header-login-link");
  login_link?.addEventListener("click", GetLogin);

  const signup_button = document.getElementById("header-signup-link");
  signup_button?.addEventListener("click", GetSignup);
});

const reactionForms = document.querySelectorAll("form.reaction-form");
reactionForms.forEach((form) => {
  form.addEventListener("submit", async function (event) {
    event.preventDefault();

    const postId = this.querySelector('input[name="postId"]').value;
    const reaction = this.querySelector('input[name="reaction"]').value;

    try {
      const data = await UpdateReaction(postId, reaction);

      // Update the HTML counters with the new data
      const likeCountSpan = document.querySelector(`.like-count[data-postid="${postId}"]`);
      const dislikeCountSpan = document.querySelector(`.dislike-count[data-postid="${postId}"]`);

      if (likeCountSpan) likeCountSpan.textContent = data.likes;
      if (dislikeCountSpan) dislikeCountSpan.textContent = data.dislikes;
    } catch (error) {
      console.error("Error:", error);
      ShowError("You are not connected");
    }
  });
});

async function UpdateReaction(postId, reaction) {
  const response = await fetch("/react", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ postId: parseInt(postId), reaction }),
  });

  if (!response.ok) {
    const errorMessage = await response.text();
    throw new Error(`Error: ${response.status} - ${errorMessage}`);
  }

  // Ensure JSON response before returning
  // location.reload();
  const contentType = response.headers.get("content-type");
  if (contentType && contentType.includes("application/json")) {
    return response.json();
  } else {
    throw new Error("Invalid JSON response");
  }
}

// ╔════════════════════ search bar ════════════════════╗
document.getElementById("search_bar").addEventListener("keypress", function (event) {
  if (event.key === "Enter") {
    event.preventDefault();
    document.getElementById("SearchForm").submit();
  }
});
