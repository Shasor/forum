export function PollForNotifications() {

    
if (window.userData.role === "visitor"){
  return
}

    fetch("https://localhost:8080/notifications")
        .then((response) => response.json())
        .then((data) => {
            const notification = data;
            console.log("Notification reçue :", notification);

            // Exemple : afficher la notification dans la page HTML
            const notificationElement = document.createElement("div");
            notificationElement.textContent = notification;
            document.body.appendChild(notificationElement);



            const notificationList = document.getElementById("notification-ul");
            if (notificationList) {
                // Parcourir les notifications reçues
                data.notifData.forEach((notification) => {
                    const { ID, Type, Sender, Receiver, Post } = notification;

                    // Créer un nouvel élément <li> pour la notification
                    const notificationItem = document.createElement("li");

                    // Construire un texte descriptif pour la notification
                    let notificationText;
                    if (Type === 'post') {
                        notificationText = `${Sender.Username} commented on your post ${Post.ParentID}.`;
                    } else if (Type === 'LIKE') {
                        notificationText = `${Sender.Username} liked ${Post.Title}.`;
                    } else if(Type ==='DISLIKE') {
                        notificationText = `${Sender.Username} disliked ${Post.Title}.`;
                    } else if(Type ==='category'){
                        notificationText = `${Sender.Username} posted ${Post.Title} on ${Post.Categories["0"].Name}.`;

                    } else {
                        notificationText = `Unknown type : ${Type}.`;
                    }

                    // Ajouter le texte à l'élément <li>
                    notificationItem.textContent = notificationText;

                    // Ajouter l'élément <li> en haut de la liste
                    notificationList.insertBefore(notificationItem, notificationList.firstChild);
                });
            } else {
                console.error("Conteneur de notifications introuvable !");
            }
            // Relancer le long polling pour recevoir la prochaine notification
            //PollForNotifications();
        })
        .catch((error) => {
            console.error("Erreur lors de la récupération de la notification :", error);
            // Attendre un peu avant de réessayer en cas d'erreur
            //setTimeout(PollForNotifications, 5000);
        });
}