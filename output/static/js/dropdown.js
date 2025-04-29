document.addEventListener('DOMContentLoaded', function() {
  const button = document.querySelector('.dropdown-button');
  
  const dropdown = button.parentElement;
  
  // Toggle dropdown when button is clicked/touched
  button.addEventListener('click', function(e) {
    e.preventDefault();
    e.stopPropagation();
    dropdown.classList.toggle('show');
  });
  
  // Prevent clicks/touches inside the dropdown from closing it
  dropdown.addEventListener('click', function(e) {
    e.stopPropagation();
  });
});

// Hide dropdowns when clicking/touching anywhere outside
document.addEventListener('click', function(e) {
  const dropdown = document.querySelector('.dropdown');
  
  if (!dropdown.contains(e.target)) {
    dropdown.classList.remove('show');
  }
});
