function isAuthenticated() {
    const token = localStorage.getItem('token');
    if (!token) {
      return false;
    }
    return true;
  }

  function protectRoute() {
    if (!isAuthenticated()) {
      // Redirect to login page
      window.location.href = '/';
    }
  }

  function logout() {
    localStorage.removeItem('authToken');
    window.location.href = '/login.html';
  }