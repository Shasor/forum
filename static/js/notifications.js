export function PollForNotifications() {

    
if (window.userData.role === "visitor"){
  return
}

    fetch("https://localhost:8080/notifications")
        .then((response) => response.json())
        .then((data) => {
            const notification = data.message;
            console.log("Notification reçue :", notification);

            // Exemple : afficher la notification dans la page HTML
            const notificationElement = document.createElement("div");
            notificationElement.textContent = notification;
            document.body.appendChild(notificationElement);

            // Relancer le long polling pour recevoir la prochaine notification
            //PollForNotifications();
        })
        .catch((error) => {
            console.error("Erreur lors de la récupération de la notification :", error);
            // Attendre un peu avant de réessayer en cas d'erreur
            //setTimeout(PollForNotifications, 5000);
        });
}