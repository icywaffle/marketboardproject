function removeemptyarrays() {
    for (i = 0; i < 10; i++) {
        if (document.getElementsByClassName("matMap")[i].innerHTML === "[]") {
            document.getElementsByClassName("matMap")[i].innerHTML = ""
            $("#materialrecipes li").eq(i).remove();
            i--
        }

    }
}

function externalFunction() {
    document.getElementById("external").innerHTML = "Hello World!!!";
}