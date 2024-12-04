export function PollForNotifications() {

console.log("user role : ", window.userData.role)
if (!window.userData || window.userData.role.trim().toLowerCase() === 'visitor') {
    console.log()
    console.log("Role is visitor, stopping execution.");
    return; // Arrêter si le rôle est 'visitor'
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
                        notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span>commented on your post<a href="/?postID=${Post.ParentID}">here.</a>`;
                    } else if (Type === 'LIKE') {
                        notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span>liked <a href="/?postID=${Post.ID}">${Post.Title}.</a>`;
                    } else if(Type ==='DISLIKE') {
                        notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span>disliked<a href="/?postID=${Post.ID}">${Post.Title}.</a>`;
                    } else if(Type ==='category'){
                        notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span>posted ${Post.Title} on <a href="/?catID=${Post.ParentID}">a followed category.</a>`;
                    }else if (Type === 'report'){
                        notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span>reported : <a href="/?postID=${Post.ID}">${Post.Title}.</a>`;
                    }else if (Type === 'request'){
                        notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span>ask to be Moderator.</a>`; 
                    } else {
                        notificationText = `Unknown type : ${Type}.`;
                    }

                    // Ajouter le texte à l'élément <li>
                    notificationItem.innerHTML = notificationText;

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
