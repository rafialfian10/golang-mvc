// Logic slide
const slide = document.querySelector('.slide');
const menuToggle = document.querySelector('.menu-toggle input');
const nav = document.querySelector('.navbar ul');

menuToggle.addEventListener('click', function(){
    nav.classList.toggle('slide');
});
// End logic slide


// Logic message
let message = document.querySelector(".message")
setTimeout(function() {
    message.style.display = "none"
}, 3000)
// End logic message