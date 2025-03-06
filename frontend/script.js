// DOM elements
const moviesContainer = document.getElementById('movies-container');
const errorElement = document.getElementById('error');
const loadingElement = document.getElementById('loading');

// API endpoint
const API_URL = '/api/movies';

// Fetch movies from the API
async function fetchMovies() {
    try {
        const response = await fetch(API_URL);
        
        // Check if the request was successful
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching movies:', error);
        showError(error.message);
        return null;
    } finally {
        // Hide loading indicator
        loadingElement.style.display = 'none';
    }
}

// Display error message
function showError(message) {
    errorElement.textContent = `Failed to load movies: ${message}`;
    errorElement.style.display = 'block';
}

// Create a movie card element
function createMovieCard(movie) {
    const movieCard = document.createElement('div');
    movieCard.className = 'movie-card';
    
    // Create image placeholder
    let imageContent = '';
    if (movie.poster_url) {
        imageContent = `<img src="${movie.poster_url}" alt="${movie.title}" class="movie-image">`;
    } else {
        imageContent = `<div class="movie-image placeholder">No Image Available</div>`;
    }
    
    // Format release date if available
    let releaseDate = movie.release_date ? new Date(movie.release_date).toLocaleDateString() : 'Unknown';
    
    movieCard.innerHTML = `
        ${imageContent}
        <div class="movie-details">
            <div class="movie-title">${movie.title}</div>
            <div class="movie-info">
                <strong>Director:</strong> ${movie.director || 'Unknown'}
            </div>
            <div class="movie-info">
                <strong>Release Date:</strong> ${releaseDate}
            </div>
            <div class="movie-info">
                <strong>Duration:</strong> ${movie.duration || 'Unknown'} min
            </div>
        </div>
    `;
    
    return movieCard;
}

// Display movies in the container
function displayMovies(movies) {
    // Clear container first
    moviesContainer.innerHTML = '';
    
    if (!movies || movies.length === 0) {
        moviesContainer.innerHTML = '<p>No movies found.</p>';
        return;
    }
    
    // Create and append movie cards
    movies.forEach(movie => {
        const movieCard = createMovieCard(movie);
        moviesContainer.appendChild(movieCard);
    });
}

// Initialize the application
async function initApp() {
    const movies = await fetchMovies();
    
    if (movies) {
        displayMovies(movies);
    }
}

// Start the application when the DOM is loaded
document.addEventListener('DOMContentLoaded', initApp);