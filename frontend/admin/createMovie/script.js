// Token handling
 const getAuthToken = () => {
    return localStorage.getItem('token');
};

// Image preview
document.getElementById('posterImage').addEventListener('change', function(e) {
    const file = e.target.files[0];
    if (file) {
        const reader = new FileReader();
        const preview = document.getElementById('imagePreview');
        
        reader.onload = function(e) {
            preview.src = e.target.result;
            preview.style.display = 'block';
        };
        
        reader.readAsDataURL(file);
    }
});

// Form submission
document.getElementById('createMovieForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const formData = new FormData(this);
    
    // Add JWT token
    const token = getAuthToken();
    if (!token) {
        alert("You are not logged in. Please log in as an admin user.");
        return;
    }
    
    try {
        const response = await fetch('/api/movies', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`
            },
            body: formData
        });
        
        const responseContainer = document.getElementById('responseContainer');
        const responseData = document.getElementById('responseData');
        responseContainer.style.display = 'block';
        
        if (response.ok) {
            const data = await response.json();
            // responseData.textContent = JSON.stringify(data, null, 2);
            // responseData.className = '';
            // Reset form on success
            document.getElementById('createMovieForm').reset();
            document.getElementById('imagePreview').style.display = 'none';
            window.location.href = '../movies/';
        } else {
            const errorData = await response.json();
            responseData.textContent = `Error (${response.status}): ${JSON.stringify(errorData, null, 2)}`;
            responseData.className = 'error';
        }
    } catch (error) {
        console.error('Error:', error);
        const responseContainer = document.getElementById('responseContainer');
        const responseData = document.getElementById('responseData');
        responseContainer.style.display = 'block';
        responseData.textContent = 'Network error: ' + error.message;
        responseData.className = 'error';
    }
});
