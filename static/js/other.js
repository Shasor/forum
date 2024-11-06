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
