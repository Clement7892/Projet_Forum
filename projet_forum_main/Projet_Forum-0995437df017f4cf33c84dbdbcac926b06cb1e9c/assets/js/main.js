function verif() {
    // pour récuperer la valeur du HTML
    const name = document.getElementById("last_name-input").value;
    console.log(name);
    if (name.length == 0) {
        alert("le nom ne doit pas etre vide !!");
        // document.getElementById("last_name-input").value = ""
        return false;
    }

}