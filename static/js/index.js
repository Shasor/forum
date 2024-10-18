document.addEventListener("DOMContentLoaded", function () {
  // ╔════════════════════ avatar ════════════════════╗
  const avatar = document.getElementById("avatar");
  const popup = document.getElementById("popup");
  avatar.onclick = function (event) {
    popup.style.display = popup.style.display === "block" ? "none" : "block";
    event.stopPropagation();
  };
  // ╚════════════════════════════════════════════════╝

  // ╔════════════════════ left bar ════════════════════╗
  const leftBar = document.querySelector(".left-bar");
  const toggleButton = leftBar.querySelector("button");
  const postsDiv = document.querySelector(".posts");
  toggleButton.addEventListener("click", function () {
    leftBar.classList.toggle("closed");
    if (leftBar.classList.contains("closed")) {
      toggleButton.textContent = ">>";
      postsDiv.style.marginRight = "-20vw";
    } else {
      toggleButton.textContent = "<<";
      postsDiv.style.marginRight = "0";
    }
  });
  // ╚══════════════════════════════════════════════════╝

  // ╔════════════════════ create-post ════════════════════╗
  const new_post = document.querySelector(".create-post");
  const bttn_new_post = document.querySelector(".create-post-button");
  bttn_new_post.onclick = function (event) {
    new_post.style.display = new_post.style.display === "block" ? "none" : "block";
    event.stopPropagation();
  };
  document.addEventListener("click", function (event) {
    if (!new_post.contains(event.target) && !bttn_new_post.contains(event.target)) {
      new_post.style.display = "none";
    }
  });
  new_post.onclick = function (event) {
    event.stopPropagation();
  };
  // ╚═════════════════════════════════════════════════════╝
});
