// editprofile.js
document.addEventListener("DOMContentLoaded", function() {
    const profilePictureInput = document.getElementById("profile_picture");
    const form = document.querySelector("form"); // Select the form
    const errorMessageContainer = document.getElementById("error-message");

    profilePictureInput.addEventListener("change", function(event) {
        const file = event.target.files[0];
        errorMessageContainer.textContent = ""; // Clear previous messages
        if (file) {
            const img = new Image();
            const objectUrl = URL.createObjectURL(file);
            img.src = objectUrl;

            img.onload = function() {
                // Check dimensions
                const maxWidth = 150; // Max width in pixels
                const maxHeight = 150; // Max height in pixels
                if (img.width > maxWidth || img.height > maxHeight) {
                    errorMessageContainer.textContent = `Image dimensions exceed the maximum allowed size of ${maxWidth}x${maxHeight} pixels.`;
                    profilePictureInput.value = ""; // Clear the input
                }
                URL.revokeObjectURL(objectUrl); // Clean up
            };
        }
    });

    // Optional: Form submission can be prevented if needed
    form.addEventListener("submit", function(event) {
        const file = profilePictureInput.files[0];
        errorMessageContainer.textContent = ""; // Clear previous messages
        if (file) {
            const img = new Image();
            const objectUrl = URL.createObjectURL(file);
            img.src = objectUrl;

            img.onload = function() {
                if (img.width > maxWidth || img.height > maxHeight) {
                    event.preventDefault(); // Prevent form submission
                    errorMessageContainer.textContent = `Please upload an image smaller than ${maxWidth}x${maxHeight} pixels.`;
                }
                URL.revokeObjectURL(objectUrl); // Clean up
            };
        }
    });
});
