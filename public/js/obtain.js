function removeemptyarrays() {
    var listlength = 10
    for (i = 0; i < listlength; i++) {
        if (document.getElementsByClassName("matMap")[i].innerHTML === "[]") {
            document.getElementsByClassName("matMap")[i].innerHTML = ""
            $("#materialrecipes li").eq(i).remove();
            // When we remove an element from the list, the index and length will change.
            i--
            listlength--
        }

    }
}

function externalFunction() {
    document.getElementById("external").innerHTML = "Hello World!!!";
}