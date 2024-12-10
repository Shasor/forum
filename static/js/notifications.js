export function PollForNotifications() {
  //console.log("user role : ", window.userData.role);
  if (!window.userData || window.userData.role.trim().toLowerCase() === "'visitor'") {
    console.log("Role is visitor, stopping execution.");
    return;
  }

  fetch("https://forum.shasor.fr/notifications")
    .then((response) => response.json())
    .then((data) => {
      //console.log("Notification reçue :", data);

      const notificationList = document?.getElementById("notification-ul");
      if (!notificationList) {
        console.error("Conteneur de notifications introuvable !");
        return;
      }

      notificationList.innerHTML = "";

      if (Array.isArray(data.notifData) && data.notifData.length > 0) {
        data.notifData.forEach((notification) => {
          const { ID, Sort, Sender, Receiver, Post } = notification;

          const notificationItem = document.createElement("li");

          let notificationText;
          if (Sort === "post") {
            notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span> commented on your post <a href="/?postID=${Post.ParentID}">here.</a>`;
          } else if (Sort === "LIKE") {
            notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span> liked : <a href="/?postID=${Post.ID}"> ${Post.Title}.</a>`;
          } else if (Sort === "DISLIKE") {
            notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span> disliked <a href="/?postID=${Post.ID}">${Post.Title}.</a>`;
          } else if (Sort === "category") {
            console.log(Post);
            notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span> posted ${Post.Title} on <a href="/?catID=${Post.Categories[0].ID}">a followed category.</a>`;
          } else if (Sort === "report") {
            notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span> reported: <a href="/?postID=${Post.ID}">${Post.Title}.</a>`;
          } else if (Sort === "request") {
            notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span> asked to be Moderator.</a>`;
          } else if (Sort === "reportdone") {
            notificationText = `The moderators have indeed deleted the message you reported.`;
          } else {
            notificationText = `Unknown Sort: ${Sort}.`;
          }

          notificationItem.innerHTML = notificationText;

          notificationList.insertBefore(notificationItem, notificationList.firstChild);
        });
      } else {
        const noNotificationItem = document.createElement("li");
        noNotificationItem.textContent = "Aucune notification pour le moment.";
        notificationList.appendChild(noNotificationItem);
      }

      //setTimeout(PollForNotifications, 15000);
    })
    .catch((error) => {
      console.error("Erreur lors de la récupération de la notification :", error);
      //setTimeout(PollForNotifications, 15000);
    });
}

export function ClearNotifications() {
  const notificationList = document.getElementById("notification-ul");

  if (!notificationList) {
    console.error("Conteneur de notifications introuvable !");
    return;
  }

  while (notificationList.firstChild) {
    notificationList.removeChild(notificationList.firstChild);
  }

  console.log("Toutes les notifications ont été supprimées de l'interface utilisateur.");
  console.log(window.userData.id);

  fetch("https://forum.shasor.fr/notifications/clear", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      userID: String(window.userData.id),
    }),
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error("Erreur lors de l'envoi de la mise à jour au serveur.");
      }
      return response.json();
    })
    .then((data) => {
      console.log("Réponse du serveur après suppression des notifications :", data);
    })
    .catch((error) => {
      console.error("Erreur lors de la suppression des notifications sur le serveur :", error);
    });
}
