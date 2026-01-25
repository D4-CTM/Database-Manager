var modal = document.getElementById("myModal");

var btn = document.getElementById("CreateConnectionBTN");

var span = document.getElementById("CloseConnectionFormBTN");

btn.onclick = function() {
  modal.style.display = "block";
}

span.onclick = function() {
  modal.style.display = "none";
}

window.onclick = function(event) {
  if (event.target == modal) {
    modal.style.display = "none";
  }
} 
