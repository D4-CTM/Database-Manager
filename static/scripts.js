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

function Collapsible(btn) {
	btn.classList.toggle("active");
	var content = btn.nextElementSibling;
	if (content.style.display === "block") {
	  content.style.display = "none";
	} else {
	  content.style.display = "block";
	}
}

function handleResponse(event, btn) {
	console.log("Handle")
    let xhr = event.detail.xhr;
    let status = xhr.getResponseHeader("HX-Status");
    let message = xhr.getResponseHeader("HX-Message");
    
    if (status != "100") {
        alert(message);
    }

	Collapsible(btn)
}
