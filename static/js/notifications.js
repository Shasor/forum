export function PollForNotifications() {

console.log("user role : ", window.userData.role)
if (!window.userData || window.userData.role.trim().toLowerCase() === "'visitor'") {
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

export function ClearNotifications() {
    const notificationList = document.getElementById("notification-ul");

    if (!notificationList) {
        console.error("Conteneur de notifications introuvable !");
        return;
    }

    // Supprimer toutes les notifications de la liste
    while (notificationList.firstChild) {
        notificationList.removeChild(notificationList.firstChild);
    }

    console.log("Toutes les notifications ont été supprimées de l'interface utilisateur.");
    console.log(window.userData.id)

    // Envoyer une mise à jour au serveur pour notifier la suppression
    fetch("https://localhost:8080/notifications/clear", {
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