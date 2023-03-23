// Logic slide navbar
const slide = document.querySelector('.slide');
const menuToggle = document.querySelector('.menu-toggle input');
const nav = document.querySelector('.navbar ul');

menuToggle.addEventListener('click', function(){
    nav.classList.toggle('slide');
});
// End logic end navbar

// Logic Input Color
let inputBorderColor1 = document.querySelectorAll(".border-color");
let inputBorderColor2 = document.querySelector(".login-container form");
let input = document.querySelectorAll("input");
let btnSubmit = document.querySelector(".btn-submit");

inputBorderColor1.forEach(function(ibc1) {
    ibc1.addEventListener('click', function(event) {
        
        inputBorderColor1.forEach(function(ibc1) {
            if(ibc1.classList.contains("violet-border")){
                ibc1.classList.remove("violet-border");
            }
        });    
        event.target.classList.add("violet-border");
    });

    btnSubmit.addEventListener('mouseover', function() {
        inputBorderColor1.forEach(function(ibc1) {
            ibc1.classList.remove("violet-border");
        });  
    });  
});

inputBorderColor2.addEventListener('mouseenter', function(){
    inputBorderColor1.forEach(function(ibc1){
        ibc1.classList.add("black-border");
    });
    
    inputBorderColor2.addEventListener('mouseleave', function(){
        inputBorderColor1.forEach(function(ibc1){
            ibc1.classList.remove("black-border");
            ibc1.classList.remove("violet-border");
        });
    });
});
// End logic input color

// Logic message
let message = document.querySelector(".message");

setTimeout(function() {
    message.style.display = "none";
}, 3000)
// End logic message