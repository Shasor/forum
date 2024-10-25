import { GetLogin, GetSignup, ShowError } from "./auth.js";

// ╔════════════════════ avatar ════════════════════╗
document.addEventListener("DOMContentLoaded", function () {
  const avatar = document.getElementById("avatar");
  const popup = document.getElementById("popup");
  avatar.onclick = function (event) {
    popup.style.display = popup.style.display === "block" ? "none" : "block";
    event.stopPropagation();
  };
});

// ╔════════════════════ left bar ════════════════════╗
document.addEventListener("DOMContentLoaded", function () {
  const leftBar = document.querySelector(".left-bar");
  const toggleButton = leftBar.querySelector("button");
  const postsDiv = document.querySelector(".posts-container");
  toggleButton.addEventListener("click", function () {
    leftBar.classList.toggle("closed");
    if (leftBar.classList.contains("closed")) {
      toggleButton.textContent = ">>";
      postsDiv.style.marginRight = "-20vw";
    } else {
      toggleButton.textContent = "<<";
      postsDiv.style.marginRight = "0vw";
    }
  });
});

// ╔════════════════════ create-post ════════════════════╗
document.addEventListener("DOMContentLoaded", function () {
  const new_post = document.querySelector(".create-post");
  const bttn_new_post = document.querySelector(".create-post-button");

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
});

// ╔════════════════════ Login-Signup ════════════════════╗
document.addEventListener("DOMContentLoaded", function () {
  const login_link = document.getElementById("header-login-link");
  login_link?.addEventListener("click", GetLogin);

  const signup_button = document.getElementById("header-signup-link");
  signup_button?.addEventListener("click", GetSignup);
});

document.addEventListener("DOMContentLoaded", function () {
  const reactionForms = document.querySelectorAll("form.reaction-form");

  reactionForms.forEach((form) => {
    form.addEventListener("submit", function (event) {
      event.preventDefault();

      const postId = this.querySelector('input[name="postId"]').value;
      const reaction = this.querySelector('input[name="reaction"]').value;
      const errorCount = 0;

      fetch("/react", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ postId: parseInt(postId), reaction }),
      })
        .then((response) => {
          // Vérifier si la réponse contient du contenu et est de type JSON
          const contentType = response.headers.get("content-type");
          if (contentType && contentType.includes("application/json")) {
            return response.json();
          }
          // Si ce n'est pas du JSON, retourner null ou une promesse résolue
          return null;
        })
        .then((data) => {
          // Vérifier si des données JSON ont été reçues
          if (data) {
            // Mettre à jour les compteurs dans le HTML
            const likeCountSpan = document.querySelector(`.like-count[data-postid="${postId}"]`);
            const dislikeCountSpan = document.querySelector(`.dislike-count[data-postid="${postId}"]`);

            if (likeCountSpan) likeCountSpan.textContent = data.likes;
            if (dislikeCountSpan) dislikeCountSpan.textContent = data.dislikes;
          } else {
            ShowError("You are not connected");
          }
        })
        .catch((error) => {
          console.error("Error:", error);
        });
    });
  });
});
