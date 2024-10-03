const avatar = document.getElementById('avatar');
const popup = document.getElementById('popup');
const new_post = document.getElementById('new_post');
const bttn_new_post = document.getElementById('create_post');

document.addEventListener('DOMContentLoaded', function() {
    // Afficher ou cacher le pop-up lorsque l'avatar est cliqué
    avatar.onclick = function(event) {
        popup.style.display = popup.style.display === 'block' ? 'none' : 'block';
        event.stopPropagation(); // Empêcher la propagation de l'événement
    };

    // Cacher le pop-up lorsque l'on clique en dehors
    document.addEventListener('click', function() {
        popup.style.display = 'none';
    });

    // Afficher ou cacher le formulaire de création de post
    bttn_new_post.onclick = function(event) {
        new_post.style.display = new_post.style.display === 'block' ? 'none' : 'block';
        event.stopPropagation(); // Empêcher la propagation de l'événement
    };

    // Cacher le formulaire lorsque l'on clique en dehors
    document.addEventListener('click', function(event) {
        // Vérifier si le clic a été fait à l'extérieur du formulaire
        if (!new_post.contains(event.target) && !bttn_new_post.contains(event.target)) {
            new_post.style.display = 'none';
        }
    });

    // Cacher le pop-up lorsque l'on clique à l'intérieur du formulaire
    new_post.onclick = function(event) {
        event.stopPropagation(); // Empêcher la propagation de l'événement
    };
});


document.addEventListener("DOMContentLoaded", function() {
    const settingsPanel = document.getElementById("settings_panel");
    const toggleButton = document.getElementById("toggle_settings");
    const sidebar = document.getElementById("side_bar");
    
    toggleButton.addEventListener("click", function() {
        settingsPanel.classList.toggle("collapsed");
        sidebar.classList.toggle("collapsed");
        toggleButton.textContent = settingsPanel.classList.contains("collapsed") ? ">>" : "<<";
    });
});