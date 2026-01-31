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

function OpenFolder(btn) {	
	let sourceFolderIcon = "fa-folder"
	let targetFolderIcon = "fa-folder-open-o"
	if (!btn.classList.contains("active")) {
		sourceFolderIcon = "fa-folder-open-o"
		targetFolderIcon = "fa-folder"
	}

	var lbl = btn.childNodes[1].childNodes[0];
	console.log(lbl)
	lbl.classList.remove(sourceFolderIcon);
	lbl.classList.add(targetFolderIcon);
}

function handleResponse(event, btn) {
	let xhr = event.detail.xhr;
	if (xhr.status == 200) {
		return
	}

	let message = xhr.getResponseHeader("HX-Message");

	alert(message);

	if (btn.classList.contains("Collapsible")) {
		Collapsible(btn);
	}
}

function IsCollapsed(btn, event) {
	if (!btn.classList.contains("active")) {
		event.preventDefault();
	}
}
