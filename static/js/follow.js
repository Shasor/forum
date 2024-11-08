document.addEventListener("DOMContentLoaded", function () {
  const followButtons = document.querySelectorAll("#follow-button");
  followButtons.forEach((button) => {
    button.addEventListener("click", function () {
      console.log("Le bouton a été cliqué !"); // Vérifier si l'événement est déclenché
      const categorieID = this.getAttribute("data-categorieid");
      console.log("CategorieID:", categorieID);
      // Envoyer la requête POST pour suivre/désuivre la catégorie
      fetch("/follow", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ categorieId: parseInt(categorieID) }),
      })
        .then((response) => response.json())
        .then((data) => {
          console.log("data => ", data);
          // Met à jour l'état de suivi en fonction de la réponse
          if (data.isFollowing) {
            this.innerHTML = '<img src="/static/assets/img/notif-on.svg" alt="Désuivre cette catégorie"/>'; // Changer le bouton en "désuivre"
            location.reload();
          } else {
            this.innerHTML = '<img src="/static/assets/img/notif-off.svg" alt="Suivre cette catégorie"/>'; // Changer le bouton en "suivre"
            location.reload();
          }
        })
        .catch((error) => {
          console.error("Erreur:", error);
        });
    });
  });
});
