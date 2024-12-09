export function PollForNotifications() {
    //console.log("user role : ", window.userData.role);
    if (!window.userData || window.userData.role.trim().toLowerCase() === "'visitor'") {
        console.log("Role is visitor, stopping execution.");
        return; 
    }

    fetch("https://localhost:8080/notifications")
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
                    const { ID, Type, Sender, Receiver, Post } = notification;

                    const notificationItem = document.createElement("li");

                    let notificationText;
                    if (Type === 'post') {
                        notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span> commented on your post <a href="/?postID=${Post.ParentID}">here.</a>`;
                    } else if (Type === 'LIKE') {
                        notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span> liked <a href="/?postID=${Post.ID}">${Post.Title}.</a>`;
                    } else if (Type === 'DISLIKE') {
                        notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span> disliked <a href="/?postID=${Post.ID}">${Post.Title}.</a>`;
                    } else if (Type === 'category') {
                        notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span> posted ${Post.Title} on <a href="/?catID=${Post.ParentID}">a followed category.</a>`;
                    } else if (Type === 'report') {
                        notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span> reported: <a href="/?postID=${Post.ID}">${Post.Title}.</a>`;
                    } else if (Type === 'request') {
                        notificationText = `<span class="sender" other_id="${Sender.ID}">${Sender.Username}</span> asked to be Moderator.</a>`;
                    } else {
                        notificationText = `Unknown type: ${Type}.`;
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
    console.log(window.userData.id)

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